---
subcategory: "Workloads"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_workloads_sap_monitor"
description: |-
  Manages a SAP Monitor.
---

# azurerm_workloads_sap_monitor

Manages a SAP Monitor.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_workloads_sap_monitor" "example" {
  name                = "example-sapmonitor"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this SAP Monitor. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the SAP Monitor should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the SAP Monitor should exist. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the SAP Monitor.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the SAP Monitor.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the SAP Monitor.
* `read` - (Defaults to 5 minutes) Used when retrieving the SAP Monitor.
* `update` - (Defaults to 60 minutes) Used when updating the SAP Monitor.
* `delete` - (Defaults to 60 minutes) Used when deleting the SAP Monitor.

## Import

SAP Monitors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_workloads_sap_monitor.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Workloads/monitors/monitor1
```
