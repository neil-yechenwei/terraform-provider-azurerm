---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub"
sidebar_current: "docs-azurerm-resource-network-virtual-hub"
description: |-
  Manages a Virtual Hub.

---

# azurerm_virtual_hub

Manages a Virtual Hub.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_virtual_hub" "test" {
   name                = "acctestvhub%d"
   resource_group_name = "${azurerm_resource_group.test.name}"
   location            = "${azurerm_resource_group.test.location}"
   address_prefix      = "10.168.0.0/24"
   virtual_wan_id      = "${azurerm_virtual_wan.test.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Virtual Hub. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Virtual Hub. Changing this forces a new resource to be created.

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

The following attributes are exported:

* `id` - The ID of the Virtual Hub.

## Import

Virtual Hub can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_hub.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/virtualHubs/testvhub
```
