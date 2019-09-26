resource "azurerm_resource_group" "test" {
  name     = "${var.resource_group_name}"
  location = "${var.location}"
}

resource "azurerm_netapp_account" "test" {
  name                = "acctestaccount"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctestnetapppool"
  account_name        = "${azurerm_netapp_account.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  service_level       = "Premium"
  size                = "4398046511104"

  tags = {
    env = "test"
  }
}
