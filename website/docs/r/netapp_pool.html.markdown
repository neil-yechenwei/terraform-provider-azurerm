---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_pool"
sidebar_current: "docs-azurerm-resource-netapp-pool"
description: |-
  Manage Azure NetApp Pool instance.
---

# azurerm_netapp_pool

Manage Azure NetApp Pool instance.


## NetApp Pool Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acctestRG"
  location = "Eastus2"
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

  tags = {
    env = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Pool. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The resource group name of the NetApp Pool. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the NetApp account.

* `location` - (Required) Resource location. Changing this forces a new resource to be created.

* `service_level` - (Required) The service level of the file system.

* `size` - (Required) Provisioned size of the pool (in bytes). Allowed values are in 4TiB chunks (value must be multiply of 4398046511104).

* `tags` - (Optional) Resource tags. Changing this forces a new resource to be created.

---

## Attributes Reference

The following attributes are exported:

* `id` - Resource id.

## Import

NetApp Pool can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_netapp_pool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG/providers/Microsoft.NetApp/netAppAccounts/acctestnetappaccount/capacityPools/
```