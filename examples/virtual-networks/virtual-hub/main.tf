resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_virtual_network" "example" {
  name                = "acctest-VirtualNetwork-neil"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_network_security_group" "example" {
  name                = "acctest-NetworkSecurityGroup-neil"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

resource "azurerm_subnet" "example" {
  name                      = "acctest-Subnet-neil"
  resource_group_name       = "${azurerm_resource_group.example.name}"
  virtual_network_name      = "${azurerm_virtual_network.example.name}"
  address_prefix            = "10.5.1.0/24"
  network_security_group_id = "${azurerm_network_security_group.example.id}"
}

resource "azurerm_virtual_wan" "example" {
  name                = "${var.prefix}-virtualwan"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
}

resource "azurerm_virtual_hub" "example" {
  name                = "${var.prefix}-virtualhub"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  address_prefix      = "10.0.1.0/24"
  virtual_wan_id      = "${azurerm_virtual_wan.example.id}"

  virtual_network_connection {
    name                                       = "testConnection2"
    remote_virtual_network_id                  = "${azurerm_virtual_network.example.id}"
    allow_hub_to_remote_vnet_transit           = false
    allow_remote_vnet_to_use_hub_vnet_gateways = false
    enable_internet_security                   = false
  }
}
