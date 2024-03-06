package config

type CpiConfig struct {
	Openstack OpenstackConfig `json:"openstack"`
}

type OpenstackConfig struct {
	AuthURL                      string   `json:"auth_url"`
	Username                     string   `json:"username"`
	APIKey                       string   `json:"api_key"`
	Region                       string   `json:"region"`
	EndpointType                 string   `json:"endpoint_type"`
	DefaultKeyName               string   `json:"default_key_name"`
	DefaultSecurityGroups        []string `json:"default_security_groups"`
	DefaultVolumeType            string   `json:"default_volume_type"`
	WaitResourcePollInterval     int      `json:"wait_resource_poll_interval"`
	BootFromVolume               string   `json:"boot_from_volume"`
	ConfigDrive                  string   `json:"config_drive"`
	UseDhcp                      string   `json:"use_dhcp"`
	IgnoreServerAvailabilityZone bool     `json:"ignore_server_availability_zone"`
	HumanReadableVMNames         string   `json:"human_readable_vm_names"`
	UseNovaNetworking            string   `json:"use_nova_networking"`
	ConnectionOptions            string   `json:"connection_options"`
	DomainName                   string   `json:"domain"`
	ProjectName                  string   `json:"project"`
	Tenant                       string   `json:"tenant"`
	VM                           struct {
		Stemcell struct {
			APIVersion int `json:"api_version"`
		} `json:"stemcell"`
	} `json:"vm"`
}

func (cpiConfig CpiConfig) Validate() error {
	// do validation here

	return nil
}

func (openstackConfig OpenstackConfig) Validate() error {
	// do validation here

	return nil
}
