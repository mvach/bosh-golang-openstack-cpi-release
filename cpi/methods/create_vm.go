package methods

import (
	"github.com/cloudfoundry/bosh-cpi-go/apiv1"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

type CreateVMMethod struct {
	logger boshlog.Logger
}

func NewCreateVMMethod(logger boshlog.Logger) CreateVMMethod {
	return CreateVMMethod{logger}
}

func (method CreateVMMethod) CreateVM(
	agentID apiv1.AgentID, stemcellCID apiv1.StemcellCID, cloudProps apiv1.VMCloudProps,
	networks apiv1.Networks, diskCIDs []apiv1.DiskCID, env apiv1.VMEnv) (apiv1.VMCID, error) {

	return apiv1.VMCID{}, nil
}

func (method CreateVMMethod) CreateVMV2(
	agentID apiv1.AgentID, stemcellCID apiv1.StemcellCID, cloudProps apiv1.VMCloudProps,
	networks apiv1.Networks, diskCIDs []apiv1.DiskCID, env apiv1.VMEnv) (apiv1.VMCID, apiv1.Networks, error) {

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "https://identity.../v3",
		Username:         "the_user",
		Password:         "the_password",
		DomainName:       "the_domain",
		TenantName:       "the_project",
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return apiv1.VMCID{}, apiv1.Networks{}, err
	}

	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "the_region",
	})
	if err != nil {
		return apiv1.VMCID{}, apiv1.Networks{}, err
	}

	server, err := servers.Create(client, servers.CreateOpts{
		Name:      "My new server!",
		FlavorRef: "22",
		ImageRef:  "5bba0da5-dfb3-49d8-a005-d799507518f7",
		Networks:  []servers.Network{{UUID: "577a35e5-6b5f-47fa-b8ac-0d38b337f358"}},
	}).Extract()
	if err != nil {
		return apiv1.VMCID{}, apiv1.Networks{}, err
	}

	return apiv1.NewVMCID(server.ID), apiv1.Networks{}, nil
}
