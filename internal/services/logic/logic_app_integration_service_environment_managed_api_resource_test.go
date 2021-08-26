package logic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogicAppIntegrationServiceEnvironmentManagedApiResource struct{}

func TestAccLogicAppIntegrationServiceEnvironmentManagedApi_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_service_environment_managed_api", "test")
	r := LogicAppIntegrationServiceEnvironmentManagedApiResource{}

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

func TestAccLogicAppIntegrationServiceEnvironmentManagedApi_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_service_environment_managed_api", "test")
	r := LogicAppIntegrationServiceEnvironmentManagedApiResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r LogicAppIntegrationServiceEnvironmentManagedApiResource) Exists(ctx context.Context, client *clients.Client, state *acceptance.InstanceState) (*bool, error) {
	id, err := parse.IntegrationServiceEnvironmentManagedApiID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Logic.IntegrationServiceEnvironmentManagedApiClient.Get(ctx, id.ResourceGroup, id.IntegrationServiceEnvironmentName, id.ManagedApiName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(resp.Properties != nil), nil
}

func (r LogicAppIntegrationServiceEnvironmentManagedApiResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-logic-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["172.16.0.0/16"]
}

resource "azurerm_subnet" "isesubnet1" {
  name                 = "isesubnet1"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["172.16.0.0/24"]

  delegation {
    name = "integrationServiceEnvironments"
    service_delegation {
      name    = "Microsoft.Logic/integrationServiceEnvironments"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_subnet" "isesubnet2" {
  name                 = "isesubnet2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["172.16.1.0/24"]
}

resource "azurerm_subnet" "isesubnet3" {
  name                 = "isesubnet3"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["172.16.2.0/24"]
}

resource "azurerm_subnet" "isesubnet4" {
  name                 = "isesubnet4"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["172.16.3.0/24"]
}

resource "azurerm_integration_service_environment" "test" {
  name                 = "acctest-ise-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  sku_name             = "Developer_0"
  access_endpoint_type = "Internal"

  virtual_network_subnet_ids = [
    azurerm_subnet.isesubnet1.id,
    azurerm_subnet.isesubnet2.id,
    azurerm_subnet.isesubnet3.id,
    azurerm_subnet.isesubnet4.id
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r LogicAppIntegrationServiceEnvironmentManagedApiResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_service_environment_managed_api" "test" {
  name                                 = "servicebus"
  resource_group_name                  = azurerm_resource_group.test.name
  integration_service_environment_name = azurerm_integration_service_environment.test.name
}
`, r.template(data))
}

func (r LogicAppIntegrationServiceEnvironmentManagedApiResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_service_environment_managed_api" "import" {
  name                                 = azurerm_logic_app_integration_service_environment_managed_api.test.name
  resource_group_name                  = azurerm_logic_app_integration_service_environment_managed_api.test.resource_group_name
  integration_service_environment_name = azurerm_logic_app_integration_service_environment_managed_api.test.integration_service_environment_name
}
`, r.basic(data))
}
