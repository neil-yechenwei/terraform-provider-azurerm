// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backupinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataProtectionBackupInstanceDataLakeStorageResource struct{}

func TestAccDataProtectionBackupInstanceDataLakeStorage_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_data_lake_storage", "test")
	r := DataProtectionBackupInstanceDataLakeStorageResource{}

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

func TestAccDataProtectionBackupInstanceDataLakeStorage_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_data_lake_storage", "test")
	r := DataProtectionBackupInstanceDataLakeStorageResource{}

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

func TestAccDataProtectionBackupInstanceDataLakeStorage_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_instance_data_lake_storage", "test")
	r := DataProtectionBackupInstanceDataLakeStorageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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

func (r DataProtectionBackupInstanceDataLakeStorageResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := backupinstances.ParseBackupInstanceID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DataProtection.BackupInstanceClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-dataprotection-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
  is_hns_enabled           = true
}

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-dataprotection-vault-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
  soft_delete         = "Off"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Account Backup Contributor"
  principal_id         = azurerm_data_protection_backup_vault.test.identity.0.principal_id
}

resource "azurerm_data_protection_backup_policy_data_lake_storage" "test" {
  name                            = "acctest-dp-%d"
  vault_id                        = azurerm_data_protection_backup_vault.test.id
  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]

  default_retention_rule {
    life_cycle {
      duration        = "P4M"
      data_store_type = "VaultStore"
    }
  }

  depends_on = [azurerm_role_assignment.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_data_protection_backup_instance_data_lake_storage" "test" {
  name               = "acctest-dbi-%d"
  location           = azurerm_resource_group.test.location
  vault_id           = azurerm_data_protection_backup_vault.test.id
  storage_account_id = azurerm_storage_account.test.id
  backup_policy_id   = azurerm_data_protection_backup_policy_data_lake_storage.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_instance_data_lake_storage" "import" {
  name               = azurerm_data_protection_backup_instance_data_lake_storage.test.name
  location           = azurerm_data_protection_backup_instance_data_lake_storage.test.location
  vault_id           = azurerm_data_protection_backup_instance_data_lake_storage.test.vault_id
  storage_account_id = azurerm_data_protection_backup_instance_data_lake_storage.test.storage_account_id
  backup_policy_id   = azurerm_data_protection_backup_instance_data_lake_storage.test.backup_policy_id
}
`, r.basic(data))
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_data_protection_backup_policy_data_lake_storage" "test2" {
  name                            = "acctest-dp2-%d"
  vault_id                        = azurerm_data_protection_backup_vault.test.id
  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]

  default_retention_rule {
    life_cycle {
      duration        = "P4M"
      data_store_type = "VaultStore"
    }
  }

  depends_on = [azurerm_role_assignment.test]
}

resource "azurerm_data_protection_backup_instance_data_lake_storage" "test" {
  name               = "acctest-dbi-%d"
  location           = azurerm_resource_group.test.location
  vault_id           = azurerm_data_protection_backup_vault.test.id
  storage_account_id = azurerm_storage_account.test.id
  backup_policy_id   = azurerm_data_protection_backup_policy_data_lake_storage.test2.id
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}
