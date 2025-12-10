// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backupinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backuppolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type BackupInstanceDataLakeStorageModel struct {
	Name             string   `tfschema:"name"`
	Location         string   `tfschema:"location"`
	VaultId          string   `tfschema:"vault_id"`
	BackupPolicyId   string   `tfschema:"backup_policy_id"`
	StorageAccountId string   `tfschema:"storage_account_id"`
	ContainerNames   []string `tfschema:"container_names"`
	ProtectionState  string   `tfschema:"protection_state"`
}

type DataProtectionBackupInstanceDataLakeStorageResource struct{}

var _ sdk.Resource = DataProtectionBackupInstanceDataLakeStorageResource{}

func (r DataProtectionBackupInstanceDataLakeStorageResource) ResourceType() string {
	return "azurerm_data_protection_backup_instance_data_lake_storage"
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) ModelObject() interface{} {
	return &BackupInstanceDataLakeStorageModel{}
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return backupinstances.ValidateBackupInstanceID
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&backuppolicies.BackupVaultId{}),

		"backup_policy_id": commonschema.ResourceIDReferenceRequired(&backuppolicies.BackupPolicyId{}),

		"storage_account_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.StorageAccountId{}),

		"container_names": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"protection_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model BackupInstanceDataLakeStorageModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.DataProtection.BackupInstanceClient

			vaultId, err := backupinstances.ParseBackupVaultID(model.VaultId)
			if err != nil {
				return err
			}

			id := backupinstances.NewBackupInstanceID(vaultId.SubscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			storageAccountId, err := commonids.ParseStorageAccountID(model.StorageAccountId)
			if err != nil {
				return err
			}

			policyId, err := backuppolicies.ParseBackupPolicyID(model.BackupPolicyId)
			if err != nil {
				return err
			}

			parameters := backupinstances.BackupInstanceResource{
				Properties: &backupinstances.BackupInstance{
					DataSourceInfo: backupinstances.Datasource{
						DatasourceType:   pointer.To("Microsoft.Storage/storageAccounts/adlsBlobServices"),
						ObjectType:       pointer.To("Datasource"),
						ResourceID:       storageAccountId.ID(),
						ResourceLocation: pointer.To(location.Normalize(model.Location)),
						ResourceName:     pointer.To(storageAccountId.StorageAccountName),
						ResourceType:     pointer.To("Microsoft.Storage/storageAccounts"),
						ResourceUri:      pointer.To(storageAccountId.ID()),
					},
					DataSourceSetInfo: &backupinstances.DatasourceSet{
						DatasourceType:   pointer.To("Microsoft.Storage/storageAccounts/adlsBlobServices"),
						ObjectType:       pointer.To("DatasourceSet"),
						ResourceID:       storageAccountId.ID(),
						ResourceLocation: pointer.To(location.Normalize(model.Location)),
						ResourceName:     pointer.To(storageAccountId.StorageAccountName),
						ResourceType:     pointer.To("Microsoft.Storage/storageAccounts"),
						ResourceUri:      pointer.To(storageAccountId.ID()),
					},
					FriendlyName: pointer.To(id.BackupInstanceName),
					PolicyInfo: backupinstances.PolicyInfo{
						PolicyId: policyId.ID(),
						PolicyParameters: &backupinstances.PolicyParameters{
							BackupDatasourceParametersList: &[]backupinstances.BackupDatasourceParameters{
								backupinstances.BlobBackupDatasourceParameters{
									ContainersList: model.ContainerNames,
								},
							},
						},
					},
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, parameters, backupinstances.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			// Service will continue to configure the protection after the resource is created and `provisioningState` returns `Succeeded`. At this time, service doesn't allow to change the resource until it is configured completely
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}

			stateConf := &pluginsdk.StateChangeConf{
				Delay:        5 * time.Second,
				Pending:      []string{string(backupinstances.CurrentProtectionStateConfiguringProtection)},
				Target:       []string{string(backupinstances.CurrentProtectionStateProtectionConfigured)},
				Refresh:      dataProtectionBackupInstanceDataLakeStorageStateRefreshFunc(ctx, client, id),
				PollInterval: 1 * time.Minute,
				Timeout:      time.Until(deadline),
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to become available: %s", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient

			id, err := backupinstances.ParseBackupInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			vaultId := backupinstances.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName)

			state := BackupInstanceDataLakeStorageModel{
				Name:    id.BackupInstanceName,
				VaultId: vaultId.ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Location = location.NormalizeNilable(props.DataSourceInfo.ResourceLocation)

					storageAccountId, err := commonids.ParseStorageAccountID(props.DataSourceInfo.ResourceID)
					if err != nil {
						return err
					}
					state.StorageAccountId = storageAccountId.ID()

					backupPolicyId, err := backuppolicies.ParseBackupPolicyID(props.PolicyInfo.PolicyId)
					if err != nil {
						return err
					}
					state.BackupPolicyId = backupPolicyId.ID()

					state.ProtectionState = pointer.FromEnum(props.CurrentProtectionState)

					if policyParams := props.PolicyInfo.PolicyParameters; policyParams != nil {
						if dataStoreParams := policyParams.BackupDatasourceParametersList; dataStoreParams != nil {
							if dsp := pointer.From(dataStoreParams); len(dsp) > 0 {
								if parameter, ok := dsp[0].(backupinstances.BlobBackupDatasourceParameters); ok {
									state.ContainerNames = parameter.ContainersList
								}
							}
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient

			id, err := backupinstances.ParseBackupInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model BackupInstanceDataLakeStorageModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			parameters := *existing.Model

			if metadata.ResourceData.HasChange("backup_policy_id") {
				policyId, err := backuppolicies.ParseBackupPolicyID(model.BackupPolicyId)
				if err != nil {
					return err
				}
				parameters.Properties.PolicyInfo.PolicyId = policyId.ID()
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, parameters, backupinstances.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			// Service will update the protection after the resource is updated and `provisioningState` returns `Succeeded`. At this time, service doesn't allow to change the resource until it is updated completely
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}

			stateConf := &pluginsdk.StateChangeConf{
				Delay:        5 * time.Second,
				Pending:      []string{string(backupinstances.CurrentProtectionStateUpdatingProtection)},
				Target:       []string{string(backupinstances.CurrentProtectionStateProtectionConfigured)},
				Refresh:      dataProtectionBackupInstanceDataLakeStorageStateRefreshFunc(ctx, client, *id),
				PollInterval: 1 * time.Minute,
				Timeout:      time.Until(deadline),
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to become available: %s", id, err)
			}

			return nil
		},
	}
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient

			id, err := backupinstances.ParseBackupInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			err = client.DeleteThenPoll(ctx, *id, backupinstances.DefaultDeleteOperationOptions())
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func dataProtectionBackupInstanceDataLakeStorageStateRefreshFunc(ctx context.Context, client *backupinstances.BackupInstancesClient, id backupinstances.BackupInstanceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("polling for %s: %+v", id, err)
		}

		if resp.Model == nil {
			return nil, "", fmt.Errorf("polling for %s: `model` was nil", id)
		}

		if resp.Model.Properties == nil {
			return nil, "", fmt.Errorf("polling for %s: `properties` was nil", id)
		}

		return resp, string(pointer.From(resp.Model.Properties.CurrentProtectionState)), nil
	}
}
