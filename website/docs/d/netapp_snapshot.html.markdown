---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_snapshot"
sidebar_current: "docs-azurerm-datasource-netapp-snapshot"
description: |-
Gets information about an existing NetApp Snapshot
---

# Data Source: azurerm_netapp_snapshot

Use this data source to access information about an existing NetApp Snapshot.


## NetApp Snapshot Usage

```hcl
data "azurerm_netapp_snapshot" "test" {
  resource_group_name = "${azurerm_netapp_snapshot.test.resource_group_name}"
  account_name        = "${azurerm_netapp_snapshot.test.account_name}"
  pool_name           = "${azurerm_netapp_snapshot.test.pool_name}"
  volume_name         = "${azurerm_netapp_snapshot.test.volume_name}"
  name                = "${azurerm_netapp_snapshot.test.name}"
}

output "netapp_snapshot_id" {
  value = "${data.azurerm_netapp_snapshot.example.id}"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Snapshot.

* `account_name` - (Required) The name of the NetApp account.

* `pool_name` - (Required) The name of the NetApp Pool.

* `volume_name` - (Required) The name of the NetApp Volume.

* `resource_group_name` - (Required) The Name of the Resource Group where the NetApp Snapshot exists.

## Attributes Reference

The following attributes are exported:

* `location` - Resource location.

* `snapshot_id` - UUID v4 used to identify the Snapshot.

* `file_system_id` - UUID v4 used to identify the FileSystem.

---
