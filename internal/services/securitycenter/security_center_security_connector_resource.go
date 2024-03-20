// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securitycenter

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2023-10-01-preview/securityconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SecurityCenterSecurityConnectorModel struct {
	Name                      string                      `tfschema:"name"`
	ResourceGroupName         string                      `tfschema:"resource_group_name"`
	Location                  string                      `tfschema:"location"`
	EnvironmentName           string                      `tfschema:"environment_name"`
	HierarchyIdentifier       string                      `tfschema:"hierarchy_identifier"`
	AwsEnvironmentData        []AwsEnvironmentData        `tfschema:"aws_environment_data"`
	GcpProjectEnvironmentData []GcpProjectEnvironmentData `tfschema:"gcp_project_environment_data"`
	Offerings                 []Offering                  `tfschema:"offering"`
	Tags                      map[string]string           `tfschema:"tags"`
}

type AwsEnvironmentData struct {
	OrganizationalDataMaster                  []AwsOrganizationalDataMaster `tfschema:"organizational_data_master"`
	OrganizationalDataMemberParentHierarchyId string                        `tfschema:"organizational_data_member_parent_hierarchy_id"`
	Regions                                   []string                      `tfschema:"regions"`
	ScanInterval                              int                           `tfschema:"scan_interval"`
}

type AwsOrganizationalDataMaster struct {
	StacksetName       string   `tfschema:"stackset_name"`
	ExcludedAccountIds []string `tfschema:"excluded_account_ids"`
}

type GcpProjectEnvironmentData struct {
	OrganizationalDataMaster []GcpOrganizationalDataMaster `tfschema:"organizational_data_master"`
	OrganizationalDataMember []GcpOrganizationalDataMember `tfschema:"organizational_data_member"`
	ProjectDetails           []GcpProjectDetails           `tfschema:"project_details"`
	ScanInterval             int                           `tfschema:"scan_interval"`
}

type GcpOrganizationalDataMaster struct {
	ExcludedProjectNumbers     []string `tfschema:"excluded_project_numbers"`
	ServiceAccountEmailAddress string   `tfschema:"service_account_email_address"`
	WorkloadIdentityProviderId string   `tfschema:"workload_identity_provider_id"`
}

type GcpOrganizationalDataMember struct {
	ManagementProjectNumber string `tfschema:"management_project_number"`
	ParentHierarchyId       string `tfschema:"parent_hierarchy_id"`
}

type GcpProjectDetails struct {
	ProjectId     string `tfschema:"project_id"`
	ProjectNumber string `tfschema:"project_number"`
}

type Offering struct {
	Type           string           `tfschema:"type"`
	CspmMonitorAws []CspmMonitorAws `tfschema:"cspm_monitor_aws"`
	CspmMonitorGcp []CspmMonitorGcp `tfschema:"cspm_monitor_gcp"`
}

type CspmMonitorAws struct {
	NativeCloudConnectionCloudRoleArn string `tfschema:"native_cloud_connection_cloud_role_arn"`
}

type CspmMonitorGcp struct {
	ServiceAccountEmailAddress string `tfschema:"service_account_email_address"`
	WorkloadIdentityProviderId string `tfschema:"workload_identity_provider_id"`
}

var _ sdk.Resource = SecurityCenterSecurityConnectorResource{}
var _ sdk.ResourceWithUpdate = SecurityCenterSecurityConnectorResource{}

type SecurityCenterSecurityConnectorResource struct{}

func (r SecurityCenterSecurityConnectorResource) ModelObject() interface{} {
	return &SecurityCenterSecurityConnectorModel{}
}

func (r SecurityCenterSecurityConnectorResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return securityconnectors.ValidateSecurityConnectorID
}

func (r SecurityCenterSecurityConnectorResource) ResourceType() string {
	return "azurerm_security_center_security_connector"
}

func (r SecurityCenterSecurityConnectorResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"environment_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(securityconnectors.CloudNameAWS),
				string(securityconnectors.CloudNameAzureDevOps),
				string(securityconnectors.CloudNameGCP),
				string(securityconnectors.CloudNameGitLab),
				string(securityconnectors.CloudNameGithub),
			}, false),
		},

		"hierarchy_identifier": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"aws_environment_data": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"organizational_data_master": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"stackset_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"excluded_account_ids": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									ForceNew: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
						ConflictsWith: []string{"aws_environment_data.0.organizational_data_member_parent_hierarchy_id"},
					},

					"organizational_data_member_parent_hierarchy_id": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ForceNew:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						ConflictsWith: []string{"aws_environment_data.0.organizational_data_master"},
					},

					"regions": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"scan_interval": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      4,
						ValidateFunc: validation.IntBetween(1, 24),
					},
				},
			},
			ConflictsWith: []string{"gcp_project_environment_data"},
		},

		"gcp_project_environment_data": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"organizational_data_master": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"service_account_email_address": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"workload_identity_provider_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"excluded_project_numbers": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									ForceNew: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
						ConflictsWith: []string{"gcp_project_environment_data.0.organizational_data_member"},
					},

					"organizational_data_member": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"parent_hierarchy_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"management_project_number": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
						ConflictsWith: []string{"gcp_project_environment_data.0.organizational_data_master"},
					},

					"project_details": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"project_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"project_number": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"scan_interval": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      4,
						ValidateFunc: validation.IntBetween(1, 24),
					},
				},
			},
			ConflictsWith: []string{"aws_environment_data"},
		},

		"offering": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeInt,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(securityconnectors.OfferingTypeCspmMonitorAws),
							string(securityconnectors.OfferingTypeCspmMonitorAzureDevOps),
							string(securityconnectors.OfferingTypeCspmMonitorGcp),
							string(securityconnectors.OfferingTypeCspmMonitorGitLab),
							string(securityconnectors.OfferingTypeCspmMonitorGithub),
						}, false),
					},

					"cspm_monitor_aws": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"native_cloud_connection_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"cspm_monitor_gcp": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"service_account_email_address": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"workload_identity_provider_id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r SecurityCenterSecurityConnectorResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r SecurityCenterSecurityConnectorResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.SecurityCenter.SecurityConnectorsClient

			var model SecurityCenterSecurityConnectorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := securityconnectors.NewSecurityConnectorID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := securityconnectors.SecurityConnector{
				Location: pointer.To(location.Normalize(model.Location)),
				Properties: &securityconnectors.SecurityConnectorProperties{
					EnvironmentName:     pointer.To(securityconnectors.CloudName(model.EnvironmentName)),
					HierarchyIdentifier: pointer.To(model.HierarchyIdentifier),
					Offerings:           expandOfferings(model.Offerings),
				},
				Tags: pointer.To(model.Tags),
			}

			if model.EnvironmentName == string(securityconnectors.CloudNameAzureDevOps) {
				if model.AwsEnvironmentData != nil {
					return fmt.Errorf("`aws_environment_data` only can be set when `environment_name` is `AWS`")
				}

				if model.GcpProjectEnvironmentData != nil {
					return fmt.Errorf("`gcp_project_environment_data` only can be set when `environment_name` is `GCP`")
				}

				parameters.Properties.EnvironmentData = securityconnectors.AzureDevOpsScopeEnvironmentData{}
			} else if model.EnvironmentName == string(securityconnectors.CloudNameGithub) {
				if model.AwsEnvironmentData != nil {
					return fmt.Errorf("`aws_environment_data` only can be set when `environment_name` is `AWS`")
				}

				if model.GcpProjectEnvironmentData != nil {
					return fmt.Errorf("`gcp_project_environment_data` only can be set when `environment_name` is `GCP`")
				}

				parameters.Properties.EnvironmentData = securityconnectors.GithubScopeEnvironmentData{}
			} else if model.EnvironmentName == string(securityconnectors.CloudNameGitLab) {
				if model.AwsEnvironmentData != nil {
					return fmt.Errorf("`aws_environment_data` only can be set when `environment_name` is `AWS`")
				}

				if model.GcpProjectEnvironmentData != nil {
					return fmt.Errorf("`gcp_project_environment_data` only can be set when `environment_name` is `GCP`")
				}

				parameters.Properties.EnvironmentData = securityconnectors.GitlabScopeEnvironmentData{}
			} else if model.EnvironmentName == string(securityconnectors.CloudNameAWS) {
				if model.GcpProjectEnvironmentData != nil {
					return fmt.Errorf("`gcp_project_environment_data` only can be set when `environment_name` is `GCP`")
				}

				parameters.Properties.EnvironmentData = expandAwsEnvironmentData(model.AwsEnvironmentData)
			} else if model.EnvironmentName == string(securityconnectors.CloudNameGCP) {
				if model.AwsEnvironmentData != nil {
					return fmt.Errorf("`aws_environment_data` only can be set when `environment_name` is `AWS`")
				}

				parameters.Properties.EnvironmentData = expandGcpProjectEnvironmentData(model.GcpProjectEnvironmentData)
			}

			if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r SecurityCenterSecurityConnectorResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SecurityCenter.SecurityConnectorsClient

			id, err := securityconnectors.ParseSecurityConnectorID(metadata.ResourceData.Id())
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

			state := SecurityCenterSecurityConnectorModel{}
			if model := resp.Model; model != nil {
				state.Name = id.SecurityConnectorName
				state.ResourceGroupName = id.ResourceGroupName
				state.Location = location.Normalize(pointer.From(model.Location))
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.EnvironmentName = string(pointer.From(props.EnvironmentName))
					state.HierarchyIdentifier = pointer.From(props.HierarchyIdentifier)
					state.Offerings = flattenOfferings(props.Offerings)

					if v, ok := props.EnvironmentData.(securityconnectors.AwsEnvironmentData); ok {
						state.AwsEnvironmentData = flattenAwsEnvironmentData(v)
					} else if v, ok := props.EnvironmentData.(securityconnectors.GcpProjectEnvironmentData); ok {
						state.GcpProjectEnvironmentData = flattenGcpProjectEnvironmentData(v)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r SecurityCenterSecurityConnectorResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SecurityCenter.SecurityConnectorsClient

			id, err := securityconnectors.ParseSecurityConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SecurityCenterSecurityConnectorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := securityconnectors.SecurityConnector{
				Properties: &securityconnectors.SecurityConnectorProperties{},
			}

			if metadata.ResourceData.HasChange("aws_environment_data") {
				if model.EnvironmentName != string(securityconnectors.CloudNameAWS) {
					return fmt.Errorf("`aws_environment_data` only can be set when `environment_name` is `AWS`")
				}

				parameters.Properties.EnvironmentData = expandAwsEnvironmentData(model.AwsEnvironmentData)
			}

			if metadata.ResourceData.HasChange("gcp_project_environment_data") {
				if model.EnvironmentName != string(securityconnectors.CloudNameGCP) {
					return fmt.Errorf("`gcp_project_environment_data` only can be set when `environment_name` is `GCP`")
				}

				parameters.Properties.EnvironmentData = expandGcpProjectEnvironmentData(model.GcpProjectEnvironmentData)
			}

			if metadata.ResourceData.HasChange("offering") {
				parameters.Properties.Offerings = expandOfferings(model.Offerings)
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

func (r SecurityCenterSecurityConnectorResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SecurityCenter.SecurityConnectorsClient

			id, err := securityconnectors.ParseSecurityConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandAwsEnvironmentData(input []AwsEnvironmentData) *securityconnectors.AwsEnvironmentData {
	if len(input) == 0 {
		return &securityconnectors.AwsEnvironmentData{}
	}

	awsEnvironmentData := input[0]

	result := &securityconnectors.AwsEnvironmentData{
		ScanInterval: pointer.To(int64(awsEnvironmentData.ScanInterval)),
		Regions:      pointer.To(awsEnvironmentData.Regions),
	}

	if v := awsEnvironmentData.OrganizationalDataMaster; v != nil {
		result.OrganizationalData = expandAwsOrganizationalDataMaster(v)
	} else if v := awsEnvironmentData.OrganizationalDataMemberParentHierarchyId; v != "" {
		result.OrganizationalData = expandAwsOrganizationalDataMember(v)
	}

	return result
}

func expandAwsOrganizationalDataMaster(input []AwsOrganizationalDataMaster) *securityconnectors.AwsOrganizationalDataMaster {
	awsOrganizationalDataMaster := input[0]

	result := &securityconnectors.AwsOrganizationalDataMaster{
		ExcludedAccountIds: pointer.To(awsOrganizationalDataMaster.ExcludedAccountIds),
		StacksetName:       pointer.To(awsOrganizationalDataMaster.StacksetName),
	}

	return result
}

func expandAwsOrganizationalDataMember(input string) *securityconnectors.AwsOrganizationalDataMember {
	result := &securityconnectors.AwsOrganizationalDataMember{
		ParentHierarchyId: pointer.To(input),
	}

	return result
}

func expandGcpProjectEnvironmentData(input []GcpProjectEnvironmentData) *securityconnectors.GcpProjectEnvironmentData {
	if len(input) == 0 {
		return &securityconnectors.GcpProjectEnvironmentData{}
	}

	gcpProjectEnvironmentData := input[0]

	result := &securityconnectors.GcpProjectEnvironmentData{
		ProjectDetails: expandGcpProjectDetails(gcpProjectEnvironmentData.ProjectDetails),
		ScanInterval:   pointer.To(int64(gcpProjectEnvironmentData.ScanInterval)),
	}

	if v := gcpProjectEnvironmentData.OrganizationalDataMaster; v != nil {
		result.OrganizationalData = expandGcpProjectOrganizationalDataMaster(v)
	} else if v := gcpProjectEnvironmentData.OrganizationalDataMember; v != nil {
		result.OrganizationalData = expandGcpProjectOrganizationalDataMember(v)
	}

	return result
}

func expandGcpProjectOrganizationalDataMaster(input []GcpOrganizationalDataMaster) *securityconnectors.GcpOrganizationalDataOrganization {
	gcpProjectOrganizationalDataMaster := input[0]

	result := &securityconnectors.GcpOrganizationalDataOrganization{
		ExcludedProjectNumbers:     pointer.To(gcpProjectOrganizationalDataMaster.ExcludedProjectNumbers),
		ServiceAccountEmailAddress: pointer.To(gcpProjectOrganizationalDataMaster.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.To(gcpProjectOrganizationalDataMaster.WorkloadIdentityProviderId),
	}

	return result
}

func expandGcpProjectOrganizationalDataMember(input []GcpOrganizationalDataMember) *securityconnectors.GcpOrganizationalDataMember {
	gcpProjectOrganizationalDataMember := input[0]

	result := &securityconnectors.GcpOrganizationalDataMember{
		ManagementProjectNumber: pointer.To(gcpProjectOrganizationalDataMember.ManagementProjectNumber),
		ParentHierarchyId:       pointer.To(gcpProjectOrganizationalDataMember.ParentHierarchyId),
	}

	return result
}

func expandGcpProjectDetails(input []GcpProjectDetails) *securityconnectors.GcpProjectDetails {
	if len(input) == 0 {
		return nil
	}

	projectDetails := input[0]

	result := &securityconnectors.GcpProjectDetails{
		ProjectId:     pointer.To(projectDetails.ProjectId),
		ProjectNumber: pointer.To(projectDetails.ProjectNumber),
	}

	return result
}

func expandOfferings(input []Offering) *[]securityconnectors.CloudOffering {
	result := make([]securityconnectors.CloudOffering, 0)
	if len(input) == 0 {
		return &result
	}

	for _, item := range input {
		if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeCspmMonitorAzureDevOps) {
			result = append(result, securityconnectors.CspmMonitorAzureDevOpsOffering{})
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeCspmMonitorGithub) {
			result = append(result, securityconnectors.CspmMonitorGithubOffering{})
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeCspmMonitorGitLab) {
			result = append(result, securityconnectors.CspmMonitorGitLabOffering{})
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeCspmMonitorAws) {
			cspmMonitorAwsOffering := securityconnectors.CspmMonitorAwsOffering{}

			if v := item.CspmMonitorAws; len(v) != 0 {
				cspmMonitorAws := v[0]

				if cloudRoleArn := cspmMonitorAws.NativeCloudConnectionCloudRoleArn; cloudRoleArn != "" {
					cspmMonitorAwsOffering.NativeCloudConnection = &securityconnectors.CspmMonitorAwsOfferingNativeCloudConnection{
						CloudRoleArn: pointer.To(cloudRoleArn),
					}
				}
			}

			result = append(result, cspmMonitorAwsOffering)
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeCspmMonitorGcp) {
			cspmMonitorGcpOffering := securityconnectors.CspmMonitorGcpOffering{}

			if v := item.CspmMonitorGcp; len(v) != 0 {
				cspmMonitorGcp := v[0]
				cspmMonitorGcpOffering.NativeCloudConnection = &securityconnectors.CspmMonitorGcpOfferingNativeCloudConnection{}

				if serviceAccountEmailAddress := cspmMonitorGcp.ServiceAccountEmailAddress; serviceAccountEmailAddress != "" {
					cspmMonitorGcpOffering.NativeCloudConnection.ServiceAccountEmailAddress = pointer.To(serviceAccountEmailAddress)
				}

				if workloadIdentityProviderId := cspmMonitorGcp.WorkloadIdentityProviderId; workloadIdentityProviderId != "" {
					cspmMonitorGcpOffering.NativeCloudConnection.WorkloadIdentityProviderId = pointer.To(workloadIdentityProviderId)
				}
			}

			result = append(result, cspmMonitorGcpOffering)
		}
	}

	return &result
}

func flattenAwsEnvironmentData(input securityconnectors.AwsEnvironmentData) []AwsEnvironmentData {
	result := make([]AwsEnvironmentData, 0)

	awsEnvironmentData := AwsEnvironmentData{
		Regions:      pointer.From(input.Regions),
		ScanInterval: int(pointer.From(input.ScanInterval)),
	}

	if v, ok := input.OrganizationalData.(securityconnectors.AwsOrganizationalDataMaster); ok {
		awsEnvironmentData.OrganizationalDataMaster = flattenAwsOrganizationalDataMaster(v)
	} else if v, ok := input.OrganizationalData.(securityconnectors.AwsOrganizationalDataMember); ok {
		awsEnvironmentData.OrganizationalDataMemberParentHierarchyId = flattenAwsOrganizationalDataMember(v)
	}

	return append(result, awsEnvironmentData)
}

func flattenAwsOrganizationalDataMaster(input securityconnectors.AwsOrganizationalDataMaster) []AwsOrganizationalDataMaster {
	result := make([]AwsOrganizationalDataMaster, 0)

	awsOrganizationalDataMaster := AwsOrganizationalDataMaster{
		ExcludedAccountIds: pointer.From(input.ExcludedAccountIds),
		StacksetName:       pointer.From(input.StacksetName),
	}

	return append(result, awsOrganizationalDataMaster)
}

func flattenAwsOrganizationalDataMember(input securityconnectors.AwsOrganizationalDataMember) string {
	organizationalDataMemberParentHierarchyId := pointer.From(input.ParentHierarchyId)

	return organizationalDataMemberParentHierarchyId
}

func flattenGcpProjectEnvironmentData(input securityconnectors.GcpProjectEnvironmentData) []GcpProjectEnvironmentData {
	result := make([]GcpProjectEnvironmentData, 0)

	gcpProjectEnvironmentData := GcpProjectEnvironmentData{
		ScanInterval:   int(pointer.From(input.ScanInterval)),
		ProjectDetails: flattenGcpProjectDetails(input.ProjectDetails),
	}

	if v, ok := input.OrganizationalData.(securityconnectors.GcpOrganizationalDataOrganization); ok {
		gcpProjectEnvironmentData.OrganizationalDataMaster = flattenGcpProjectOrganizationalDataMaster(v)
	} else if v, ok := input.OrganizationalData.(securityconnectors.GcpOrganizationalDataMember); ok {
		gcpProjectEnvironmentData.OrganizationalDataMember = flattenGcpProjectOrganizationalDataMember(v)
	}

	return append(result, gcpProjectEnvironmentData)
}

func flattenGcpProjectOrganizationalDataMaster(input securityconnectors.GcpOrganizationalDataOrganization) []GcpOrganizationalDataMaster {
	result := make([]GcpOrganizationalDataMaster, 0)

	gcpOrganizationalDataMaster := GcpOrganizationalDataMaster{
		ExcludedProjectNumbers:     pointer.From(input.ExcludedProjectNumbers),
		ServiceAccountEmailAddress: pointer.From(input.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.From(input.WorkloadIdentityProviderId),
	}

	return append(result, gcpOrganizationalDataMaster)
}

func flattenGcpProjectOrganizationalDataMember(input securityconnectors.GcpOrganizationalDataMember) []GcpOrganizationalDataMember {
	result := make([]GcpOrganizationalDataMember, 0)

	gcpOrganizationalDataMember := GcpOrganizationalDataMember{
		ManagementProjectNumber: pointer.From(input.ManagementProjectNumber),
		ParentHierarchyId:       pointer.From(input.ParentHierarchyId),
	}

	return append(result, gcpOrganizationalDataMember)
}

func flattenGcpProjectDetails(input *securityconnectors.GcpProjectDetails) []GcpProjectDetails {
	result := make([]GcpProjectDetails, 0)
	if input == nil {
		return result
	}

	gcpProjectDetails := GcpProjectDetails{
		ProjectId:     pointer.From(input.ProjectId),
		ProjectNumber: pointer.From(input.ProjectNumber),
	}

	return append(result, gcpProjectDetails)
}

func flattenOfferings(input *[]securityconnectors.CloudOffering) []Offering {
	result := make([]Offering, 0)
	if input == nil {
		return result
	}

	for _, item := range *input {
		if _, ok := item.(securityconnectors.CspmMonitorAzureDevOpsOffering); ok {
			result = append(result, Offering{
				Type: string(securityconnectors.OfferingTypeCspmMonitorAzureDevOps),
			})
		} else if _, ok := item.(securityconnectors.CspmMonitorGithubOffering); ok {
			result = append(result, Offering{
				Type: string(securityconnectors.OfferingTypeCspmMonitorGithub),
			})
		} else if _, ok := item.(securityconnectors.CspmMonitorGitLabOffering); ok {
			result = append(result, Offering{
				Type: string(securityconnectors.OfferingTypeCspmMonitorGitLab),
			})
		} else if v, ok := item.(securityconnectors.CspmMonitorAwsOffering); ok {
			cspmMonitorAwsOffering := Offering{
				Type: string(securityconnectors.OfferingTypeCspmMonitorAws),
			}

			if nativeCloudConnection := v.NativeCloudConnection; nativeCloudConnection != nil {
				cspmMonitorAwsOffering.CspmMonitorAws = []CspmMonitorAws{
					{
						NativeCloudConnectionCloudRoleArn: pointer.From(nativeCloudConnection.CloudRoleArn),
					},
				}
			}

			result = append(result, cspmMonitorAwsOffering)
		} else if v, ok := item.(securityconnectors.CspmMonitorGcpOffering); ok {
			cspmMonitorGcpOffering := Offering{
				Type: string(securityconnectors.OfferingTypeCspmMonitorGcp),
			}

			if nativeCloudConnection := v.NativeCloudConnection; nativeCloudConnection != nil {
				cspmMonitorGcpOffering.CspmMonitorGcp = []CspmMonitorGcp{
					{
						ServiceAccountEmailAddress: pointer.From(nativeCloudConnection.ServiceAccountEmailAddress),
						WorkloadIdentityProviderId: pointer.From(nativeCloudConnection.WorkloadIdentityProviderId),
					},
				}
			}

			result = append(result, cspmMonitorGcpOffering)
		}
	}

	return result
}
