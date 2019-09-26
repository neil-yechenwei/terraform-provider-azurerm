---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_pool"
sidebar_current: "docs-azurerm-datasource-netapp-pool"
description: |-
Gets information about an existing NetApp Pool
---

# Data Source: azurerm_netapp_pool

Use this data source to access information about an existing NetApp Pool.


## NetApp Pool Usage

```hcl
data "azurerm_netapp_pool" "test" {
  resource_group_name = "${azurerm_netapp_pool.test.resource_group_name}"
  account_name        = "${azurerm_netapp_pool.test.account_name}"
  name                = "${azurerm_netapp_pool.test.name}"
}

output "netapp_pool_id" {
  value = "${data.azurerm_netapp_pool.example.id}"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Pool.

* `account_name` - (Required) The name of the NetApp account.

* `resource_group_name` - (Required) The Name of the Resource Group where the NetApp Pool exists.


## Attributes Reference

The following attributes are exported:

* `location` - Resource location.

* `service_level` - The service level of the file system.

* `size` - Provisioned size of the pool (in bytes). Allowed values are in 4TiB chunks (value must be multiply of 4398046511104).

* `tags` - Resource tags.

---
