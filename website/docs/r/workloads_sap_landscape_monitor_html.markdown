---
subcategory: "Workloads"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_workloads_sap_landscape_monitor"
description: |-
  Manages a SAP Landscape Monitor.
---

# azurerm_workloads_sap_landscape_monitor

Manages a SAP Landscape Monitor.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_workloads_sap_landscape_monitor" "example" {
  name                = "example-sapmonitor"
  resource_group_name = azurerm_resource_group.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this SAP Landscape Monitor. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the SAP Landscape Monitor should exist. Changing this forces a new resource to be created.

* `grouping` - (Optional) A `grouping` block as defined below.

* `top_metrics_thresholds` - (Optional) A `top_metrics_thresholds` block as defined below.

---

A `grouping` block supports the following:

* `landscape` - (Optional) A `landscape` block as defined below.

* `sap_application` - (Optional) A `sap_application` block as defined below.

---

A `landscape` block supports the following:

* `name` - (Optional) Gets or sets the name of the grouping.

* `top_sid` - (Optional) Gets or sets the list of SID's.

---

A `sap_application` block supports the following:

* `name` - (Optional) Gets or sets the name of the grouping.

* `top_sid` - (Optional) Gets or sets the list of SID's.

---

A `top_metrics_thresholds` block supports the following:

* `green` - (Optional) Gets or sets the threshold value for Green.

* `name` - (Optional) Gets or sets the name of the threshold.

* `red` - (Optional) Gets or sets the threshold value for Red.

* `yellow` - (Optional) Gets or sets the threshold value for Yellow.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the SAP Landscape Monitor.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the SAP Landscape Monitor.
* `read` - (Defaults to 5 minutes) Used when retrieving the SAP Landscape Monitor.
* `update` - (Defaults to 60 minutes) Used when updating the SAP Landscape Monitor.
* `delete` - (Defaults to 60 minutes) Used when deleting the SAP Landscape Monitor.

## Import

SAP Landscape Monitors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_workloads_sap_landscape_monitor.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Workloads/monitors/monitor1
```
