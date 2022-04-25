package logic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LogicAppActionCustomResource struct{}

func TestAccLogicAppActionCustom_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_action_custom", "test")
	r := LogicAppActionCustomResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppActionCustom_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_action_custom", "test")
	r := LogicAppActionCustomResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_logic_app_action_custom"),
		},
	})
}

func TestAccLogicAppActionCustom_multipleActionsWithIntegrationAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_action_custom", "test")
	r := LogicAppActionCustomResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multipleActionsWithIntegrationAccount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (LogicAppActionCustomResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	return actionExists(ctx, clients, state)
}

func (r LogicAppActionCustomResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_custom" "test" {
  name         = "action%d"
  logic_app_id = azurerm_logic_app_workflow.test.id

  body = <<BODY
{
    "description": "A variable to configure the auto expiration age in days. Configured in negative number. Default is -30 (30 days old).",
    "inputs": {
        "variables": [
            {
                "name": "ExpirationAgeInDays",
                "type": "Integer",
                "value": -30
            }
        ]
    },
    "runAfter": {},
    "type": "InitializeVariable"
}
BODY

}
`, r.template(data), data.RandomInteger)
}

func (r LogicAppActionCustomResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_action_custom" "import" {
  name         = azurerm_logic_app_action_custom.test.name
  logic_app_id = azurerm_logic_app_action_custom.test.logic_app_id
  body         = azurerm_logic_app_action_custom.test.body
}
`, r.basic(data))
}

func (LogicAppActionCustomResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlaw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r LogicAppActionCustomResource) multipleActionsWithIntegrationAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_integration_account" "test" {
  name                = "acctestlaia-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_logic_app_workflow" "test" {
  name                             = "acctestlaw-%d"
  location                         = azurerm_resource_group.test.location
  resource_group_name              = azurerm_resource_group.test.name
  logic_app_integration_account_id = azurerm_logic_app_integration_account.test.id
}

resource "azurerm_logic_app_action_custom" "test" {
  name         = "acctestlaac1-%d"
  logic_app_id = azurerm_logic_app_workflow.test.id

  body = <<BODY
{
    "description": "Decode the AS2 payload",
    "inputs": {
        "messageHeaders": "@triggerOutputs()['headers']",
        "messageToDecode": "@triggerBody()"
    },
    "runAfter": {},
    "type": "As2Decode"
}
BODY
}

resource "azurerm_logic_app_action_custom" "test2" {
  name         = "acctestlaac2-%d"
  logic_app_id = azurerm_logic_app_workflow.test.id

  body = <<BODY
{
    "description": "Decode the EDIFACT message content.",
    "inputs": {
        "body": "@body('${azurerm_logic_app_action_custom.test.name}')?['messageToDecode']",
        "host": {
            "connection": {
                "name": "@parameters('$connections')['edifact']['connectionId']"
            }
        },
        "method": "post",
        "path": "/decode",
        "queries": {
            "componentSeparator": 58,
            "dataElementSeparator": 43,
            "decimalIndicator": "Comma",
            "payloadCharacterSet": "Legacy",
            "releaseIndicator": 63,
            "repetitionSeparator": 42,
            "segmentTerminator": 39,
            "segmentTerminatorSuffix": "None"
        }
    },
    "runAfter": {
        "${azurerm_logic_app_action_custom.test.name}": [
            "Succeeded"
        ]
    },
    "type": "ApiConnection"
}
BODY
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
