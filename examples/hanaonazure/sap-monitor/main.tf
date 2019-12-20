data "azurerm_resource_group" "example" {
  name     = "HanaOnAzure-SapMonitor-RG5"
}

data "azurerm_subnet" "example" {
  name                 = "hdb-subnet"
  resource_group_name  = "${data.azurerm_resource_group.example.name}"
  virtual_network_name = "PV1-vnet"
}

resource "azurerm_hanaonazure_sap_monitor" "example" {
  name                = "${var.prefix}hanaonazuresapmonitor"
  resource_group_name = "${data.azurerm_resource_group.example.name}"
  location            = "${data.azurerm_resource_group.example.location}"
  hana_host_name      = "10.0.0.6"
  hana_subnet_id      = "${data.azurerm_subnet.example.id}"
  hana_db_name        = "SYSTEMDB"
  hana_db_sql_port    = 30215
  hana_db_username    = "SYSTEM"
  hana_db_password    = "${var.hana_db_password}"
}
