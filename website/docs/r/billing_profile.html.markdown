---
subcategory: "Billing"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_billing_profile"
description: |-
  Manages a Billing Profile.
---

# azurerm_billing_profile

Manages a Billing Profile.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_billing_profile" "example" {
  name                 = "example-bp"
  billing_account_name = "exampleBA"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Billing Profile. Changing this forces a new resource to be created.

* `billing_account_name` - (Required) The name of the Billing Account. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Billing Account.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Billing Profile.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating this Billing Profile.
* `read` - (Defaults to 5 minutes) Used when retrieving this Billing Profile.
* `update` - (Defaults to 30 minutes) Used when updating this Billing Profile.
* `delete` - (Defaults to 30 minutes) Used when deleting this Billing Profile.

## Import

Billing Profiles can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_billing_profile.example /providers/Microsoft.Billing/billingAccounts/ba1/billingProfiles/bp1
```
