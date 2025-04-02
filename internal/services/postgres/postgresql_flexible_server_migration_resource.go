// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2024-08-01/migrations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PostgresqlFlexibleServerMigrationModel struct {
	Name                          string            `tfschema:"name"`
	Location                      string            `tfschema:"location"`
	ServerId                      string            `tfschema:"server_id"`
	CancelEnabled                 bool              `tfschema:"cancel_enabled"`
	DbsToCancelMigrationOn        []string          `tfschema:"dbs_to_cancel_migration_on"`
	DbsToMigrate                  []string          `tfschema:"dbs_to_migrate"`
	DbsToTriggerCutoverOn         []string          `tfschema:"dbs_to_trigger_cutover_on"`
	MigrateRolesEnabled           bool              `tfschema:"migrate_roles_enabled"`
	MigrationInstanceResourceId   string            `tfschema:"migration_instance_resource_id"`
	MigrationMode                 string            `tfschema:"migration_mode"`
	MigrationOption               string            `tfschema:"migration_option"`
	MigrationWindowEndTimeInUtc   string            `tfschema:"migration_window_end_time_in_utc"`
	MigrationWindowStartTimeInUtc string            `tfschema:"migration_window_start_time_in_utc"`
	OverwriteDbsInTargetEnabled   bool              `tfschema:"overwrite_dbs_in_target_enabled"`
	TriggerCutoverEnabled         bool              `tfschema:"trigger_cutover_enabled"`
	Tags                          map[string]string `tfschema:"tags"`
}

var _ sdk.ResourceWithUpdate = PostgresqlFlexibleServerMigrationResource{}

type PostgresqlFlexibleServerMigrationResource struct{}

func (r PostgresqlFlexibleServerMigrationResource) ModelObject() interface{} {
	return &PostgresqlFlexibleServerMigrationModel{}
}

func (r PostgresqlFlexibleServerMigrationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return migrations.ValidateMigrationID
}

func (r PostgresqlFlexibleServerMigrationResource) ResourceType() string {
	return "azurerm_postgresql_flexible_server_migration"
}

func (r PostgresqlFlexibleServerMigrationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FlexibleServerMigrationName,
		},

		"location": commonschema.Location(),

		"server_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: migrations.ValidateFlexibleServerID,
		},

		"cancel_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"dbs_to_cancel_migration_on": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.FlexibleServerDatabaseName,
			},
		},

		"dbs_to_migrate": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 50,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.FlexibleServerDatabaseName,
			},
		},

		"dbs_to_trigger_cutover_on": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.FlexibleServerDatabaseName,
			},
		},

		"migrate_roles_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"migration_instance_resource_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: migrations.ValidateFlexibleServerID,
		},

		"migration_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(migrations.PossibleValuesForMigrationMode(), false),
		},

		"migration_option": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(migrations.PossibleValuesForMigrationOption(), false),
		},

		"migration_window_end_time_in_utc": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"migration_window_start_time_in_utc": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"overwrite_dbs_in_target_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"trigger_cutover_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"tags": commonschema.Tags(),
	}
}

func (r PostgresqlFlexibleServerMigrationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r PostgresqlFlexibleServerMigrationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.Postgres.MigrationsClient

			var model PostgresqlFlexibleServerMigrationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			serverId, err := migrations.ParseFlexibleServerID(model.ServerId)
			if err != nil {
				return err
			}

			id := migrations.NewMigrationID(subscriptionId, serverId.ResourceGroupName, serverId.FlexibleServerName, model.Name)

			locks.ByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)
			defer locks.UnlockByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := migrations.MigrationResource{
				Location:   location.Normalize(model.Location),
				Properties: &migrations.MigrationResourceProperties{},
				Tags:       pointer.To(model.Tags),
			}

			cancelEnabled := migrations.CancelEnumFalse
			if model.CancelEnabled {
				cancelEnabled = migrations.CancelEnumTrue
			}
			parameters.Properties.Cancel = pointer.To(cancelEnabled)

			if v := model.DbsToCancelMigrationOn; v != nil {
				parameters.Properties.DbsToCancelMigrationOn = pointer.To(v)
			}

			if v := model.DbsToMigrate; v != nil {
				parameters.Properties.DbsToMigrate = pointer.To(v)
			}

			if v := model.DbsToTriggerCutoverOn; v != nil {
				parameters.Properties.DbsToTriggerCutoverOn = pointer.To(v)
			}

			migrateRolesEnabled := migrations.MigrateRolesEnumFalse
			if model.MigrateRolesEnabled {
				migrateRolesEnabled = migrations.MigrateRolesEnumTrue
			}
			parameters.Properties.MigrateRoles = pointer.To(migrateRolesEnabled)

			if v := model.MigrationInstanceResourceId; v != "" {
				parameters.Properties.MigrationInstanceResourceId = pointer.To(v)
			}

			if v := model.MigrationMode; v != "" {
				parameters.Properties.MigrationMode = pointer.To(migrations.MigrationMode(v))
			}

			if v := model.MigrationOption; v != "" {
				parameters.Properties.MigrationOption = pointer.To(migrations.MigrationOption(v))
			}

			if v := model.MigrationWindowEndTimeInUtc; v != "" {
				parameters.Properties.MigrationWindowEndTimeInUtc = pointer.To(v)
			}

			if v := model.MigrationWindowStartTimeInUtc; v != "" {
				parameters.Properties.MigrationWindowStartTimeInUtc = pointer.To(v)
			}

			overwriteDbsInTargetEnabled := migrations.OverwriteDbsInTargetEnumFalse
			if model.OverwriteDbsInTargetEnabled {
				overwriteDbsInTargetEnabled = migrations.OverwriteDbsInTargetEnumTrue
			}
			parameters.Properties.OverwriteDbsInTarget = pointer.To(overwriteDbsInTargetEnabled)

			triggerCutoverEnabled := migrations.TriggerCutoverEnumFalse
			if model.TriggerCutoverEnabled {
				triggerCutoverEnabled = migrations.TriggerCutoverEnumTrue
			}
			parameters.Properties.TriggerCutover = pointer.To(triggerCutoverEnabled)

			if _, err := client.Create(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r PostgresqlFlexibleServerMigrationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Postgres.MigrationsClient

			id, err := migrations.ParseMigrationID(metadata.ResourceData.Id())
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

			state := PostgresqlFlexibleServerMigrationModel{
				Name:     id.MigrationName,
				ServerId: migrations.NewFlexibleServerID(id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName).ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				if props := model.Properties; props != nil {
					state.CancelEnabled = pointer.From(props.Cancel) == migrations.CancelEnumTrue
					state.DbsToCancelMigrationOn = pointer.From(props.DbsToCancelMigrationOn)
					state.DbsToMigrate = pointer.From(props.DbsToMigrate)
					state.DbsToTriggerCutoverOn = pointer.From(props.DbsToTriggerCutoverOn)
					state.MigrateRolesEnabled = pointer.From(props.MigrateRoles) == migrations.MigrateRolesEnumTrue
					state.MigrationInstanceResourceId = pointer.From(props.MigrationInstanceResourceId)
					state.MigrationMode = string(pointer.From(props.MigrationMode))
					state.MigrationOption = string(pointer.From(props.MigrationOption))
					state.MigrationWindowEndTimeInUtc = pointer.From(props.MigrationWindowEndTimeInUtc)
					state.MigrationWindowStartTimeInUtc = pointer.From(props.MigrationWindowStartTimeInUtc)
					state.OverwriteDbsInTargetEnabled = pointer.From(props.OverwriteDbsInTarget) == migrations.OverwriteDbsInTargetEnumTrue
					state.TriggerCutoverEnabled = pointer.From(props.TriggerCutover) == migrations.TriggerCutoverEnumTrue
				}

				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r PostgresqlFlexibleServerMigrationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Postgres.MigrationsClient

			id, err := migrations.ParseMigrationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)
			defer locks.UnlockByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)

			var model PostgresqlFlexibleServerMigrationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := migrations.MigrationResourceForPatch{
				Properties: &migrations.MigrationResourcePropertiesForPatch{},
			}

			if metadata.ResourceData.HasChange("cancel_enabled") {
				cancelEnabled := migrations.CancelEnumFalse
				if model.CancelEnabled {
					cancelEnabled = migrations.CancelEnumTrue
				}
				parameters.Properties.Cancel = pointer.To(cancelEnabled)
			}

			if metadata.ResourceData.HasChange("dbs_to_cancel_migration_on") {
				parameters.Properties.DbsToCancelMigrationOn = pointer.To(model.DbsToCancelMigrationOn)
			}

			if metadata.ResourceData.HasChange("dbs_to_migrate") {
				parameters.Properties.DbsToMigrate = pointer.To(model.DbsToMigrate)
			}

			if metadata.ResourceData.HasChange("dbs_to_trigger_cutover_on") {
				parameters.Properties.DbsToTriggerCutoverOn = pointer.To(model.DbsToTriggerCutoverOn)
			}

			if metadata.ResourceData.HasChange("migrate_roles_enabled") {
				migrateRolesEnabled := migrations.MigrateRolesEnumFalse
				if model.MigrateRolesEnabled {
					migrateRolesEnabled = migrations.MigrateRolesEnumTrue
				}
				parameters.Properties.MigrateRoles = pointer.To(migrateRolesEnabled)
			}

			if metadata.ResourceData.HasChange("migration_mode") {
				parameters.Properties.MigrationMode = pointer.To(migrations.MigrationMode(model.MigrationMode))
			}

			if metadata.ResourceData.HasChange("migration_window_start_time_in_utc") {
				parameters.Properties.MigrationWindowStartTimeInUtc = pointer.To(model.MigrationWindowStartTimeInUtc)
			}

			if metadata.ResourceData.HasChange("overwrite_dbs_in_target_enabled") {
				overwriteDbsInTargetEnabled := migrations.OverwriteDbsInTargetEnumFalse
				if model.OverwriteDbsInTargetEnabled {
					overwriteDbsInTargetEnabled = migrations.OverwriteDbsInTargetEnumTrue
				}
				parameters.Properties.OverwriteDbsInTarget = pointer.To(overwriteDbsInTargetEnabled)
			}

			if metadata.ResourceData.HasChange("trigger_cutover_enabled") {
				triggerCutoverEnabled := migrations.TriggerCutoverEnumFalse
				if model.TriggerCutoverEnabled {
					triggerCutoverEnabled = migrations.TriggerCutoverEnumTrue
				}
				parameters.Properties.TriggerCutover = pointer.To(triggerCutoverEnabled)
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = pointer.To(model.Tags)
			}

			if _, err := client.Update(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r PostgresqlFlexibleServerMigrationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Postgres.MigrationsClient

			id, err := migrations.ParseMigrationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)
			defer locks.UnlockByName(id.FlexibleServerName, postgresqlFlexibleServerResourceName)

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
