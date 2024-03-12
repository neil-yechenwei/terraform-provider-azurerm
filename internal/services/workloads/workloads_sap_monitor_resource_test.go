package workloads_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkloadsSAPMonitorResource struct{}

func TestAccWorkloadsSAPMonitor_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_monitor", "test")
	r := WorkloadsSAPMonitorResource{}

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

func TestAccWorkloadsSAPMonitor_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_monitor", "test")
	r := WorkloadsSAPMonitorResource{}

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

func TestAccWorkloadsSAPMonitor_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_monitor", "test")
	r := WorkloadsSAPMonitorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWorkloadsSAPMonitor_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_monitor", "test")
	r := WorkloadsSAPMonitorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r WorkloadsSAPMonitorResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := monitors.ParseMonitorID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Workloads.Monitors
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r WorkloadsSAPMonitorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sapmonitor-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-subnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r WorkloadsSAPMonitorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_workloads_sap_monitor" "test" {
  name                        = "acctest-sapmonitor-%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  app_location                = "%s"
  managed_resource_group_name = "managedRGSAPMonitor"
  routing_preference          = "RouteAll"
  subnet_id                   = azurerm_subnet.test.id
}
`, r.template(data), data.RandomInteger, data.Locations.Secondary)
}

func (r WorkloadsSAPMonitorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_workloads_sap_monitor" "import" {
  name                        = azurerm_workloads_sap_monitor.test.name
  resource_group_name         = azurerm_workloads_sap_monitor.test.resource_group_name
  location                    = azurerm_workloads_sap_monitor.test.location
  app_location                = azurerm_workloads_sap_monitor.test.app_location
  managed_resource_group_name = azurerm_workloads_sap_monitor.test.managed_resource_group_name
  routing_preference          = azurerm_workloads_sap_monitor.test.routing_preference
  subnet_id                   = azurerm_workloads_sap_monitor.test.subnet_id
}
`, r.basic(data))
}

func (r WorkloadsSAPMonitorResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_workloads_sap_monitor" "test" {
  name                        = "acctest-sapmonitor-%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  app_location                = "%s"
  managed_resource_group_name = "managedRGSAPMonitor"
  routing_preference          = "RouteAll"
  subnet_id                   = azurerm_subnet.test.id
  log_analytics_workspace_id  = azurerm_log_analytics_workspace.test.id
  zone_redundancy_preference  = "ZoneRedundantApp"

  tags = {
    Env = "Test"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.Locations.Secondary)
}

func (r WorkloadsSAPMonitorResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_workloads_sap_monitor" "test" {
  name                        = "acctest-sapmonitor-%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  app_location                = "%s"
  managed_resource_group_name = "managedRGSAPMonitor"
  routing_preference          = "RouteAll"
  subnet_id                   = azurerm_subnet.test.id
  log_analytics_workspace_id  = azurerm_log_analytics_workspace.test.id
  zone_redundancy_preference  = "ZoneRedundantApp"

  tags = {
    Env = "Test2"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.Locations.Secondary)
}
