package workloads_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapapplicationserverinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkloadsSAPApplicationServerInstanceResource struct{}

func TestAccWorkloadsSAPApplicationServerInstance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_application_server_instance", "test")
	r := WorkloadsSAPApplicationServerInstanceResource{}
	sapVISNameSuffix := RandomInt()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, sapVISNameSuffix),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWorkloadsSAPApplicationServerInstance_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_application_server_instance", "test")
	r := WorkloadsSAPApplicationServerInstanceResource{}
	sapVISNameSuffix := RandomInt()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, sapVISNameSuffix),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, sapVISNameSuffix),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
	})
}

func TestAccWorkloadsSAPApplicationServerInstance_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_application_server_instance", "test")
	r := WorkloadsSAPApplicationServerInstanceResource{}
	sapVISNameSuffix := RandomInt()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, sapVISNameSuffix),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWorkloadsSAPApplicationServerInstance_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_application_server_instance", "test")
	r := WorkloadsSAPApplicationServerInstanceResource{}
	sapVISNameSuffix := RandomInt()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, sapVISNameSuffix),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, sapVISNameSuffix),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r WorkloadsSAPApplicationServerInstanceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := sapapplicationserverinstances.ParseApplicationInstanceID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Workloads.SAPApplicationServerInstances
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r WorkloadsSAPApplicationServerInstanceResource) basic(data acceptance.TestData, sapVISNameSuffix int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_workloads_sap_application_server_instance" "test" {
  name                    = "acctest-sapapp-%d"
  resource_group_name     = azurerm_resource_group.test.name
  location                = azurerm_resource_group.test.location
  sap_virtual_instance_id = azurerm_workloads_sap_three_tier_virtual_instance.test.id
}
`, WorkloadsSAPThreeTierVirtualInstanceResource{}.basic(data, sapVISNameSuffix), data.RandomInteger)
}

func (r WorkloadsSAPApplicationServerInstanceResource) requiresImport(data acceptance.TestData, sapVISNameSuffix int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_workloads_sap_application_server_instance" "import" {
  name                    = azurerm_workloads_sap_application_server_instance.test.name
  resource_group_name     = azurerm_workloads_sap_application_server_instance.test.resource_group_name
  location                = azurerm_workloads_sap_application_server_instance.test.location
  sap_virtual_instance_id = azurerm_workloads_sap_application_server_instance.test.sap_virtual_instance_id
}
`, r.basic(data, sapVISNameSuffix))
}

func (r WorkloadsSAPApplicationServerInstanceResource) complete(data acceptance.TestData, sapVISNameSuffix int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_workloads_sap_application_server_instance" "test" {
  name                    = "acctest-sapapp-%d"
  resource_group_name     = azurerm_resource_group.test.name
  location                = azurerm_resource_group.test.location
  sap_virtual_instance_id = azurerm_workloads_sap_three_tier_virtual_instance.test.id

  tags = {
    Env = "Test"
  }
}
`, WorkloadsSAPThreeTierVirtualInstanceResource{}.basic(data, sapVISNameSuffix), data.RandomInteger)
}

func (r WorkloadsSAPApplicationServerInstanceResource) update(data acceptance.TestData, sapVISNameSuffix int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_workloads_sap_application_server_instance" "test" {
  name                    = "acctest-sapapp-%d"
  resource_group_name     = azurerm_resource_group.test.name
  location                = azurerm_resource_group.test.location
  sap_virtual_instance_id = azurerm_workloads_sap_three_tier_virtual_instance.test.id

  tags = {
    Env = "Test2"
  }
}
`, WorkloadsSAPThreeTierVirtualInstanceResource{}.basic(data, sapVISNameSuffix), data.RandomInteger)
}
