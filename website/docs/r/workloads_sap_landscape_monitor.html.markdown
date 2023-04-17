---
subcategory: "workloads"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_workloads_sap_landscape_monitor"
description: |-
  Manages a Workloads Sap Landscape Monitor.
---

# azurerm_workloads_sap_landscape_monitor

Manages a Workloads Sap Landscape Monitor.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_workloads_monitor" "example" {
  name                = "example-wm"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_workloads_sap_landscape_monitor" "example" {
  name                 = "example-wslm"
  workloads_monitor_id = azurerm_workloads_monitor.test.id
  grouping {
    landscape {
      name    = ""
      top_sid = []
    }
    sap_application {
      name    = ""
      top_sid = []
    }
  }
  top_metrics_thresholds {
    green  = 0.0
    name   = ""
    red    = 0.0
    yellow = 0.0
  }

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Workloads Sap Landscape Monitor. Changing this forces a new Workloads Sap Landscape Monitor to be created.

* `workloads_monitor_id` - (Required) Specifies the ID of the Workloads Sap Landscape Monitor. Changing this forces a new Workloads Sap Landscape Monitor to be created.

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

* `id` - The ID of the Workloads Sap Landscape Monitor.



## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Workloads Sap Landscape Monitor.
* `read` - (Defaults to 5 minutes) Used when retrieving the Workloads Sap Landscape Monitor.
* `update` - (Defaults to 30 minutes) Used when updating the Workloads Sap Landscape Monitor.
* `delete` - (Defaults to 30 minutes) Used when deleting the Workloads Sap Landscape Monitor.

## Import

Workloads Sap Landscape Monitor can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_workloads_sap_landscape_monitor.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Workloads/monitors/monitor1/sapLandscapeMonitor/default
```
