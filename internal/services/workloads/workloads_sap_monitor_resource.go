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
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/saplandscapemonitor"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkloadsSAPMonitorModel struct {
	Name                     string                 `tfschema:"name"`
	ResourceGroupName        string                 `tfschema:"resource_group_name"`
	Location                 string                 `tfschema:"location"`
	AppLocation              string                 `tfschema:"app_location"`
	Grouping                 []Grouping             `tfschema:"grouping"`
	LogAnalyticsWorkspaceId  string                 `tfschema:"log_analytics_workspace_id"`
	ManagedResourceGroupName string                 `tfschema:"managed_resource_group_name"`
	RoutingPreference        string                 `tfschema:"routing_preference"`
	SubnetId                 string                 `tfschema:"subnet_id"`
	TopMetricsThresholds     []TopMetricsThresholds `tfschema:"top_metrics_thresholds"`
	ZoneRedundancyPreference string                 `tfschema:"zone_redundancy_preference"`
	Tags                     map[string]string      `tfschema:"tags"`
}

type Grouping struct {
	Landscape      []SAPLandscapeMonitorSIDMapping `tfschema:"landscape"`
	SAPApplication []SAPLandscapeMonitorSIDMapping `tfschema:"sap_application"`
}

type SAPLandscapeMonitorSIDMapping struct {
	Name   string   `tfschema:"name"`
	TopSid []string `tfschema:"top_sid"`
}

type TopMetricsThresholds struct {
	Green  float64 `tfschema:"green"`
	Name   string  `tfschema:"name"`
	Red    float64 `tfschema:"red"`
	Yellow float64 `tfschema:"yellow"`
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

		"grouping": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"landscape": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"top_sid": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"sap_application": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"top_sid": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},

		"log_analytics_workspace_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"top_metrics_thresholds": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"green": {
						Type:     pluginsdk.TypeFloat,
						Optional: true,
					},

					"red": {
						Type:     pluginsdk.TypeFloat,
						Optional: true,
					},

					"yellow": {
						Type:     pluginsdk.TypeFloat,
						Optional: true,
					},
				},
			},
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

			if model.Grouping != nil || model.TopMetricsThresholds != nil {
				landscapeMonitorClient := metadata.Client.Workloads.SapLandscapeMonitor
				landscapeMonitorId := saplandscapemonitor.NewMonitorID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName)

				properties := &saplandscapemonitor.SapLandscapeMonitor{
					Properties: &saplandscapemonitor.SapLandscapeMonitorProperties{},
				}

				grouping, err := expandGrouping(model.Grouping)
				if err != nil {
					return err
				}
				properties.Properties.Grouping = grouping

				topMetricsThresholds, err := expandTopMetricsThresholds(model.TopMetricsThresholds)
				if err != nil {
					return err
				}
				properties.Properties.TopMetricsThresholds = topMetricsThresholds

				if _, err := landscapeMonitorClient.Create(ctx, landscapeMonitorId, *properties); err != nil {
					return fmt.Errorf("creating %s: %+v", landscapeMonitorId, err)
				}
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

			if metadata.ResourceData.HasChange("grouping") || metadata.ResourceData.HasChange("top_metrics_thresholds") {
				landscapeMonitorClient := metadata.Client.Workloads.SapLandscapeMonitor
				landscapeMonitorId, err := saplandscapemonitor.ParseMonitorID(metadata.ResourceData.Id())
				if err != nil {
					return err
				}

				resp, err := landscapeMonitorClient.Get(ctx, *landscapeMonitorId)
				if err != nil {
					return fmt.Errorf("retrieving %s: %+v", *id, err)
				}

				properties := resp.Model
				if properties == nil {
					return fmt.Errorf("retrieving %s: properties was nil", id)
				}

				if metadata.ResourceData.HasChange("grouping") {
					grouping, err := expandGrouping(model.Grouping)
					if err != nil {
						return err
					}
					properties.Properties.Grouping = grouping
				}

				if metadata.ResourceData.HasChange("top_metrics_thresholds") {
					topMetricsThresholds, err := expandTopMetricsThresholds(model.TopMetricsThresholds)
					if err != nil {
						return err
					}
					properties.Properties.TopMetricsThresholds = topMetricsThresholds
				}

				if _, err := landscapeMonitorClient.Update(ctx, *landscapeMonitorId, *properties); err != nil {
					return fmt.Errorf("updating %s: %+v", *landscapeMonitorId, err)
				}
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

			landscapeMonitorClient := metadata.Client.Workloads.SapLandscapeMonitor
			landscapeMonitorId, err := saplandscapemonitor.ParseMonitorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			landscapeMonitorResp, err := landscapeMonitorClient.Get(ctx, *landscapeMonitorId)
			if err != nil {
				if !response.WasNotFound(landscapeMonitorResp.HttpResponse) {
					return fmt.Errorf("retrieving %s: %+v", *landscapeMonitorId, err)
				}
			} else {
				landscapeMonitorModel := landscapeMonitorResp.Model
				if landscapeMonitorModel == nil {
					return fmt.Errorf("retrieving %s: model was nil", landscapeMonitorId)
				}

				if properties := landscapeMonitorModel.Properties; properties != nil {
					grouping, err := flattenGrouping(properties.Grouping)
					if err != nil {
						return err
					}
					state.Grouping = grouping

					topMetricsThresholds, err := flattenTopMetricsThresholds(properties.TopMetricsThresholds)
					if err != nil {
						return err
					}
					state.TopMetricsThresholds = topMetricsThresholds
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r WorkloadsSAPMonitorResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			landscapeMonitorClient := metadata.Client.Workloads.SapLandscapeMonitor
			landscapeMonitorId, err := saplandscapemonitor.ParseMonitorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			landscapeMonitorResp, err := landscapeMonitorClient.Get(ctx, *landscapeMonitorId)
			if err != nil {
				if !response.WasNotFound(landscapeMonitorResp.HttpResponse) {
					return fmt.Errorf("retrieving %s: %+v", *landscapeMonitorId, err)
				}
			} else {
				if _, err := landscapeMonitorClient.Delete(ctx, *landscapeMonitorId); err != nil {
					return fmt.Errorf("deleting %s: %+v", landscapeMonitorId, err)
				}
			}

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

func expandGrouping(input []Grouping) (*saplandscapemonitor.SapLandscapeMonitorPropertiesGrouping, error) {
	if len(input) == 0 {
		return nil, nil
	}

	grouping := &input[0]
	result := saplandscapemonitor.SapLandscapeMonitorPropertiesGrouping{}

	landscape, err := expandSAPLandscapeMonitorSIDMapping(grouping.Landscape)
	if err != nil {
		return nil, err
	}
	result.Landscape = landscape

	sapApplication, err := expandSAPLandscapeMonitorSIDMapping(grouping.SAPApplication)
	if err != nil {
		return nil, err
	}
	result.SapApplication = sapApplication

	return &result, nil
}

func expandSAPLandscapeMonitorSIDMapping(input []SAPLandscapeMonitorSIDMapping) (*[]saplandscapemonitor.SapLandscapeMonitorSidMapping, error) {
	var result []saplandscapemonitor.SapLandscapeMonitorSidMapping

	for _, v := range input {
		output := saplandscapemonitor.SapLandscapeMonitorSidMapping{
			TopSid: &v.TopSid,
		}

		if v.Name != "" {
			output.Name = &v.Name
		}

		result = append(result, output)
	}

	return &result, nil
}

func expandTopMetricsThresholds(input []TopMetricsThresholds) (*[]saplandscapemonitor.SapLandscapeMonitorMetricThresholds, error) {
	var result []saplandscapemonitor.SapLandscapeMonitorMetricThresholds

	for _, v := range input {
		output := saplandscapemonitor.SapLandscapeMonitorMetricThresholds{
			Green:  &v.Green,
			Red:    &v.Red,
			Yellow: &v.Yellow,
		}

		if v.Name != "" {
			output.Name = &v.Name
		}

		result = append(result, output)
	}

	return &result, nil
}

func flattenGrouping(input *saplandscapemonitor.SapLandscapeMonitorPropertiesGrouping) ([]Grouping, error) {
	var result []Grouping
	if input == nil {
		return result, nil
	}

	output := Grouping{}

	landscape, err := flattenSAPLandscapeMonitorSIDMapping(input.Landscape)
	if err != nil {
		return nil, err
	}
	output.Landscape = landscape

	sapApplication, err := flattenSAPLandscapeMonitorSIDMapping(input.SapApplication)
	if err != nil {
		return nil, err
	}
	output.SAPApplication = sapApplication

	return append(result, output), nil
}

func flattenSAPLandscapeMonitorSIDMapping(input *[]saplandscapemonitor.SapLandscapeMonitorSidMapping) ([]SAPLandscapeMonitorSIDMapping, error) {
	var result []SAPLandscapeMonitorSIDMapping
	if input == nil {
		return result, nil
	}

	for _, v := range *input {
		sapLandscapeMonitorSIDMapping := SAPLandscapeMonitorSIDMapping{}

		if v.Name != nil {
			sapLandscapeMonitorSIDMapping.Name = *v.Name
		}

		if v.TopSid != nil {
			sapLandscapeMonitorSIDMapping.TopSid = *v.TopSid
		}

		result = append(result, sapLandscapeMonitorSIDMapping)
	}

	return result, nil
}

func flattenTopMetricsThresholds(input *[]saplandscapemonitor.SapLandscapeMonitorMetricThresholds) ([]TopMetricsThresholds, error) {
	var result []TopMetricsThresholds
	if input == nil {
		return result, nil
	}

	for _, v := range *input {
		output := TopMetricsThresholds{}

		if v.Green != nil {
			output.Green = *v.Green
		}

		if v.Name != nil {
			output.Name = *v.Name
		}

		if v.Red != nil {
			output.Red = *v.Red
		}

		if v.Yellow != nil {
			output.Yellow = *v.Yellow
		}

		result = append(result, output)
	}

	return result, nil
}
