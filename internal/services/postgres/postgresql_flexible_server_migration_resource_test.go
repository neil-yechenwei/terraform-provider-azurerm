// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

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
  name                            = "acctest-pfsm-%d"
  location                        = azurerm_resource_group.test.location
  server_id                       = azurerm_postgresql_flexible_server.test.id
  source_db_server_resource_id    = azurerm_postgresql_server.test.id
  dbs_to_migrate                  = [azurerm_postgresql_database.test.name]
  overwrite_dbs_in_target_enabled = true

  secrets {
    admin_credentials {
      source_server_password = azurerm_postgresql_server.test.administrator_login_password
      target_server_password = azurerm_postgresql_flexible_server.test.administrator_password
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerMigrationTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server_migration" "import" {
  name                            = azurerm_postgresql_flexible_server_migration.test.name
  location                        = azurerm_postgresql_flexible_server_migration.test.location
  server_id                       = azurerm_postgresql_flexible_server_migration.test.server_id
  source_db_server_resource_id    = azurerm_postgresql_flexible_server_migration.test.source_db_server_resource_id
  dbs_to_migrate                  = azurerm_postgresql_flexible_server_migration.test.dbs_to_migrate
  overwrite_dbs_in_target_enabled = azurerm_postgresql_flexible_server_migration.test.overwrite_dbs_in_target_enabled

  secrets {
    admin_credentials {
      source_server_password = azurerm_postgresql_server.test.administrator_login_password
      target_server_password = azurerm_postgresql_flexible_server.test.administrator_password
    }
  }
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
  name                                                     = "acctest-pfsm-%d"
  location                                                 = azurerm_resource_group.test.location
  server_id                                                = azurerm_postgresql_flexible_server.test.id
  cancel_enabled                                           = false
  dbs_to_migrate                                           = [azurerm_postgresql_database.test.name]
  migration_instance_resource_id                           = azurerm_postgresql_flexible_server.test.id
  migration_mode                                           = "Online"
  migration_option                                         = "ValidateAndMigrate"
  migration_window_start_time_in_utc                       = "%s"
  migration_window_end_time_in_utc                         = "%s"
  overwrite_dbs_in_target_enabled                          = true
  setup_logical_replication_on_source_db_if_needed_enabled = true
  source_db_server_resource_id                             = azurerm_postgresql_server.test.id
  source_type                                              = "PostgreSQLSingleServer"
  ssl_mode                                                 = "Prefer"
  start_data_migration_enabled                             = true

  secrets {
    admin_credentials {
      source_server_password = azurerm_postgresql_server.test.administrator_login_password
      target_server_password = azurerm_postgresql_flexible_server.test.administrator_password
    }

    source_server_username = azurerm_postgresql_server.test.administrator_login
    target_server_username = azurerm_postgresql_flexible_server.test.administrator_login
  }

  tags = {
    Env = "Test"
  }
}
`, r.template(data), data.RandomInteger, time.Now().UTC().Format(time.RFC3339), time.Now().Add(time.Duration(15)*time.Minute).UTC().Format(time.RFC3339))
}

func (r PostgresqlFlexibleServerMigrationTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_postgresql_server" "test2" {
  name                         = "acctest-psql-server2-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  sku_name                     = "B_Gen5_1"
  version                      = "11"
  storage_mb                   = 51200
  ssl_enforcement_enabled      = true
}

resource "azurerm_postgresql_database" "test2" {
  name                = "acctest-postgresqldb2-%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_postgresql_server.test2.name
  charset             = "UTF8"
  collation           = "English_United States.1252"
}

resource "azurerm_postgresql_flexible_server_migration" "test" {
  name                                                     = "acctest-pfsm-%d"
  location                                                 = azurerm_resource_group.test.location
  server_id                                                = azurerm_postgresql_flexible_server.test.id
  cancel_enabled                                           = true
  dbs_to_cancel_migration_on                               = [azurerm_postgresql_database.test.name]
  trigger_cutover_enabled                                  = true
  dbs_to_trigger_cutover_on                                = [azurerm_postgresql_database.test.name]
  migrate_roles_enabled                                    = true
  migration_mode                                           = "Offline"
  overwrite_dbs_in_target_enabled                          = false
  setup_logical_replication_on_source_db_if_needed_enabled = false
  start_data_migration_enabled                             = false
  source_db_server_resource_id                             = azurerm_postgresql_server.test2.id
  dbs_to_migrate                                           = [azurerm_postgresql_database.test2.name]

  secrets {
    admin_credentials {
      source_server_password = "H@Sh1CoR34!"
      target_server_password = "QAZwsx1234"
    }

    source_server_username = "acctestun2"
    target_server_username = "adminTerraform2"
  }

  tags = {
    Env = "Test2"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r PostgresqlFlexibleServerMigrationTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-postgresql-%d"
  location = "%s"
}

resource "azurerm_postgresql_server" "test" {
  name                         = "acctest-psql-server-%d"
  location                     = azurerm_resource_group.test.location
  resource_group_name          = azurerm_resource_group.test.name
  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  sku_name                     = "B_Gen5_1"
  version                      = "11"
  storage_mb                   = 51200
  ssl_enforcement_enabled      = true
}

resource "azurerm_postgresql_database" "test" {
  name                = "acctest-postgresqldb-%d"
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_postgresql_server.test.name
  charset             = "UTF8"
  collation           = "English_United States.1252"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
