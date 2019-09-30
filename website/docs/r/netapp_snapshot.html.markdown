---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_snapshot"
sidebar_current: "docs-azurerm-resource-netapp-snapshot"
description: |-
  Manage Azure NetApp Snapshot instance.
---

# azurerm_netapp_snapshot

Manage Azure NetApp Snapshot instance.


## NetApp Snapshot Usage

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

  delegation {
    name = "acctestdelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_netapp_account" "test" {
  name                = "acctestaccount"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctestpool"
  account_name        = "${azurerm_netapp_account.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  service_level       = "Premium"
  size                = "4398046511104"
}

resource "azurerm_netapp_volume" "test" {
  name                = "acctestvolume"
  account_name        = "${azurerm_netapp_account.test.name}"
  pool_name           = "${azurerm_netapp_pool.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  creation_token      = "my-unique-file-path"
  service_level       = "Premium"
  subnet_id           = "${azurerm_subnet.test.id}"
  usage_threshold     = "107374182400"
}

resource "azurerm_netapp_snapshot" "test" {
  name                = "acctestsnapshot"
  account_name        = "${azurerm_netapp_account.test.name}"
  pool_name           = "${azurerm_netapp_pool.test.name}"
  volume_name         = "${azurerm_netapp_volume.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  file_system_id      = "${azurerm_netapp_volume.test.file_system_id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Snapshot. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The resource group name of the NetApp Snapshot. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the NetApp Account.

* `pool_name` - (Required) The name of the NetApp Pool.

* `volume_name` - (Required) The name of the NetApp Volume.

* `location` - (Required) Resource location. Changing this forces a new resource to be created.

* `file_system_id` - (Optional) UUID v4 used to identify the FileSystem.

---

## Attributes Reference

The following attributes are exported:

* `id` - Resource id.

* `snapshot_id` - UUID v4 used to identify the Snapshot.

## Import

NetApp Snapshot can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_netapp_snapshot.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG/providers/Microsoft.NetApp/netAppAccounts/acctestaccount/capacityPools/acctestpool/volumes/acctestpool/snapshots/
```