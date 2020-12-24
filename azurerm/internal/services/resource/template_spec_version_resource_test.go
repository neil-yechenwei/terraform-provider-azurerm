package resource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type TemplateSpecVersionResource struct{}

func TestAccTemplateSpecVersion_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_template_spec_version", "test")
	r := TemplateSpecVersionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r TemplateSpecVersionResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	versionClient := client.Resource.TemplateSpecVersionsClient
	id, err := parse.TemplateSpecVersionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := versionClient.Get(ctx, id.ResourceGroup, id.TemplateSpecName, id.VersionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving Template Spec Version %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.VersionProperties != nil), nil
}

func (r TemplateSpecVersionResource) basic() string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_resource_group" "test" {
  name     = "acctestRG-templatespec-test01"
}

resource "azurerm_template_spec_version" "test" {
  name                = "acctest-TemplateSpecVersion-test01"
  resource_group_name = data.azurerm_resource_group.test.name
  location            = data.azurerm_resource_group.test.location
  template_spec_name  = "acctest-TemplateSpec-test01"
  description         = "test description version"

  template_content = <<DEPLOY
{
  "$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    "location": {
      "type": "string",
      "defaultValue": "eastus",
      "metadata":{
        "description": "Specify the location for the resources."
      }
    },
    "storageAccountType": {
      "type": "string",
      "defaultValue": "Standard_LRS",
      "metadata":{
        "description": "Specify the storage account type."
      }
    }
  },
  "variables": {
    "appServicePlanName": "[concat('plan', uniquestring(resourceGroup().id))]"
  },
  "resources": [
    {
      "type": "Microsoft.Web/serverfarms",
      "apiVersion": "2016-09-01",
      "name": "[variables('appServicePlanName')]",
      "location": "[parameters('location')]",
      "sku": {
        "name": "B1",
        "tier": "Basic",
        "size": "B1",
        "family": "B",
        "capacity": 1
      },
      "kind": "linux",
      "properties": {
        "perSiteScaling": false,
        "reserved": true,
        "targetWorkerCount": 0,
        "targetWorkerSizeId": 0
      }
    },
    {
      "type": "Microsoft.Resources/deployments",
      "apiVersion": "2020-06-01",
      "name": "createStorage",
      "properties": {
        "mode": "Incremental",
        "templateLink": {
          "relativePath": "./artifacts/linkedTemplate.json"
        },
        "parameters": {
          "storageAccountType": {
            "value": "[parameters('storageAccountType')]"
          }
        }
      }
    }
  ]
}
DEPLOY

  tags = {
    abc = "test"
  }
}
`)
}
