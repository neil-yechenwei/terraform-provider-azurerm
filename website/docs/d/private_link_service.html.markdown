---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_link_service"
sidebar_current: "docs-azurerm-datasource-private-link-service"
description: |-
  Gets information about an existing Private Link Service.
---

# Data Source: azurerm_private_link_service

Use this data source to access information about an existing Private Link Service.

## Example Usage

```hcl
data "azurerm_private_link_service" "test" {
  name                = "myprivatelinkservice"
  resource_group_name = "myprivatelinkservicerg"
}

output "private_link_service_id" {
  value = "${data.azurerm_private_link_service.test.id}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Private Link Service.

* `resource_group_name` - (Required) Specifies the name of the resource group the Private Link Service is located in.

## Attributes Reference

* `id` - The ID of the Private Link Service.
