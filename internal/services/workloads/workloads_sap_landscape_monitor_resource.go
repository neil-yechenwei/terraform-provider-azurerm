package workloads

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/saplandscapemonitor"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type WorkloadsSAPLandscapeMonitorModel struct {
	Name                 string                 `tfschema:"name"`
	ResourceGroupName    string                 `tfschema:"resource_group_name"`
	Grouping             []Grouping             `tfschema:"grouping"`
	TopMetricsThresholds []TopMetricsThresholds `tfschema:"top_metrics_thresholds"`
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

type WorkloadsSAPLandscapeMonitorResource struct{}

var _ sdk.ResourceWithUpdate = WorkloadsSAPLandscapeMonitorResource{}

func (r WorkloadsSAPLandscapeMonitorResource) ResourceType() string {
	return "azurerm_workloads_sap_landscape_monitor"
}

func (r WorkloadsSAPLandscapeMonitorResource) ModelObject() interface{} {
	return &WorkloadsSAPLandscapeMonitorModel{}
}

func (r WorkloadsSAPLandscapeMonitorResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return saplandscapemonitor.ValidateMonitorID
}

func (r WorkloadsSAPLandscapeMonitorResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

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
	}
}

func (r WorkloadsSAPLandscapeMonitorResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r WorkloadsSAPLandscapeMonitorResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkloadsSAPLandscapeMonitorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Workloads.SapLandscapeMonitor
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := saplandscapemonitor.NewMonitorID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			if model.Grouping != nil || model.TopMetricsThresholds != nil {
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

				if _, err := client.Create(ctx, id, *properties); err != nil {
					return fmt.Errorf("creating %s: %+v", id, err)
				}
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r WorkloadsSAPLandscapeMonitorResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SapLandscapeMonitor
			id, err := saplandscapemonitor.ParseMonitorIDInsensitively(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model WorkloadsSAPLandscapeMonitorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("grouping") || metadata.ResourceData.HasChange("top_metrics_thresholds") {
				resp, err := client.Get(ctx, *id)
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

				if _, err := client.Update(ctx, *id, *properties); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}

func (r WorkloadsSAPLandscapeMonitorResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SapLandscapeMonitor
			id, err := saplandscapemonitor.ParseMonitorIDInsensitively(metadata.ResourceData.Id())
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

			state := WorkloadsSAPLandscapeMonitorModel{
				Name:              id.MonitorName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if properties := model.Properties; properties != nil {
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

			return metadata.Encode(&state)
		},
	}
}

func (r WorkloadsSAPLandscapeMonitorResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SapLandscapeMonitor
			id, err := saplandscapemonitor.ParseMonitorIDInsensitively(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
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
