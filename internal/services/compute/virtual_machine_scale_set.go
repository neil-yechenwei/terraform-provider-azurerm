package compute

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/virtualmachinescalesets"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-11-01/compute"
	identity "github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func VirtualMachineScaleSetAdditionalCapabilitiesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				// NOTE: requires registration to use:
				// $ az feature show --namespace Microsoft.Compute --name UltraSSDWithVMSS
				// $ az provider register -n Microsoft.Compute
				"ultra_ssd_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
					ForceNew: true,
				},
			},
		},
	}
}

func ExpandVirtualMachineScaleSetAdditionalCapabilities(input []interface{}) *virtualmachinescalesets.AdditionalCapabilities {
	capabilities := virtualmachinescalesets.AdditionalCapabilities{}

	if len(input) > 0 {
		raw := input[0].(map[string]interface{})

		capabilities.UltraSSDEnabled = utils.Bool(raw["ultra_ssd_enabled"].(bool))
	}

	return &capabilities
}

func FlattenVirtualMachineScaleSetAdditionalCapabilities(input *virtualmachinescalesets.AdditionalCapabilities) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	ultraSsdEnabled := false

	if input.UltraSSDEnabled != nil {
		ultraSsdEnabled = *input.UltraSSDEnabled
	}

	return []interface{}{
		map[string]interface{}{
			"ultra_ssd_enabled": ultraSsdEnabled,
		},
	}
}

func expandVirtualMachineScaleSetIdentity(input []interface{}) (*identity.SystemAndUserAssignedMap, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}
	out := identity.SystemAndUserAssignedMap{
		Type: expanded.Type,
	}
	if expanded.Type == identity.TypeUserAssigned || expanded.Type == identity.TypeSystemAssignedUserAssigned {
		out.IdentityIds = make(map[string]identity.UserAssignedIdentityDetails)
		for k := range expanded.IdentityIds {
			out.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				// intentionally empty
			}
		}
	}

	return &out, nil
}

func flattenVirtualMachineScaleSetIdentity(input *identity.SystemAndUserAssignedMap) (*[]interface{}, error) {
	var transform *identity.SystemAndUserAssignedMap

	if input != nil {
		transform = &identity.SystemAndUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		if input.PrincipalId != "" {
			transform.PrincipalId = input.PrincipalId
		}
		if input.TenantId != "" {
			transform.TenantId = input.TenantId
		}
		for k, v := range input.IdentityIds {
			transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				ClientId:    v.ClientId,
				PrincipalId: v.PrincipalId,
			}
		}
	}

	return identity.FlattenSystemAndUserAssignedMap(transform)
}

func flattenOrchestratedVirtualMachineScaleSetIdentity(input *compute.VirtualMachineScaleSetIdentity) (*[]interface{}, error) {
	var transform *identity.UserAssignedMap

	if input != nil {
		transform = &identity.UserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		for k, v := range input.UserAssignedIdentities {
			transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				ClientId:    v.ClientID,
				PrincipalId: v.PrincipalID,
			}
		}
	}

	return identity.FlattenUserAssignedMap(transform)
}

func VirtualMachineScaleSetNetworkInterfaceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"ip_configuration": virtualMachineScaleSetIPConfigurationSchema(),

				"dns_servers": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
				// TODO 4.0: change this from enable_* to *_enabled
				"enable_accelerated_networking": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
				// TODO 4.0: change this from enable_* to *_enabled
				"enable_ip_forwarding": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
				"network_security_group_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: azure.ValidateResourceIDOrEmpty,
				},
				"primary": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func VirtualMachineScaleSetNetworkInterfaceSchemaForDataSource() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"ip_configuration": virtualMachineScaleSetIPConfigurationSchemaForDataSource(),

				"dns_servers": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},
				// TODO 4.0: change this from enable_* to *_enabled
				"enable_accelerated_networking": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
				// TODO 4.0: change this from enable_* to *_enabled
				"enable_ip_forwarding": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
				"network_security_group_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"primary": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
			},
		},
	}
}

func virtualMachineScaleSetIPConfigurationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				// Optional
				"application_gateway_backend_address_pool_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Set:      pluginsdk.HashString,
				},

				"application_security_group_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: azure.ValidateResourceID,
					},
					Set:      pluginsdk.HashString,
					MaxItems: 20,
				},

				"load_balancer_backend_address_pool_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Set:      pluginsdk.HashString,
				},

				"load_balancer_inbound_nat_rules_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Set:      pluginsdk.HashString,
				},

				"primary": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"public_ip_address": virtualMachineScaleSetPublicIPAddressSchema(),

				"subnet_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: azure.ValidateResourceID,
				},

				"version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(compute.IPVersionIPv4),
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.IPVersionIPv4),
						string(compute.IPVersionIPv6),
					}, false),
				},
			},
		},
	}
}

func virtualMachineScaleSetIPConfigurationSchemaForDataSource() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"application_gateway_backend_address_pool_ids": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"application_security_group_ids": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"load_balancer_backend_address_pool_ids": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"load_balancer_inbound_nat_rules_ids": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"primary": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"public_ip_address": virtualMachineScaleSetPublicIPAddressSchemaForDataSource(),

				"subnet_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func virtualMachineScaleSetPublicIPAddressSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				// Optional
				"domain_name_label": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"idle_timeout_in_minutes": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(4, 32),
				},
				"ip_tag": {
					// TODO: does this want to be a Set?
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"tag": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"type": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
				// TODO: preview feature
				// $ az feature register --namespace Microsoft.Network --name AllowBringYourOwnPublicIpAddress
				// $ az provider register -n Microsoft.Network
				"public_ip_prefix_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: azure.ValidateResourceIDOrEmpty,
				},
			},
		},
	}
}

func virtualMachineScaleSetPublicIPAddressSchemaForDataSource() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"domain_name_label": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"idle_timeout_in_minutes": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"ip_tag": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"tag": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
							"type": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},

				"public_ip_prefix_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func ExpandVirtualMachineScaleSetNetworkInterface(input []interface{}) (*[]virtualmachinescalesets.VirtualMachineScaleSetNetworkConfiguration, error) {
	output := make([]virtualmachinescalesets.VirtualMachineScaleSetNetworkConfiguration, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		dnsServers := utils.ExpandStringSlice(raw["dns_servers"].([]interface{}))

		ipConfigurations := make([]virtualmachinescalesets.VirtualMachineScaleSetIPConfiguration, 0)
		ipConfigurationsRaw := raw["ip_configuration"].([]interface{})
		for _, configV := range ipConfigurationsRaw {
			configRaw := configV.(map[string]interface{})
			ipConfiguration, err := expandVirtualMachineScaleSetIPConfiguration(configRaw)
			if err != nil {
				return nil, err
			}

			ipConfigurations = append(ipConfigurations, *ipConfiguration)
		}

		config := virtualmachinescalesets.VirtualMachineScaleSetNetworkConfiguration{
			Name: raw["name"].(string),
			Properties: &virtualmachinescalesets.VirtualMachineScaleSetNetworkConfigurationProperties{
				DnsSettings: &virtualmachinescalesets.VirtualMachineScaleSetNetworkConfigurationDnsSettings{
					DnsServers: dnsServers,
				},
				EnableAcceleratedNetworking: utils.Bool(raw["enable_accelerated_networking"].(bool)),
				EnableIPForwarding:          utils.Bool(raw["enable_ip_forwarding"].(bool)),
				IPConfigurations:            ipConfigurations,
				Primary:                     utils.Bool(raw["primary"].(bool)),
			},
		}

		if nsgId := raw["network_security_group_id"].(string); nsgId != "" {
			config.Properties.NetworkSecurityGroup = &virtualmachinescalesets.SubResource{
				Id: utils.String(nsgId),
			}
		}

		output = append(output, config)
	}

	return &output, nil
}

func expandVirtualMachineScaleSetIPConfiguration(raw map[string]interface{}) (*virtualmachinescalesets.VirtualMachineScaleSetIPConfiguration, error) {
	applicationGatewayBackendAddressPoolIdsRaw := raw["application_gateway_backend_address_pool_ids"].(*pluginsdk.Set).List()
	applicationGatewayBackendAddressPoolIds := expandIDsToSubResources(applicationGatewayBackendAddressPoolIdsRaw)

	applicationSecurityGroupIdsRaw := raw["application_security_group_ids"].(*pluginsdk.Set).List()
	applicationSecurityGroupIds := expandIDsToSubResources(applicationSecurityGroupIdsRaw)

	loadBalancerBackendAddressPoolIdsRaw := raw["load_balancer_backend_address_pool_ids"].(*pluginsdk.Set).List()
	loadBalancerBackendAddressPoolIds := expandIDsToSubResources(loadBalancerBackendAddressPoolIdsRaw)

	loadBalancerInboundNatPoolIdsRaw := raw["load_balancer_inbound_nat_rules_ids"].(*pluginsdk.Set).List()
	loadBalancerInboundNatPoolIds := expandIDsToSubResources(loadBalancerInboundNatPoolIdsRaw)

	primary := raw["primary"].(bool)
	version := virtualmachinescalesets.IPVersion(raw["version"].(string))
	if primary && version == virtualmachinescalesets.IPVersionIPvSix {
		return nil, fmt.Errorf("an IPv6 Primary IP Configuration is unsupported - instead add a IPv4 IP Configuration as the Primary and make the IPv6 IP Configuration the secondary")
	}

	ipConfiguration := virtualmachinescalesets.VirtualMachineScaleSetIPConfiguration{
		Name: raw["name"].(string),
		Properties: &virtualmachinescalesets.VirtualMachineScaleSetIPConfigurationProperties{
			Primary:                               utils.Bool(primary),
			PrivateIPAddressVersion:               &version,
			ApplicationGatewayBackendAddressPools: applicationGatewayBackendAddressPoolIds,
			ApplicationSecurityGroups:             applicationSecurityGroupIds,
			LoadBalancerBackendAddressPools:       loadBalancerBackendAddressPoolIds,
			LoadBalancerInboundNatPools:           loadBalancerInboundNatPoolIds,
		},
	}

	if subnetId := raw["subnet_id"].(string); subnetId != "" {
		ipConfiguration.Properties.Subnet = &virtualmachinescalesets.ApiEntityReference{
			Id: utils.String(subnetId),
		}
	}

	publicIPConfigsRaw := raw["public_ip_address"].([]interface{})
	if len(publicIPConfigsRaw) > 0 {
		publicIPConfigRaw := publicIPConfigsRaw[0].(map[string]interface{})
		publicIPAddressConfig := expandVirtualMachineScaleSetPublicIPAddress(publicIPConfigRaw)
		ipConfiguration.Properties.PublicIPAddressConfiguration = publicIPAddressConfig
	}

	return &ipConfiguration, nil
}

func expandVirtualMachineScaleSetPublicIPAddress(raw map[string]interface{}) *virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfiguration {
	ipTagsRaw := raw["ip_tag"].([]interface{})
	ipTags := make([]virtualmachinescalesets.VirtualMachineScaleSetIPTag, 0)
	for _, ipTagV := range ipTagsRaw {
		ipTagRaw := ipTagV.(map[string]interface{})
		ipTags = append(ipTags, virtualmachinescalesets.VirtualMachineScaleSetIPTag{
			Tag:       utils.String(ipTagRaw["tag"].(string)),
			IPTagType: utils.String(ipTagRaw["type"].(string)),
		})
	}

	publicIPAddressConfig := virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfiguration{
		Name: raw["name"].(string),
		Properties: &virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfigurationProperties{
			IPTags: &ipTags,
		},
	}

	if domainNameLabel := raw["domain_name_label"].(string); domainNameLabel != "" {
		dns := &virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettings{
			DomainNameLabel: domainNameLabel,
		}
		publicIPAddressConfig.Properties.DnsSettings = dns
	}

	if idleTimeout := raw["idle_timeout_in_minutes"].(int); idleTimeout > 0 {
		publicIPAddressConfig.Properties.IdleTimeoutInMinutes = utils.Int64(int64(raw["idle_timeout_in_minutes"].(int)))
	}

	if publicIPPrefixID := raw["public_ip_prefix_id"].(string); publicIPPrefixID != "" {
		publicIPAddressConfig.Properties.PublicIPPrefix = &virtualmachinescalesets.SubResource{
			Id: utils.String(publicIPPrefixID),
		}
	}

	return &publicIPAddressConfig
}

func ExpandVirtualMachineScaleSetNetworkInterfaceUpdate(input []interface{}) (*[]virtualmachinescalesets.VirtualMachineScaleSetUpdateNetworkConfiguration, error) {
	output := make([]virtualmachinescalesets.VirtualMachineScaleSetUpdateNetworkConfiguration, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		dnsServers := utils.ExpandStringSlice(raw["dns_servers"].([]interface{}))

		ipConfigurations := make([]virtualmachinescalesets.VirtualMachineScaleSetUpdateIPConfiguration, 0)
		ipConfigurationsRaw := raw["ip_configuration"].([]interface{})
		for _, configV := range ipConfigurationsRaw {
			configRaw := configV.(map[string]interface{})
			ipConfiguration, err := expandVirtualMachineScaleSetIPConfigurationUpdate(configRaw)
			if err != nil {
				return nil, err
			}

			ipConfigurations = append(ipConfigurations, *ipConfiguration)
		}

		config := virtualmachinescalesets.VirtualMachineScaleSetUpdateNetworkConfiguration{
			Name: utils.String(raw["name"].(string)),
			Properties: &virtualmachinescalesets.VirtualMachineScaleSetUpdateNetworkConfigurationProperties{
				DnsSettings: &virtualmachinescalesets.VirtualMachineScaleSetNetworkConfigurationDnsSettings{
					DnsServers: dnsServers,
				},
				EnableAcceleratedNetworking: utils.Bool(raw["enable_accelerated_networking"].(bool)),
				EnableIPForwarding:          utils.Bool(raw["enable_ip_forwarding"].(bool)),
				IPConfigurations:            &ipConfigurations,
				Primary:                     utils.Bool(raw["primary"].(bool)),
			},
		}

		if nsgId := raw["network_security_group_id"].(string); nsgId != "" {
			config.Properties.NetworkSecurityGroup = &virtualmachinescalesets.SubResource{
				Id: utils.String(nsgId),
			}
		}

		output = append(output, config)
	}

	return &output, nil
}

func expandVirtualMachineScaleSetIPConfigurationUpdate(raw map[string]interface{}) (*virtualmachinescalesets.VirtualMachineScaleSetUpdateIPConfiguration, error) {
	applicationGatewayBackendAddressPoolIdsRaw := raw["application_gateway_backend_address_pool_ids"].(*pluginsdk.Set).List()
	applicationGatewayBackendAddressPoolIds := expandIDsToSubResources(applicationGatewayBackendAddressPoolIdsRaw)

	applicationSecurityGroupIdsRaw := raw["application_security_group_ids"].(*pluginsdk.Set).List()
	applicationSecurityGroupIds := expandIDsToSubResources(applicationSecurityGroupIdsRaw)

	loadBalancerBackendAddressPoolIdsRaw := raw["load_balancer_backend_address_pool_ids"].(*pluginsdk.Set).List()
	loadBalancerBackendAddressPoolIds := expandIDsToSubResources(loadBalancerBackendAddressPoolIdsRaw)

	loadBalancerInboundNatPoolIdsRaw := raw["load_balancer_inbound_nat_rules_ids"].(*pluginsdk.Set).List()
	loadBalancerInboundNatPoolIds := expandIDsToSubResources(loadBalancerInboundNatPoolIdsRaw)

	primary := raw["primary"].(bool)
	version := virtualmachinescalesets.IPVersion(raw["version"].(string))

	if primary && version == virtualmachinescalesets.IPVersionIPvSix {
		return nil, fmt.Errorf("An IPv6 Primary IP Configuration is unsupported - instead add a IPv4 IP Configuration as the Primary and make the IPv6 IP Configuration the secondary")
	}

	ipConfiguration := virtualmachinescalesets.VirtualMachineScaleSetUpdateIPConfiguration{
		Name: utils.String(raw["name"].(string)),
		Properties: &virtualmachinescalesets.VirtualMachineScaleSetUpdateIPConfigurationProperties{
			Primary:                               utils.Bool(primary),
			PrivateIPAddressVersion:               &version,
			ApplicationGatewayBackendAddressPools: applicationGatewayBackendAddressPoolIds,
			ApplicationSecurityGroups:             applicationSecurityGroupIds,
			LoadBalancerBackendAddressPools:       loadBalancerBackendAddressPoolIds,
			LoadBalancerInboundNatPools:           loadBalancerInboundNatPoolIds,
		},
	}

	if subnetId := raw["subnet_id"].(string); subnetId != "" {
		ipConfiguration.Properties.Subnet = &virtualmachinescalesets.ApiEntityReference{
			Id: utils.String(subnetId),
		}
	}

	publicIPConfigsRaw := raw["public_ip_address"].([]interface{})
	if len(publicIPConfigsRaw) > 0 {
		publicIPConfigRaw := publicIPConfigsRaw[0].(map[string]interface{})
		publicIPAddressConfig := expandVirtualMachineScaleSetPublicIPAddressUpdate(publicIPConfigRaw)
		ipConfiguration.Properties.PublicIPAddressConfiguration = publicIPAddressConfig
	}

	return &ipConfiguration, nil
}

func expandVirtualMachineScaleSetPublicIPAddressUpdate(raw map[string]interface{}) *virtualmachinescalesets.VirtualMachineScaleSetUpdatePublicIPAddressConfiguration {
	publicIPAddressConfig := virtualmachinescalesets.VirtualMachineScaleSetUpdatePublicIPAddressConfiguration{
		Name:       utils.String(raw["name"].(string)),
		Properties: &virtualmachinescalesets.VirtualMachineScaleSetUpdatePublicIPAddressConfigurationProperties{},
	}

	if domainNameLabel := raw["domain_name_label"].(string); domainNameLabel != "" {
		dns := &virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettings{
			DomainNameLabel: domainNameLabel,
		}
		publicIPAddressConfig.Properties.DnsSettings = dns
	}

	if idleTimeout := raw["idle_timeout_in_minutes"].(int); idleTimeout > 0 {
		publicIPAddressConfig.Properties.IdleTimeoutInMinutes = utils.Int64(int64(raw["idle_timeout_in_minutes"].(int)))
	}

	return &publicIPAddressConfig
}

func FlattenVirtualMachineScaleSetNetworkInterface(input *[]virtualmachinescalesets.VirtualMachineScaleSetNetworkConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, v := range *input {
		var name, networkSecurityGroupId string
		if v.Name != "" {
			name = v.Name
		}
		if v.Properties.NetworkSecurityGroup != nil && v.Properties.NetworkSecurityGroup.Id != nil {
			networkSecurityGroupId = *v.Properties.NetworkSecurityGroup.Id
		}

		var enableAcceleratedNetworking, enableIPForwarding, primary bool
		if v.Properties.EnableAcceleratedNetworking != nil {
			enableAcceleratedNetworking = *v.Properties.EnableAcceleratedNetworking
		}
		if v.Properties.EnableIPForwarding != nil {
			enableIPForwarding = *v.Properties.EnableIPForwarding
		}
		if v.Properties.Primary != nil {
			primary = *v.Properties.Primary
		}

		var dnsServers []interface{}
		if settings := v.Properties.DnsSettings; settings != nil {
			dnsServers = utils.FlattenStringSlice(v.Properties.DnsSettings.DnsServers)
		}

		var ipConfigurations []interface{}
		if v.Properties.IPConfigurations != nil {
			for _, configRaw := range v.Properties.IPConfigurations {
				config := flattenVirtualMachineScaleSetIPConfiguration(configRaw)
				ipConfigurations = append(ipConfigurations, config)
			}
		}

		results = append(results, map[string]interface{}{
			"name":                          name,
			"dns_servers":                   dnsServers,
			"enable_accelerated_networking": enableAcceleratedNetworking,
			"enable_ip_forwarding":          enableIPForwarding,
			"ip_configuration":              ipConfigurations,
			"network_security_group_id":     networkSecurityGroupId,
			"primary":                       primary,
		})
	}

	return results
}

func flattenVirtualMachineScaleSetIPConfiguration(input virtualmachinescalesets.VirtualMachineScaleSetIPConfiguration) map[string]interface{} {
	var name, subnetId string
	if input.Name != "" {
		name = input.Name
	}
	if input.Properties.Subnet != nil && input.Properties.Subnet.Id != nil {
		subnetId = *input.Properties.Subnet.Id
	}

	var primary bool
	if input.Properties.Primary != nil {
		primary = *input.Properties.Primary
	}

	var publicIPAddresses []interface{}
	if input.Properties.PublicIPAddressConfiguration != nil {
		publicIPAddresses = append(publicIPAddresses, flattenVirtualMachineScaleSetPublicIPAddress(*input.Properties.PublicIPAddressConfiguration))
	}

	applicationGatewayBackendAddressPoolIds := flattenSubResourcesToIDs(input.Properties.ApplicationGatewayBackendAddressPools)
	applicationSecurityGroupIds := flattenSubResourcesToIDs(input.Properties.ApplicationSecurityGroups)
	loadBalancerBackendAddressPoolIds := flattenSubResourcesToIDs(input.Properties.LoadBalancerBackendAddressPools)
	loadBalancerInboundNatRuleIds := flattenSubResourcesToIDs(input.Properties.LoadBalancerInboundNatPools)

	return map[string]interface{}{
		"name":              name,
		"primary":           primary,
		"public_ip_address": publicIPAddresses,
		"subnet_id":         subnetId,
		"version":           string(*input.Properties.PrivateIPAddressVersion),
		"application_gateway_backend_address_pool_ids": applicationGatewayBackendAddressPoolIds,
		"application_security_group_ids":               applicationSecurityGroupIds,
		"load_balancer_backend_address_pool_ids":       loadBalancerBackendAddressPoolIds,
		"load_balancer_inbound_nat_rules_ids":          loadBalancerInboundNatRuleIds,
	}
}

func flattenVirtualMachineScaleSetPublicIPAddress(input virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfiguration) map[string]interface{} {
	ipTags := make([]interface{}, 0)
	if input.Properties.IPTags != nil {
		for _, rawTag := range *input.Properties.IPTags {
			var tag, tagType string

			if rawTag.IPTagType != nil {
				tagType = *rawTag.IPTagType
			}

			if rawTag.Tag != nil {
				tag = *rawTag.Tag
			}

			ipTags = append(ipTags, map[string]interface{}{
				"tag":  tag,
				"type": tagType,
			})
		}
	}

	var domainNameLabel, name, publicIPPrefixId string
	if input.Properties.DnsSettings != nil && input.Properties.DnsSettings.DomainNameLabel != "" {
		domainNameLabel = input.Properties.DnsSettings.DomainNameLabel
	}
	if input.Name != "" {
		name = input.Name
	}
	if input.Properties.PublicIPPrefix != nil && input.Properties.PublicIPPrefix.Id != nil {
		publicIPPrefixId = *input.Properties.PublicIPPrefix.Id
	}

	var idleTimeoutInMinutes int
	if input.Properties.IdleTimeoutInMinutes != nil {
		idleTimeoutInMinutes = int(*input.Properties.IdleTimeoutInMinutes)
	}

	return map[string]interface{}{
		"name":                    name,
		"domain_name_label":       domainNameLabel,
		"idle_timeout_in_minutes": idleTimeoutInMinutes,
		"ip_tag":                  ipTags,
		"public_ip_prefix_id":     publicIPPrefixId,
	}
}

func VirtualMachineScaleSetDataDiskSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		// TODO: does this want to be a Set?
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"caching": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.CachingTypesNone),
						string(compute.CachingTypesReadOnly),
						string(compute.CachingTypesReadWrite),
					}, false),
				},

				"create_option": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.DiskCreateOptionTypesEmpty),
						string(compute.DiskCreateOptionTypesFromImage),
					}, false),
					Default: string(compute.DiskCreateOptionTypesEmpty),
				},

				"disk_encryption_set_id": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					// whilst the API allows updating this value, it's never actually set at Azure's end
					// presumably this'll take effect once key rotation is supported a few months post-GA?
					// however for now let's make this ForceNew since it can't be (successfully) updated
					ForceNew:     true,
					ValidateFunc: validate.DiskEncryptionSetID,
				},

				"disk_size_gb": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(1, 32767),
				},

				"lun": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(0, 2000), // TODO: confirm upper bounds
				},

				"storage_account_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.StorageAccountTypesPremiumLRS),
						string(compute.StorageAccountTypesStandardLRS),
						string(compute.StorageAccountTypesStandardSSDLRS),
						string(compute.StorageAccountTypesUltraSSDLRS),
					}, false),
				},

				"write_accelerator_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"ultra_ssd_disk_iops_read_write": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Computed: true,
				},

				"ultra_ssd_disk_mbps_read_write": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func ExpandVirtualMachineScaleSetDataDisk(input []interface{}, ultraSSDEnabled bool) (*[]virtualmachinescalesets.VirtualMachineScaleSetDataDisk, error) {
	disks := make([]virtualmachinescalesets.VirtualMachineScaleSetDataDisk, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		cachingType := virtualmachinescalesets.CachingTypes(raw["caching"].(string))
		storageAccountType := virtualmachinescalesets.StorageAccountTypes(raw["storage_account_type"].(string))
		disk := virtualmachinescalesets.VirtualMachineScaleSetDataDisk{
			Caching:    &cachingType,
			DiskSizeGB: utils.Int64(int64(raw["disk_size_gb"].(int))),
			Lun:        int64(raw["lun"].(int)),
			ManagedDisk: &virtualmachinescalesets.VirtualMachineScaleSetManagedDiskParameters{
				StorageAccountType: &storageAccountType,
			},
			WriteAcceleratorEnabled: utils.Bool(raw["write_accelerator_enabled"].(bool)),
			CreateOption:            virtualmachinescalesets.DiskCreateOptionTypes(raw["create_option"].(string)),
		}

		if id := raw["disk_encryption_set_id"].(string); id != "" {
			disk.ManagedDisk.DiskEncryptionSet = &virtualmachinescalesets.SubResource{
				Id: utils.String(id),
			}
		}

		var iops int
		if diskIops, ok := raw["disk_iops_read_write"]; ok && diskIops.(int) > 0 {
			iops = diskIops.(int)
		} else if ssdIops, ok := raw["ultra_ssd_disk_iops_read_write"]; ok && ssdIops.(int) > 0 {
			iops = ssdIops.(int)
		}

		if iops > 0 && !ultraSSDEnabled {
			return nil, fmt.Errorf("disk_iops_read_write and ultra_ssd_disk_iops_read_write are only available for UltraSSD disks")
		}

		var mbps int
		if diskMbps, ok := raw["disk_mbps_read_write"]; ok && diskMbps.(int) > 0 {
			mbps = diskMbps.(int)
		} else if ssdMbps, ok := raw["ultra_ssd_disk_mbps_read_write"]; ok && ssdMbps.(int) > 0 {
			mbps = ssdMbps.(int)
		}

		if mbps > 0 && !ultraSSDEnabled {
			return nil, fmt.Errorf("disk_mbps_read_write and ultra_ssd_disk_mbps_read_write are only available for UltraSSD disks")
		}

		// Do not set value unless value is greater than 0 - issue 15516
		if iops > 0 {
			disk.DiskIOPSReadWrite = utils.Int64(int64(iops))
		}

		// Do not set value unless value is greater than 0 - issue 15516
		if mbps > 0 {
			disk.DiskMBpsReadWrite = utils.Int64(int64(mbps))
		}

		disks = append(disks, disk)
	}

	return &disks, nil
}

func FlattenVirtualMachineScaleSetDataDisk(input *[]virtualmachinescalesets.VirtualMachineScaleSetDataDisk) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, v := range *input {
		diskSizeGb := 0
		if v.DiskSizeGB != nil && *v.DiskSizeGB != 0 {
			diskSizeGb = int(*v.DiskSizeGB)
		}

		lun := 0
		if v.Lun != 0 {
			lun = int(v.Lun)
		}

		storageAccountType := ""
		diskEncryptionSetId := ""
		if v.ManagedDisk != nil {
			storageAccountType = string(*v.ManagedDisk.StorageAccountType)
			if v.ManagedDisk.DiskEncryptionSet != nil && v.ManagedDisk.DiskEncryptionSet.Id != nil {
				diskEncryptionSetId = *v.ManagedDisk.DiskEncryptionSet.Id
			}
		}

		writeAcceleratorEnabled := false
		if v.WriteAcceleratorEnabled != nil {
			writeAcceleratorEnabled = *v.WriteAcceleratorEnabled
		}

		iops := 0
		if v.DiskIOPSReadWrite != nil {
			iops = int(*v.DiskIOPSReadWrite)
		}

		mbps := 0
		if v.DiskMBpsReadWrite != nil {
			mbps = int(*v.DiskMBpsReadWrite)
		}

		dataDisk := map[string]interface{}{
			"caching":                   string(*v.Caching),
			"create_option":             string(v.CreateOption),
			"lun":                       lun,
			"disk_encryption_set_id":    diskEncryptionSetId,
			"disk_size_gb":              diskSizeGb,
			"storage_account_type":      storageAccountType,
			"write_accelerator_enabled": writeAcceleratorEnabled,
		}

		// Do not set value unless value is greater than 0 - issue 15516
		if iops > 0 {
			dataDisk["ultra_ssd_disk_iops_read_write"] = iops
		}

		// Do not set value unless value is greater than 0 - issue 15516
		if mbps > 0 {
			dataDisk["ultra_ssd_disk_mbps_read_write"] = mbps
		}

		output = append(output, dataDisk)
	}

	return output
}

func VirtualMachineScaleSetOSDiskSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"caching": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.CachingTypesNone),
						string(compute.CachingTypesReadOnly),
						string(compute.CachingTypesReadWrite),
					}, false),
				},
				"storage_account_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					// whilst this appears in the Update block the API returns this when changing:
					// Changing property 'osDisk.managedDisk.storageAccountType' is not allowed
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						// note: OS Disks don't support Ultra SSDs
						string(compute.StorageAccountTypesPremiumLRS),
						string(compute.StorageAccountTypesStandardLRS),
						string(compute.StorageAccountTypesStandardSSDLRS),
					}, false),
				},

				"diff_disk_settings": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"option": {
								Type:     pluginsdk.TypeString,
								Required: true,
								ForceNew: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(compute.DiffDiskOptionsLocal),
								}, false),
							},
							"placement": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								ForceNew: true,
								Default:  string(compute.DiffDiskPlacementCacheDisk),
								ValidateFunc: validation.StringInSlice([]string{
									string(compute.DiffDiskPlacementCacheDisk),
									string(compute.DiffDiskPlacementResourceDisk),
								}, false),
							},
						},
					},
				},

				"disk_encryption_set_id": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					// whilst the API allows updating this value, it's never actually set at Azure's end
					// presumably this'll take effect once key rotation is supported a few months post-GA?
					// however for now let's make this ForceNew since it can't be (successfully) updated
					ForceNew:      true,
					ValidateFunc:  validate.DiskEncryptionSetID,
					ConflictsWith: []string{"os_disk.0.secure_vm_disk_encryption_set_id"},
				},

				"disk_size_gb": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(0, 4095),
				},

				"secure_vm_disk_encryption_set_id": {
					Type:          pluginsdk.TypeString,
					Optional:      true,
					ForceNew:      true,
					ValidateFunc:  validate.DiskEncryptionSetID,
					ConflictsWith: []string{"os_disk.0.disk_encryption_set_id"},
				},

				"security_encryption_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(compute.SecurityEncryptionTypesVMGuestStateOnly),
						string(compute.SecurityEncryptionTypesDiskWithVMGuestState),
					}, false),
				},

				"write_accelerator_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func ExpandVirtualMachineScaleSetOSDisk(input []interface{}, osType virtualmachinescalesets.OperatingSystemTypes) (*virtualmachinescalesets.VirtualMachineScaleSetOSDisk, error) {
	raw := input[0].(map[string]interface{})
	caching := raw["caching"].(string)
	cachingType := virtualmachinescalesets.CachingTypes(caching)
	storageAccountType := virtualmachinescalesets.StorageAccountTypes(raw["storage_account_type"].(string))
	disk := virtualmachinescalesets.VirtualMachineScaleSetOSDisk{
		Caching: &cachingType,
		ManagedDisk: &virtualmachinescalesets.VirtualMachineScaleSetManagedDiskParameters{
			StorageAccountType: &storageAccountType,
		},
		WriteAcceleratorEnabled: utils.Bool(raw["write_accelerator_enabled"].(bool)),

		// these have to be hard-coded so there's no point exposing them
		CreateOption: virtualmachinescalesets.DiskCreateOptionTypesFromImage,
		OsType:       &osType,
	}

	securityEncryptionType := virtualmachinescalesets.SecurityEncryptionTypes(raw["security_encryption_type"].(string))
	if securityEncryptionType != "" {
		disk.ManagedDisk.SecurityProfile = &virtualmachinescalesets.VMDiskSecurityProfile{
			SecurityEncryptionType: &securityEncryptionType,
		}
	}
	if secureVMDiskEncryptionId := raw["secure_vm_disk_encryption_set_id"].(string); secureVMDiskEncryptionId != "" {
		if virtualmachinescalesets.SecurityEncryptionTypesDiskWithVMGuestState != virtualmachinescalesets.SecurityEncryptionTypes(securityEncryptionType) {
			return nil, fmt.Errorf("`secure_vm_disk_encryption_set_id` can only be specified when `security_encryption_type` is set to `DiskWithVMGuestState`")
		}
		disk.ManagedDisk.SecurityProfile.DiskEncryptionSet = &virtualmachinescalesets.SubResource{
			Id: utils.String(secureVMDiskEncryptionId),
		}
	}

	if diskEncryptionSetId := raw["disk_encryption_set_id"].(string); diskEncryptionSetId != "" {
		disk.ManagedDisk.DiskEncryptionSet = &virtualmachinescalesets.SubResource{
			Id: utils.String(diskEncryptionSetId),
		}
	}

	if osDiskSize := raw["disk_size_gb"].(int); osDiskSize > 0 {
		disk.DiskSizeGB = utils.Int64(int64(osDiskSize))
	}

	if diffDiskSettingsRaw := raw["diff_disk_settings"].([]interface{}); len(diffDiskSettingsRaw) > 0 {
		if caching != string(virtualmachinescalesets.CachingTypesReadOnly) {
			// Restriction per https://docs.microsoft.com/azure/virtual-machines/ephemeral-os-disks-deploy#vm-template-deployment
			return nil, fmt.Errorf("`diff_disk_settings` can only be set when `caching` is set to `ReadOnly`")
		}

		diffDiskRaw := diffDiskSettingsRaw[0].(map[string]interface{})
		diffDiskOption := virtualmachinescalesets.DiffDiskOptions(diffDiskRaw["option"].(string))
		diffDiskPlacement := virtualmachinescalesets.DiffDiskPlacement(diffDiskRaw["placement"].(string))
		disk.DiffDiskSettings = &virtualmachinescalesets.DiffDiskSettings{
			Option:    &diffDiskOption,
			Placement: &diffDiskPlacement,
		}
	}

	return &disk, nil
}

func ExpandVirtualMachineScaleSetOSDiskUpdate(input []interface{}) *virtualmachinescalesets.VirtualMachineScaleSetUpdateOSDisk {
	raw := input[0].(map[string]interface{})
	cachingType := virtualmachinescalesets.CachingTypes(raw["caching"].(string))
	storageAccountType := virtualmachinescalesets.StorageAccountTypes(raw["storage_account_type"].(string))
	disk := virtualmachinescalesets.VirtualMachineScaleSetUpdateOSDisk{
		Caching: &cachingType,
		ManagedDisk: &virtualmachinescalesets.VirtualMachineScaleSetManagedDiskParameters{
			StorageAccountType: &storageAccountType,
		},
		WriteAcceleratorEnabled: utils.Bool(raw["write_accelerator_enabled"].(bool)),
	}

	if diskEncryptionSetId := raw["disk_encryption_set_id"].(string); diskEncryptionSetId != "" {
		disk.ManagedDisk.DiskEncryptionSet = &virtualmachinescalesets.SubResource{
			Id: utils.String(diskEncryptionSetId),
		}
	}

	if osDiskSize := raw["disk_size_gb"].(int); osDiskSize > 0 {
		disk.DiskSizeGB = utils.Int64(int64(osDiskSize))
	}

	return &disk
}

func FlattenVirtualMachineScaleSetOSDisk(input *virtualmachinescalesets.VirtualMachineScaleSetOSDisk) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	diffDiskSettings := make([]interface{}, 0)
	if input.DiffDiskSettings != nil {
		diffDiskSettings = append(diffDiskSettings, map[string]interface{}{
			"option":    input.DiffDiskSettings.Option,
			"placement": input.DiffDiskSettings.Placement,
		})
	}

	diskSizeGb := 0
	if input.DiskSizeGB != nil && *input.DiskSizeGB != 0 {
		diskSizeGb = int(*input.DiskSizeGB)
	}

	storageAccountType := ""
	diskEncryptionSetId := ""
	secureVMDiskEncryptionSetId := ""
	securityEncryptionType := ""
	if input.ManagedDisk != nil {
		storageAccountType = string(*input.ManagedDisk.StorageAccountType)
		if input.ManagedDisk.DiskEncryptionSet != nil && input.ManagedDisk.DiskEncryptionSet.Id != nil {
			diskEncryptionSetId = *input.ManagedDisk.DiskEncryptionSet.Id
		}

		if securityProfile := input.ManagedDisk.SecurityProfile; securityProfile != nil {
			securityEncryptionType = string(*securityProfile.SecurityEncryptionType)
			if securityProfile.DiskEncryptionSet != nil && securityProfile.DiskEncryptionSet.Id != nil {
				secureVMDiskEncryptionSetId = *securityProfile.DiskEncryptionSet.Id
			}
		}
	}

	writeAcceleratorEnabled := false
	if input.WriteAcceleratorEnabled != nil {
		writeAcceleratorEnabled = *input.WriteAcceleratorEnabled
	}

	return []interface{}{
		map[string]interface{}{
			"caching":                          input.Caching,
			"disk_size_gb":                     diskSizeGb,
			"diff_disk_settings":               diffDiskSettings,
			"storage_account_type":             storageAccountType,
			"write_accelerator_enabled":        writeAcceleratorEnabled,
			"disk_encryption_set_id":           diskEncryptionSetId,
			"secure_vm_disk_encryption_set_id": secureVMDiskEncryptionSetId,
			"security_encryption_type":         securityEncryptionType,
		},
	}
}

func VirtualMachineScaleSetAutomatedOSUpgradePolicySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				// TODO: should these be optional + defaulted?
				"disable_automatic_rollback": {
					Type:     pluginsdk.TypeBool,
					Required: true,
				},
				// TODO 4.0: change this from enable_* to *_enabled
				"enable_automatic_os_upgrade": {
					Type:     pluginsdk.TypeBool,
					Required: true,
				},
			},
		},
	}
}

func ExpandVirtualMachineScaleSetAutomaticUpgradePolicy(input []interface{}) *virtualmachinescalesets.AutomaticOSUpgradePolicy {
	if len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]interface{})
	return &virtualmachinescalesets.AutomaticOSUpgradePolicy{
		DisableAutomaticRollback: utils.Bool(raw["disable_automatic_rollback"].(bool)),
		EnableAutomaticOSUpgrade: utils.Bool(raw["enable_automatic_os_upgrade"].(bool)),
	}
}

func FlattenVirtualMachineScaleSetAutomaticOSUpgradePolicy(input *virtualmachinescalesets.AutomaticOSUpgradePolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	disableAutomaticRollback := false
	if input.DisableAutomaticRollback != nil {
		disableAutomaticRollback = *input.DisableAutomaticRollback
	}

	enableAutomaticOSUpgrade := false
	if input.EnableAutomaticOSUpgrade != nil {
		enableAutomaticOSUpgrade = *input.EnableAutomaticOSUpgrade
	}

	return []interface{}{
		map[string]interface{}{
			"disable_automatic_rollback":  disableAutomaticRollback,
			"enable_automatic_os_upgrade": enableAutomaticOSUpgrade,
		},
	}
}

func VirtualMachineScaleSetRollingUpgradePolicySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"max_batch_instance_percent": {
					Type:     pluginsdk.TypeInt,
					Required: true,
				},
				"max_unhealthy_instance_percent": {
					Type:     pluginsdk.TypeInt,
					Required: true,
				},
				"max_unhealthy_upgraded_instance_percent": {
					Type:     pluginsdk.TypeInt,
					Required: true,
				},
				"pause_time_between_batches": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: azValidate.ISO8601Duration,
				},
			},
		},
	}
}

func ExpandVirtualMachineScaleSetRollingUpgradePolicy(input []interface{}) *virtualmachinescalesets.RollingUpgradePolicy {
	if len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]interface{})

	return &virtualmachinescalesets.RollingUpgradePolicy{
		MaxBatchInstancePercent:             utils.Int64(int64(raw["max_batch_instance_percent"].(int))),
		MaxUnhealthyInstancePercent:         utils.Int64(int64(raw["max_unhealthy_instance_percent"].(int))),
		MaxUnhealthyUpgradedInstancePercent: utils.Int64(int64(raw["max_unhealthy_upgraded_instance_percent"].(int))),
		PauseTimeBetweenBatches:             utils.String(raw["pause_time_between_batches"].(string)),
	}
}

func FlattenVirtualMachineScaleSetRollingUpgradePolicy(input *virtualmachinescalesets.RollingUpgradePolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	maxBatchInstancePercent := 0
	if input.MaxBatchInstancePercent != nil {
		maxBatchInstancePercent = int(*input.MaxBatchInstancePercent)
	}

	maxUnhealthyInstancePercent := 0
	if input.MaxUnhealthyInstancePercent != nil {
		maxUnhealthyInstancePercent = int(*input.MaxUnhealthyInstancePercent)
	}

	maxUnhealthyUpgradedInstancePercent := 0
	if input.MaxUnhealthyUpgradedInstancePercent != nil {
		maxUnhealthyUpgradedInstancePercent = int(*input.MaxUnhealthyUpgradedInstancePercent)
	}

	pauseTimeBetweenBatches := ""
	if input.PauseTimeBetweenBatches != nil {
		pauseTimeBetweenBatches = *input.PauseTimeBetweenBatches
	}

	return []interface{}{
		map[string]interface{}{
			"max_batch_instance_percent":              maxBatchInstancePercent,
			"max_unhealthy_instance_percent":          maxUnhealthyInstancePercent,
			"max_unhealthy_upgraded_instance_percent": maxUnhealthyUpgradedInstancePercent,
			"pause_time_between_batches":              pauseTimeBetweenBatches,
		},
	}
}

// TODO remove VirtualMachineScaleSetTerminateNotificationSchema in 4.0
func VirtualMachineScaleSetTerminateNotificationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:       pluginsdk.TypeList,
		Optional:   true,
		Computed:   true,
		MaxItems:   1,
		Deprecated: "`terminate_notification` has been renamed to `termination_notification` and will be removed in 4.0.",
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"enabled": {
					Type:     pluginsdk.TypeBool,
					Required: true,
				},
				"timeout": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: azValidate.ISO8601DurationBetween("PT5M", "PT15M"),
					Default:      "PT5M",
				},
			},
		},
		ConflictsWith: []string{"termination_notification"},
	}
}

func VirtualMachineScaleSetTerminationNotificationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"enabled": {
					Type:     pluginsdk.TypeBool,
					Required: true,
				},
				"timeout": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: azValidate.ISO8601DurationBetween("PT5M", "PT15M"),
					Default:      "PT5M",
				},
			},
		},
		// TODO remove ConflictsWith in 4.0
		ConflictsWith: []string{"terminate_notification"},
	}
}

func ExpandVirtualMachineScaleSetScheduledEventsProfile(input []interface{}) *virtualmachinescalesets.ScheduledEventsProfile {
	if len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]interface{})
	enabled := raw["enabled"].(bool)
	timeout := raw["timeout"].(string)

	return &virtualmachinescalesets.ScheduledEventsProfile{
		TerminateNotificationProfile: &virtualmachinescalesets.TerminateNotificationProfile{
			Enable:           &enabled,
			NotBeforeTimeout: &timeout,
		},
	}
}

func FlattenVirtualMachineScaleSetScheduledEventsProfile(input *virtualmachinescalesets.ScheduledEventsProfile) []interface{} {
	// if enabled is set to false, there will be no ScheduledEventsProfile in response, to avoid plan non empty when
	// a user explicitly set enabled to false, we need to assign a default block to this field

	enabled := false
	if input != nil && input.TerminateNotificationProfile != nil && input.TerminateNotificationProfile.Enable != nil {
		enabled = *input.TerminateNotificationProfile.Enable
	}

	timeout := "PT5M"
	if input != nil && input.TerminateNotificationProfile != nil && input.TerminateNotificationProfile.NotBeforeTimeout != nil {
		timeout = *input.TerminateNotificationProfile.NotBeforeTimeout
	}

	return []interface{}{
		map[string]interface{}{
			"enabled": enabled,
			"timeout": timeout,
		},
	}
}

func VirtualMachineScaleSetAutomaticRepairsPolicySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"enabled": {
					Type:     pluginsdk.TypeBool,
					Required: true,
				},
				"grace_period": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "PT30M",
					// this field actually has a range from 30m to 90m, is there a function that can do this validation?
					ValidateFunc: azValidate.ISO8601Duration,
				},
			},
		},
	}
}

func ExpandVirtualMachineScaleSetAutomaticRepairsPolicy(input []interface{}) *virtualmachinescalesets.AutomaticRepairsPolicy {
	if len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]interface{})

	return &virtualmachinescalesets.AutomaticRepairsPolicy{
		Enabled:     utils.Bool(raw["enabled"].(bool)),
		GracePeriod: utils.String(raw["grace_period"].(string)),
	}
}

func FlattenVirtualMachineScaleSetAutomaticRepairsPolicy(input *virtualmachinescalesets.AutomaticRepairsPolicy) []interface{} {
	// if enabled is set to false, there will be no AutomaticRepairsPolicy in response, to avoid plan non empty when
	// a user explicitly set enabled to false, we need to assign a default block to this field

	enabled := false
	if input != nil && input.Enabled != nil {
		enabled = *input.Enabled
	}

	gracePeriod := "PT30M"
	if input != nil && input.GracePeriod != nil {
		gracePeriod = *input.GracePeriod
	}

	return []interface{}{
		map[string]interface{}{
			"enabled":      enabled,
			"grace_period": gracePeriod,
		},
	}
}

func VirtualMachineScaleSetExtensionsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"publisher": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"type": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"type_handler_version": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"auto_upgrade_minor_version": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
				},

				"automatic_upgrade_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},

				"force_update_tag": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"protected_settings": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsJSON,
				},

				"provision_after_extensions": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"settings": {
					Type:             pluginsdk.TypeString,
					Optional:         true,
					ValidateFunc:     validation.StringIsJSON,
					DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
				},
			},
		},
		Set: virtualMachineScaleSetExtensionHash,
	}
}

func virtualMachineScaleSetExtensionHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["publisher"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["type"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["type_handler_version"].(string)))
		buf.WriteString(fmt.Sprintf("%t-", m["auto_upgrade_minor_version"].(bool)))

		if v, ok = m["force_update_tag"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}

		if v, ok := m["provision_after_extensions"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}

		// we need to ensure the whitespace is consistent
		settings := m["settings"].(string)
		if settings != "" {
			expandedSettings, err := pluginsdk.ExpandJsonFromString(settings)
			if err == nil {
				serializedSettings, err := pluginsdk.FlattenJsonToString(expandedSettings)
				if err == nil {
					buf.WriteString(fmt.Sprintf("%s-", serializedSettings))
				}
			}
		}

		if v, ok := m["protected_settings"]; ok {
			settings := v.(string)
			if settings != "" {
				expandedSettings, err := pluginsdk.ExpandJsonFromString(settings)
				if err == nil {
					serializedSettings, err := pluginsdk.FlattenJsonToString(expandedSettings)
					if err == nil {
						buf.WriteString(fmt.Sprintf("%s-", serializedSettings))
					}
				}
			}
		}
	}

	return pluginsdk.HashString(buf.String())
}

func expandVirtualMachineScaleSetExtensions(input []interface{}) (extensionProfile *virtualmachinescalesets.VirtualMachineScaleSetExtensionProfile, hasHealthExtension bool, err error) {
	extensionProfile = &virtualmachinescalesets.VirtualMachineScaleSetExtensionProfile{}
	if len(input) == 0 {
		return extensionProfile, false, nil
	}

	extensions := make([]virtualmachinescalesets.VirtualMachineScaleSetExtension, 0)
	for _, v := range input {
		extensionRaw := v.(map[string]interface{})
		extension := virtualmachinescalesets.VirtualMachineScaleSetExtension{
			Name: utils.String(extensionRaw["name"].(string)),
		}
		extensionType := extensionRaw["type"].(string)

		extensionProps := virtualmachinescalesets.VirtualMachineScaleSetExtensionProperties{
			Publisher:                utils.String(extensionRaw["publisher"].(string)),
			Type:                     &extensionType,
			TypeHandlerVersion:       utils.String(extensionRaw["type_handler_version"].(string)),
			AutoUpgradeMinorVersion:  utils.Bool(extensionRaw["auto_upgrade_minor_version"].(bool)),
			EnableAutomaticUpgrade:   utils.Bool(extensionRaw["automatic_upgrade_enabled"].(bool)),
			ProvisionAfterExtensions: utils.ExpandStringSlice(extensionRaw["provision_after_extensions"].([]interface{})),
		}

		if extensionType == "ApplicationHealthLinux" || extensionType == "ApplicationHealthWindows" {
			hasHealthExtension = true
		}

		if forceUpdateTag := extensionRaw["force_update_tag"]; forceUpdateTag != nil {
			extensionProps.ForceUpdateTag = utils.String(forceUpdateTag.(string))
		}

		if val, ok := extensionRaw["settings"]; ok && val.(string) != "" {
			var result interface{}
			err = json.Unmarshal([]byte(val.(string)), &result)
			if err != nil {
				return nil, false, fmt.Errorf("failed to parse JSON from `settings`: %+v", err)
			}
			extensionProps.Settings = &result
		}

		if val, ok := extensionRaw["protected_settings"]; ok && val.(string) != "" {
			var result interface{}
			err = json.Unmarshal([]byte(val.(string)), &result)
			if err != nil {
				return nil, false, fmt.Errorf("failed to parse JSON from `protected_settings`: %+v", err)
			}
			extensionProps.ProtectedSettings = &result
		}

		extension.Properties = &extensionProps
		extensions = append(extensions, extension)
	}
	extensionProfile.Extensions = &extensions

	return extensionProfile, hasHealthExtension, nil
}

func flattenVirtualMachineScaleSetExtensions(input *virtualmachinescalesets.VirtualMachineScaleSetExtensionProfile, d *pluginsdk.ResourceData) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)
	if input == nil || input.Extensions == nil {
		return result, nil
	}

	// extensionsFromState holds the "extension" block, which is used to retrieve the "protected_settings" to fill it back the state,
	// since it is not returned from the API.
	extensionsFromState := map[string]map[string]interface{}{}
	if extSet, ok := d.GetOk("extension"); ok && extSet != nil {
		extensions := extSet.(*pluginsdk.Set).List()
		for _, ext := range extensions {
			if ext == nil {
				continue
			}
			ext := ext.(map[string]interface{})
			extensionsFromState[ext["name"].(string)] = ext
		}
	}

	for _, v := range *input.Extensions {
		name := ""
		if v.Name != nil {
			name = *v.Name
		}

		autoUpgradeMinorVersion := false
		enableAutomaticUpgrade := false
		forceUpdateTag := ""
		provisionAfterExtension := make([]interface{}, 0)
		protectedSettings := ""
		extPublisher := ""
		extSettings := ""
		extType := ""
		extTypeVersion := ""

		if props := v.Properties; props != nil {
			if props.Publisher != nil {
				extPublisher = *props.Publisher
			}

			if props.Type != nil {
				extType = *props.Type
			}

			if props.TypeHandlerVersion != nil {
				extTypeVersion = *props.TypeHandlerVersion
			}

			if props.AutoUpgradeMinorVersion != nil {
				autoUpgradeMinorVersion = *props.AutoUpgradeMinorVersion
			}

			if props.EnableAutomaticUpgrade != nil {
				enableAutomaticUpgrade = *props.EnableAutomaticUpgrade
			}

			if props.ForceUpdateTag != nil {
				forceUpdateTag = *props.ForceUpdateTag
			}

			if props.ProvisionAfterExtensions != nil {
				provisionAfterExtension = utils.FlattenStringSlice(props.ProvisionAfterExtensions)
			}

			if props.Settings != nil {
				extSettingsRaw, err := json.Marshal(props.Settings)
				if err != nil {
					return nil, err
				}
				extSettings = string(extSettingsRaw)
			}
		}
		// protected_settings isn't returned, so we attempt to get it from state otherwise set to empty string
		if ext, ok := extensionsFromState[name]; ok {
			if protectedSettingsFromState, ok := ext["protected_settings"]; ok {
				if protectedSettingsFromState.(string) != "" && protectedSettingsFromState.(string) != "{}" {
					protectedSettings = protectedSettingsFromState.(string)
				}
			}
		}

		result = append(result, map[string]interface{}{
			"name":                       name,
			"auto_upgrade_minor_version": autoUpgradeMinorVersion,
			"automatic_upgrade_enabled":  enableAutomaticUpgrade,
			"force_update_tag":           forceUpdateTag,
			"provision_after_extensions": provisionAfterExtension,
			"protected_settings":         protectedSettings,
			"publisher":                  extPublisher,
			"settings":                   extSettings,
			"type":                       extType,
			"type_handler_version":       extTypeVersion,
		})
	}
	return result, nil
}
