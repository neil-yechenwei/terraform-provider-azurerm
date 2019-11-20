resource "azurerm_resource_group" "test" {
  name     = "${var.resource_group_name}"
  location = "${var.location}"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.5.1.0/24"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_virtual_hub" "test" {
  name           = "acctestvirtualhub"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location       = "${azurerm_resource_group.test.location}"
  address_prefix = "10.0.1.0/24"

  virtual_wan_id = "${azurerm_virtual_wan.test.id}"

  virtual_network_connections {
	  name = "testConnection"
	
    remote_virtual_network_id = "${azurerm_virtual_network.test.id}"

    allow_hub_to_remote_vnet_transit           = "false"
    allow_remote_vnet_to_use_hub_vnet_gateways = "false"
    enable_internet_security                   = "false"
  }

  route_table {
    routes {
      address_prefixes    = ["10.0.2.0/24", "10.0.3.0/24"]
      next_hop_ip_address = "10.0.4.6"
    }
  }

  tags = {
    env = "test"
  }
}