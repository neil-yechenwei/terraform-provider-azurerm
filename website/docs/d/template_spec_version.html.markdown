---
subcategory: "Template"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_template_spec_version"
description: |-
  Gets information about an existing Template Spec Version
---

# Data Source: azurerm_template_spec_version

Use this data source to access information about an existing Template Spec Version

## Example Usage

```hcl
data "azurerm_template_spec_version" "example" {
  name                = "example-templatespecversion"
  resource_group_name = "example-resources"
  template_spec_name  = "example-templatespec"
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the Template Spec Version resource.

* `resource_group_name` - The name of the Resource Group where the Template Spec Version exists.

* `template_spec_name` - The name of the Template Spec resource.

## Attributes Reference

The following attributes are exported:

* `id` - The Template Spec ID.

* `artifact` - An `artifact` block as defined below.

---

An `artifact` block exports the following:

* `path` - The safe relative path of the template artifact in File System.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Template Spec Version.
