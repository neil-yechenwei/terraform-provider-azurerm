package msi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/sdk/managedidentity"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type UserAssignedIdentityResource struct{}

func TestAccAzureRMUserAssignedIdentity_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_user_assigned_identity", "test")
	r := UserAssignedIdentityResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("principal_id").MatchesRegex(validate.UUIDRegExp),
				check.That(data.ResourceName).Key("client_id").MatchesRegex(validate.UUIDRegExp),
				check.That(data.ResourceName).Key("tenant_id").MatchesRegex(validate.UUIDRegExp),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMUserAssignedIdentity_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_user_assigned_identity", "test")
	r := UserAssignedIdentityResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("principal_id").MatchesRegex(validate.UUIDRegExp),
				check.That(data.ResourceName).Key("client_id").MatchesRegex(validate.UUIDRegExp),
				check.That(data.ResourceName).Key("tenant_id").MatchesRegex(validate.UUIDRegExp),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r UserAssignedIdentityResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := managedidentity.ParseUserAssignedIdentitiesID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSI.UserAssignedIdentitiesClient.UserAssignedIdentitiesGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r UserAssignedIdentityResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r UserAssignedIdentityResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "import" {
  name                = azurerm_user_assigned_identity.test.name
  resource_group_name = azurerm_user_assigned_identity.test.resource_group_name
  location            = azurerm_user_assigned_identity.test.location
}
`, template)
}
