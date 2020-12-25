---
subcategory: "Template"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_template_spec_version"
description: |-
  Manages a Template Spec Version.
---

# azurerm_template_spec_version

Manages a Template Spec Version.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_template_spec" "example" {
  name                = "example-templatespec"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_template_spec_version" "example" {
  name                = "example-templatespecversion"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  template_spec_name  = azurerm_template_spec.example.name

  template_content = <<DEPLOY
  {
      "$schema": "http://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
      "contentVersion": "1.0.0.0",
      "parameters": {},
      "resources": []
  }
  DEPLOY
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Template Spec Version. Changing this forces a new Template Spec Version to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Template Spec Version should exist. Changing this forces a new Template Spec Version to be created.

* `location` - (Required) The Azure Region where the Template Spec Version should exist. Changing this forces a new Template Spec Version to be created.

* `template_content` - (Optional) The content of the main ARM template. Changing this forces a new Template Spec Version to be created.

* `template_spec_name` - (Required) The name of the Template Spec. Changing this forces a new Template Spec Version to be created.

---

* `artifact` - (Optional)  An `artifact` block as defined below. Changing this forces a new Template Spec Version to be created.

* `description` - (Optional) The description of the Template Spec Version. Changing this forces a new Template Spec Version to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Template Spec Version.

---

An `artifact` block exports the following:

* `path` - (Required) The relative path of the template artifact. Changing this forces a new Template Spec Version to be created.

* `template_content` - (Required) The content of the linked template which is associated with main template. Changing this forces a new Template Spec Version to be created.

* `kind` - (Optional) The kind of the template artifact. Possible value is `template`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Template Spec Version.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Template Spec Version.
* `read` - (Defaults to 5 minutes) Used when retrieving the Template Spec Version.
* `update` - (Defaults to 30 minutes) Used when updating the Template Spec Version.
* `delete` - (Defaults to 30 minutes) Used when deleting the Template Spec Version.

## Import

Template Spec Versions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_template_spec_version.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Resources/templateSpecs/spec1/versions/v1
```
