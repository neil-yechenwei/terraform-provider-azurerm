resource "azurerm_resource_group" "test" {
  name     = "${var.resource_group_name}"
  location = "${var.location}"
}

resource "azurerm_virtual_network" "test" {
  name                = "${var.name_prefix}-vnet"
  address_space       = ["10.5.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  subnet {
    name           = "${var.name_prefix}-vsubnet"
    address_prefix = "10.5.1.0/24"
  }
}

resource "azurerm_virtual_wan" "test" {
  name                = "${var.name_prefix}-vwan"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_virtual_hub" "test" {
  name                = "${var.name_prefix}-vhub"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  address_prefix      = "10.0.1.0/24"
  virtual_wan_id      = "${azurerm_virtual_wan.test.id}"
  virtual_network_connections {
    name                                       = "testConnection"
    remote_virtual_network_id                  = "${azurerm_virtual_network.test.id}"
    allow_hub_to_remote_vnet_transit           = "false"
    allow_remote_vnet_to_use_hub_vnet_gateways = "false"
    enable_internet_security                   = "false"
  }
  route_table {
    address_prefixes    = ["10.0.2.0/24", "10.0.3.0/24"]
    next_hop_ip_address = "10.0.4.5"
  }
}
