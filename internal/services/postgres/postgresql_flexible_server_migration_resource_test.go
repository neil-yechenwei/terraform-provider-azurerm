// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2024-08-01/migrations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PostgresqlFlexibleServerMigrationTestResource struct{}

func TestAccPostgresqlFlexibleServerMigration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_migration", "test")
	r := PostgresqlFlexibleServerMigrationTestResource{}

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

func TestAccPostgresqlFlexibleServerMigration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_migration", "test")
	r := PostgresqlFlexibleServerMigrationTestResource{}

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

func TestAccPostgresqlFlexibleServerMigration_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_migration", "test")
	r := PostgresqlFlexibleServerMigrationTestResource{}

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

func TestAccPostgresqlFlexibleServerMigration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server_migration", "test")
	r := PostgresqlFlexibleServerMigrationTestResource{}

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

func (r PostgresqlFlexibleServerMigrationTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := migrations.ParseMigrationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Postgres.MigrationsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r PostgresqlFlexibleServerMigrationTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_postgresql_flexible_server_migration" "test" {
  name      = "acctest-pfsm-%d"
  location  = azurerm_resource_group.test.location
  server_id = azurerm_postgresql_flexible_server.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerMigrationTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server_migration" "import" {
  name      = azurerm_postgresql_flexible_server_migration.test.name
  location  = azurerm_postgresql_flexible_server_migration.test.location
  server_id = azurerm_postgresql_flexible_server_migration.test.server_id
}
`, r.basic(data))
}

func (r PostgresqlFlexibleServerMigrationTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_postgresql_flexible_server_migration" "test" {
  name      = "acctest-pfsm-%d"
  location  = azurerm_resource_group.test.location
  server_id = azurerm_postgresql_flexible_server.test.id

  tags = {
    Env = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerMigrationTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_postgresql_flexible_server_migration" "test" {
  name      = "acctest-pfsm-%d"
  location  = azurerm_resource_group.test.location
  server_id = azurerm_postgresql_flexible_server.test.id

  tags = {
    Env = "Test2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerMigrationTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-postgresql-%d"
  location = "%s"
}

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  version                = "12"
  sku_name               = "GP_Standard_D2s_v3"
  zone                   = "2"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
