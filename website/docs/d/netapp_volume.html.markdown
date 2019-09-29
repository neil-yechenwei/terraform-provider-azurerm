---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_volume"
sidebar_current: "docs-azurerm-datasource-netapp-volume"
description: |-
Gets information about an existing NetApp Volume
---

# Data Source: azurerm_netapp_volume

Use this data source to access information about an existing NetApp Volume.


## NetApp Volume Usage

```hcl
data "azurerm_netapp_volume" "test" {
  resource_group_name = "${azurerm_netapp_volume.test.resource_group_name}"
  account_name        = "${azurerm_netapp_volume.test.account_name}"
  pool_name           = "${azurerm_netapp_volume.test.pool_name}"
  name                = "${azurerm_netapp_volume.test.name}"
}

output "netapp_volume_id" {
  value = "${data.azurerm_netapp_volume.example.id}"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Pool.

* `account_name` - (Required) The name of the NetApp account.

* `pool_name` - (Required) The name of the NetApp Pool.

* `resource_group_name` - (Required) The Name of the Resource Group where the NetApp Pool exists.


## Attributes Reference

The following attributes are exported:

* `location` - Resource location.

* `creation_token` - A unique file path for the volume. Used when creating mount targets.

* `export_policy` - One `export_policy` block defined below.

* `service_level` - The service level of the file system.

* `subnet_id` - The Azure Resource URI for a delegated subnet. Must have the delegation Microsoft.NetApp/volumes.

* `usage_threshold` - Maximum storage quota allowed for a file system in bytes. This is a soft quota used for alerting only. Minimum size is 100 GiB. Upper limit is 100TiB.

* `tags` - Resource tags.

---

The `export_policy` block contains the following:

* `rules` - One or more `rules` block defined below.

---

The `rules` block contains the following:

* `allowed_clients` - Client ingress specification as comma separated string with IPv4 CIDRs, IPv4 host addresses and host names.

* `cifs` - Allows CIFS protocol.

* `nfsv3` - Allows NFSv3 protocol.

* `nfsv4` - Allows NFSv4 protocol.

* `rule_index` - Order index.

* `unix_read_only` - Read only access.

* `unix_read_write` - Read and write access.

---
