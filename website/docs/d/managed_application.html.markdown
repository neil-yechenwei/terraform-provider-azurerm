---
subcategory: "Solutions"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_application"
sidebar_current: "docs-azurerm-datasource-managed-application"
description: |-
Gets information about an existing Managed Application
---

# Data Source: azurerm_managed_application

Uses this data source to access information about an existing Managed Application.

## Managed Application Usage

```hcl
data "azurerm_managed_application" "test" {
  resource_group_name = "acctestRG"
  name                = "acctestnetappsnapshot"
}

output "managed_application_id" {
  value = "${data.azurerm_managed_application.test.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Manage Application.

* `resource_group_name` - (Required) The Name of the Resource Group where the Managed Application exists.

## Attributes Reference

The following attributes are exported:

* `location` - The Azure Region where the Managed Application exists.

* `kind` - The kind of the managed application. Possible values are `MarketPlace` and `ServiceCatalog`.

* `managed_resource_group_id` - The managed resource group Id.

* `application_definition_id` - The fully qualified path of managed application definition Id.

* `parameters` - Name and value pairs that define the managed application parameters.

* `plan` - One or more `plan` blocks as defined below.

* `tags` - A mapping of tags assigned to the resource.

---

The `plan` block exports the following:

* `name` - The name of the plan from the marketplace.

* `publisher` - The publisher of the plan.

* `product` - The product of the plan from the marketplace.

* `version` - The version of the plan from the marketplace.

* `promotion_code` - The promotion code of the plan.
