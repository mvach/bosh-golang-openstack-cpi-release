package methods

import (
	"github.com/cloudfoundry/bosh-cpi-go/apiv1"
	"github.com/cloudfoundry/bosh-golang-openstack-cpi-go/config"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/google/uuid"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/pagination"
)

type CreateVMMethod struct {
	config config.OpenstackConfig
	logger boshlog.Logger
}

func NewCreateVMMethod(openstackConfig config.OpenstackConfig, logger boshlog.Logger) CreateVMMethod {
	return CreateVMMethod{openstackConfig, logger}
}

func (method CreateVMMethod) CreateVM(
	agentID apiv1.AgentID, stemcellCID apiv1.StemcellCID, cloudProps apiv1.VMCloudProps,
	networks apiv1.Networks, diskCIDs []apiv1.DiskCID, env apiv1.VMEnv) (apiv1.VMCID, error) {

	return apiv1.VMCID{}, nil
}

type VmCloudProps struct {
	AvailabilityZone string `json:"availability_zone"`
	EphemeralDisk    string `json:"ephemeral_disk"`
	InstanceType     string `json:"instance_type"`
}

type NetworkCloudProps struct {
	NetID          string   `json:"net_id"`
	SecurityGroups []string `json:"security_groups"`
}

func (method CreateVMMethod) CreateVMV2(
	agentID apiv1.AgentID, stemcellCID apiv1.StemcellCID, cloudProps apiv1.VMCloudProps,
	networks apiv1.Networks, diskCIDs []apiv1.DiskCID, env apiv1.VMEnv) (apiv1.VMCID, apiv1.Networks, error) {

	var vmCloudProps = VmCloudProps{}
	cloudProps.As(&vmCloudProps)

	var networkCloudProps = NetworkCloudProps{}
	networks.Default().CloudProps().As(&networkCloudProps)

	method.logger.Error("create_vm", ">>>>> agentID: %v", agentID)
	method.logger.Error("create_vm", ">>>>> stemcellCID: %v", stemcellCID)
	method.logger.Error("create_vm", ">>>>> cloudProps: %v", vmCloudProps)
	method.logger.Error("create_vm", ">>>>> networks: %v", networkCloudProps.NetID)
	method.logger.Error("create_vm", ">>>>> IP: %v", networks.Default().IP())
	method.logger.Error("create_vm", ">>>>> diskCIDs: %v", diskCIDs)
	method.logger.Error("create_vm", ">>>>> env: %v", env)

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: method.config.AuthURL,
		Username:         method.config.Username,
		Password:         method.config.APIKey,
		DomainName:       method.config.DomainName,
		TenantName:       method.config.Tenant,
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return apiv1.VMCID{}, apiv1.Networks{}, err
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{Region: method.config.Region})
	if err != nil {
		return apiv1.VMCID{}, apiv1.Networks{}, err
	}

	var flavorRef string
	flavorsDetails := flavors.ListDetail(client, nil)
	flavorsDetails.EachPage(func(page pagination.Page) (bool, error) {
		pageElements, err := flavors.ExtractFlavors(page)
		if err != nil {
			return false, err
		}

		for _, element := range pageElements {
			if element.Name == vmCloudProps.InstanceType {
				flavorRef = element.ID
			}
		}
		return true, nil
	})

	serverCreateOpts := servers.CreateOpts{
		Name:             "vm-" + uuid.New().String(),
		ImageRef:         stemcellCID.AsString(),
		Networks:         []servers.Network{{UUID: networkCloudProps.NetID, FixedIP: networks.Default().IP()}},
		SecurityGroups:   networkCloudProps.SecurityGroups,
		AvailabilityZone: vmCloudProps.AvailabilityZone,
		FlavorRef:        flavorRef,
	}

	createOpts := keypairs.CreateOptsExt{
		CreateOptsBuilder: serverCreateOpts,
		KeyName:           method.config.DefaultKeyName,
	}

	server, err := servers.Create(client, createOpts).Extract()
	if err != nil {
		return apiv1.VMCID{}, apiv1.Networks{}, err
	}

	return apiv1.NewVMCID(server.ID), apiv1.Networks{}, nil
}
