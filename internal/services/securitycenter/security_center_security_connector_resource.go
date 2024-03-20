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
	Name                string            `tfschema:"name"`
	ResourceGroupName   string            `tfschema:"resource_group_name"`
	Location            string            `tfschema:"location"`
	EnvironmentName     string            `tfschema:"environment_name"`
	HierarchyIdentifier string            `tfschema:"hierarchy_identifier"`
	EnvironmentData     []EnvironmentData `tfschema:"environment_data"`
	Offerings           []Offering        `tfschema:"offering"`
	Tags                map[string]string `tfschema:"tags"`
}

type EnvironmentData struct {
	EnvironmentType string       `tfschema:"environment_type"`
	AwsAccount      []AwsAccount `tfschema:"aws_account"`
	GcpProject      []GcpProject `tfschema:"gcp_project"`
}

type AwsAccount struct {
	OrganizationalDataMaster                  []AwsOrganizationalDataMaster `tfschema:"organizational_data_master"`
	OrganizationalDataMemberParentHierarchyId string                        `tfschema:"organizational_data_member_parent_hierarchy_id"`
	Regions                                   []string                      `tfschema:"regions"`
	ScanInterval                              int                           `tfschema:"scan_interval"`
}

type AwsOrganizationalDataMaster struct {
	StacksetName       string   `tfschema:"stackset_name"`
	ExcludedAccountIds []string `tfschema:"excluded_account_ids"`
}

type GcpProject struct {
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

		"environment_data": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"environment_type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(securityconnectors.PossibleValuesForEnvironmentType(), false),
					},

					"aws_account": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
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
									ConflictsWith: []string{"environment_data.0.aws_account.0.organizational_data_member_parent_hierarchy_id"},
								},

								"organizational_data_member_parent_hierarchy_id": {
									Type:          pluginsdk.TypeString,
									Optional:      true,
									ForceNew:      true,
									ValidateFunc:  validation.StringIsNotEmpty,
									ConflictsWith: []string{"environment_data.0.aws_account.0.organizational_data_master"},
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
					},

					"gcp_project": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
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
									ConflictsWith: []string{"environment_data.0.gcp_project.0.organizational_data_member"},
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
									ConflictsWith: []string{"environment_data.0.gcp_project.0.organizational_data_master"},
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
					},
				},
			},
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
					EnvironmentData:     expandEnvironmentData(model.EnvironmentData),
				},
				Tags: pointer.To(model.Tags),
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
					state.EnvironmentData = flattenEnvironmentData(props.EnvironmentData)
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

			if metadata.ResourceData.HasChange("environment_data") {
				parameters.Properties.EnvironmentData = expandEnvironmentData(model.EnvironmentData)
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

func expandEnvironmentData(input []EnvironmentData) *securityconnectors.EnvironmentData {
	if len(input) == 0 {
		return nil
	}

	environmentData := input[0]
	var result securityconnectors.EnvironmentData

	if environmentData.EnvironmentType == string(securityconnectors.EnvironmentTypeAzureDevOpsScope) {
		result = securityconnectors.AzureDevOpsScopeEnvironmentData{}
	} else if environmentData.EnvironmentType == string(securityconnectors.EnvironmentTypeGithubScope) {
		result = securityconnectors.GithubScopeEnvironmentData{}
	} else if environmentData.EnvironmentType == string(securityconnectors.EnvironmentTypeGitlabScope) {
		result = securityconnectors.GitlabScopeEnvironmentData{}
	} else if environmentData.EnvironmentType == string(securityconnectors.EnvironmentTypeAwsAccount) {
		result = expandAwsAccount(environmentData.AwsAccount)
	} else if environmentData.EnvironmentType == string(securityconnectors.EnvironmentTypeGcpProject) {
		result = expandGcpProject(environmentData.GcpProject)
	}

	return &result
}

func expandAwsAccount(input []AwsAccount) *securityconnectors.AwsEnvironmentData {
	if len(input) == 0 {
		return &securityconnectors.AwsEnvironmentData{}
	}

	awsAccount := input[0]

	result := &securityconnectors.AwsEnvironmentData{
		ScanInterval: pointer.To(int64(awsAccount.ScanInterval)),
		Regions:      pointer.To(awsAccount.Regions),
	}

	if v := awsAccount.OrganizationalDataMaster; v != nil {
		result.OrganizationalData = expandAwsOrganizationalDataMaster(v)
	} else if v := awsAccount.OrganizationalDataMemberParentHierarchyId; v != "" {
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

func expandGcpProject(input []GcpProject) *securityconnectors.GcpProjectEnvironmentData {
	if len(input) == 0 {
		return &securityconnectors.GcpProjectEnvironmentData{}
	}

	gcpProject := input[0]

	result := &securityconnectors.GcpProjectEnvironmentData{
		ProjectDetails: expandGcpProjectDetails(gcpProject.ProjectDetails),
		ScanInterval:   pointer.To(int64(gcpProject.ScanInterval)),
	}

	if v := gcpProject.OrganizationalDataMaster; v != nil {
		result.OrganizationalData = expandGcpProjectOrganizationalDataMaster(v)
	} else if v := gcpProject.OrganizationalDataMember; v != nil {
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

func flattenEnvironmentData(input securityconnectors.EnvironmentData) []EnvironmentData {
	result := make([]EnvironmentData, 0)

	if _, ok := input.(securityconnectors.AzureDevOpsScopeEnvironmentData); ok {
		result = append(result, EnvironmentData{
			EnvironmentType: string(securityconnectors.EnvironmentTypeAzureDevOpsScope),
		})
	} else if _, ok := input.(securityconnectors.GithubScopeEnvironmentData); ok {
		result = append(result, EnvironmentData{
			EnvironmentType: string(securityconnectors.EnvironmentTypeGithubScope),
		})
	} else if _, ok := input.(securityconnectors.GitlabScopeEnvironmentData); ok {
		result = append(result, EnvironmentData{
			EnvironmentType: string(securityconnectors.EnvironmentTypeGitlabScope),
		})
	} else if _, ok := input.(securityconnectors.AwsEnvironmentData); ok {
		awsEnvironmentData := input.(securityconnectors.AwsEnvironmentData)

		result = append(result, EnvironmentData{
			EnvironmentType: string(securityconnectors.EnvironmentTypeAwsAccount),
			AwsAccount:      flattenAwsAccount(awsEnvironmentData),
		})
	} else if _, ok := input.(securityconnectors.GcpProjectEnvironmentData); ok {
		gcpEnvironmentData := input.(securityconnectors.GcpProjectEnvironmentData)

		result = append(result, EnvironmentData{
			EnvironmentType: string(securityconnectors.EnvironmentTypeGcpProject),
			GcpProject:      flattenGcpProject(gcpEnvironmentData),
		})
	}

	return result
}

func flattenAwsAccount(input securityconnectors.AwsEnvironmentData) []AwsAccount {
	result := make([]AwsAccount, 0)

	awsAccount := AwsAccount{
		Regions:      pointer.From(input.Regions),
		ScanInterval: int(pointer.From(input.ScanInterval)),
	}

	if v, ok := input.OrganizationalData.(securityconnectors.AwsOrganizationalDataMaster); ok {
		awsAccount.OrganizationalDataMaster = flattenAwsOrganizationalDataMaster(v)
	} else if v, ok := input.OrganizationalData.(securityconnectors.AwsOrganizationalDataMember); ok {
		awsAccount.OrganizationalDataMemberParentHierarchyId = flattenAwsOrganizationalDataMember(v)
	}

	return append(result, awsAccount)
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

func flattenGcpProject(input securityconnectors.GcpProjectEnvironmentData) []GcpProject {
	result := make([]GcpProject, 0)

	gcpProject := GcpProject{
		ScanInterval:   int(pointer.From(input.ScanInterval)),
		ProjectDetails: flattenGcpProjectDetails(input.ProjectDetails),
	}

	if v, ok := input.OrganizationalData.(securityconnectors.GcpOrganizationalDataOrganization); ok {
		gcpProject.OrganizationalDataMaster = flattenGcpProjectOrganizationalDataMaster(v)
	} else if v, ok := input.OrganizationalData.(securityconnectors.GcpOrganizationalDataMember); ok {
		gcpProject.OrganizationalDataMember = flattenGcpProjectOrganizationalDataMember(v)
	}

	return append(result, gcpProject)
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
