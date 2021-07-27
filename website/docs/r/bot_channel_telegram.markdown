---
subcategory: "Bot"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_bot_channel_telegram"
description: |-
  Manages a Telegram integration for a Bot Channel
---

# azurerm_bot_channel_telegram

Manages a Telegram integration for a Bot Channel

~> **Note** A bot can only have a single Telegram Channel associated with it.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_bot_channels_registration" "example" {
  name                = "example-bcr"
  location            = "global"
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "F0"
  microsoft_app_id    = data.azurerm_client_config.current.client_id
}

resource "azurerm_bot_channel_telegram" "example" {
  bot_name                      = azurerm_bot_channels_registration.example.name
  location                      = azurerm_bot_channels_registration.example.location
  resource_group_name           = azurerm_resource_group.example.name
  telegram_channel_access_token = "TestAccessToken"
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group where the Telegram Channel should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `bot_name` - (Required) The name of the Bot Resource this channel will be associated with. Changing this forces a new resource to be created.

* `telegram_channel_access_token` - (Required) The access token for the Telegram Channel.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Telegram Integration for a Bot Channel.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Telegram Integration for a Bot Channel.
* `update` - (Defaults to 30 minutes) Used when updating the Telegram Integration for a Bot Channel.
* `read` - (Defaults to 5 minutes) Used when retrieving the Telegram Integration for a Bot Channel.
* `delete` - (Defaults to 30 minutes) Used when deleting the Telegram Integration for a Bot Channel.

## Import

The Telegram Integration for a Bot Channel can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_bot_channel_telegram.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.BotService/botServices/botService1/channels/TelegramChannel
```
