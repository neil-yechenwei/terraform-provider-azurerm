---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub"
sidebar_current: "docs-azurerm-resource-virtual-hub"
description: |-
  Manage Azure VirtualHub instance.
---

# azurerm_virtual_hub

Manage Azure VirtualHub instance.


## Vitual Hub Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "acctestRG"
  location = "Eastus2"
}

resource "azurerm_virtual_wan" "example" {
  name                = "acctestvwan-%d"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
}

resource "azurerm_virtual_hub" "example" {
  name                = "acctestvirtualhub-%d"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  address_prefix      = "10.0.1.0/24"

  virtual_wan_id = "${azurerm_virtual_wan.example.id}"

  route_table {
    routes {
      address_prefixes    = ["10.0.2.0/24", "10.0.3.0/24"]
      next_hop_ip_address = "10.0.4.5"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the VirtualHub. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The resource group name of the VirtualHub. Changing this forces a new resource to be created.

* `location` - (Required) Resource location. Changing this forces a new resource to be created.

* `address_prefix` - (Required) Address-prefix for this VirtualHub.

* `express_route_gateway_id` - (Optional) The resource id of express route gateway.

* `p2svpn_gateway_id` - (Optional) The resource id of p2svpn gateway.

* `route_table` - (Optional) One `route_table` block defined below.

* `virtual_network_connections` - (Optional) One or more `virtual_network_connection` block defined below.

* `virtual_wan_id` - (Required) The resource id of virtual wan.

* `vpn_gateway_id` - (Optional) The resource id of vpn gateway.

* `tags` - (Optional) Resource tags. Changing this forces a new resource to be created.

---

The `route_table` block supports the following:

* `routes` - (Optional) One or more `route` block defined below.

---

The `route` block supports the following:

* `address_prefixes` - (Optional) List of all addressPrefixes.

* `next_hop_ip_address` - (Optional) NextHop ip address.

---

The `virtual_network_connection` block supports the following:

* `id` - (Optional) Resource ID.

* `remote_virtual_network_id` - (Optional) The resource id of remote virtual network.

* `allow_hub_to_remote_vnet_transit` - (Optional) VirtualHub to RemoteVnet transit to enabled or not.

* `allow_remote_vnet_to_use_hub_vnet_gateways` - (Optional) Allow RemoteVnet to use Virtual Hub's gateways.

* `enable_internet_security` - (Optional) Enable internet security.

* `name` - (Optional) The name of the resource that is unique within a resource group. This name can be used to access the resource.

---

## Attributes Reference

The following attributes are exported:

* `id` - Resource id.

## Import

Virtual Hub can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_virtual_hub.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG/providers/Microsoft.Network/virtualHubs/
```
