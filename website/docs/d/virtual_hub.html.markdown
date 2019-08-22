---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_hub"
sidebar_current: "docs-azurerm-datasource-virtual-hub"
description: |-
  Gets information about an existing Virtual Hub.
---

# Data Source: azurerm_virtual_hub

Use this data source to access information about an existing Virtual Hub.

## Example Usage

```hcl
data "azurerm_virtual_hub" "test" {
  name                = "testvhub"
  resource_group_name = "testRG"
  location            = "eastus2"
}

output "virtual_hub_id" {
  value = "${data.azurerm_virtual_hub.test.id}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Virtual Hub.

* `resource_group_name` - (Required) Specifies the name of the resource group the Virtual Hub is located in.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `address_prefix` - (Optional) Address-prefix for this VirtualHub.

* `virtual_wan_id` - (Optional) The VirtualWAN to which the VirtualHub belongs.

* `tags` - (Optional) A mapping of tags to assign to the Virtual Hub.

## Attributes Reference

* `id` - The ID of the Virtual Hub.
