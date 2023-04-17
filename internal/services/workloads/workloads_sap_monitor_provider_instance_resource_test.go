package workloads_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/providerinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkloadsSAPMonitorProviderInstanceResource struct{}

func TestAccWorkloadsSAPMonitorProviderInstanceResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_monitor_provider_instance", "test")
	r := WorkloadsSAPMonitorProviderInstanceResource{}

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

func TestAccWorkloadsSAPMonitorProviderInstanceResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_monitor_provider_instance", "test")
	r := WorkloadsSAPMonitorProviderInstanceResource{}

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

func (r WorkloadsSAPMonitorProviderInstanceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := providerinstances.ParseProviderInstanceID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Workloads.ProviderInstances
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r WorkloadsSAPMonitorProviderInstanceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mpi-%d"
  location = "%s"
}

resource "azurerm_workloads_sap_monitor" "test" {
  name                = "acctest-sapmonitor-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r WorkloadsSAPMonitorProviderInstanceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_workloads_sap_monitor_provider_instance" "test" {
  name       = "acctest-mpi-%d"
  monitor_id = azurerm_workloads_sap_monitor.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r WorkloadsSAPMonitorProviderInstanceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_workloads_sap_monitor_provider_instance" "import" {
  name       = azurerm_workloads_sap_monitor_provider_instance.test.name
  monitor_id = azurerm_workloads_sap_monitor_provider_instance.test.monitor_id
}
`, r.basic(data))
}
