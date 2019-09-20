package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmVirtualHub() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmVirtualHubRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"address_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"express_route_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"p2svpn_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"route_table": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"routes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address_prefixes": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"next_hop_ip_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"virtual_network_connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_hub_to_remote_vnet_transit": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"allow_remote_vnet_to_use_hub_vnet_gateways": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_internet_security": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_virtual_network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"virtual_wan_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmVirtualHubRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.VirtualHubClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Virtual Hub %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if virtualHubProperties := resp.VirtualHubProperties; virtualHubProperties != nil {
		d.Set("address_prefix", virtualHubProperties.AddressPrefix)
		if virtualHubProperties.ExpressRouteGateway != nil {
			if err := d.Set("express_route_gateway_id", virtualHubProperties.ExpressRouteGateway.ID); err != nil {
				return fmt.Errorf("Error setting `express_route_gateway`: %+v", err)
			}
		}
		if virtualHubProperties.P2SVpnGateway != nil {
			if err := d.Set("p2svpn_gateway_id", virtualHubProperties.P2SVpnGateway.ID); err != nil {
				return fmt.Errorf("Error setting `p2svpn_gateway`: %+v", err)
			}
		}
		if err := d.Set("route_table", flattenArmVirtualHubVirtualHubRouteTable(virtualHubProperties.RouteTable)); err != nil {
			return fmt.Errorf("Error setting `route_table`: %+v", err)
		}
		if err := d.Set("virtual_network_connections", flattenArmVirtualHubHubVirtualNetworkConnection(virtualHubProperties.VirtualNetworkConnections)); err != nil {
			return fmt.Errorf("Error setting `virtual_network_connections`: %+v", err)
		}
		if virtualHubProperties.VirtualWan != nil {
			if err := d.Set("virtual_wan_id", virtualHubProperties.VirtualWan.ID); err != nil {
				return fmt.Errorf("Error setting `virtual_wan`: %+v", err)
			}
		}
		if virtualHubProperties.VpnGateway != nil {
			if err := d.Set("vpn_gateway_id", virtualHubProperties.VpnGateway.ID); err != nil {
				return fmt.Errorf("Error setting `vpn_gateway`: %+v", err)
			}
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
