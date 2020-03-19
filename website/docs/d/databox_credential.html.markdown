---
subcategory: "DataBox"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_databox_credential"
description: |-
  Get information about an existing DataBox Credential.
---

# Data Source: azurerm_databox_credential

Use this data source to access information about an existing DataBox Credential.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_databox_credential" "existing" {
  databox_credential_name = "example-databoxjob"
  resource_group_name     = "example-resources"
}

output "id" {
  value = data.azurerm_databox_credential.existing.id
}
```

## Argument Reference

* `databox_credential_name` - Specifies the name of the DataBox.

* `resource_group_name` - Specifies the name of the Resource Group where this DataBox exists.

## Attributes Reference

* `location` - The Azure location where the resource exists.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the DataBox Credential.
