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
)

type PostgresqlFlexibleServerMigrationModel struct {
	Name     string            `tfschema:"name"`
	Location string            `tfschema:"location"`
	ServerId string            `tfschema:"server_id"`
	Tags     map[string]string `tfschema:"tags"`
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
				Tags: pointer.To(model.Tags),
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
