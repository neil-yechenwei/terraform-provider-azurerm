---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_flexible_server_migration"
description: |-
  Manages a PostgreSQL Flexible Server Migration.
---

# azurerm_postgresql_flexible_server_migration

Manages a PostgreSQL Flexible Server Migration.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_postgresql_flexible_server" "example" {
  name                   = "example-fs"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  version                = "12"
  sku_name               = "GP_Standard_D2s_v3"
  zone                   = "2"
}

resource "azurerm_postgresql_flexible_server_migration" "example" {
  name      = "example-pfsm"
  location  = azurerm_resource_group.example.location
  server_id = azurerm_postgresql_flexible_server.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this PostgreSQL Flexible Server Migration. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the PostgreSQL Flexible Server Migration should exist. Changing this forces a new resource to be created.

* `server_id` - (Required) The ID of the PostgreSQL Flexible Server from which to create this PostgreSQL Flexible Server Migration. Changing this forces a new resource to be created.

* `cancel_enabled` - (Optional) Should cancel be enabled for entire migration?

* `dbs_to_cancel_migration_on` - (Optional) A list of databases to trigger cancel.

* `dbs_to_migrate` - (Optional) A list of databases to migrate.

* `dbs_to_trigger_cutover_on` - (Optional) A list of databases to trigger cutover.

* `migrate_roles_enabled` - (Optional) Should roles and permissions migration be enabled for entire migration?

* `migration_instance_resource_id` - (Optional) The ID of the Private Endpoint migration instance. Changing this forces a new resource to be created.

* `migration_mode` - (Optional) The migration mode. Possible values are `Online` and `Offline`.

* `migration_option` - (Optional) The migration option for the migration. Possible values are `Validate`, `Migrate` and `ValidateAndMigrate`.

* `migration_window_end_time_in_utc` - (Optional) The end time in UTC for migration window. Changing this forces a new resource to be created.

* `migration_window_start_time_in_utc` - (Optional) The start time in UTC for migration window.

* `overwrite_dbs_in_target_enabled` - (Optional) Should the databases on the target server can be overwritten?

* `trigger_cutover_enabled` - (Optional) Should cutover be enabled for entire migration?

* `tags` - (Optional) A mapping of tags which should be assigned to the PostgreSQL Flexible Server Migration.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the PostgreSQL Flexible Server Migration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating this PostgreSQL Flexible Server Migration.
* `delete` - (Defaults to 30 minutes) Used when deleting this PostgreSQL Flexible Server Migration.
* `update` - (Defaults to 30 minutes) Used when updating the PostgreSQL Flexible Server Migration.
* `read` - (Defaults to 5 minutes) Used when retrieving this PostgreSQL Flexible Server Migration.

## Import

An existing PostgreSQL Flexible Server Migration can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_postgresql_flexible_server_migration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DBforPostgreSQL/flexibleServers/fs1/migrations/migration1
```
