// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securitycenter_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2023-10-01-preview/securityconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SecurityCenterSecurityConnectorResource struct{}

func TestAccSecurityCenterSecurityConnector_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_security_connector", "test")
	r := SecurityCenterSecurityConnectorResource{}

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

func TestAccSecurityCenterSecurityConnector_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_security_connector", "test")
	r := SecurityCenterSecurityConnectorResource{}

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

func TestAccSecurityCenterSecurityConnector_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_security_connector", "test")
	r := SecurityCenterSecurityConnectorResource{}

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

func TestAccSecurityCenterSecurityConnector_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_security_connector", "test")
	r := SecurityCenterSecurityConnectorResource{}

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

func (r SecurityCenterSecurityConnectorResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := securityconnectors.ParseSecurityConnectorID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SecurityCenter.SecurityConnectorsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r SecurityCenterSecurityConnectorResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-securityconnector-%d"
  location = "%s"
}

resource "azurerm_security_center_security_connector" "test" {
  name                = "acctest-securityconnector-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r SecurityCenterSecurityConnectorResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_security_center_security_connector" "import" {
  name                = azurerm_security_center_security_connector.test.name
  resource_group_name = azurerm_security_center_security_connector.test.resource_group_name
  location            = azurerm_security_center_security_connector.test.location
}
`, r.basic(data))
}

func (r SecurityCenterSecurityConnectorResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-securityconnector-%d"
  location = "%s"
}

resource "azurerm_security_center_security_connector" "test" {
  name                = "acctest-securityconnector-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    Env = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r SecurityCenterSecurityConnectorResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-securityconnector-%d"
  location = "%s"
}

resource "azurerm_security_center_security_connector" "test" {
  name                = "acctest-securityconnector-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    Env = "Test2"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
