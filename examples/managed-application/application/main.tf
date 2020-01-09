resource "azurerm_resource_group" "app_def_group" {
  name     = "${var.prefix}-app-def-group-resources"
  location = "${var.location}"
}

resource "azurerm_resource_group" "app_group" {
  name     = "${var.prefix}-app-group-resources"
  location = "${var.location}"
}

data "azurerm_client_config" "current" {}

resource "azurerm_managed_application_definition" "example" {
  name                 = "${var.prefix}managedappdefinition"
  location             = "${azurerm_resource_group.app_def_group.location}"
  resource_group_name  = "${azurerm_resource_group.app_def_group.name}"
  lock_level           = "ReadOnly"
  package_file_uri     = "https://github.com/Azure/azure-managedapp-samples/raw/master/Managed Application Sample Packages/201-managed-storage-account/managedstorage.zip"
  display_name         = "TestManagedAppDefinition"
  description          = "Test Managed App Definition"

  authorization {
    service_principal_id = "${data.azurerm_client_config.current.object_id}"
    role_definition_id   = "b24988ac-6180-42a0-ab88-20f7382dd24c"
  }
}

resource "azurerm_managed_application" "example" {
  name                      = "${var.prefix}managedapp"
  location                  = "${azurerm_resource_group.app_group.location}"
  resource_group_name       = "${azurerm_resource_group.app_group.name}"
  kind                      = "ServiceCatalog"
  managed_resource_group_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/infrastructureGroup"
  application_definition_id = "${azurerm_managed_application_definition.example.id}"
  parameters = <<PARAMETERS
    {
        "storageAccountNamePrefix": {
            "value": "demostorags"
        },
        "storageAccountType": {
            "value": "Standard_LRS"
        }
    }
    PARAMETERS
}
