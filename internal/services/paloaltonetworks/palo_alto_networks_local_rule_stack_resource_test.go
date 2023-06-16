package paloaltonetworks_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrulestacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PaloAltoNetworksLocalRuleStackResource struct{}

func TestAccPaloAltoNetworksLocalRuleStack_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_networks_local_rule_stack", "test")
	r := PaloAltoNetworksLocalRuleStackResource{}

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

func TestAccPaloAltoNetworksLocalRuleStack_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_networks_local_rule_stack", "test")
	r := PaloAltoNetworksLocalRuleStackResource{}

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

func TestAccPaloAltoNetworksLocalRuleStack_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_networks_local_rule_stack", "test")
	r := PaloAltoNetworksLocalRuleStackResource{}

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

func TestAccPaloAltoNetworksLocalRuleStack_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_networks_local_rule_stack", "test")
	r := PaloAltoNetworksLocalRuleStackResource{}

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

func (r PaloAltoNetworksLocalRuleStackResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := localrulestacks.ParseLocalRuleStackID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.PaloAltoNetworks.LocalRuleStacks
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r PaloAltoNetworksLocalRuleStackResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-lrs-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r PaloAltoNetworksLocalRuleStackResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_palo_alto_networks_local_rule_stack" "test" {
  name                = "acctest-lrs-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, r.template(data), data.RandomInteger)
}

func (r PaloAltoNetworksLocalRuleStackResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_palo_alto_networks_local_rule_stack" "import" {
  name                = azurerm_palo_alto_networks_local_rule_stack.test.name
  resource_group_name = azurerm_palo_alto_networks_local_rule_stack.test.resource_group_name
  location            = azurerm_palo_alto_networks_local_rule_stack.test.location
}
`, r.basic(data))
}

func (r PaloAltoNetworksLocalRuleStackResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_palo_alto_networks_local_rule_stack" "test" {
  name                = "acctest-lrs-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Env = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r PaloAltoNetworksLocalRuleStackResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_palo_alto_networks_local_rule_stack" "test" {
  name                = "acctest-lrs-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    Env = "Test2"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}
