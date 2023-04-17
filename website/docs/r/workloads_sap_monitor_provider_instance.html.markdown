---
subcategory: "Workloads"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_workloads_sap_monitor_provider_instance"
description: |-
  Manages a SAP Monitor Provider Instance.
---

# azurerm_workloads_sap_monitor_provider_instance

Manages a SAP Monitor Provider Instance.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_workloads_sap_monitor" "example" {
  name                = "example-wm"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_workloads_sap_monitor_provider_instance" "example" {
  name       = "example-wpi"
  monitor_id = azurerm_workloads_sap_monitor.test.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this SAP Monitor Provider Instance. Changing this forces a new resource to be created.

* `monitor_id` - (Required) The resource ID of the SAP Monitor Provider Instance. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `prometheus_os_provider_settings` - (Optional) A `prometheus_os_provider_settings` block as defined below.

---

An `identity` block supports the following:

* `type` - (Required) The type of Managed Service Identity. Possible values are `UserAssigned`.

* `identity_ids` - (Required) A list of IDs for User Assigned Managed Identity resources to be assigned.

---

A `prometheus_os_provider_settings` block supports the following:

* `prometheus_url` - (Optional) The prometheus url.

* `sid` - (Optional) The SID.

* `ssl_certificate_uri` - (Optional) The ssl certificate uri.

* `ssl_preference` - (Optional) The ssl preference.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the SAP Monitor Provider Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the SAP Monitor Provider Instance.
* `read` - (Defaults to 5 minutes) Used when retrieving the SAP Monitor Provider Instance.
* `delete` - (Defaults to 30 minutes) Used when deleting the SAP Monitor Provider Instance.

## Import

SAP Monitor Provider Instances can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_workloads_sap_monitor_provider_instance.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Workloads/monitors/monitor1/providerInstances/providerInstance1
```
