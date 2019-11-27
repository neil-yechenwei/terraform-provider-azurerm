resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-virtualnetwork"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

resource "azurerm_subnet" "example" {
  name                 = "${var.prefix}-subnet"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  virtual_network_name = "${azurerm_virtual_network.example.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_hanaonazure_sap_monitor" "example" {
  name                = "${var.prefix}-hanaonazuresapmonitor"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  hana_db_username    = "SYSTEM"
  hana_db_sql_port    = 30815
  hana_host_name      = "10.0.0.6"
  hana_db_name        = "SYSTEMDB"
  hana_db_password    = "Manager1"
  hana_subnet_id      = "${azurerm_subnet.example.id}"
}
