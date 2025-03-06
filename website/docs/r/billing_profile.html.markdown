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
resource "azurerm_billing_profile" "example" {
  name                 = "example-bp"
  billing_account_name = "exampleBA"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Billing Profile. Changing this forces a new resource to be created.

* `billing_account_name` - (Required) The name of the Billing Account. Changing this forces a new resource to be created.

* `bill_to` - (Optional) A `bill_to` block as defined below.

* `current_payment_term` - (Optional) A `current_payment_term` block as defined below.

* `display_name` - (Optional) The display name of the Billing Profile.

* `enabled_azure_plan` - (Optional) A `enabled_azure_plan` block as defined below.

* `indirect_relationship_info` - (Optional) A `indirect_relationship_info` block as defined below.

* `invoice_email_opt_in_enabled` - (Optional) Should the invoices for the Billing Profile send through email?

* `invoice_recipients` - (Optional) A list of email addresses to receive invoices by email for the Billing Profile.

* `po_number` - (Optional) The default purchase order number that will appear on the invoices generated for the Billing Profile.

* `ship_to` - (Optional) A `ship_to` block as defined below.

* `sold_to` - (Optional) A `sold_to` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Billing Profile.

---

A `bill_to` block supports the following:

* `address_line_1` - (Required) The address line 1.

* `country` - (Required) The country code uses ISO 3166-1 Alpha-2 format.

* `address_line_2` - (Optional) The address line 2.

* `address_line_3` - (Optional) The address line 3.

* `city` - (Optional) The address city.

* `company_name` - (Optional) The company name.

* `district` - (Optional) The address district.

* `email` - (Optional) The email address.

* `first_name` - (Optional) The first name.

* `last_name` - (Optional) The last name.

* `middle_name` - (Optional) The middle name.

* `phone_number` - (Optional) The phone number.

* `postal_code` - (Optional) The postal code.

* `region` - (Optional) The address region.

* `valid_address_enabled` - (Optional) Should the address be incomplete?

---

A `current_payment_term` block supports the following:

* `end_date` - (Optional) The end date on when the defined `Payment Term` will end and is always in UTC.

* `start_date` - (Optional) The start date on when the defined `Payment Term` will be effective from and is always in UTC.

* `term` - (Optional) The term that represents duration in netXX format. Always in days.

---

A `enabled_azure_plan` block supports the following:

* `product_id` - (Optional) The ID that uniquely identifies a product.

* `sku_description` - (Optional) The SKU description.

* `sku_id` - (Optional) The ID that uniquely identifies a SKU.

---

A `indirect_relationship_info` block supports the following:

* `billing_account_name` - (Optional) The Billing Account name of the partner or the customer for an indirect motion.

* `billing_profile_name` - (Optional) The Billing Profile name of the partner or the customer for an indirect motion.

* `display_name` - (Optional) The display name of the partner or customer for an indirect motion.

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
