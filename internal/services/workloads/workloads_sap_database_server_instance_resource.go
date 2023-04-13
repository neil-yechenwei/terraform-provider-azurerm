package workloads

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapdatabaseinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type WorkloadsSAPDatabaseServerInstanceModel struct {
	Name                 string            `tfschema:"name"`
	ResourceGroupName    string            `tfschema:"resource_group_name"`
	Location             string            `tfschema:"location"`
	SAPVirtualInstanceId string            `tfschema:"sap_virtual_instance_id"`
	Tags                 map[string]string `tfschema:"tags"`
}

type WorkloadsSAPDatabaseServerInstanceResource struct{}

var _ sdk.ResourceWithUpdate = WorkloadsSAPDatabaseServerInstanceResource{}

func (r WorkloadsSAPDatabaseServerInstanceResource) ResourceType() string {
	return "azurerm_workloads_sap_database_server_instance"
}

func (r WorkloadsSAPDatabaseServerInstanceResource) ModelObject() interface{} {
	return &WorkloadsSAPDatabaseServerInstanceModel{}
}

func (r WorkloadsSAPDatabaseServerInstanceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sapdatabaseinstances.ValidateDatabaseInstanceID
}

func (r WorkloadsSAPDatabaseServerInstanceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sap_virtual_instance_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: sapdatabaseinstances.ValidateSapVirtualInstanceID,
		},

		"tags": commonschema.Tags(),
	}
}

func (r WorkloadsSAPDatabaseServerInstanceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r WorkloadsSAPDatabaseServerInstanceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkloadsSAPDatabaseServerInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			sapVirtualInstanceId, err := sapdatabaseinstances.ParseSapVirtualInstanceID(model.SAPVirtualInstanceId)
			if err != nil {
				return err
			}

			client := metadata.Client.Workloads.SAPDatabaseInstances
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := sapdatabaseinstances.NewDatabaseInstanceID(subscriptionId, sapVirtualInstanceId.ResourceGroupName, sapVirtualInstanceId.SapVirtualInstanceName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := &sapdatabaseinstances.SAPDatabaseInstance{
				Location: location.Normalize(model.Location),
				Tags:     &model.Tags,
			}

			if err := client.CreateThenPoll(ctx, id, *parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r WorkloadsSAPDatabaseServerInstanceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPDatabaseInstances

			id, err := sapdatabaseinstances.ParseDatabaseInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model WorkloadsSAPDatabaseServerInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := &sapdatabaseinstances.UpdateSAPDatabaseInstanceRequest{}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = &model.Tags
			}

			if err := client.UpdateThenPoll(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r WorkloadsSAPDatabaseServerInstanceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPDatabaseInstances

			id, err := sapdatabaseinstances.ParseDatabaseInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := WorkloadsSAPDatabaseServerInstanceModel{
				Name:                 id.DatabaseInstanceName,
				ResourceGroupName:    id.ResourceGroupName,
				Location:             location.Normalize(model.Location),
				SAPVirtualInstanceId: sapdatabaseinstances.NewSapVirtualInstanceID(id.SubscriptionId, id.ResourceGroupName, id.SapVirtualInstanceName).ID(),
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r WorkloadsSAPDatabaseServerInstanceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPDatabaseInstances

			id, err := sapdatabaseinstances.ParseDatabaseInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
