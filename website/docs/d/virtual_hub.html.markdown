---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub"
sidebar_current: "docs-azurerm-datasource-virtual-hub"
description: |-
  Gets information about an existing Virtual Hub
---

# Data Source: azurerm_virtual_hub

Use this data source to access information about an existing Virtual Hub.


## Vitual Hub Usage

```hcl
data "azurerm_virtual_hub" "example" {
  resource_group = "acctestRG"
  name           = "acctestvhub"
}

output "virtual_hub_id" {
  value = "${data.azurerm_virtual_hub.example.id}"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the VirtualHub.

* `resource_group_name` - (Required) The Name of the Resource Group where the App Service exists.


## Attributes Reference

The following attributes are exported:

* `location` - Resource location.

* `address_prefix` - Address-prefix for this VirtualHub.

* `express_route_gateway` - One `express_route_gateway` block defined below.

* `p2svpn_gateway` - One `p2svpn_gateway` block defined below.

* `route_table` - One `route_table` block defined below.

* `type` - Resource type.

* `virtual_network_connections` - One or more `virtual_network_connection` block defined below.

* `virtual_wan` - One `virtual_wan` block defined below.

* `vpn_gateway` - One `vpn_gateway` block defined below.

* `tags` - Resource tags.


---

The `express_route_gateway` block contains the following:

* `id` - Resource ID.

---

The `p2svpn_gateway` block contains the following:

* `id` - Resource ID.

---

The `route_table` block contains the following:

* `routes` - One or more `route` block defined below.


---

The `route` block contains the following:

* `address_prefixes` - List of all addressPrefixes.

* `next_hop_ip_address` - NextHop ip address.

---

The `virtual_network_connection` block contains the following:

* `id` - Resource ID.

* `remote_virtual_network` - One `remote_virtual_network` block defined below.

* `allow_hub_to_remote_vnet_transit` - VirtualHub to RemoteVnet transit to enabled or not.

* `allow_remote_vnet_to_use_hub_vnet_gateways` - Allow RemoteVnet to use Virtual Hub's gateways.

* `enable_internet_security` - Enable internet security.

* `name` - The name of the resource that is unique within a resource group. This name can be used to access the resource.


---

The `remote_virtual_network` block contains the following:

* `id` - Resource ID.

---

The `virtual_wan` block contains the following:

* `id` - Resource ID.

---

The `vpn_gateway` block contains the following:

* `id` - Resource ID.
