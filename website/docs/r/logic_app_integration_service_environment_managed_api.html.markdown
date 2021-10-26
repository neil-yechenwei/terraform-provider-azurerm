---
subcategory: "Logic App"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_logic_app_integration_service_environment_managed_api"
description: |-
  Manages a Logic App Integration Service Environment Managed Api.
---

# azurerm_logic_app_integration_service_environment_managed_api

Manages a Logic App Integration Service Environment Managed Api.

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
  address_space       = ["10.0.0.0/22"]
}

resource "azurerm_subnet" "isesubnet1" {
  name                 = "isesubnet1"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/26"]

  delegation {
    name = "integrationServiceEnvironments"
    service_delegation {
      name    = "Microsoft.Logic/integrationServiceEnvironments"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_subnet" "isesubnet2" {
  name                 = "isesubnet2"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.64/26"]
}

resource "azurerm_subnet" "isesubnet3" {
  name                 = "isesubnet3"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.128/26"]
}

resource "azurerm_subnet" "isesubnet4" {
  name                 = "isesubnet4"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.192/26"]
}

resource "azurerm_integration_service_environment" "example" {
  name                 = "example-ise"
  location             = azurerm_resource_group.example.location
  resource_group_name  = azurerm_resource_group.example.name
  sku_name             = "Premium_0"
  access_endpoint_type = "Internal"

  virtual_network_subnet_ids = [
    azurerm_subnet.isesubnet1.id,
    azurerm_subnet.isesubnet2.id,
    azurerm_subnet.isesubnet3.id,
    azurerm_subnet.isesubnet4.id
  ]
}

resource "azurerm_logic_app_integration_service_environment_managed_api" "example" {
  name                               = "servicebus"
  location                           = azurerm_resource_group.example.location
  resource_group_name                = azurerm_resource_group.example.name
  integration_service_environment_id = azurerm_integration_service_environment.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Logic App Integration Service Environment Managed Api. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Logic App Integration Service Environment Managed Api should exist. Changing this forces a new resource to be created.

* `integration_service_environment_id` - (Required) The resource ID of the Logic App Integration Service Environment. Changing this forces a new resource to be created.

* `deployment_content_link_definition_uri` - (Optional) The content link of deployment definition for the Logic App Integration Service Environment.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Logic App Integration Service Environment Managed Api.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Logic App Integration Service Environment Managed Api.
* `update` - (Defaults to 30 minutes) Used when updating the Logic App Integration Service Environment Managed Api.
* `read` - (Defaults to 5 minutes) Used when retrieving the Logic App Integration Service Environment Managed Api.
* `delete` - (Defaults to 30 minutes) Used when deleting the Logic App Integration Service Environment Managed Api.

## Import

Logic App Integration Service Environment Managed Apis can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_logic_app_integration_service_environment_managed_api.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Logic/integrationServiceEnvironments/ise1/managedApis/servicebus
```
