package workloads

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkloadsSAPMonitorModel struct {
	Name                     string            `tfschema:"name"`
	ResourceGroupName        string            `tfschema:"resource_group_name"`
	Location                 string            `tfschema:"location"`
	AppLocation              string            `tfschema:"app_location"`
	LogAnalyticsWorkspaceId  string            `tfschema:"log_analytics_workspace_id"`
	ManagedResourceGroupName string            `tfschema:"managed_resource_group_name"`
	RoutingPreference        string            `tfschema:"routing_preference"`
	SubnetId                 string            `tfschema:"subnet_id"`
	ZoneRedundancyPreference string            `tfschema:"zone_redundancy_preference"`
	Tags                     map[string]string `tfschema:"tags"`
}

type WorkloadsSAPMonitorResource struct{}

var _ sdk.ResourceWithUpdate = WorkloadsSAPMonitorResource{}

func (r WorkloadsSAPMonitorResource) ResourceType() string {
	return "azurerm_workloads_sap_monitor"
}

func (r WorkloadsSAPMonitorResource) ModelObject() interface{} {
	return &WorkloadsSAPMonitorModel{}
}

func (r WorkloadsSAPMonitorResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return monitors.ValidateMonitorID
}

func (r WorkloadsSAPMonitorResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"app_location": commonschema.Location(),

		"managed_resource_group_name": commonschema.ResourceGroupName(),

		"routing_preference": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(monitors.RoutingPreferenceDefault),
				string(monitors.RoutingPreferenceRouteAll),
			}, false),
		},

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: networkValidate.SubnetID,
		},

		"log_analytics_workspace_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"zone_redundancy_preference": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (r WorkloadsSAPMonitorResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r WorkloadsSAPMonitorResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkloadsSAPMonitorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Workloads.Monitors
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := monitors.NewMonitorID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			routingPreference := monitors.RoutingPreference(model.RoutingPreference)

			parameters := &monitors.Monitor{
				Location: location.Normalize(model.Location),
				Tags:     &model.Tags,
				Properties: &monitors.MonitorProperties{
					AppLocation: utils.String(model.AppLocation),
					ManagedResourceGroupConfiguration: &monitors.ManagedRGConfiguration{
						Name: utils.String(model.ManagedResourceGroupName),
					},
					MonitorSubnet:     utils.String(model.SubnetId),
					RoutingPreference: &routingPreference,
				},
			}

			if v := model.LogAnalyticsWorkspaceId; v != "" {
				parameters.Properties.LogAnalyticsWorkspaceArmId = utils.String(model.LogAnalyticsWorkspaceId)
			}

			if v := model.ZoneRedundancyPreference; v != "" {
				parameters.Properties.ZoneRedundancyPreference = utils.String(model.ZoneRedundancyPreference)
			}

			if err := client.CreateThenPoll(ctx, id, *parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r WorkloadsSAPMonitorResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.Monitors

			id, err := monitors.ParseMonitorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model WorkloadsSAPMonitorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := &monitors.UpdateMonitorRequest{}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = &model.Tags
			}

			if _, err := client.Update(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r WorkloadsSAPMonitorResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.Monitors

			id, err := monitors.ParseMonitorID(metadata.ResourceData.Id())
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

			state := WorkloadsSAPMonitorModel{
				Name:              id.MonitorName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			if props := model.Properties; props != nil {
				state.AppLocation = *props.AppLocation
				state.RoutingPreference = string(*props.RoutingPreference)
				state.SubnetId = *props.MonitorSubnet

				if v := props.ManagedResourceGroupConfiguration; v != nil {
					state.ManagedResourceGroupName = *v.Name
				}

				if v := props.LogAnalyticsWorkspaceArmId; v != nil {
					state.LogAnalyticsWorkspaceId = *v
				}

				if v := props.ZoneRedundancyPreference; v != nil {
					state.ZoneRedundancyPreference = *v
				}
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r WorkloadsSAPMonitorResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.Monitors
			id, err := monitors.ParseMonitorID(metadata.ResourceData.Id())
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
