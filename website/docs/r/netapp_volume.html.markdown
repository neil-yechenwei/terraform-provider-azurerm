---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_volume"
sidebar_current: "docs-azurerm-resource-netapp-volume"
description: |-
  Manage Azure NetApp Volume instance.
---

# azurerm_netapp_volume

Manage Azure NetApp Volume instance.


## NetApp Volume Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acctestRG"
  location = "Eastus2"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctvn"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctsub"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_netapp_account" "test" {
  name                = "acctestnetappaccount"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctestnetapppool"
  account_name        = "${azurerm_netapp_account.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  service_level       = "Premium"
  size                = "4398046511104"
}

resource "azurerm_netapp_volume" "test" {
  name                = "acctestnetappvolume"
  account_name        = "${azurerm_netapp_account.test.name}"
  pool_name           = "${azurerm_netapp_pool.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  creation_token      = "my-unique-file-path"
  service_level       = "Premium"
  subnet_id           = "${azurerm_subnet.test.id}"
  usage_threshold     = "107374182400"

  export_policy {
    rules {
      unix_read_write = "true"
    }
  }

  tags = {
    env = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Volume. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The resource group name of the NetApp Volume. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the NetApp Account.

* `pool_name` - (Required) The name of the NetApp Pool.

* `location` - (Required) Resource location. Changing this forces a new resource to be created.

* `creation_token` - (Required) A unique file path for the volume. Used when creating mount targets.

* `export_policy` - (Optional) One `export_policy` block defined below.

* `service_level` - (Required) The service level of the file system.

* `subnet_id` - (Required) The Azure Resource URI for a delegated subnet. Must have the delegation Microsoft.NetApp/volumes.

* `usage_threshold` - (Required) Maximum storage quota allowed for a file system in bytes. This is a soft quota used for alerting only. Minimum size is 100 GiB. Upper limit is 100TiB.

* `tags` - Resource tags.

---

The `export_policy` block contains the following:

* `rules` - (Optional) One or more `rules` block defined below.

---

The `rules` block contains the following:

* `allowed_clients` - (Optional) Client ingress specification as comma separated string with IPv4 CIDRs, IPv4 host addresses and host names.

* `cifs` - (Optional) Allows CIFS protocol.

* `nfsv3` - (Optional) Allows NFSv3 protocol.

* `nfsv4` - (Optional) Allows NFSv4 protocol.

* `rule_index` - (Optional) Order index.

* `unix_read_only` - (Optional) Read only access.

* `unix_read_write` - (Optional) Read and write access.

---

## Attributes Reference

The following attributes are exported:

* `id` - Resource id.

## Import

NetApp Volume can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_netapp_volume.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG/providers/Microsoft.NetApp/netAppAccounts/acctestnetappaccount/capacityPools/acctestnetapppool/volumes/
```