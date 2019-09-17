package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-06-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVirtualHub() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualHubCreateUpdate,
		Read:   resourceArmVirtualHubRead,
		Update: resourceArmVirtualHubCreateUpdate,
		Delete: resourceArmVirtualHubDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"address_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"express_route_gateway": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"p2svpn_gateway": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"route_table": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"routes": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address_prefixes": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"next_hop_ip_address": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"tags": tags.Schema(),

			"virtual_network_connections": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_hub_to_remote_vnet_transit": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"allow_remote_vnet_to_use_hub_vnet_gateways": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"enable_internet_security": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"remote_virtual_network": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"virtual_wan": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"vpn_gateway": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmVirtualHubCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.VirtualHubClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for present of existing Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_virtual_hub", *resp.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	addressPrefix := d.Get("address_prefix").(string)
	expressRouteGateway := d.Get("express_route_gateway").([]interface{})
	p2svpnGateway := d.Get("p2svpn_gateway").([]interface{})
	routeTable := d.Get("route_table").([]interface{})
	virtualNetworkConnections := d.Get("virtual_network_connections").([]interface{})
	virtualWan := d.Get("virtual_wan").([]interface{})
	vpnGateway := d.Get("vpn_gateway").([]interface{})
	t := d.Get("tags").(map[string]interface{})

	virtualHubParameters := network.VirtualHub{
		Location: utils.String(location),
		VirtualHubProperties: &network.VirtualHubProperties{
			AddressPrefix:             utils.String(addressPrefix),
			ExpressRouteGateway:       expandArmVirtualHubSubResource(expressRouteGateway),
			P2SVpnGateway:             expandArmVirtualHubSubResource(p2svpnGateway),
			RouteTable:                expandArmVirtualHubVirtualHubRouteTable(routeTable),
			VirtualNetworkConnections: expandArmVirtualHubHubVirtualNetworkConnection(virtualNetworkConnections),
			VirtualWan:                expandArmVirtualHubSubResource(virtualWan),
			VpnGateway:                expandArmVirtualHubSubResource(vpnGateway),
		},
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, virtualHubParameters)
	if err != nil {
		return fmt.Errorf("Error creating Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Virtual Hub %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmVirtualHubRead(d, meta)
}

func resourceArmVirtualHubRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.VirtualHubClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["virtualHubs"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Virtual Hub %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if virtualHubProperties := resp.VirtualHubProperties; virtualHubProperties != nil {
		d.Set("address_prefix", virtualHubProperties.AddressPrefix)
		if err := d.Set("express_route_gateway", flattenArmVirtualHubSubResource(virtualHubProperties.ExpressRouteGateway)); err != nil {
			return fmt.Errorf("Error setting `express_route_gateway`: %+v", err)
		}
		if err := d.Set("p2svpn_gateway", flattenArmVirtualHubSubResource(virtualHubProperties.P2SVpnGateway)); err != nil {
			return fmt.Errorf("Error setting `p2svpn_gateway`: %+v", err)
		}
		if err := d.Set("route_table", flattenArmVirtualHubVirtualHubRouteTable(virtualHubProperties.RouteTable)); err != nil {
			return fmt.Errorf("Error setting `route_table`: %+v", err)
		}
		if err := d.Set("virtual_network_connections", flattenArmVirtualHubHubVirtualNetworkConnection(virtualHubProperties.VirtualNetworkConnections)); err != nil {
			return fmt.Errorf("Error setting `virtual_network_connections`: %+v", err)
		}
		if err := d.Set("virtual_wan", flattenArmVirtualHubSubResource(virtualHubProperties.VirtualWan)); err != nil {
			return fmt.Errorf("Error setting `virtual_wan`: %+v", err)
		}
		if err := d.Set("vpn_gateway", flattenArmVirtualHubSubResource(virtualHubProperties.VpnGateway)); err != nil {
			return fmt.Errorf("Error setting `vpn_gateway`: %+v", err)
		}
	}
	d.Set("type", resp.Type)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmVirtualHubDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.VirtualHubClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["virtualHubs"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandArmVirtualHubSubResource(input []interface{}) *network.SubResource {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	id := v["id"].(string)

	result := network.SubResource{
		ID: utils.String(id),
	}
	return &result
}

func expandArmVirtualHubVirtualHubRouteTable(input []interface{}) *network.VirtualHubRouteTable {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	routes := v["routes"].([]interface{})

	result := network.VirtualHubRouteTable{
		Routes: expandArmVirtualHubVirtualHubRoute(routes),
	}
	return &result
}

func expandArmVirtualHubHubVirtualNetworkConnection(input []interface{}) *[]network.HubVirtualNetworkConnection {
	results := make([]network.HubVirtualNetworkConnection, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		id := v["id"].(string)
		remoteVirtualNetwork := v["remote_virtual_network"].([]interface{})
		allowHubToRemoteVnetTransit := v["allow_hub_to_remote_vnet_transit"].(bool)
		allowRemoteVnetToUseHubVnetGateways := v["allow_remote_vnet_to_use_hub_vnet_gateways"].(bool)
		enableInternetSecurity := v["enable_internet_security"].(bool)
		name := v["name"].(string)

		result := network.HubVirtualNetworkConnection{
			ID:   utils.String(id),
			Name: utils.String(name),
			HubVirtualNetworkConnectionProperties: &network.HubVirtualNetworkConnectionProperties{
				AllowHubToRemoteVnetTransit:         utils.Bool(allowHubToRemoteVnetTransit),
				AllowRemoteVnetToUseHubVnetGateways: utils.Bool(allowRemoteVnetToUseHubVnetGateways),
				EnableInternetSecurity:              utils.Bool(enableInternetSecurity),
				RemoteVirtualNetwork:                expandArmVirtualHubSubResource(remoteVirtualNetwork),
			},
		}

		results = append(results, result)
	}
	return &results
}

func expandArmVirtualHubVirtualHubRoute(input []interface{}) *[]network.VirtualHubRoute {
	results := make([]network.VirtualHubRoute, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		addressPrefixes := v["address_prefixes"].([]interface{})
		nextHopIpAddress := v["next_hop_ip_address"].(string)

		result := network.VirtualHubRoute{
			AddressPrefixes:  utils.ExpandStringSlice(addressPrefixes),
			NextHopIPAddress: utils.String(nextHopIpAddress),
		}

		results = append(results, result)
	}
	return &results
}

func flattenArmVirtualHubSubResource(input *network.SubResource) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if id := input.ID; id != nil {
		result["id"] = *id
	}

	return []interface{}{result}
}

func flattenArmVirtualHubVirtualHubRouteTable(input *network.VirtualHubRouteTable) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	result["routes"] = flattenArmVirtualHubVirtualHubRoute(input.Routes)

	return []interface{}{result}
}

func flattenArmVirtualHubHubVirtualNetworkConnection(input *[]network.HubVirtualNetworkConnection) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		if id := item.ID; id != nil {
			v["id"] = *id
		}
		if name := item.Name; name != nil {
			v["name"] = *name
		}
		if hubVirtualNetworkConnectionProperties := item.HubVirtualNetworkConnectionProperties; hubVirtualNetworkConnectionProperties != nil {
			if allowHubToRemoteVnetTransit := hubVirtualNetworkConnectionProperties.AllowHubToRemoteVnetTransit; allowHubToRemoteVnetTransit != nil {
				v["allow_hub_to_remote_vnet_transit"] = *allowHubToRemoteVnetTransit
			}
			if allowRemoteVnetToUseHubVnetGateways := hubVirtualNetworkConnectionProperties.AllowRemoteVnetToUseHubVnetGateways; allowRemoteVnetToUseHubVnetGateways != nil {
				v["allow_remote_vnet_to_use_hub_vnet_gateways"] = *allowRemoteVnetToUseHubVnetGateways
			}
			if enableInternetSecurity := hubVirtualNetworkConnectionProperties.EnableInternetSecurity; enableInternetSecurity != nil {
				v["enable_internet_security"] = *enableInternetSecurity
			}
			v["remote_virtual_network"] = flattenArmVirtualHubSubResource(hubVirtualNetworkConnectionProperties.RemoteVirtualNetwork)
		}

		results = append(results, v)
	}

	return results
}

func flattenArmVirtualHubVirtualHubRoute(input *[]network.VirtualHubRoute) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		v["address_prefixes"] = utils.FlattenStringSlice(item.AddressPrefixes)
		if nextHopIpAddress := item.NextHopIPAddress; nextHopIpAddress != nil {
			v["next_hop_ip_address"] = *nextHopIpAddress
		}

		results = append(results, v)
	}

	return results
}
