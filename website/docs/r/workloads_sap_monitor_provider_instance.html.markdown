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

* `prometheus_ha_cluster_provider_settings` - (Optional) A `prometheus_ha_cluster_provider_settings` block as defined below.

* `mssql_server_provider_settings` - (Optional) A `mssql_server_provider_settings` block as defined below.

* `db2_provider_settings` - (Optional) A `db2_provider_settings` block as defined below.

* `sap_hana_provider_settings` - (Optional) A `sap_hana_provider_settings` block as defined below.

* `sap_net_weaver_provider_settings` - (Optional) A `sap_net_weaver_provider_settings` block as defined below.

---

An `identity` block supports the following:

* `type` - (Required) The type of Managed Service Identity. Possible values are `UserAssigned`.

* `identity_ids` - (Required) A list of IDs for User Assigned Managed Identity resources to be assigned.

---

A `prometheus_os_provider_settings` block supports the following:

* `prometheus_url` - (Optional) The URL of the Node Exporter endpoint.

* `sid` - (Optional) The SAP System Identifier.

* `ssl_certificate_uri` - (Optional) The Blob URI to SSL certificate for the prometheus node exporter.

* `ssl_preference` - (Optional) The certificate preference.

---

A `prometheus_ha_cluster_provider_settings` block supports the following:

* `cluster_name` - (Optional) The name of the cluster.

* `host_name` - (Optional) The target machine name.

* `prometheus_url` - (Optional) The URL of the Node Exporter endpoint.

* `sid` - (Optional) The SID of the cluster.

* `ssl_certificate_uri` - (Optional) The Blob URI to SSL certificate for the HA cluster exporter.

* `ssl_preference` - (Optional) The certificate preference.

---

A `mssql_server_provider_settings` block supports the following:

* `db_password` - (Optional) The password of the database.

* `db_password_uri` - (Optional) The Key Vault URI to secret with the database password.

* `db_port` - (Optional) The SQL port of the databse.

* `db_username` - (Optional) The username of the database.

* `host_name` - (Optional) The host name of the SQL Server.

* `sid` - (Optional) The SAP System Identifier.

* `ssl_certificate_uri` - (Optional) The Blob URI to SSL certificate for the SQL Database.

* `ssl_preference` - (Optional) The certificate preference.

---

A `db2_provider_settings` block supports the following:

* `db_name` - (Optional) The name of the db2 database.

* `db_password` - (Optional) The password of the db2 database.

* `db_password_uri` - (Optional) The Key Vault URI to secret with the database password.

* `db_port` - (Optional) The SQL port of the db2 database.

* `db_username` - (Optional) The username of the db2 database.

* `host_name` - (Optional) The name of the target Virtual Machine.

* `sid` - (Optional) The SAP System Identifier.

* `ssl_certificate_uri` - (Optional) The Blob URI to SSL certificate for the SQL Database.

* `ssl_preference` - (Optional) The certificate preference.

---

A `sap_hana_provider_settings` block supports the following:

* `db_name` - (Optional) The name of the HANA database.

* `db_password` - (Optional) The password of the database.

* `db_password_uri` - (Optional) The Key Vault URI to secret with the database password.

* `db_username` - (Optional) The username of the database.

* `host_name` - (Optional) The name of the target Virtual Machine.

* `instance_number` - (Optional) The instance number of the database.

* `sid` - (Optional) The SAP System Identifier.

* `sql_port` - (Optional) The SQL port of the database.

* `ssl_certificate_uri` - (Optional) The Blob URI to SSL certificate for the DB.

* `ssl_host_name_in_certificate` - (Optional) The host name in the SSL certificate.

* `ssl_preference` - (Optional) The certificate preference.

---

A `sap_net_weaver_provider_settings` block supports the following:

* `client_id` - (Optional) The client ID of the SAP.

* `host_file_entries` - (Optional) A list of the host file entries.

* `host_name` - (Optional) The target Virtual Machine IP Address/FQDN.

* `instance_nr` - (Optional) The instance number of the SAP NetWeaver.

* `password` - (Optional) The password of the SAP.

* `password_uri` - (Optional) The Key Vault URI to secret with the SAP password.

* `port_number` - (Optional) The HTTP port number of the SAP.

* `sid` - (Optional) The SAP System Identifier.

* `username` - (Optional) The username of the SAP.

* `ssl_certificate_uri` - (Optional) The Blob URI to SSL certificate for the SAP system.

* `ssl_preference` - (Optional) The certificate preference.

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
