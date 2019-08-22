package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-06-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmPrivateLinkService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPrivateLinkServiceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"auto_approval_subscription_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"fqdns": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"visibility_subscription_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"ip_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"private_ip_address_version": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"private_ip_address_allocation": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"load_balancer_frontend_ip_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmPrivateLinkServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.PrivateLinkServiceClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Private Link Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)

	d.Set("resource_group_name", resGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if fqdns := resp.Fqdns; fqdns != nil {
		d.Set("fqdns", fqdns)
	}

	if autoApproval := resp.AutoApproval; autoApproval != nil {
		if autoApprovalSubs := autoApproval.Subscriptions; autoApprovalSubs != nil {
			d.Set("auto_approval_subscription_names", autoApprovalSubs)
		}
	}

	if visibility := resp.Visibility; visibility != nil {
		if visibilitySubs := visibility.Subscriptions; visibilitySubs != nil {
			d.Set("visibility_subscription_names", visibilitySubs)
		}
	}

	frontendIPConfigs := flattenLoadBalancerFrontendIPConfigurations(resp.LoadBalancerFrontendIPConfigurations)
	if err := d.Set("load_balancer_frontend_ip_configuration", frontendIPConfigs); err != nil {
		return err
	}

	ipConfigs := flattenIPConfigurations(resp.IPConfigurations)
	if err := d.Set("ip_configuration", ipConfigs); err != nil {
		return err
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func flattenLoadBalancerFrontendIPConfigurations(input *[]network.FrontendIPConfiguration) []interface{} {
	frontendIPConfigs := make([]interface{}, 0)
	if input == nil {
		return frontendIPConfigs
	}

	for _, frontendIPConfig := range *input {
		ipConfig := make(map[string]interface{})
		if frontendIPConfig.ID != nil {
			ipConfig["id"] = *frontendIPConfig.ID
		}
		frontendIPConfigs = append(frontendIPConfigs, ipConfig)
	}

	return frontendIPConfigs
}

func flattenIPConfigurations(input *[]network.PrivateLinkServiceIPConfiguration) []interface{} {
	ipConfigs := make([]interface{}, 0)
	if input == nil {
		return ipConfigs
	}

	for _, ipConfig := range *input {
		data := make(map[string]interface{})

		data["private_ip_address_allocation"] = ipConfig.PrivateIPAllocationMethod

		data["private_ip_address_version"] = ipConfig.PrivateIPAddressVersion

		if ipConfig.Name != nil {
			data["name"] = *ipConfig.Name
		}

		if ipConfig.PrivateIPAddress != nil {
			data["private_ip_address"] = *ipConfig.PrivateIPAddress
		}

		if ipConfig.Subnet != nil {
			if ipConfig.Subnet.ID != nil {
				data["subnet_id"] = *ipConfig.Subnet.ID
			}
		}

		ipConfigs = append(ipConfigs, data)
	}

	return ipConfigs
}
