---
subcategory: "HanaOnAzure"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hanaonazure_sap_monitor"
sidebar_current: "docs-azurerm-datasource-hanaonazure-sap-monitor"
description: |-
Gets information about an existing HanaOnAzure Sap Monitor
---

# Data Source: azurerm_hanaonazure_sap_monitor

Uses this data source to access information about an existing HanaOnAzure Sap Monitor.

## NetApp Pool Usage

```hcl
data "azurerm_hanaonazure_sap_monitor" "example" {
  resource_group_name = "acctestRG"
  name                = "acctesthanaonazuresapmonitor"
}
output "hanaonazure_sap_monitor_id" {
  value = "${data.azurerm_hanaonazure_sap_monitor.example.id}"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the HanaOnAzure Sap Monitor.

* `resource_group_name` - (Required) The Name of the Resource Group where the HanaOnAzure Sap Monitor exists.


## Attributes Reference

The following attributes are exported:

* `location` - The Azure Region where the HanaOnAzure Sap Monitor exists.

* `hana_db_username` - The database username of the HANA instance.

* `hana_db_sql_port` - The database port of the HANA instance.

* `hana_host_name` - The hostname of the HANA instance.

* `hana_subnet` - The subnet to create Azure Monitor for SAP Solutions into.

* `hana_db_password` - The database password of the HANA instance.

* `hana_db_name` - The database name of the HANA instance.

* `key_vault_id` - The key vault id containing customer's HANA credentials.

* `hana_db_password_key_vault_url` - The key vault url link to the password for the HANA database.

* `hana_db_credentials_msi_id` - The msi id passed by customer which has access to customer's key vault and to be assigned to the collector vm.

* `log_analytics_workspace_arm_id` - The ARM ID of the Log Analytics Workspace that is used for monitoring.

* `managed_resource_group_name` - The name of the resource group the SAP Monitor resources get deployed into.

* `tags` - A mapping of tags to assign to the resource.
