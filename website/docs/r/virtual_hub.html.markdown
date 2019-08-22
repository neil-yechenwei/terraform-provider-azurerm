---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub"
sidebar_current: "docs-azurerm-resource-network-virtual-hub"
description: |-
  Manages a Virtual Hub.

---

# azurerm_virtual_hub

Manages a Virtual Hub.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_virtual_hub" "test" {
   name                = "acctestvhub%d"
   resource_group_name = "${azurerm_resource_group.test.name}"
   location            = "${azurerm_resource_group.test.location}"
   address_prefix      = "10.168.0.0/24"
   virtual_wan_id      = "${azurerm_virtual_wan.test.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Virtual Hub. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Virtual Hub. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `address_prefix` - (Optional) Address-prefix for this VirtualHub.

* `virtual_wan_id` - (Optional) The VirtualWAN to which the VirtualHub belongs.

* `tags` - (Optional) A mapping of tags to assign to the Virtual Hub.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual Hub.

## Import

Virtual Hub can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_hub.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/virtualHubs/testvhub
```
