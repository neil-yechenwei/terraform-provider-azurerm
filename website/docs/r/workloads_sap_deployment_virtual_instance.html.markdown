---
subcategory: "Workloads"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_workloads_sap_deployment_virtual_instance"
description: |-
  Manages a SAP Virtual Instance with new SAP System.
---

# azurerm_workloads_sap_deployment_virtual_instance

Manages a SAP Virtual Instance with new SAP System.

-> **Note:** Before using this resource, it's required to submit the request of registering the Resource Provider with Azure CLI `az provider register --namespace "Microsoft.Workloads"`. The Resource Provider can take a while to register, you can check the status by running `az provider show --namespace "Microsoft.Workloads" --query "registrationState"`. Once this outputs "Registered" the Resource Provider is available for use.

## Example Usage

```hcl
data "azurerm_subscription" "current" {}

resource "tls_private_key" "example" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

data "tls_public_key" "example" {
  private_key_pem = tls_private_key.example.private_key_pem
}

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

resource "azurerm_resource_group" "app" {
  name     = "example-sapapp"
  location = "West Europe"

  depends_on = [
    azurerm_subnet.example
  ]
}

resource "azurerm_workloads_sap_deployment_virtual_instance" "example" {
  name                        = "X05"
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  environment                 = "NonProd"
  sap_product                 = "S4HANA"
  managed_resource_group_name = "managedTestRG"

  deployment_with_os_configuration {
    app_location = azurerm_resource_group.app.location

    os_sap_configuration {
      sap_fqdn = "sap.bpaas.com"
    }

    single_server_configuration {
      app_resource_group_name = azurerm_resource_group.app.name
      subnet_id               = azurerm_subnet.example.id
      database_type           = "HANA"
      secondary_ip_enabled    = true

      virtual_machine_configuration {
        virtual_machine_size = "Standard_E32ds_v4"

        image {
          offer     = "RHEL-SAP-HA"
          publisher = "RedHat"
          sku       = "82sapha-gen2"
          version   = "latest"
        }

        os_profile {
          admin_username  = "testAdmin"
          ssh_private_key = tls_private_key.example.private_key_pem
          ssh_public_key  = data.tls_public_key.example.public_key_openssh
        }
      }

      disk_volume_configuration {
        volume_name     = "hana/data"
        number_of_disks = 3
        size_in_gb      = 128
        sku_name        = "Premium_LRS"
      }

      disk_volume_configuration {
        volume_name     = "hana/log"
        number_of_disks = 3
        size_in_gb      = 128
        sku_name        = "Premium_LRS"
      }

      disk_volume_configuration {
        volume_name     = "hana/shared"
        number_of_disks = 1
        size_in_gb      = 256
        sku_name        = "Premium_LRS"
      }

      disk_volume_configuration {
        volume_name     = "usr/sap"
        number_of_disks = 1
        size_in_gb      = 128
        sku_name        = "Premium_LRS"
      }

      disk_volume_configuration {
        volume_name     = "backup"
        number_of_disks = 2
        size_in_gb      = 256
        sku_name        = "StandardSSD_LRS"
      }

      disk_volume_configuration {
        volume_name     = "os"
        number_of_disks = 1
        size_in_gb      = 64
        sku_name        = "StandardSSD_LRS"
      }

      virtual_machine_full_resource_names {
        host_name               = "apphostName0"
        os_disk_name            = "app0osdisk"
        virtual_machine_name    = "appvm0"
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

  depends_on = [
    azurerm_role_assignment.example
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this SAP Virtual Instance. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the SAP Virtual Instance should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the SAP Virtual Instance should exist. Changing this forces a new resource to be created.

* `deployment_with_os_configuration` - (Required) A `deployment_with_os_configuration` block as defined below. Changing this forces a new resource to be created.

* `environment` - (Required) The environment type for the SAP Virtual Instance. Possible values are `NonProd` and `Prod`. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `sap_product` - (Required) The SAP Product type for the SAP Virtual Instance. Possible values are `ECC`, `Other` and `S4HANA`. Changing this forces a new resource to be created.

* `managed_resource_group_name` - (Optional) The name of the managed Resource Group for the SAP Virtual Instance. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the SAP Virtual Instance.

---

A `deployment_with_os_configuration` block supports the following:

* `app_location` - (Required) The Geo-Location where the SAP system is to be created. Changing this forces a new resource to be created.

* `os_sap_configuration` - (Required) An `os_sap_configuration` block as defined below. Changing this forces a new resource to be created.

* `single_server_configuration` - (Optional) A `single_server_configuration` block as defined below. Changing this forces a new resource to be created.

* `three_tier_configuration` - (Optional) A `three_tier_configuration` block as defined below. Changing this forces a new resource to be created.

---

An `os_sap_configuration` block supports the following:

* `sap_fqdn` - (Required) The FQDN of the SAP system. Changing this forces a new resource to be created.

* `deployer_virtual_machine_packages` - (Optional) A `deployer_virtual_machine_packages` block as defined below. Changing this forces a new resource to be created.

---

A `deployer_virtual_machine_packages` block supports the following:

* `storage_account_id` - (Required) A `deployer_virtual_machine_packages` block as defined below. Changing this forces a new resource to be created.

* `url` - (Required) The URL of the deployer VM packages file. Changing this forces a new resource to be created.

---

A `single_server_configuration` block supports the following:

* `app_resource_group_name` - (Required) The name of the application Resource Group where SAP system resources will be deployed. Changing this forces a new resource to be created.

~> **Note:** While creating SAP Virtual Instance, service would provision the extra SAP system/component in `app_resource_group_name` which aren't defined in tf config. At this time, when `app_resource_group_name` is different with the Resource Group where SAP Virtual Instance exists, we can set `prevent_deletion_if_contains_resources` to `false` to delete all resources defined in tf config and the resources created in `app_resource_group_name` with `tf destroy`. But when `app_resource_group_name` is same with the Resource Group where SAP Virtual Instance exists, some resources like the subnet defined in tf config cannot be deleted with `tf destroy` since the resources defined in tf config are being referenced by the SAP system/component created in `app_resource_group_name`. So it has to manually delete the resources in `app_resource_group_name` first after the SAP Virtual Instance is deleted, and then the resources in tf config can be deleted successfully for this situation.

* `subnet_id` - (Required) The resource ID of the Subnet for the SAP Virtual Instance. Changing this forces a new resource to be created.

* `virtual_machine_configuration` - (Required) A `virtual_machine_configuration` block as defined below. Changing this forces a new resource to be created.

* `database_type` - (Optional) The supported SAP database type. Possible values are `DB2` and `HANA`. Changing this forces a new resource to be created.

* `disk_volume_configuration` - (Optional) One or more `disk_volume_configuration` blocks as defined below. Changing this forces a new resource to be created.

* `secondary_ip_enabled` - (Optional) Is a secondary IP Address that should be added to the Network Interface on all VMs of the SAP system being deployed enabled? Defaults to `false`. Changing this forces a new resource to be created.

* `virtual_machine_full_resource_names` - (Optional) A `virtual_machine_full_resource_names` block as defined below. Changing this forces a new resource to be created.

---

A `disk_volume_configuration` block supports the following:

* `volume_name` - (Required) The name of the DB volume of the disk configuration. Possible values are `backup`, `hana/data`, `hana/log`, `hana/shared`, `os` and `usr/sap`. Changing this forces a new resource to be created.

* `number_of_disks` - (Required) The total number of disks required for the concerned volume. Changing this forces a new resource to be created.

* `size_in_gb` - (Required) The size of the Disk in GB. Changing this forces a new resource to be created.

* `sku_name` - (Required) The name of the Disk SKU. Possible values are `Premium_LRS`, `PremiumV2_LRS`, `Premium_ZRS`, `Standard_LRS`, `StandardSSD_LRS`, `StandardSSD_ZRS` and `UltraSSD_LRS`. Changing this forces a new resource to be created.

---

A `virtual_machine_configuration` block supports the following:

* `image` - (Required) An `image` block as defined below. Changing this forces a new resource to be created.

* `os_profile` - (Required) An `os_profile` block as defined below. Changing this forces a new resource to be created.

* `virtual_machine_size` - (Required) The size of the Virtual Machine. Changing this forces a new resource to be created.

---

An `image` block supports the following:

* `offer` - (Required) The offer of the platform image or marketplace image used to create the Virtual Machine. Changing this forces a new resource to be created.

* `publisher` - (Required) The publisher of the Image. Possible values are `RedHat` and `SUSE`. Changing this forces a new resource to be created.

* `sku` - (Required) The SKU of the Image. Changing this forces a new resource to be created.

* `version` - (Required) The version of the platform image or marketplace image used to create the Virtual Machine. Changing this forces a new resource to be created.

---

An `os_profile` block supports the following:

* `admin_username` - (Required) The name of the administrator account. Changing this forces a new resource to be created.

* `ssh_private_key` - (Required) The SSH public key that is used to authenticate with the VM. Changing this forces a new resource to be created.

* `ssh_public_key` - (Required) The SSH private key that is used to authenticate with the VM. Changing this forces a new resource to be created.

---

A `virtual_machine_full_resource_names` block supports the following:

* `data_disk_names` - (Optional) A mapping of Data Disk names to pass to the backend host. The keys are Volume names and the values are a comma separated string of full names for Data Disks belonging to the specific Volume. This is converted to a list before being passed to the API. Changing this forces a new resource to be created.

* `host_name` - (Optional) The full name of the host of the Virtual Machine. Changing this forces a new resource to be created.

* `network_interface_names` - (Optional) A list of full names for the Network Interface of the Virtual Machine. Changing this forces a new resource to be created.

* `os_disk_name` - (Optional) The full name of the OS Disk attached to the VM. Changing this forces a new resource to be created.

* `virtual_machine_name` - (Optional) The full name of the Virtual Machine in a single server SAP system. Changing this forces a new resource to be created.

---

A `three_tier_configuration` block supports the following:

* `app_resource_group_name` - (Required) The name of the application Resource Group where SAP system resources will be deployed. Changing this forces a new resource to be created.

* `application_server_configuration` - (Required) An `application_server_configuration` block as defined below. Changing this forces a new resource to be created.

* `central_server_configuration` - (Required) A `central_server_configuration` block as defined below. Changing this forces a new resource to be created.

* `database_server_configuration` - (Required) A `database_server_configuration` block as defined below. Changing this forces a new resource to be created.

* `full_resource_names` - (Optional) A `full_resource_names` block as defined below. Changing this forces a new resource to be created.

* `high_availability_type` - (Optional) The high availability type for the three tier configuration. Possible values are `AvailabilitySet` and `AvailabilityZone`. Changing this forces a new resource to be created.

* `secondary_ip_enabled` - (Optional) Is a secondary IP Address that should be added to the Network Interface on all VMs of the SAP system being deployed enabled? Defaults to `false`. Changing this forces a new resource to be created.

* `transport_create_and_mount` - (Optional) A `transport_create_and_mount` block as defined below. Changing this forces a new resource to be created.

* `transport_mount` - (Optional) A `transport_mount` block as defined below. Changing this forces a new resource to be created.

~> **Note:** The `Skip` configuration type would be enabled when the `transport_create_and_mount` and the `transport_mount` aren't set.

---

A `transport_create_and_mount` block supports the following:

* `resource_group_id` - (Optional) The ID of the Resource Group of the transport File Share. Changing this forces a new resource to be created.

* `storage_account_name` - (Optional) The name of the Storage Account of the File Share. Changing this forces a new resource to be created.

---

A `transport_mount` block supports the following:

* `file_share_id` - (Required) The resource ID of the Storage File Share. Changing this forces a new resource to be created.

* `private_endpoint_id` - (Required) The resource ID of the Private Endpoint. Changing this forces a new resource to be created.

---

An `application_server_configuration` block supports the following:

* `instance_count` - (Required) The number of instances for the Application Server. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The resource ID of the Subnet for the Application Server. Changing this forces a new resource to be created.

* `virtual_machine_configuration` - (Required) A `virtual_machine_configuration` block as defined below. Changing this forces a new resource to be created.

---

A `central_server_configuration` block supports the following:

* `instance_count` - (Required) The number of instances for the Central Server. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The resource ID of the Subnet for the Central Server. Changing this forces a new resource to be created.

* `virtual_machine_configuration` - (Required) A `virtual_machine_configuration` block as defined below. Changing this forces a new resource to be created.

---

A `database_server_configuration` block supports the following:

* `instance_count` - (Required) The number of instances for the Database Server. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The resource ID of the Subnet for the Database Server. Changing this forces a new resource to be created.

* `virtual_machine_configuration` - (Required) A `virtual_machine_configuration` block as defined below. Changing this forces a new resource to be created.

* `database_type` - (Optional) The database type for the Database Server. Possible values are `DB2` and `HANA`. Changing this forces a new resource to be created.

* `disk_volume_configuration` - (Optional) One or more `disk_volume_configuration` blocks as defined below. Changing this forces a new resource to be created.

---

A `full_resource_names` block supports the following:

* `application_server` - (Optional) An `application_server` block as defined below. Changing this forces a new resource to be created.

* `central_server` - (Optional) A `central_server` block as defined below. Changing this forces a new resource to be created.

* `database_server` - (Optional) A `database_server` block as defined below. Changing this forces a new resource to be created.

* `shared_storage` - (Optional) A `shared_storage` block as defined below. Changing this forces a new resource to be created.

---

An `application_server` block supports the following:

* `availability_set_name` - (Optional) The full name for the availability set. Changing this forces a new resource to be created.

* `virtual_machine` - (Optional) One or more `virtual_machine` blocks as defined below. Changing this forces a new resource to be created.

---

A `virtual_machine` block supports the following:

* `data_disk_names` - (Optional) A mapping of Data Disk names to pass to the backend host. The keys are Volume names and the values are a comma separated string of full names for Data Disks belonging to the specific Volume. This is converted to a list before being passed to the API. Changing this forces a new resource to be created.

* `host_name` - (Optional) The full name of the host of the Virtual Machine. Changing this forces a new resource to be created.

* `network_interface_names` - (Optional) A list of full names for the Network Interface of the Virtual Machine. Changing this forces a new resource to be created.

* `os_disk_name` - (Optional) The full name of the OS Disk attached to the VM. Changing this forces a new resource to be created.

* `virtual_machine_name` - (Optional) The full name of the Virtual Machine in a single server SAP system. Changing this forces a new resource to be created.

---

A `central_server` block supports the following:

* `availability_set_name` - (Optional) The full name for the availability set. Changing this forces a new resource to be created.

* `load_balancer` - (Optional) A `load_balancer` block as defined below. Changing this forces a new resource to be created.

* `virtual_machine` - (Optional) One or more `virtual_machine` blocks as defined below. Changing this forces a new resource to be created.

---

A `load_balancer` block supports the following:

* `name` - (Optional) The full resource name of the Load Balancer. Changing this forces a new resource to be created.

* `backend_pool_names` - (Optional) A list of Backend Pool names for the Load Balancer. Changing this forces a new resource to be created.

* `frontend_ip_configuration_names` - (Optional) A list of Frontend IP Configuration names. Changing this forces a new resource to be created.

* `health_probe_names` - (Optional) A list of Health Probe names. Changing this forces a new resource to be created.

---

A `database_server` block supports the following:

* `availability_set_name` - (Optional) The full name for the availability set. Changing this forces a new resource to be created.

* `load_balancer` - (Optional) A `load_balancer` block as defined below. Changing this forces a new resource to be created.

* `virtual_machine` - (Optional) One or more `virtual_machine` blocks as defined below. Changing this forces a new resource to be created.

---

A `shared_storage` block supports the following:

* `account_name` - (Optional) The full name of the Shared Storage Account. Changing this forces a new resource to be created.

* `private_endpoint_name` - (Optional) The full name of Private Endpoint for the Shared Storage Account. Changing this forces a new resource to be created.

---

An `identity` block supports the following:

* `type` - (Required) The type of Managed Service Identity that should be configured on this SAP Virtual Instance. Only possible value is `UserAssigned`.

* `identity_ids` - (Required) A list of User Assigned Managed Identity IDs to be assigned to this SAP Virtual Instance.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the SAP Virtual Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the SAP Virtual Instance with new SAP System.
* `read` - (Defaults to 5 minutes) Used when retrieving the SAP Virtual Instance with new SAP System.
* `update` - (Defaults to 60 minutes) Used when updating the SAP Virtual Instance with new SAP System.
* `delete` - (Defaults to 60 minutes) Used when deleting the SAP Virtual Instance with new SAP System.

## Import

SAP Virtual Instances with new SAP Systems can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_workloads_sap_deployment_virtual_instance.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Workloads/sapVirtualInstances/vis1
```
