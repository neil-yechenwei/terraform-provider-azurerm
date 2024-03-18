---
subcategory: "Security Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_security_center_security_connector"
description: |-
  Manages a Security Connector for Azure Security Center.
---

# azurerm_security_center_security_connector

Manages a Security Connector for Azure Security Center.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_security_center_security_connector" "example" {
  name                = "example-securityconnector"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Security Connector. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Security Connector should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Security Connector should exist. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Security Connector.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Security Connector.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating this Security Connector.
* `read` - (Defaults to 5 minutes) Used when retrieving this Security Connector.
* `update` - (Defaults to 30 minutes) Used when updating this Security Connector.
* `delete` - (Defaults to 30 minutes) Used when deleting this Security Connector.

## Import

Security Connectors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_security_center_security_connector.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Security/securityConnectors/securityConnector1
```
