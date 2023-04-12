package workloads

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapapplicationserverinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type WorkloadsSAPApplicationServerInstanceModel struct {
	Name                 string            `tfschema:"name"`
	ResourceGroupName    string            `tfschema:"resource_group_name"`
	Location             string            `tfschema:"location"`
	SAPVirtualInstanceId string            `tfschema:"sap_virtual_instance_id"`
	Tags                 map[string]string `tfschema:"tags"`
}

type WorkloadsSAPApplicationServerInstanceResource struct{}

var _ sdk.ResourceWithUpdate = WorkloadsSAPApplicationServerInstanceResource{}

func (r WorkloadsSAPApplicationServerInstanceResource) ResourceType() string {
	return "azurerm_workloads_sap_application_server_instance"
}

func (r WorkloadsSAPApplicationServerInstanceResource) ModelObject() interface{} {
	return &WorkloadsSAPApplicationServerInstanceModel{}
}

func (r WorkloadsSAPApplicationServerInstanceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sapapplicationserverinstances.ValidateApplicationInstanceID
}

func (r WorkloadsSAPApplicationServerInstanceResource) Arguments() map[string]*pluginsdk.Schema {
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
			ValidateFunc: sapapplicationserverinstances.ValidateSapVirtualInstanceID,
		},

		"tags": commonschema.Tags(),
	}
}

func (r WorkloadsSAPApplicationServerInstanceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r WorkloadsSAPApplicationServerInstanceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkloadsSAPApplicationServerInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			sapVirtualInstanceId, err := sapapplicationserverinstances.ParseSapVirtualInstanceID(model.SAPVirtualInstanceId)
			if err != nil {
				return err
			}

			client := metadata.Client.Workloads.SAPApplicationServerInstances
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := sapapplicationserverinstances.NewApplicationInstanceID(subscriptionId, sapVirtualInstanceId.ResourceGroupName, sapVirtualInstanceId.SapVirtualInstanceName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := &sapapplicationserverinstances.SAPApplicationServerInstance{
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

func (r WorkloadsSAPApplicationServerInstanceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPApplicationServerInstances

			id, err := sapapplicationserverinstances.ParseApplicationInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model WorkloadsSAPApplicationServerInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := &sapapplicationserverinstances.UpdateSAPApplicationInstanceRequest{}

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

func (r WorkloadsSAPApplicationServerInstanceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPApplicationServerInstances

			id, err := sapapplicationserverinstances.ParseApplicationInstanceID(metadata.ResourceData.Id())
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

			state := WorkloadsSAPApplicationServerInstanceModel{
				Name:                 id.ApplicationInstanceName,
				ResourceGroupName:    id.ResourceGroupName,
				Location:             location.Normalize(model.Location),
				SAPVirtualInstanceId: sapapplicationserverinstances.NewSapVirtualInstanceID(id.SubscriptionId, id.ResourceGroupName, id.SapVirtualInstanceName).ID(),
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r WorkloadsSAPApplicationServerInstanceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPApplicationServerInstances

			id, err := sapapplicationserverinstances.ParseApplicationInstanceID(metadata.ResourceData.Id())
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
