---
subcategory: "Solutions"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_application"
sidebar_current: "docs-azurerm-resource-managed-application"
description: |-
  Manages a Managed Application.
---

# azurerm_managed_application

Manages a Managed Application.

## Managed Application Usage

```hcl
resource "azurerm_resource_group" "app_def_group" {
  name     = "example-app-def-group-resources"
  location = "West Europe"
}

resource "azurerm_resource_group" "app_group" {
  name     = "example-app-group-resources"
  location = "West Europe"
}

data "azurerm_client_config" "example" {}

resource "azurerm_managed_application_definition" "example" {
  name                 = "example-managedappdef"
  location             = "${azurerm_resource_group.app_def_group.location}"
  resource_group_name  = "${azurerm_resource_group.app_def_group.name}"
  lock_level           = "ReadOnly"
  package_file_uri     = "https://github.com/Azure/azure-managedapp-samples/raw/master/Managed Application Sample Packages/201-managed-storage-account/managedstorage.zip"
  display_name         = "TestManagedAppDefinition"
  description          = "Test Managed App Definition"

  authorization {
    service_principal_id = "${data.azurerm_client_config.example.object_id}"
    role_definition_id   = "b24988ac-6180-42a0-ab88-20f7382dd24c"
  }
}

resource "azurerm_managed_application" "example" {
  name                      = "example-managedapp"
  location                  = "${azurerm_resource_group.app_group.location}"
  resource_group_name       = "${azurerm_resource_group.app_group.name}"
  kind                      = "ServiceCatalog"
  managed_resource_group_id = "/subscriptions/${data.azurerm_client_config.example.subscription_id}/resourceGroups/infrastructureGroup"
  application_definition_id = "${azurerm_managed_application_definition.example.id}"
  parameters                = <<PARAMETERS
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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Snapshot. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the NetApp Snapshot should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `kind` - (Required) The kind of the managed application. Possible values are `MarketPlace` and `ServiceCatalog`.

* `managed_resource_group_id` - (Required) The managed resource group Id.

* `application_definition_id` - (Optional) The fully qualified path of managed application definition Id.

* `parameters` - (Optional) The name and value pairs that define the managed application parameters.

* `plan` - (Optional) One or more `plan` blocks as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `plan` block exports the following:

* `name` - (Required) Specifies the name of the plan from the marketplace.

* `publisher` - (Required) Specifies the publisher of the plan.

* `product` - (Required) Specifies the product of the plan from the marketplace.

* `version` - (Required) Specifies the version of the plan from the marketplace.

* `promotion_code` - (Optional) Specifies the promotion code of the plan. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Managed Application.

## Import

Managed Application can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_managed_application.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Solutions/applications/app1
```
