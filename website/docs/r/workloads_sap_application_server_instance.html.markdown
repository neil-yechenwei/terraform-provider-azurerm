---
subcategory: "Workloads"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_workloads_sap_application_server_instance"
description: |-
  Manages a SAP Application Server Instance.
---

# azurerm_workloads_sap_application_server_instance

Manages a SAP Application Server Instance.

## Example Usage

```hcl
data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example-uai"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_role_assignment" "example" {
  scope                = data.azurerm_subscription.current.id
  role_definition_name = "Azure Center for SAP solutions service role"
  principal_id         = azurerm_user_assigned_identity.example.principal_id
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_workloads_sap_virtual_instance" "example" {
  name                        = "X01"
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  environment                 = "NonProd"
  sap_product                 = "S4HANA"
  managed_resource_group_name = "managedTestRG"

  deployment_with_os_configuration {
    app_location = azurerm_resource_group.example.location

    os_sap_configuration {
        sap_fqdn = "sap.bpaas.com"
    }

    single_server_configuration {
        app_resource_group_name = azurerm_resource_group.example.name
        subnet_id               = azurerm_subnet.example.id
        database_type           = "HANA"
        is_secondary_ip_enabled = true

        virtual_machine_configuration {
            vm_size = "Standard_E32ds_v4"

            image_reference {
                offer     = "RHEL-SAP-HA"
                publisher = "RedHat"
                sku       = "82sapha-gen2"
                version   = "latest"
            }

            os_profile {
                admin_username = "testAdmin"

                ssh_key_pair {
                    private_key = "-----BEGIN RSA PRIVATE KEY-----\n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\n-----END RSA PRIVATE KEY-----\n"
                    public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC8k1K0mjpBtyBRcr+74p6zShq09hMXmQZJJem41RlqV132kbGUy/EHcQLoOmp2qdv6Kz/jgz49YWM2WxxhB8G+Uh8G58feYFaGNkT+WTrfE/XLHxFnXMghZgCEiV8vfNUROuHux7UfsxiDt6Nw2215fkJ8lQjQEduKPfS7OmoQQDUvGeg8rxUbwouYdR9z5MmcNKSiBC0Z3EEBWVkY0jsbASarE9iRLDT5JrXeeZWCSe+/lg3Ofe7b28f+nXo9LH/voRVc+eKsx6T6reW9wVtm6ygEBXZbqvuEeqISgLvSFWXD7SsO2CbyJpmaOkzjIMbzlaHV126+PYdAY2vf0gVhCwRt5wU/ueFloMmTb3HbLkyVczluBUsSfcJ3H5DYPECiCy5qKaPsGHDzFGByc/0Yy+fNn2L/WrNgNuC3LaeszTZ7YqNVo/6tXTGSqOMrv8AD1rP71sX2HDSdlLf/HyLRgNjfJn/KR3ZE9QZ2Hl9TRcZ9DN+V/wvc/lci/7nsSx0= generated-by-azure"
                }
            }
        }

        disk_volume_configuration {
            volume_name = "hana/data"
            count       = 3
            size_gb     = 128
            sku_name    = "Premium_LRS"
        }

        disk_volume_configuration {
            volume_name = "hana/log"
            count       = 3
            size_gb     = 128
            sku_name    = "Premium_LRS"
        }

        disk_volume_configuration {
            volume_name = "hana/shared"
            count       = 1
            size_gb     = 256
            sku_name    = "Premium_LRS"
        }

        disk_volume_configuration {
            volume_name = "usr/sap"
            count       = 1
            size_gb     = 128
            sku_name    = "Premium_LRS"
        }

        disk_volume_configuration {
            volume_name = "backup"
            count       = 2
            size_gb     = 256
            sku_name    = "StandardSSD_LRS"
        }

        disk_volume_configuration {
            volume_name = "os"
            count       = 1
            size_gb     = 64
            sku_name    = "StandardSSD_LRS"
        }

        virtual_machine_full_resource_names {
            host_name               = "apphostName0"
            os_disk_name            = "app0osdisk"
            vm_name                 = "appvm0"
            network_interface_names = ["appnic0"]

            data_disk_names = {
                default = "app0disk0"
            }
        }
    }
  }

  identity {
    type = "UserAssigned"

    identity_ids = [
        azurerm_user_assigned_identity.example.id,
    ]
  }

  tags = {
    Env = "Test"
  }

  depends_on = [
    azurerm_role_assignment.example
  ]
}

resource "azurerm_workloads_sap_application_server_instance" "example" {
  name                    = "example-sapapp"
  resource_group_name     = azurerm_resource_group.example.name
  location                = azurerm_resource_group.example.location
  sap_virtual_instance_id = azurerm_workloads_sap_virtual_instance.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this SAP Application Server Instance. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the SAP Application Server Instance should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the SAP Application Server Instance should exist. Changing this forces a new resource to be created.

* `sap_virtual_instance_id` - (Required) The resource ID of the SAP Virtual Instance. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the SAP Application Server Instance.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the SAP Application Server Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the SAP Application Server Instance.
* `read` - (Defaults to 5 minutes) Used when retrieving the SAP Application Server Instance.
* `update` - (Defaults to 60 minutes) Used when updating the SAP Application Server Instance.
* `delete` - (Defaults to 60 minutes) Used when deleting the SAP Application Server Instance.

## Import

SAP Application Server Instances can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_workloads_sap_application_server_instance.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Workloads/sapVirtualInstances/sapVirtualInstance1/applicationInstances/applicationInstance1
```
