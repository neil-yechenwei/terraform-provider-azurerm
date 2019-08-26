---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub"
sidebar_current: "docs-azurerm-datasource-virtual-hub"
description: |-
  Gets information about an existing Virtual Hub.
---

# Data Source: azurerm_virtual_hub

Use this data source to access information about an existing Virtual Hub.

## Example Usage

```hcl
data "azurerm_virtual_hub" "test" {
  name                = "testvhub"
  resource_group_name = "testRG"
  location            = "eastus2"
}

output "virtual_hub_id" {
  value = "${data.azurerm_virtual_hub.test.id}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Virtual Hub.

* `resource_group_name` - (Required) Specifies the name of the resource group the Virtual Hub is located in.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `address_prefix` - (Optional) Address-prefix for this VirtualHub.

* `express_route_gateway_id` - (Optional) The expressRouteGateway associated with this VirtualHub.

* `p2s_vpn_gateway_id` - (Optional) The P2SVpnGateway associated with this VirtualHub.

* `route_table` - (Optional) The routeTable associated with this virtual hub.

* `virtual_wan_id` - (Optional) The VirtualWAN to which the VirtualHub belongs.

* `vpn_gateway_id` - (Optional) The VpnGateway associated with this VirtualHub.

* `tags` - (Optional) A mapping of tags to assign to the Virtual Hub.

---

A `route_table` block supports the following:

* `address_prefixes` - (Optional) List of all addressPrefixes.

* `next_hop_ip_address` - (Optional) NextHop ip address.

* `virtual_network_connections` - List of all vnet connections with this VirtualHub.

---

A `virtual_network_connections` block supports the following:

* `id` - Resource ID of virtual network connection.

* `name` - The name of the resource that is unique within a resource group. This name can be used to access the resource.

* `allow_hub_to_remote_vnet_transit` - VirtualHub to RemoteVnet transit to enabled or not.

* `allow_remote_vnet_to_use_hub_vnet_gateways` - Allow RemoteVnet to use Virtual Hub's gateways.

* `enable_internet_security` - Enable internet security.

* `remote_virtual_network_id` - Reference to the remote virtual network.

## Attributes Reference

* `id` - The ID of the Virtual Hub.
