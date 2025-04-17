---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_network_gateway_learned_routes"
description: |-
  Gets information about an existing Virtual Network Gateway Learned Routes.
---

# Data Source: azurerm_virtual_network_gateway_learned_routes

Use this data source to access information about an existing Virtual Network Gateway Learned Routes.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "example" {
  name                = "example-pip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Dynamic"
  sku                 = "Basic"
}

resource "azurerm_virtual_network_gateway" "example" {
  name                = "example-vng"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "Basic"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.example.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.example.id
  }
}

data "azurerm_virtual_network_gateway_learned_routes" "example" {
  virtual_network_gateway_id = azurerm_virtual_network_gateway.example.id
}

output "virtual_network_gateway_learned_routes" {
  value = data.azurerm_virtual_network_gateway_learned_routes.example.gateway_routes
}
```

## Argument Reference

* `virtual_network_gateway_id` - The ID of the Virtual Network Gateway.

## Attributes Reference

* `as_path` - The route's AS path sequence.

* `local_address` - The gateway's local address.

* `network` - The route's network prefix.

* `next_hop` - The route's next hop.

* `origin` - The source this route was learned from.

* `source_peer` - The peer this route was learned from.

* `weight` - The route's weight.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Network Gateway Learned Routes.
