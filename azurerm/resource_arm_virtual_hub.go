package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
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

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"address_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"express_route_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"p2s_vpn_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"route_table": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
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

			"virtual_network_connections": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},

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

						"remote_virtual_network_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"virtual_wan_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmVirtualHubCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.VirtualHubClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Virtual Hub creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	addressPrefix := d.Get("address_prefix").(string)
	virtualWanId := d.Get("virtual_wan_id").(string)
	tags := d.Get("tags").(map[string]interface{})

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_virtual_hub", *existing.ID)
		}
	}

	hub := network.VirtualHub{
		Location: utils.String(location),
		Tags:     expandTags(tags),
		VirtualHubProperties: &network.VirtualHubProperties{
			AddressPrefix: &addressPrefix,
			VirtualWan: &network.SubResource{
				ID: utils.String(virtualWanId),
			},
		},
	}

	expressRouteGatewayId := d.Get("express_route_gateway_id").(string)
	if expressRouteGatewayId != "" {
		hub.VirtualHubProperties.ExpressRouteGateway.ID = &expressRouteGatewayId
	}

	p2SVpnGatewayId := d.Get("p2s_vpn_gateway_id").(string)
	if p2SVpnGatewayId != "" {
		hub.VirtualHubProperties.P2SVpnGateway.ID = &p2SVpnGatewayId
	}

	vpnGatewayId := d.Get("vpn_gateway_id").(string)
	if vpnGatewayId != "" {
		hub.VirtualHubProperties.VpnGateway.ID = &vpnGatewayId
	}

	virtualHubRoutes := expandArmRouteTable(d)
	if virtualHubRoutes != nil {
		hub.RouteTable = &network.VirtualHubRouteTable{
			Routes: virtualHubRoutes,
		}
	}

	hubVirtualNetworkConnections := expandArmHubVirtualNetworkConnections(d)
	if hubVirtualNetworkConnections != nil {
		hub.VirtualNetworkConnections = hubVirtualNetworkConnections
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, hub)
	if err != nil {
		return fmt.Errorf("Error creating/updating Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation/update of Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Virtual Hub %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmVirtualHubRead(d, meta)
}

func resourceArmVirtualHubRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.VirtualHubClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["virtualHubs"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Virtual Hub %q (Resource Group %q) was not found - removing from state", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
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
			virtualHubRouteTable := flattenArmVirtualHubRouteTable(props.RouteTable.Routes)
			if virtualHubRouteTable != nil {
				d.Set("route_table", virtualHubRouteTable)
			}
		}

		if props.VirtualNetworkConnections != nil {
			virtualNetworkConnections := flattenArmHubVirtualNetworkConnections(props.VirtualNetworkConnections)
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

func resourceArmVirtualHubDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.VirtualHubClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["virtualHubs"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		// deleted outside of Terraform
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error deleting Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for the deletion of Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandArmRouteTable(d *schema.ResourceData) *[]network.VirtualHubRoute {
	data := d.Get("route_table").([]interface{})
	if len(data) == 0 {
		return nil
	}

	virtualHubRouteTable := make([]network.VirtualHubRoute, 0)
	for _, v := range data {
		route := v.(map[string]interface{})
		virtualHubRoute := network.VirtualHubRoute{}

		nextHopIpAddress := route["next_hop_ip_address"].(string)
		if nextHopIpAddress != "" {
			virtualHubRoute.NextHopIPAddress = &nextHopIpAddress
		}

		addressPrefixes := route["address_prefixes"].([]interface{})
		if addressPrefixes != nil {
			addressPrefixeList := make([]string, len(addressPrefixes))
			for i, v := range addressPrefixes {
				addressPrefixeList[i] = fmt.Sprint(v)
			}
			virtualHubRoute.AddressPrefixes = &addressPrefixeList
		}

		virtualHubRouteTable = append(virtualHubRouteTable, virtualHubRoute)
	}

	return &virtualHubRouteTable
}

func expandArmHubVirtualNetworkConnections(d *schema.ResourceData) *[]network.HubVirtualNetworkConnection {
	data := d.Get("virtual_network_connections").([]interface{})
	if len(data) == 0 {
		return nil
	}

	hubVirtualNetworkConnections := make([]network.HubVirtualNetworkConnection, 0)
	for _, v := range data {
		connection := v.(map[string]interface{})
		hubVirtualNetworkConnection := network.HubVirtualNetworkConnection{}

		id := connection["id"].(string)
		if id != "" {
			hubVirtualNetworkConnection.ID = &id
		}

		name := connection["name"].(string)
		if name != "" {
			hubVirtualNetworkConnection.Name = &name
		}

		allowHubToRemoteVnetTransit := connection["allow_hub_to_remote_vnet_transit"].(bool)
		hubVirtualNetworkConnection.AllowHubToRemoteVnetTransit = &allowHubToRemoteVnetTransit

		allowRemoteVnetToUseHubVnetGateways := connection["allow_remote_vnet_to_use_hub_vnet_gateways"].(bool)
		hubVirtualNetworkConnection.AllowRemoteVnetToUseHubVnetGateways = &allowRemoteVnetToUseHubVnetGateways

		enableInternetSecurity := connection["enable_internet_security"].(bool)
		hubVirtualNetworkConnection.EnableInternetSecurity = &enableInternetSecurity

		remoteVirtualNetworkId := connection["remote_virtual_network_id"].(string)
		log.Print("[neil] %s", remoteVirtualNetworkId)
		if remoteVirtualNetworkId != "" {
			hubVirtualNetworkConnection.RemoteVirtualNetwork.ID = &remoteVirtualNetworkId
		}

		hubVirtualNetworkConnections = append(hubVirtualNetworkConnections, hubVirtualNetworkConnection)
	}

	return &hubVirtualNetworkConnections
}

func flattenArmVirtualHubRouteTable(routes *[]network.VirtualHubRoute) []interface{} {
	results := make([]interface{}, 0)

	if routes == nil {
		return results
	}

	for _, v := range *routes {
		route := make(map[string]interface{})
		if v.NextHopIPAddress != nil {
			route["next_hop_ip_address"] = v.NextHopIPAddress
		}

		if v.AddressPrefixes != nil {
			route["address_prefixes"] = v.AddressPrefixes
		}

		results = append(results, route)
	}

	return results
}

func flattenArmHubVirtualNetworkConnections(input *[]network.HubVirtualNetworkConnection) []interface{} {
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
