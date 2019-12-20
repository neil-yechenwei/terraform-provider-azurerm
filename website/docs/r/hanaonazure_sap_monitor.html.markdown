---
subcategory: "HanaOnAzure"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hanaonazure_sap_monitor"
sidebar_current: "docs-azurerm-resource-hanaonazure-sap-monitor"
description: |-
  Manages a HanaOnAzure Sap Monitor.
---

# azurerm_hanaonazure_sap_monitor

Manages a HanaOnAzure Sap Monitor.

## HanaOnAzure Sap Monitor Usage

```hcl
data "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = "${data.azurerm_resource_group.example.name}"
  virtual_network_name = "PV1-vnet"
}

resource "azurerm_hanaonazure_sap_monitor" "example" {
  name                = "examplehanaonazuresapmonitor"
  resource_group_name = "${data.azurerm_resource_group.example.name}"
  location            = "${data.azurerm_resource_group.example.location}"
  hana_host_name      = "10.0.0.6"
  hana_subnet_id      = "${data.azurerm_subnet.example.id}"
  hana_db_name        = "SYSTEMDB"
  hana_db_sql_port    = 30215
  hana_db_username    = "SYSTEM"
  hana_db_password    = "TestPWD"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the HanaOnAzure Sap Monitor. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the HanaOnAzure Sap Monitor should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `hana_host_name` - (Required) The hostname of the HANA instance.

* `hana_subnet_id` - (Required) The subnet to create Azure Monitor for SAP Solutions into.

* `hana_db_name` - (Required) The database name of the HANA instance.

* `hana_db_sql_port` - (Required) The database port of the HANA instance.

* `hana_db_username` - (Required) The database username of the HANA instance.

* `hana_db_password` - (Optional) The database password of the HANA instance.

* `key_vault_id` - (Optional) The key vault id containing customer's HANA credentials. It conflicts with `hana_db_password`.

* `hana_db_password_key_vault_url` - (Optional) The key vault url link to the password for the HANA database. It conflicts with `hana_db_password`.

* `hana_db_credentials_msi_id` - (Optional) The msi id passed by customer which has access to customer's key vault and to be assigned to the collector vm. It conflicts with `hana_db_password`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the HanaOnAzure Sap Monitor.

* `log_analytics_workspace_arm_id` - The ARM ID of the Log Analytics Workspace that is used for monitoring.

## Import

HanaOnAzure Sap Monitor can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_hanaonazure_sap_monitor.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.HanaOnAzure/sapMonitors/monitor1
```
