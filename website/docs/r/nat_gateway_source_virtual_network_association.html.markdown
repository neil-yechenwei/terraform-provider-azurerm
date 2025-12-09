---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_nat_gateway_source_virtual_network_association"
description: |-
  Manages the association between a NAT Gateway and a Virtual Network.
---

# azurerm_nat_gateway_source_virtual_network_association

Manages the association between a NAT Gateway and a Virtual Network as a source virtual network.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_nat_gateway" "example" {
  name                = "example-NatGateway"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "StandardV2"
  zones               = ["1", "2", "3"]
}

resource "azurerm_nat_gateway_source_virtual_network_association" "example" {
  nat_gateway_id     = azurerm_nat_gateway.example.id
  virtual_network_id = azurerm_virtual_network.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `nat_gateway_id` - (Required) The ID of the NAT Gateway. Changing this forces a new resource to be created.

* `virtual_network_id` - (Required) The ID of the Virtual Network which should be associated with the NAT Gateway as a source virtual network. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The (Terraform specific) ID of the Association between the NAT Gateway and the Virtual Network.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the association between the NAT Gateway and the Virtual Network.
* `read` - (Defaults to 5 minutes) Used when retrieving the association between the NAT Gateway and the Virtual Network.
* `delete` - (Defaults to 30 minutes) Used when deleting the association between the NAT Gateway and the Virtual Network.

## Import

Associations between NAT Gateway and Virtual Networks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_nat_gateway_source_virtual_network_association.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/natGateways/gateway1|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myVirtualNetwork1"
```

-> **Note:** This is a Terraform Specific ID in the format `{natGatewayID}|{virtualNetworkID}`

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network` - 2025-01-01
