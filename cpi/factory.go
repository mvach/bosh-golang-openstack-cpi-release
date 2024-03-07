package cpi

import (
	"github.com/cloudfoundry/bosh-cpi-go/apiv1"
	"github.com/cloudfoundry/bosh-golang-openstack-cpi-go/config"
	"github.com/cloudfoundry/bosh-golang-openstack-cpi-go/cpi/methods"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
)

type Factory struct {
	fs              boshsys.FileSystem
	uuidGen         boshuuid.Generator
	openstackConfig config.OpenstackConfig
	logger          boshlog.Logger
}

type CPI struct {
	methods.InfoMethod

	methods.CreateStemcellMethod
	methods.DeleteStemcellMethod

	methods.CreateVMMethod
	methods.DeleteVMMethod
	methods.CalculateVMCloudPropertiesMethod
	methods.HasVMMethod
	methods.RebootVMMethod
	methods.SetVMMetadataMethod
	methods.GetDisksMethod

	methods.CreateDiskMethod
	methods.DeleteDiskMethod
	methods.AttachDiskMethod
	methods.DetachDiskMethod
	methods.HasDiskMethod
	methods.ResizeDiskMethod
	methods.SetDiskMetadataMethod

	methods.DeleteSnapshotMethod
	methods.SnapshotDiskMethod
}

func NewFactory(
	fs boshsys.FileSystem,
	uuidGen boshuuid.Generator,
	openstackConfig config.OpenstackConfig,
	logger boshlog.Logger,
) Factory {
	return Factory{fs, uuidGen, openstackConfig, logger}
}

func (cpiFactory Factory) New(ctx apiv1.CallContext) (apiv1.CPI, error) {
	_, err := cpiFactory.getValidatedConfig(ctx)
	if err != nil {
		return CPI{}, err
	}

	return CPI{
		methods.NewInfoMethod(),

		methods.NewCreateStemcellMethod(),
		methods.NewDeleteStemcellMethod(),

		methods.NewCreateVMMethod(cpiFactory.openstackConfig, cpiFactory.logger),
		methods.NewDeleteVMMethod(),
		methods.NewCalculateVMCloudPropertiesMethod(),
		methods.NewHasVMMethod(),
		methods.NewRebootVMMethod(),
		methods.NewSetVMMetadataMethod(),
		methods.NewGetDisksMethod(),

		methods.NewCreateDiskMethod(),
		methods.NewDeleteDiskMethod(),
		methods.NewAttachDiskMethod(),
		methods.NewDetachDiskMethod(),
		methods.NewHasDiskMethod(),
		methods.NewResizeDiskMethod(),
		methods.NewSetDiskMetadataMethod(),
		methods.NewDeleteSnapshotMethod(),
		methods.NewSnapshotDiskMethod(),
	}, nil
}

func (cpiFactory Factory) getValidatedConfig(ctx apiv1.CallContext) (config.OpenstackConfig, error) {
	var validateConfig config.OpenstackConfig

	err := ctx.As(&validateConfig)
	if err != nil {
		return config.OpenstackConfig{}, bosherr.WrapError(err, "Parsing CPI context")
	}

	err = validateConfig.Validate()
	if err != nil {
		return config.OpenstackConfig{}, bosherr.WrapError(err, "Validating CPI context")
	}

	return validateConfig, nil
}
