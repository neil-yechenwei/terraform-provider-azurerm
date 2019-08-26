package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmVirtualHub() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmVirtualHubRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"address_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"express_route_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"p2s_vpn_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"route_table": {
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

			"virtual_network_connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

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
			return fmt.Errorf("Virtual Hub %q was not found in Resource Group %q", name, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	if id := resp.ID; id != nil {
		d.SetId(*resp.ID)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.VirtualHubProperties; props != nil {
		if props.AddressPrefix != nil {
			d.Set("address_prefix", props.AddressPrefix)
		}

		if props.ExpressRouteGateway != nil {
			d.Set("express_route_gateway_id", props.ExpressRouteGateway.ID)
		}

		if props.P2SVpnGateway != nil {
			d.Set("p2s_vpn_gateway_id", props.P2SVpnGateway.ID)
		}

		if props.RouteTable != nil {
			virtualHubRouteTable := flattenVirtualHubRouteTable(props.RouteTable.Routes)
			if virtualHubRouteTable != nil {
				d.Set("route_table", virtualHubRouteTable)
			}
		}

		if props.VirtualNetworkConnections != nil {
			virtualNetworkConnections := flattenHubVirtualNetworkConnections(props.VirtualNetworkConnections)
			if virtualNetworkConnections != nil {
				d.Set("virtual_network_connections", virtualNetworkConnections)
			}
		}

		if props.VirtualWan != nil {
			d.Set("virtual_wan_id", props.VirtualWan.ID)
		}

		if props.VpnGateway != nil {
			d.Set("vpn_gateway_id", props.VpnGateway.ID)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func flattenVirtualHubRouteTable(routes *[]network.VirtualHubRoute) []interface{} {
	results := make([]interface{}, 0)

	if routes == nil {
		return results
	}

	for _, v := range *routes {
		route := make(map[string]interface{})
		if v.AddressPrefixes != nil {
			route["address_prefixes"] = v.AddressPrefixes
		}

		if v.NextHopIPAddress != nil {
			route["next_hop_ip_address"] = v.NextHopIPAddress
		}

		results = append(results, route)
	}

	return results
}

func flattenHubVirtualNetworkConnections(input *[]network.HubVirtualNetworkConnection) []interface{} {
	connections := make([]interface{}, 0)

	if input == nil {
		return connections
	}

	for _, v := range *input {
		connection := make(map[string]interface{})
		if v.ID != nil {
			connection["id"] = v.ID
		}

		if v.Name != nil {
			connection["name"] = v.Name
		}

		if v.AllowHubToRemoteVnetTransit != nil {
			connection["allow_hub_to_remote_vnet_transit"] = v.AllowHubToRemoteVnetTransit
		}

		if v.AllowRemoteVnetToUseHubVnetGateways != nil {
			connection["allow_remote_vnet_to_use_hub_vnet_gateways"] = v.AllowRemoteVnetToUseHubVnetGateways
		}

		if v.EnableInternetSecurity != nil {
			connection["enable_internet_security"] = v.EnableInternetSecurity
		}

		if v.RemoteVirtualNetwork != nil {
			connection["remote_virtual_network_id"] = v.RemoteVirtualNetwork.ID
		}

		connections = append(connections, connection)
	}

	return connections
}
