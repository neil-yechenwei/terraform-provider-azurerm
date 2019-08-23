resource "azurerm_resource_group" "test" {
  name     = "${var.resource_group_name}"
  location = "${var.location}"
}

 resource "azurerm_private_link_service" "test" {
    name = "${var.name}"
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    fqdns = ["testFqdns"]
    ip_configuration {
      name = "${azurerm_public_ip.test.name}"
      subnet_id = "${azurerm_subnet.test.id}"
      private_ip_address = "10.0.1.17"
      private_ip_address_version = "IPv4"
      private_ip_address_allocation = "Static"
    }
    load_balancer_frontend_ip_configuration {
       id = "${azurerm_lb.test.frontend_ip_configuration.0.id}"
    }
    tags = {
      env = "test"
    }
}

 resource "azurerm_virtual_network" "test" {
    name                = "testVnet"
    address_space       = ["10.0.0.0/16"]
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

 resource "azurerm_subnet" "test" {
    name                 = "testSubnet"
    resource_group_name  = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix       = "10.0.1.0/24"
    private_link_service_network_policies = "Disabled"
}

 resource "azurerm_public_ip" "test" {
    name                = "testPip"
    sku                 = "Standard"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    allocation_method   = "Static"
}

 resource "azurerm_lb" "test" {
    name                = "testLb"
    sku                 = "Standard"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    frontend_ip_configuration {
        name                 = "${azurerm_public_ip.test.name}"
        public_ip_address_id = "${azurerm_public_ip.test.id}"
    }
} 