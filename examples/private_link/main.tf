resource "azurerm_resource_group" "example" {
  name     = "neilrgforpep01"
  location = "eastus"
}

resource "azurerm_virtual_network" "example" {
  name                = "neilvnetforpep01"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  address_space       = ["172.22.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                                          = "neilsubnetforpep01"
  virtual_network_name                          = "${azurerm_virtual_network.example.name}"
  resource_group_name                           = "${azurerm_resource_group.example.name}"
  address_prefix                                = "172.22.0.0/24"
  disable_private_link_service_network_policies = false
  disable_private_endpoint_network_policies     = true
}

resource "azurerm_public_ip" "example" {
  name                = "neilpipforpep01"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  allocation_method   = "Dynamic"
  sku                 = "Basic"
}

resource "azurerm_network_security_group" "example" {
  name                = "neilnsgforpep01"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

resource "azurerm_network_security_rule" "example1" {
  name                        = "HTTP"
  priority                    = 300
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "TCP"
  source_port_range           = "*"
  destination_port_range      = "80"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = "${azurerm_resource_group.example.name}"
  network_security_group_name = "${azurerm_network_security_group.example.name}"
}

resource "azurerm_network_security_rule" "example2" {
  name                        = "RDP"
  priority                    = 320
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "3389"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = "${azurerm_resource_group.example.name}"
  network_security_group_name = "${azurerm_network_security_group.example.name}"
}

resource "azurerm_network_interface" "example" {
  name                = "neilnicforpep01"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  ip_configuration {
    name                          = "testipconfig1"
    subnet_id                     = "${azurerm_subnet.example.id}"
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = "${azurerm_public_ip.example.id}"
  }

  network_security_group_id = "${azurerm_network_security_group.example.id}"
}

resource "azurerm_storage_account" "example" {
  name                = "neilstorageacctforpep01"
  resource_group_name = "${azurerm_resource_group.example.name}"

  location                 = "${azurerm_resource_group.example.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "Storage"
}

resource "azurerm_virtual_machine" "example" {
  name                  = "neilvmforpep01"
  location              = "${azurerm_resource_group.example.location}"
  resource_group_name   = "${azurerm_resource_group.example.name}"
  network_interface_ids = ["${azurerm_network_interface.example.id}"]
  vm_size               = "Standard_DS1_v2"

  storage_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osdisk"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Premium_LRS"
  }

  os_profile {
    computer_name  = "neilvmforpep01"
    admin_username = "v-cheye"
    admin_password = "Password1234!"
  }

  os_profile_windows_config {
    enable_automatic_upgrades = true
    provision_vm_agent = true
  }

  boot_diagnostics {
    enabled     = true
    storage_uri = "${azurerm_storage_account.example.primary_blob_endpoint}"
  }
}

resource "azurerm_sql_server" "example" {
  name                         = "neilsqlserverforpep01"
  resource_group_name          = "${azurerm_resource_group.example.name}"
  location                     = "${azurerm_resource_group.example.location}"
  version                      = "12.0"
  administrator_login          = "v-cheye"
  administrator_login_password = "Password2345!"
}

resource "azurerm_sql_database" "example" {
  name                             = "neilsqldbforpep01"
  resource_group_name              = "${azurerm_resource_group.example.name}"
  server_name                      = "${azurerm_sql_server.example.name}"
  location                         = "${azurerm_resource_group.example.location}"
  edition                          = "GeneralPurpose"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "34359738368"
  read_scale                       = false

  threat_detection_policy {
    state                      = "Enabled"
    email_account_admins       = "Enabled"
  }
}

resource "azurerm_private_link_endpoint" "example" {
  name                = "neiltestpep01"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  subnet_id           = "${azurerm_subnet.example.id}"

  private_link_service_connection {
    name = "testplsconnection"
    private_link_service_id = "${azurerm_sql_server.example.id}"
    group_ids               = ["sqlServer"]
    request_message         = "Please approve neil connection"
  }
}

resource "azurerm_private_dns_zone" "example" {
  name                = "privatelink.database.windows.net"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

resource "azurerm_private_dns_zone_virtual_network_link" "example" {
  name                	= "neilpdzvnetlink01"
  private_dns_zone_name = "${azurerm_private_dns_zone.example.name}"
  virtual_network_id 	= "${azurerm_virtual_network.example.id}"
  resource_group_name 	= "${azurerm_resource_group.example.name}"
  registration_enabled  = false
}

resource "azurerm_private_dns_a_record" "example" {
  name                = "neilsqlserverforpep01"
  resource_group_name = "${azurerm_resource_group.example.name}"
  zone_name           = "${azurerm_private_dns_zone.example.name}"
  ttl                 = 300
  records             = ["172.22.0.5"]
}