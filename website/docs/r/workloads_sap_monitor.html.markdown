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

resource "azurerm_workloads_sap_monitor" "example" {
  name                        = "example-sapmonitor"
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  app_location                = "East US"
  managed_resource_group_name = "managedRGSAPMonitor"
  routing_preference          = "RouteAll"
  subnet_id                   = azurerm_subnet.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this SAP Monitor. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the SAP Monitor should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the SAP Monitor should exist. Changing this forces a new resource to be created.

* `app_location` - (Required) The app location where the SAP Monitor will be deployed in. Changing this forces a new resource to be created.

* `managed_resource_group_name` - (Required) The name of the managed Resource Group for this SAP Monitor. Changing this forces a new resource to be created.

* `routing_preference` - (Required) The routing preference of the SAP Monitor. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The subnet which the SAP Monitor will be deployed in. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `log_analytics_workspace_id` - (Optional) The resource ID of the Log Analytics Workspace that is used for SAP Monitor. Changing this forces a new resource to be created.

* `zone_redundancy_preference` - (Optional) The preference for zone redundancy on resources created for the SAP Monitor. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the SAP Monitor.

---

An `identity` block supports the following:

* `type` - (Required) The type of Managed Service Identity. Possible values are `UserAssigned`.

* `identity_ids` - (Required) A list of IDs for User Assigned Managed Identity resources to be assigned.

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
