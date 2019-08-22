resource "azurerm_resource_group" "test" {
  name     = "${var.resource_group_name}"
  location = "${var.location}"
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
	address_prefix      = "10.168.0.0/24"
	virtual_wan_id      = "${azurerm_virtual_wan.test.id}"
}