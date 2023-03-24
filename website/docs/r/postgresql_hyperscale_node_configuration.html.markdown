---
subcategory: "PostgreSQL HyperScale"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_hyperscale_node_configuration"
description: |-
  Sets a PostgreSQL HyperScale Node Configuration value on a Azure PostgreSQL HyperScale Cluster.
---

# azurerm_postgresql_hyperscale_node_configuration

Sets a PostgreSQL HyperScale Node Configuration value on a Azure PostgreSQL HyperScale Cluster.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_postgresql_hyperscale_cluster" "example" {
  name                = "example-cluster"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_postgresql_hyperscale_node_configuration" "example" {
  name       = "array_nulls"
  cluster_id = azurerm_postgresql_hyperscale_cluster.example.id
  value      = "on"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the PostgreSQL HyperScale Node Configuration. Changing this forces a new resource to be created.

* `cluster_id` - (Required) The resource ID of the PostgreSQL HyperScale Cluster where we want to change configuration. Changing this forces a new resource to be created.

* `value` - (Required) The value of the PostgreSQL HyperScale Node Configuration.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the PostgreSQL HyperScale Node Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the PostgreSQL HyperScale Node Configuration.
* `update` - (Defaults to 30 minutes) Used when updating the PostgreSQL HyperScale Node Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the PostgreSQL HyperScale Node Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the PostgreSQL HyperScale Node Configuration.

## Import

PostgreSQL HyperScale Node Configurations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_postgresql_hyperscale_node_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DBforPostgreSQL/serverGroupsv2/cluster1/nodeConfigurations/array_nulls
```
