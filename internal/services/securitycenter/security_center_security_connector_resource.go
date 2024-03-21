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
	Type                     string                     `tfschema:"type"`
	CspmMonitorAws           []CspmMonitorAws           `tfschema:"cspm_monitor_aws"`
	CspmMonitorGcp           []CspmMonitorGcp           `tfschema:"cspm_monitor_gcp"`
	DefenderForDatabasesAws  []DefenderForDatabasesAws  `tfschema:"defender_for_databases_aws"`
	DefenderForDatabasesGcp  []DefenderForDatabasesGcp  `tfschema:"defender_for_databases_gcp"`
	DefenderForContainersAws []DefenderForContainersAws `tfschema:"defender_for_containers_aws"`
	DefenderForContainersGcp []DefenderForContainersGcp `tfschema:"defender_for_containers_gcp"`
}

type CspmMonitorAws struct {
	NativeCloudConnectionCloudRoleArn string `tfschema:"native_cloud_connection_cloud_role_arn"`
}

type CspmMonitorGcp struct {
	ServiceAccountEmailAddress string `tfschema:"service_account_email_address"`
	WorkloadIdentityProviderId string `tfschema:"workload_identity_provider_id"`
}

type DefenderForDatabasesAws struct {
	ArcAutoProvisioning []DatabasesAwsArcAutoProvisioning `tfschema:"arc_auto_provisioning"`
	DatabasesDspm       []DatabasesAwsDspm                `tfschema:"databases_dspm"`
	Rds                 []DatabasesAwsRds                 `tfschema:"rds"`
}

type DatabasesAwsArcAutoProvisioning struct {
	CloudRoleArn  string                      `tfschema:"cloud_role_arn"`
	Configuration []DatabasesAwsConfiguration `tfschema:"configuration"`
	Enabled       bool                        `tfschema:"enabled"`
}

type DatabasesAwsConfiguration struct {
	PrivateLinkScope string `tfschema:"private_link_scope"`
	Proxy            string `tfschema:"proxy"`
}

type DatabasesAwsDspm struct {
	CloudRoleArn string `tfschema:"cloud_role_arn"`
	Enabled      bool   `tfschema:"enabled"`
}

type DatabasesAwsRds struct {
	CloudRoleArn string `tfschema:"cloud_role_arn"`
	Enabled      bool   `tfschema:"enabled"`
}

type DefenderForDatabasesGcp struct {
	ArcAutoProvisioning                     []DatabasesGcpArcAutoProvisioning                     `tfschema:"arc_auto_provisioning"`
	DefenderForDatabasesArcAutoProvisioning []DatabasesGcpDefenderForDatabasesArcAutoProvisioning `tfschema:"defender_for_databases_arc_auto_provisioning"`
}

type DatabasesGcpArcAutoProvisioning struct {
	Configuration []DatabasesGcpConfiguration `tfschema:"configuration"`
	Enabled       bool                        `tfschema:"enabled"`
}

type DatabasesGcpConfiguration struct {
	PrivateLinkScope string `tfschema:"private_link_scope"`
	Proxy            string `tfschema:"proxy"`
}

type DatabasesGcpDefenderForDatabasesArcAutoProvisioning struct {
	ServiceAccountEmailAddress string `tfschema:"service_account_email_address"`
	WorkloadIdentityProviderId string `tfschema:"workload_identity_provider_id"`
}

type DefenderForContainersAws struct {
	AutoProvisioningEnabled                          bool                                              `tfschema:"auto_provisioning_enabled"`
	CloudWatchToKinesisCloudRoleArn                  string                                            `tfschema:"cloud_watch_to_kinesis_cloud_role_arn"`
	ContainerVulnerabilityAssessmentCloudRoleArn     string                                            `tfschema:"container_vulnerability_assessment_cloud_role_arn"`
	ContainerVulnerabilityAssessmentTaskCloudRoleArn string                                            `tfschema:"container_vulnerability_assessment_task_cloud_role_arn"`
	ContainerVulnerabilityAssessmentEnabled          bool                                              `tfschema:"container_vulnerability_assessment_enabled"`
	KinesisToS3CloudRoleArn                          string                                            `tfschema:"kinesis_to_s3_cloud_role_arn"`
	KubeAuditRetentionTime                           int                                               `tfschema:"kube_audit_retention_time"`
	KubernetesScubaReaderCloudRoleArn                string                                            `tfschema:"kubernetes_scuba_reader_cloud_role_arn"`
	KubernetesServiceCloudRoleArn                    string                                            `tfschema:"kubernetes_service_cloud_role_arn"`
	MdcContainersAgentlessDiscoveryK8s               []ContainersAwsMdcContainersAgentlessDiscoveryK8s `tfschema:"mdc_containers_agentless_discovery_k8s"`
	MdcContainersImageAssessment                     []ContainersAwsMdcContainersImageAssessment       `tfschema:"mdc_containers_image_assessment"`
	ScubaExternalId                                  string                                            `tfschema:"scuba_external_id"`
}

type ContainersAwsMdcContainersAgentlessDiscoveryK8s struct {
	CloudRoleArn string `tfschema:"cloud_role_arn"`
	Enabled      bool   `tfschema:"enabled"`
}

type ContainersAwsMdcContainersImageAssessment struct {
	CloudRoleArn string `tfschema:"cloud_role_arn"`
	Enabled      bool   `tfschema:"enabled"`
}

type DefenderForContainersGcp struct {
	AuditLogsAutoProvisioningFlagEnabled     bool                                              `tfschema:"audit_logs_auto_provisioning_flag_enabled"`
	DataPipelineNativeCloudConnection        []ContainersGcpDataPipelineNativeCloudConnection  `tfschema:"data_pipeline_native_cloud_connection"`
	DefenderAgentAutoProvisioningFlagEnabled bool                                              `tfschema:"defender_agent_auto_provisioning_flag_enabled"`
	MdcContainersAgentlessDiscoveryK8s       []ContainersGcpMdcContainersAgentlessDiscoveryK8s `tfschema:"mdc_containers_agentless_discovery_k8s"`
	MdcContainersImageAssessment             []ContainersGcpMdcContainersImageAssessment       `tfschema:"mdc_containers_image_assessment"`
	NativeCloudConnection                    []ContainersGcpNativeCloudConnection              `tfschema:"native_cloud_connection"`
	PolicyAgentAutoProvisioningFlagEnabled   bool                                              `tfschema:"policy_agent_auto_provisioning_flag_enabled"`
}

type ContainersGcpDataPipelineNativeCloudConnection struct {
	ServiceAccountEmailAddress string `tfschema:"service_account_email_address"`
	WorkloadIdentityProviderId string `tfschema:"workload_identity_provider_id"`
}

type ContainersGcpMdcContainersAgentlessDiscoveryK8s struct {
	ServiceAccountEmailAddress string `tfschema:"service_account_email_address"`
	WorkloadIdentityProviderId string `tfschema:"workload_identity_provider_id"`
	Enabled                    bool   `tfschema:"enabled"`
}

type ContainersGcpMdcContainersImageAssessment struct {
	ServiceAccountEmailAddress string `tfschema:"service_account_email_address"`
	WorkloadIdentityProviderId string `tfschema:"workload_identity_provider_id"`
	Enabled                    bool   `tfschema:"enabled"`
}

type ContainersGcpNativeCloudConnection struct {
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
							string(securityconnectors.OfferingTypeDefenderCspmAws),
							string(securityconnectors.OfferingTypeDefenderCspmGcp),
							string(securityconnectors.OfferingTypeDefenderForContainersAws),
							string(securityconnectors.OfferingTypeDefenderForContainersGcp),
							string(securityconnectors.OfferingTypeDefenderForDatabasesAws),
							string(securityconnectors.OfferingTypeDefenderForDatabasesGcp),
							string(securityconnectors.OfferingTypeDefenderForServersAws),
							string(securityconnectors.OfferingTypeDefenderForServersGcp),
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

					"defender_for_databases_aws": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"arc_auto_provisioning": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									ForceNew: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"cloud_role_arn": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"configuration": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												ForceNew: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"private_link_scope": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ForceNew:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"proxy": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ForceNew:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},
												},
											},

											"enabled": {
												Type:     pluginsdk.TypeBool,
												Required: true,
												ForceNew: true,
											},
										},
									},
								},

								"databases_dspm": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									ForceNew: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"cloud_role_arn": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"enabled": {
												Type:     pluginsdk.TypeBool,
												Required: true,
												ForceNew: true,
											},
										},
									},
								},

								"rds": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									ForceNew: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"cloud_role_arn": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"enabled": {
												Type:     pluginsdk.TypeBool,
												Required: true,
												ForceNew: true,
											},
										},
									},
								},
							},
						},
					},

					"defender_for_databases_gcp": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"arc_auto_provisioning": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									ForceNew: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"configuration": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												ForceNew: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"private_link_scope": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ForceNew:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"proxy": {
															Type:         pluginsdk.TypeString,
															Optional:     true,
															ForceNew:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},
													},
												},
											},

											"enabled": {
												Type:     pluginsdk.TypeBool,
												Required: true,
												ForceNew: true,
											},
										},
									},
								},

								"defender_for_databases_arc_auto_provisioning": {
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
										},
									},
								},
							},
						},
					},

					"defender_for_containers_aws": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"auto_provisioning_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									ForceNew: true,
								},

								"cloud_watch_to_kinesis_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"container_vulnerability_assessment_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"container_vulnerability_assessment_task_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"container_vulnerability_assessment_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									ForceNew: true,
								},

								"kinesis_to_s3_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"kube_audit_retention_time": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									ForceNew: true,
								},

								"kubernetes_scuba_reader_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"kubernetes_service_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"mdc_containers_agentless_discovery_k8s": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									ForceNew: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"cloud_role_arn": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"enabled": {
												Type:     pluginsdk.TypeBool,
												Required: true,
												ForceNew: true,
											},
										},
									},
								},

								"mdc_containers_image_assessment": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									ForceNew: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"cloud_role_arn": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ForceNew:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"enabled": {
												Type:     pluginsdk.TypeBool,
												Required: true,
												ForceNew: true,
											},
										},
									},
								},

								"scuba_external_id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"defender_for_containers_gcp": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"audit_logs_auto_provisioning_flag_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									ForceNew: true,
								},

								"data_pipeline_native_cloud_connection": {
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
										},
									},
								},

								"defender_agent_auto_provisioning_flag_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									ForceNew: true,
								},

								"mdc_containers_agentless_discovery_k8s": {
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

											"enabled": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
												ForceNew: true,
											},
										},
									},
								},

								"mdc_containers_image_assessment": {
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

											"enabled": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
												ForceNew: true,
											},
										},
									},
								},

								"native_cloud_connection": {
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
										},
									},
								},

								"policy_agent_auto_provisioning_flag_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									ForceNew: true,
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

				cspmMonitorAwsOffering.NativeCloudConnection = &securityconnectors.CspmMonitorAwsOfferingNativeCloudConnection{
					CloudRoleArn: pointer.To(cspmMonitorAws.NativeCloudConnectionCloudRoleArn),
				}
			}

			result = append(result, cspmMonitorAwsOffering)
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeCspmMonitorGcp) {
			cspmMonitorGcpOffering := securityconnectors.CspmMonitorGcpOffering{}

			if v := item.CspmMonitorGcp; len(v) != 0 {
				cspmMonitorGcp := v[0]

				cspmMonitorGcpOffering.NativeCloudConnection = &securityconnectors.CspmMonitorGcpOfferingNativeCloudConnection{
					ServiceAccountEmailAddress: pointer.To(cspmMonitorGcp.ServiceAccountEmailAddress),
					WorkloadIdentityProviderId: pointer.To(cspmMonitorGcp.WorkloadIdentityProviderId),
				}
			}

			result = append(result, cspmMonitorGcpOffering)
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeDefenderForDatabasesAws) {
			defenderFoDatabasesAwsOffering := securityconnectors.DefenderFoDatabasesAwsOffering{}

			if v := item.DefenderForDatabasesAws; len(v) != 0 {
				defenderForDatabasesAws := v[0]

				defenderFoDatabasesAwsOffering.ArcAutoProvisioning = expandDatabasesAwsArcAutoProvisioning(defenderForDatabasesAws.ArcAutoProvisioning)
				defenderFoDatabasesAwsOffering.DatabasesDspm = expandDatabasesAwsDspm(defenderForDatabasesAws.DatabasesDspm)
				defenderFoDatabasesAwsOffering.Rds = expandDatabasesAwsRds(defenderForDatabasesAws.Rds)
			}

			result = append(result, defenderFoDatabasesAwsOffering)
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeDefenderForDatabasesGcp) {
			defenderForDatabasesGcpOffering := securityconnectors.DefenderForDatabasesGcpOffering{}

			if v := item.DefenderForDatabasesGcp; len(v) != 0 {
				defenderForDatabasesGcp := v[0]

				defenderForDatabasesGcpOffering.ArcAutoProvisioning = expandDatabasesGcpArcAutoProvisioning(defenderForDatabasesGcp.ArcAutoProvisioning)
				defenderForDatabasesGcpOffering.DefenderForDatabasesArcAutoProvisioning = expandDatabasesGcpDefenderForDatabasesArcAutoProvisioning(defenderForDatabasesGcp.DefenderForDatabasesArcAutoProvisioning)
			}

			result = append(result, defenderForDatabasesGcpOffering)
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeDefenderForContainersAws) {
			defenderForContainersAwsOffering := securityconnectors.DefenderForContainersAwsOffering{}

			if v := item.DefenderForContainersAws; len(v) != 0 {
				defenderForContainersAws := v[0]

				defenderForContainersAwsOffering.AutoProvisioning = pointer.To(defenderForContainersAws.AutoProvisioningEnabled)
				defenderForContainersAwsOffering.EnableContainerVulnerabilityAssessment = pointer.To(defenderForContainersAws.ContainerVulnerabilityAssessmentEnabled)
				defenderForContainersAwsOffering.MdcContainersAgentlessDiscoveryK8s = expandContainersAwsMdcContainersAgentlessDiscoveryK8s(defenderForContainersAws.MdcContainersAgentlessDiscoveryK8s)
				defenderForContainersAwsOffering.MdcContainersImageAssessment = expandContainersAwsMdcContainersImageAssessment(defenderForContainersAws.MdcContainersImageAssessment)
				defenderForContainersAwsOffering.ScubaExternalId = pointer.To(defenderForContainersAws.ScubaExternalId)

				defenderForContainersAwsOffering.CloudWatchToKinesis = &securityconnectors.DefenderForContainersAwsOfferingCloudWatchToKinesis{
					CloudRoleArn: pointer.To(defenderForContainersAws.CloudWatchToKinesisCloudRoleArn),
				}

				defenderForContainersAwsOffering.ContainerVulnerabilityAssessment = &securityconnectors.DefenderForContainersAwsOfferingContainerVulnerabilityAssessment{
					CloudRoleArn: pointer.To(defenderForContainersAws.ContainerVulnerabilityAssessmentCloudRoleArn),
				}

				defenderForContainersAwsOffering.ContainerVulnerabilityAssessmentTask = &securityconnectors.DefenderForContainersAwsOfferingContainerVulnerabilityAssessmentTask{
					CloudRoleArn: pointer.To(defenderForContainersAws.ContainerVulnerabilityAssessmentTaskCloudRoleArn),
				}

				defenderForContainersAwsOffering.KinesisToS3 = &securityconnectors.DefenderForContainersAwsOfferingKinesisToS3{
					CloudRoleArn: pointer.To(defenderForContainersAws.KinesisToS3CloudRoleArn),
				}

				if v := defenderForContainersAws.KubeAuditRetentionTime; v != 0 {
					defenderForContainersAwsOffering.KubeAuditRetentionTime = pointer.To(int64(v))
				}

				defenderForContainersAwsOffering.KubernetesScubaReader = &securityconnectors.DefenderForContainersAwsOfferingKubernetesScubaReader{
					CloudRoleArn: pointer.To(defenderForContainersAws.KubernetesScubaReaderCloudRoleArn),
				}

				defenderForContainersAwsOffering.KubernetesService = &securityconnectors.DefenderForContainersAwsOfferingKubernetesService{
					CloudRoleArn: pointer.To(defenderForContainersAws.KubernetesServiceCloudRoleArn),
				}
			}

			result = append(result, defenderForContainersAwsOffering)
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeDefenderForContainersGcp) {
			defenderForContainersGcpOffering := securityconnectors.DefenderForContainersGcpOffering{}

			if v := item.DefenderForContainersGcp; len(v) != 0 {
				defenderForContainersGcp := v[0]

				defenderForContainersGcpOffering.AuditLogsAutoProvisioningFlag = pointer.To(defenderForContainersGcp.AuditLogsAutoProvisioningFlagEnabled)
				defenderForContainersGcpOffering.DataPipelineNativeCloudConnection = expandContainersGcpDataPipelineNativeCloudConnection(defenderForContainersGcp.DataPipelineNativeCloudConnection)
				defenderForContainersGcpOffering.DefenderAgentAutoProvisioningFlag = pointer.To(defenderForContainersGcp.DefenderAgentAutoProvisioningFlagEnabled)
				defenderForContainersGcpOffering.MdcContainersAgentlessDiscoveryK8s = expandContainersGcpMdcContainersAgentlessDiscoveryK8s(defenderForContainersGcp.MdcContainersAgentlessDiscoveryK8s)
				defenderForContainersGcpOffering.MdcContainersImageAssessment = expandContainersGcpMdcContainersImageAssessment(defenderForContainersGcp.MdcContainersImageAssessment)
				defenderForContainersGcpOffering.NativeCloudConnection = expandContainersGcpNativeCloudConnection(defenderForContainersGcp.NativeCloudConnection)
				defenderForContainersGcpOffering.PolicyAgentAutoProvisioningFlag = pointer.To(defenderForContainersGcp.PolicyAgentAutoProvisioningFlagEnabled)
			}

			result = append(result, defenderForContainersGcpOffering)
		}
	}

	return &result
}

func expandDatabasesAwsArcAutoProvisioning(input []DatabasesAwsArcAutoProvisioning) *securityconnectors.DefenderFoDatabasesAwsOfferingArcAutoProvisioning {
	if len(input) == 0 {
		return nil
	}

	awsArcAutoProvisioning := input[0]

	result := &securityconnectors.DefenderFoDatabasesAwsOfferingArcAutoProvisioning{
		CloudRoleArn:  pointer.To(awsArcAutoProvisioning.CloudRoleArn),
		Configuration: expandDatabasesAwsConfiguration(awsArcAutoProvisioning.Configuration),
		Enabled:       pointer.To(awsArcAutoProvisioning.Enabled),
	}

	return result
}

func expandDatabasesAwsConfiguration(input []DatabasesAwsConfiguration) *securityconnectors.DefenderFoDatabasesAwsOfferingArcAutoProvisioningConfiguration {
	if len(input) == 0 {
		return nil
	}

	awsConfiguration := input[0]

	result := &securityconnectors.DefenderFoDatabasesAwsOfferingArcAutoProvisioningConfiguration{
		PrivateLinkScope: pointer.To(awsConfiguration.PrivateLinkScope),
		Proxy:            pointer.To(awsConfiguration.Proxy),
	}

	return result
}

func expandDatabasesAwsDspm(input []DatabasesAwsDspm) *securityconnectors.DefenderFoDatabasesAwsOfferingDatabasesDspm {
	if len(input) == 0 {
		return nil
	}

	databasesDspm := input[0]

	result := &securityconnectors.DefenderFoDatabasesAwsOfferingDatabasesDspm{
		CloudRoleArn: pointer.To(databasesDspm.CloudRoleArn),
		Enabled:      pointer.To(databasesDspm.Enabled),
	}

	return result
}

func expandDatabasesAwsRds(input []DatabasesAwsRds) *securityconnectors.DefenderFoDatabasesAwsOfferingRds {
	if len(input) == 0 {
		return nil
	}

	rds := input[0]

	result := &securityconnectors.DefenderFoDatabasesAwsOfferingRds{
		CloudRoleArn: pointer.To(rds.CloudRoleArn),
		Enabled:      pointer.To(rds.Enabled),
	}

	return result
}

func expandDatabasesGcpArcAutoProvisioning(input []DatabasesGcpArcAutoProvisioning) *securityconnectors.DefenderForDatabasesGcpOfferingArcAutoProvisioning {
	if len(input) == 0 {
		return nil
	}

	gcpArcAutoProvisioning := input[0]

	result := &securityconnectors.DefenderForDatabasesGcpOfferingArcAutoProvisioning{
		Configuration: expandDatabasesGcpConfiguration(gcpArcAutoProvisioning.Configuration),
		Enabled:       pointer.To(gcpArcAutoProvisioning.Enabled),
	}

	return result
}

func expandDatabasesGcpConfiguration(input []DatabasesGcpConfiguration) *securityconnectors.DefenderForDatabasesGcpOfferingArcAutoProvisioningConfiguration {
	if len(input) == 0 {
		return nil
	}

	gcpConfiguration := input[0]

	result := &securityconnectors.DefenderForDatabasesGcpOfferingArcAutoProvisioningConfiguration{
		PrivateLinkScope: pointer.To(gcpConfiguration.PrivateLinkScope),
		Proxy:            pointer.To(gcpConfiguration.Proxy),
	}

	return result
}

func expandDatabasesGcpDefenderForDatabasesArcAutoProvisioning(input []DatabasesGcpDefenderForDatabasesArcAutoProvisioning) *securityconnectors.DefenderForDatabasesGcpOfferingDefenderForDatabasesArcAutoProvisioning {
	if len(input) == 0 {
		return nil
	}

	defenderForDatabasesArcAutoProvisioning := input[0]

	result := &securityconnectors.DefenderForDatabasesGcpOfferingDefenderForDatabasesArcAutoProvisioning{
		ServiceAccountEmailAddress: pointer.To(defenderForDatabasesArcAutoProvisioning.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.To(defenderForDatabasesArcAutoProvisioning.WorkloadIdentityProviderId),
	}

	return result
}

func expandContainersAwsMdcContainersAgentlessDiscoveryK8s(input []ContainersAwsMdcContainersAgentlessDiscoveryK8s) *securityconnectors.DefenderForContainersAwsOfferingMdcContainersAgentlessDiscoveryK8s {
	if len(input) == 0 {
		return nil
	}

	awsMdcContainersAgentlessDiscoveryK8s := input[0]

	result := &securityconnectors.DefenderForContainersAwsOfferingMdcContainersAgentlessDiscoveryK8s{
		CloudRoleArn: pointer.To(awsMdcContainersAgentlessDiscoveryK8s.CloudRoleArn),
		Enabled:      pointer.To(awsMdcContainersAgentlessDiscoveryK8s.Enabled),
	}

	return result
}

func expandContainersAwsMdcContainersImageAssessment(input []ContainersAwsMdcContainersImageAssessment) *securityconnectors.DefenderForContainersAwsOfferingMdcContainersImageAssessment {
	if len(input) == 0 {
		return nil
	}

	awsMdcContainersImageAssessment := input[0]

	result := &securityconnectors.DefenderForContainersAwsOfferingMdcContainersImageAssessment{
		CloudRoleArn: pointer.To(awsMdcContainersImageAssessment.CloudRoleArn),
		Enabled:      pointer.To(awsMdcContainersImageAssessment.Enabled),
	}

	return result
}

func expandContainersGcpDataPipelineNativeCloudConnection(input []ContainersGcpDataPipelineNativeCloudConnection) *securityconnectors.DefenderForContainersGcpOfferingDataPipelineNativeCloudConnection {
	if len(input) == 0 {
		return nil
	}

	dataPipelineNativeCloudConnection := input[0]

	result := &securityconnectors.DefenderForContainersGcpOfferingDataPipelineNativeCloudConnection{
		ServiceAccountEmailAddress: pointer.To(dataPipelineNativeCloudConnection.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.To(dataPipelineNativeCloudConnection.WorkloadIdentityProviderId),
	}

	return result
}

func expandContainersGcpMdcContainersAgentlessDiscoveryK8s(input []ContainersGcpMdcContainersAgentlessDiscoveryK8s) *securityconnectors.DefenderForContainersGcpOfferingMdcContainersAgentlessDiscoveryK8s {
	if len(input) == 0 {
		return nil
	}

	gcpMdcContainersAgentlessDiscoveryK8s := input[0]

	result := &securityconnectors.DefenderForContainersGcpOfferingMdcContainersAgentlessDiscoveryK8s{
		ServiceAccountEmailAddress: pointer.To(gcpMdcContainersAgentlessDiscoveryK8s.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.To(gcpMdcContainersAgentlessDiscoveryK8s.WorkloadIdentityProviderId),
		Enabled:                    pointer.To(gcpMdcContainersAgentlessDiscoveryK8s.Enabled),
	}

	return result
}

func expandContainersGcpMdcContainersImageAssessment(input []ContainersGcpMdcContainersImageAssessment) *securityconnectors.DefenderForContainersGcpOfferingMdcContainersImageAssessment {
	if len(input) == 0 {
		return nil
	}

	gcpMdcContainersImageAssessment := input[0]

	result := &securityconnectors.DefenderForContainersGcpOfferingMdcContainersImageAssessment{
		ServiceAccountEmailAddress: pointer.To(gcpMdcContainersImageAssessment.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.To(gcpMdcContainersImageAssessment.WorkloadIdentityProviderId),
		Enabled:                    pointer.To(gcpMdcContainersImageAssessment.Enabled),
	}

	return result
}

func expandContainersGcpNativeCloudConnection(input []ContainersGcpNativeCloudConnection) *securityconnectors.DefenderForContainersGcpOfferingNativeCloudConnection {
	if len(input) == 0 {
		return nil
	}

	containersGcpNativeCloudConnection := input[0]

	result := &securityconnectors.DefenderForContainersGcpOfferingNativeCloudConnection{
		ServiceAccountEmailAddress: pointer.To(containersGcpNativeCloudConnection.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.To(containersGcpNativeCloudConnection.WorkloadIdentityProviderId),
	}

	return result
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
		} else if v, ok := item.(securityconnectors.DefenderFoDatabasesAwsOffering); ok {
			defenderFoDatabasesAwsOffering := Offering{
				Type: string(securityconnectors.OfferingTypeDefenderForDatabasesAws),
				DefenderForDatabasesAws: []DefenderForDatabasesAws{
					{
						ArcAutoProvisioning: flattenDatabasesAwsArcAutoProvisioning(v.ArcAutoProvisioning),
						DatabasesDspm:       flattenDatabasesAwsDspm(v.DatabasesDspm),
						Rds:                 flattenDatabasesAwsRds(v.Rds),
					},
				},
			}

			result = append(result, defenderFoDatabasesAwsOffering)
		} else if v, ok := item.(securityconnectors.DefenderForDatabasesGcpOffering); ok {
			defenderForDatabasesGcpOffering := Offering{
				Type: string(securityconnectors.OfferingTypeDefenderForDatabasesGcp),
				DefenderForDatabasesGcp: []DefenderForDatabasesGcp{
					{
						ArcAutoProvisioning:                     flattenDatabasesGcpGcpArcAutoProvisioning(v.ArcAutoProvisioning),
						DefenderForDatabasesArcAutoProvisioning: flattenDatabasesGcpDefenderForDatabasesArcAutoProvisioning(v.DefenderForDatabasesArcAutoProvisioning),
					},
				},
			}

			result = append(result, defenderForDatabasesGcpOffering)
		} else if v, ok := item.(securityconnectors.DefenderForContainersAwsOffering); ok {
			defenderForContainersAwsOffering := Offering{
				Type: string(securityconnectors.OfferingTypeDefenderForContainersAws),
				DefenderForContainersAws: []DefenderForContainersAws{
					{
						AutoProvisioningEnabled:                 pointer.From(v.AutoProvisioning),
						ContainerVulnerabilityAssessmentEnabled: pointer.From(v.EnableContainerVulnerabilityAssessment),
						KubeAuditRetentionTime:                  int(pointer.From(v.KubeAuditRetentionTime)),
						MdcContainersAgentlessDiscoveryK8s:      flattenContainersAwsMdcContainersAgentlessDiscoveryK8s(v.MdcContainersAgentlessDiscoveryK8s),
						MdcContainersImageAssessment:            flattenContainersAwsMdcContainersImageAssessment(v.MdcContainersImageAssessment),
						ScubaExternalId:                         pointer.From(v.ScubaExternalId),
					},
				},
			}

			if cloudWatchToKinesis := v.CloudWatchToKinesis; cloudWatchToKinesis != nil {
				defenderForContainersAwsOffering.DefenderForContainersAws[0].CloudWatchToKinesisCloudRoleArn = pointer.From(cloudWatchToKinesis.CloudRoleArn)
			}

			if containerVulnerabilityAssessment := v.ContainerVulnerabilityAssessment; containerVulnerabilityAssessment != nil {
				defenderForContainersAwsOffering.DefenderForContainersAws[0].ContainerVulnerabilityAssessmentCloudRoleArn = pointer.From(containerVulnerabilityAssessment.CloudRoleArn)
			}

			if containerVulnerabilityAssessmentTask := v.ContainerVulnerabilityAssessmentTask; containerVulnerabilityAssessmentTask != nil {
				defenderForContainersAwsOffering.DefenderForContainersAws[0].ContainerVulnerabilityAssessmentTaskCloudRoleArn = pointer.From(containerVulnerabilityAssessmentTask.CloudRoleArn)
			}

			if kinesisToS3 := v.KinesisToS3; kinesisToS3 != nil {
				defenderForContainersAwsOffering.DefenderForContainersAws[0].KinesisToS3CloudRoleArn = pointer.From(kinesisToS3.CloudRoleArn)
			}

			if kubernetesScubaReader := v.KubernetesScubaReader; kubernetesScubaReader != nil {
				defenderForContainersAwsOffering.DefenderForContainersAws[0].KubernetesScubaReaderCloudRoleArn = pointer.From(kubernetesScubaReader.CloudRoleArn)
			}

			if kubernetesService := v.KubernetesService; kubernetesService != nil {
				defenderForContainersAwsOffering.DefenderForContainersAws[0].KubernetesServiceCloudRoleArn = pointer.From(kubernetesService.CloudRoleArn)
			}

			result = append(result, defenderForContainersAwsOffering)
		} else if v, ok := item.(securityconnectors.DefenderForContainersGcpOffering); ok {
			defenderForContainersGcpOffering := Offering{
				Type: string(securityconnectors.OfferingTypeDefenderForContainersGcp),
				DefenderForContainersGcp: []DefenderForContainersGcp{
					{
						AuditLogsAutoProvisioningFlagEnabled:     pointer.From(v.AuditLogsAutoProvisioningFlag),
						DataPipelineNativeCloudConnection:        flattenContainersGcpDataPipelineNativeCloudConnection(v.DataPipelineNativeCloudConnection),
						DefenderAgentAutoProvisioningFlagEnabled: pointer.From(v.DefenderAgentAutoProvisioningFlag),
						MdcContainersAgentlessDiscoveryK8s:       flattenContainersGcpMdcContainersAgentlessDiscoveryK8s(v.MdcContainersAgentlessDiscoveryK8s),
						MdcContainersImageAssessment:             flattenContainersGcpMdcContainersImageAssessment(v.MdcContainersImageAssessment),
						NativeCloudConnection:                    flattenContainersGcpNativeCloudConnection(v.NativeCloudConnection),
						PolicyAgentAutoProvisioningFlagEnabled:   pointer.From(v.PolicyAgentAutoProvisioningFlag),
					},
				},
			}

			result = append(result, defenderForContainersGcpOffering)
		}
	}

	return result
}

func flattenDatabasesAwsArcAutoProvisioning(input *securityconnectors.DefenderFoDatabasesAwsOfferingArcAutoProvisioning) []DatabasesAwsArcAutoProvisioning {
	result := make([]DatabasesAwsArcAutoProvisioning, 0)
	if input == nil {
		return result
	}

	awsArcAutoProvisioning := DatabasesAwsArcAutoProvisioning{
		CloudRoleArn:  pointer.From(input.CloudRoleArn),
		Configuration: flattenDatabasesAwsConfiguration(input.Configuration),
		Enabled:       pointer.From(input.Enabled),
	}

	return append(result, awsArcAutoProvisioning)
}

func flattenDatabasesAwsConfiguration(input *securityconnectors.DefenderFoDatabasesAwsOfferingArcAutoProvisioningConfiguration) []DatabasesAwsConfiguration {
	result := make([]DatabasesAwsConfiguration, 0)
	if input == nil {
		return result
	}

	awsConfiguration := DatabasesAwsConfiguration{
		PrivateLinkScope: pointer.From(input.PrivateLinkScope),
		Proxy:            pointer.From(input.Proxy),
	}

	return append(result, awsConfiguration)
}

func flattenDatabasesAwsDspm(input *securityconnectors.DefenderFoDatabasesAwsOfferingDatabasesDspm) []DatabasesAwsDspm {
	result := make([]DatabasesAwsDspm, 0)
	if input == nil {
		return result
	}

	databasesDspm := DatabasesAwsDspm{
		CloudRoleArn: pointer.From(input.CloudRoleArn),
		Enabled:      pointer.From(input.Enabled),
	}

	return append(result, databasesDspm)
}

func flattenDatabasesAwsRds(input *securityconnectors.DefenderFoDatabasesAwsOfferingRds) []DatabasesAwsRds {
	result := make([]DatabasesAwsRds, 0)
	if input == nil {
		return result
	}

	rds := DatabasesAwsRds{
		CloudRoleArn: pointer.From(input.CloudRoleArn),
		Enabled:      pointer.From(input.Enabled),
	}

	return append(result, rds)
}

func flattenDatabasesGcpGcpArcAutoProvisioning(input *securityconnectors.DefenderForDatabasesGcpOfferingArcAutoProvisioning) []DatabasesGcpArcAutoProvisioning {
	result := make([]DatabasesGcpArcAutoProvisioning, 0)
	if input == nil {
		return result
	}

	gcpArcAutoProvisioning := DatabasesGcpArcAutoProvisioning{
		Configuration: flattenDatabasesGcpConfiguration(input.Configuration),
		Enabled:       pointer.From(input.Enabled),
	}

	return append(result, gcpArcAutoProvisioning)
}

func flattenDatabasesGcpConfiguration(input *securityconnectors.DefenderForDatabasesGcpOfferingArcAutoProvisioningConfiguration) []DatabasesGcpConfiguration {
	result := make([]DatabasesGcpConfiguration, 0)
	if input == nil {
		return result
	}

	gcpConfiguration := DatabasesGcpConfiguration{
		PrivateLinkScope: pointer.From(input.PrivateLinkScope),
		Proxy:            pointer.From(input.Proxy),
	}

	return append(result, gcpConfiguration)
}

func flattenDatabasesGcpDefenderForDatabasesArcAutoProvisioning(input *securityconnectors.DefenderForDatabasesGcpOfferingDefenderForDatabasesArcAutoProvisioning) []DatabasesGcpDefenderForDatabasesArcAutoProvisioning {
	result := make([]DatabasesGcpDefenderForDatabasesArcAutoProvisioning, 0)
	if input == nil {
		return result
	}

	defenderForDatabasesArcAutoProvisioning := DatabasesGcpDefenderForDatabasesArcAutoProvisioning{
		ServiceAccountEmailAddress: pointer.From(input.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.From(input.WorkloadIdentityProviderId),
	}

	return append(result, defenderForDatabasesArcAutoProvisioning)
}

func flattenContainersAwsMdcContainersAgentlessDiscoveryK8s(input *securityconnectors.DefenderForContainersAwsOfferingMdcContainersAgentlessDiscoveryK8s) []ContainersAwsMdcContainersAgentlessDiscoveryK8s {
	result := make([]ContainersAwsMdcContainersAgentlessDiscoveryK8s, 0)
	if input == nil {
		return result
	}

	awsMdcContainersAgentlessDiscoveryK8s := ContainersAwsMdcContainersAgentlessDiscoveryK8s{
		CloudRoleArn: pointer.From(input.CloudRoleArn),
		Enabled:      pointer.From(input.Enabled),
	}

	return append(result, awsMdcContainersAgentlessDiscoveryK8s)
}

func flattenContainersAwsMdcContainersImageAssessment(input *securityconnectors.DefenderForContainersAwsOfferingMdcContainersImageAssessment) []ContainersAwsMdcContainersImageAssessment {
	result := make([]ContainersAwsMdcContainersImageAssessment, 0)
	if input == nil {
		return result
	}

	awsMdcContainersImageAssessment := ContainersAwsMdcContainersImageAssessment{
		CloudRoleArn: pointer.From(input.CloudRoleArn),
		Enabled:      pointer.From(input.Enabled),
	}

	return append(result, awsMdcContainersImageAssessment)
}

func flattenContainersGcpDataPipelineNativeCloudConnection(input *securityconnectors.DefenderForContainersGcpOfferingDataPipelineNativeCloudConnection) []ContainersGcpDataPipelineNativeCloudConnection {
	result := make([]ContainersGcpDataPipelineNativeCloudConnection, 0)
	if input == nil {
		return result
	}

	dataPipelineNativeCloudConnection := ContainersGcpDataPipelineNativeCloudConnection{
		ServiceAccountEmailAddress: pointer.From(input.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.From(input.WorkloadIdentityProviderId),
	}

	return append(result, dataPipelineNativeCloudConnection)
}

func flattenContainersGcpMdcContainersAgentlessDiscoveryK8s(input *securityconnectors.DefenderForContainersGcpOfferingMdcContainersAgentlessDiscoveryK8s) []ContainersGcpMdcContainersAgentlessDiscoveryK8s {
	result := make([]ContainersGcpMdcContainersAgentlessDiscoveryK8s, 0)
	if input == nil {
		return result
	}

	gcpMdcContainersAgentlessDiscoveryK8s := ContainersGcpMdcContainersAgentlessDiscoveryK8s{
		ServiceAccountEmailAddress: pointer.From(input.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.From(input.WorkloadIdentityProviderId),
		Enabled:                    pointer.From(input.Enabled),
	}

	return append(result, gcpMdcContainersAgentlessDiscoveryK8s)
}

func flattenContainersGcpMdcContainersImageAssessment(input *securityconnectors.DefenderForContainersGcpOfferingMdcContainersImageAssessment) []ContainersGcpMdcContainersImageAssessment {
	result := make([]ContainersGcpMdcContainersImageAssessment, 0)
	if input == nil {
		return result
	}

	gcpMdcContainersImageAssessment := ContainersGcpMdcContainersImageAssessment{
		ServiceAccountEmailAddress: pointer.From(input.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.From(input.WorkloadIdentityProviderId),
		Enabled:                    pointer.From(input.Enabled),
	}

	return append(result, gcpMdcContainersImageAssessment)
}

func flattenContainersGcpNativeCloudConnection(input *securityconnectors.DefenderForContainersGcpOfferingNativeCloudConnection) []ContainersGcpNativeCloudConnection {
	result := make([]ContainersGcpNativeCloudConnection, 0)
	if input == nil {
		return result
	}

	containersGcpNativeCloudConnection := ContainersGcpNativeCloudConnection{
		ServiceAccountEmailAddress: pointer.From(input.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.From(input.WorkloadIdentityProviderId),
	}

	return append(result, containersGcpNativeCloudConnection)
}
