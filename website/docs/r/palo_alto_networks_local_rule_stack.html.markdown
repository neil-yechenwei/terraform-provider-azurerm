---
subcategory: "PaloAltoNetworks"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_palo_alto_networks_local_rule_stack"
description: |-
  Manages a Palo Alto Networks Local Rule Stack.
---

# azurerm_palo_alto_networks_local_rule_stack

Manages a Palo Alto Networks Local Rule Stack.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_palo_alto_networks_local_rule_stack" "example" {
  name                = "example-lrs"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Palo Alto Networks Local Rule Stack. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Palo Alto Networks Local Rule Stack should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Palo Alto Networks Local Rule Stack should exist. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Palo Alto Networks Local Rule Stack.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) A list of IDs for User Assigned Managed Identity resources to be assigned.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Palo Alto Networks Local Rule Stack.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Palo Alto Networks Local Rule Stack.
* `read` - (Defaults to 5 minutes) Used when retrieving the Palo Alto Networks Local Rule Stack.
* `update` - (Defaults to 30 minutes) Used when updating the Palo Alto Networks Local Rule Stack.
* `delete` - (Defaults to 30 minutes) Used when deleting the Palo Alto Networks Local Rule Stack.

## Import

Palo Alto Networks Local Rule Stacks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_palo_alto_networks_local_rule_stack.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/PaloAltoNetworks.CloudNGFW/localRuleStacks/localRuleStack1
```
