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
	Regions      []string `tfschema:"regions"`
	ScanInterval int      `tfschema:"scan_interval"`
}

type GcpProject struct {
	ProjectDetails []GcpProjectDetails `tfschema:"project_details"`
	ScanInterval   int                 `tfschema:"scan_interval"`
}

type GcpProjectDetails struct {
	ProjectId     string `tfschema:"project_id"`
	ProjectNumber string `tfschema:"project_number"`
}

type Offering struct {
	Type                                            string                     `tfschema:"type"`
	CspmMonitorAwsNativeCloudConnectionCloudRoleArn string                     `tfschema:"cspm_monitor_aws_native_cloud_connection_cloud_role_arn"`
	CspmMonitorGcp                                  []CspmMonitorGcp           `tfschema:"cspm_monitor_gcp"`
	DefenderForDatabasesAws                         []DefenderForDatabasesAws  `tfschema:"defender_for_databases_aws"`
	DefenderForDatabasesGcp                         []DefenderForDatabasesGcp  `tfschema:"defender_for_databases_gcp"`
	DefenderForContainersAws                        []DefenderForContainersAws `tfschema:"defender_for_containers_aws"`
	DefenderForContainersGcp                        []DefenderForContainersGcp `tfschema:"defender_for_containers_gcp"`
	DefenderForServersAws                           []DefenderForServersAws    `tfschema:"defender_for_servers_aws"`
	DefenderForServersGcp                           []DefenderForServersGcp    `tfschema:"defender_for_servers_gcp"`
	DefenderCspmAws                                 []DefenderCspmAws          `tfschema:"defender_cspm_aws"`
	DefenderCspmGcp                                 []DefenderCspmGcp          `tfschema:"defender_cspm_gcp"`
}

type CspmMonitorGcp struct {
	ServiceAccountEmailAddress string `tfschema:"service_account_email_address"`
	WorkloadIdentityProviderId string `tfschema:"workload_identity_provider_id"`
}

type DefenderForDatabasesAws struct {
	ArcAutoProvisioningCloudRoleArn string `tfschema:"arc_auto_provisioning_cloud_role_arn"`
	DatabasesDspmCloudRoleArn       string `tfschema:"databases_dspm_cloud_role_arn"`
	RdsCloudRoleArn                 string `tfschema:"rds_cloud_role_arn"`
}

type DefenderForDatabasesGcp struct {
	ArcAutoProvisioningEnabled              bool                                                  `tfschema:"arc_auto_provisioning_enabled"`
	DefenderForDatabasesArcAutoProvisioning []DatabasesGcpDefenderForDatabasesArcAutoProvisioning `tfschema:"defender_for_databases_arc_auto_provisioning"`
}

type DatabasesGcpDefenderForDatabasesArcAutoProvisioning struct {
	ServiceAccountEmailAddress string `tfschema:"service_account_email_address"`
	WorkloadIdentityProviderId string `tfschema:"workload_identity_provider_id"`
}

type DefenderForContainersAws struct {
	AutoProvisioningEnabled                          bool   `tfschema:"auto_provisioning_enabled"`
	CloudWatchToKinesisCloudRoleArn                  string `tfschema:"cloud_watch_to_kinesis_cloud_role_arn"`
	ContainerVulnerabilityAssessmentCloudRoleArn     string `tfschema:"container_vulnerability_assessment_cloud_role_arn"`
	ContainerVulnerabilityAssessmentTaskCloudRoleArn string `tfschema:"container_vulnerability_assessment_task_cloud_role_arn"`
	ContainerVulnerabilityAssessmentEnabled          bool   `tfschema:"container_vulnerability_assessment_enabled"`
	KinesisToS3CloudRoleArn                          string `tfschema:"kinesis_to_s3_cloud_role_arn"`
	KubeAuditRetentionTime                           int    `tfschema:"kube_audit_retention_time"`
	KubernetesScubaReaderCloudRoleArn                string `tfschema:"kubernetes_scuba_reader_cloud_role_arn"`
	KubernetesServiceCloudRoleArn                    string `tfschema:"kubernetes_service_cloud_role_arn"`
	MdcContainersAgentlessDiscoveryK8sCloudRoleArn   string `tfschema:"mdc_containers_agentless_discovery_k8s_cloud_role_arn"`
	MdcContainersImageAssessmentCloudRoleArn         string `tfschema:"mdc_containers_image_assessment_cloud_role_arn"`
	ScubaExternalId                                  string `tfschema:"scuba_external_id"`
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
}

type ContainersGcpMdcContainersImageAssessment struct {
	ServiceAccountEmailAddress string `tfschema:"service_account_email_address"`
	WorkloadIdentityProviderId string `tfschema:"workload_identity_provider_id"`
}

type ContainersGcpNativeCloudConnection struct {
	ServiceAccountEmailAddress string `tfschema:"service_account_email_address"`
	WorkloadIdentityProviderId string `tfschema:"workload_identity_provider_id"`
}

type DefenderForServersAws struct {
	ArcAutoProvisioningCloudRoleArn     string                `tfschema:"arc_auto_provisioning_cloud_role_arn"`
	DefenderForServersCloudRoleArn      string                `tfschema:"defender_for_servers_cloud_role_arn"`
	MdeAutoProvisioningEnabled          bool                  `tfschema:"mde_auto_provisioning_enabled"`
	SubPlanType                         string                `tfschema:"sub_plan_type"`
	VaAutoProvisioningConfigurationType string                `tfschema:"va_auto_provisioning_configuration_type"`
	VmScanner                           []ServersAwsVmScanner `tfschema:"vm_scanner"`
}

type ServersAwsVmScanner struct {
	CloudRoleArn  string            `tfschema:"cloud_role_arn"`
	ExclusionTags map[string]string `tfschema:"exclusion_tags"`
	ScanningMode  string            `tfschema:"scanning_mode"`
}

type DefenderForServersGcp struct {
	ArcAutoProvisioningEnabled          bool                           `tfschema:"arc_auto_provisioning_enabled"`
	DefenderForServers                  []ServersGcpDefenderForServers `tfschema:"defender_for_servers"`
	MdeAutoProvisioningEnabled          bool                           `tfschema:"mde_auto_provisioning_enabled"`
	SubPlanType                         string                         `tfschema:"sub_plan_type"`
	VaAutoProvisioningConfigurationType string                         `tfschema:"va_auto_provisioning_configuration_type"`
	VmScanner                           []ServersGcpVmScanner          `tfschema:"vm_scanner"`
}

type ServersGcpDefenderForServers struct {
	ServiceAccountEmailAddress string `tfschema:"service_account_email_address"`
	WorkloadIdentityProviderId string `tfschema:"workload_identity_provider_id"`
}

type ServersGcpVmScanner struct {
	ExclusionTags map[string]string `tfschema:"exclusion_tags"`
	ScanningMode  string            `tfschema:"scanning_mode"`
}

type DefenderCspmAws struct {
	Ciem                                           []CspmAwsCiem      `tfschema:"ciem"`
	DataSensitivityDiscoveryCloudRoleArn           string             `tfschema:"data_sensitivity_discovery_cloud_role_arn"`
	DatabasesDspmCloudRoleArn                      string             `tfschema:"databases_dspm_cloud_role_arn"`
	MdcContainersAgentlessDiscoveryK8sCloudRoleArn string             `tfschema:"mdc_containers_agentless_discovery_k8s_cloud_role_arn"`
	MdcContainersImageAssessmentCloudRoleArn       string             `tfschema:"mdc_containers_image_assessment_cloud_role_arn"`
	VmScanner                                      []CspmAwsVmScanner `tfschema:"vm_scanner"`
}

type CspmAwsCiem struct {
	CiemDiscoveryCloudRoleArn string            `tfschema:"ciem_discovery_cloud_role_arn"`
	CiemOidc                  []CspmAwsCiemOidc `tfschema:"ciem_oidc"`
}

type CspmAwsCiemOidc struct {
	AzureActiveDirectoryAppName string `tfschema:"azure_active_directory_app_name"`
	CloudRoleArn                string `tfschema:"cloud_role_arn"`
}

type CspmAwsVmScanner struct {
	CloudRoleArn  string            `tfschema:"cloud_role_arn"`
	ExclusionTags map[string]string `tfschema:"exclusion_tags"`
	ScanningMode  string            `tfschema:"scanning_mode"`
}

type DefenderCspmGcp struct {
	CiemDiscovery                      []CspmGcpCiemDiscovery                      `tfschema:"ciem_discovery"`
	DataSensitivityDiscovery           []CspmGcpDataSensitivityDiscovery           `tfschema:"data_sensitivity_discovery"`
	MdcContainersAgentlessDiscoveryK8s []CspmGcpMdcContainersAgentlessDiscoveryK8s `tfschema:"mdc_containers_agentless_discovery_k8s"`
	MdcContainersImageAssessment       []CspmGcpMdcContainersImageAssessment       `tfschema:"mdc_containers_image_assessment"`
	VmScanner                          []CspmGcpVmScanner                          `tfschema:"vm_scanner"`
}

type CspmGcpCiemDiscovery struct {
	AzureActiveDirectoryAppName string `tfschema:"azure_active_directory_app_name"`
	ServiceAccountEmailAddress  string `tfschema:"service_account_email_address"`
	WorkloadIdentityProviderId  string `tfschema:"workload_identity_provider_id"`
}

type CspmGcpDataSensitivityDiscovery struct {
	ServiceAccountEmailAddress string `tfschema:"service_account_email_address"`
	WorkloadIdentityProviderId string `tfschema:"workload_identity_provider_id"`
}

type CspmGcpMdcContainersAgentlessDiscoveryK8s struct {
	ServiceAccountEmailAddress string `tfschema:"service_account_email_address"`
	WorkloadIdentityProviderId string `tfschema:"workload_identity_provider_id"`
}

type CspmGcpMdcContainersImageAssessment struct {
	ServiceAccountEmailAddress string `tfschema:"service_account_email_address"`
	WorkloadIdentityProviderId string `tfschema:"workload_identity_provider_id"`
}

type CspmGcpVmScanner struct {
	ExclusionTags map[string]string `tfschema:"exclusion_tags"`
	ScanningMode  string            `tfschema:"scanning_mode"`
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
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"regions": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
										AtLeastOneOf: []string{"environment_data.0.aws_account.0.regions", "environment_data.0.aws_account.0.scan_interval"},
									},
								},

								"scan_interval": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      4,
									ValidateFunc: validation.IntBetween(1, 24),
									AtLeastOneOf: []string{"environment_data.0.aws_account.0.regions", "environment_data.0.aws_account.0.scan_interval"},
								},
							},
						},
					},

					"gcp_project": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"project_details": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"project_id": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"project_number": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
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
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
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

					"cspm_monitor_aws_native_cloud_connection_cloud_role_arn": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"cspm_monitor_gcp": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"service_account_email_address": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"workload_identity_provider_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
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
								"arc_auto_provisioning_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"databases_dspm_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"rds_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
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
								"arc_auto_provisioning_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"defender_for_databases_arc_auto_provisioning": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"service_account_email_address": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"workload_identity_provider_id": {
												Type:         pluginsdk.TypeString,
												Required:     true,
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
									Default:  true,
								},

								"cloud_watch_to_kinesis_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"kinesis_to_s3_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"kubernetes_scuba_reader_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"kubernetes_service_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"container_vulnerability_assessment_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"container_vulnerability_assessment_task_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"container_vulnerability_assessment_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"kube_audit_retention_time": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									Default:  30,
								},

								"mdc_containers_agentless_discovery_k8s_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"mdc_containers_image_assessment_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"scuba_external_id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
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
								"native_cloud_connection": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"service_account_email_address": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"workload_identity_provider_id": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},

								"data_pipeline_native_cloud_connection": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"service_account_email_address": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"workload_identity_provider_id": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},

								"audit_logs_auto_provisioning_flag_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"defender_agent_auto_provisioning_flag_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"policy_agent_auto_provisioning_flag_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"mdc_containers_agentless_discovery_k8s": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"service_account_email_address": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"workload_identity_provider_id": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},

								"mdc_containers_image_assessment": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"service_account_email_address": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"workload_identity_provider_id": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},
							},
						},
					},

					"defender_for_servers_aws": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"defender_for_servers_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"arc_auto_provisioning_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"mde_auto_provisioning_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"sub_plan_type": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      "P2",
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"va_auto_provisioning_configuration_type": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"vm_scanner": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"cloud_role_arn": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"exclusion_tags": commonschema.Tags(),

											"scanning_mode": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												Default:      "Default",
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},
							},
						},
					},

					"defender_for_servers_gcp": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"arc_auto_provisioning_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"defender_for_servers": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"service_account_email_address": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"workload_identity_provider_id": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},

								"mde_auto_provisioning_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"sub_plan_type": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      "P2",
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"va_auto_provisioning_configuration_type": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"vm_scanner": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"exclusion_tags": commonschema.Tags(),

											"scanning_mode": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												Default:      "Default",
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},
							},
						},
					},

					"defender_cspm_aws": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"ciem": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"ciem_discovery_cloud_role_arn": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"ciem_oidc": {
												Type:     pluginsdk.TypeList,
												Required: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"cloud_role_arn": {
															Type:         pluginsdk.TypeString,
															Required:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"azure_active_directory_app_name": {
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

								"data_sensitivity_discovery_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"databases_dspm_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"mdc_containers_agentless_discovery_k8s_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"mdc_containers_image_assessment_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"vm_scanner": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"cloud_role_arn": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"exclusion_tags": commonschema.Tags(),

											"scanning_mode": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												Default:      false,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},
							},
						},
					},

					"defender_cspm_gcp": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"ciem_discovery": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"service_account_email_address": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"workload_identity_provider_id": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"azure_active_directory_app_name": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},

								"data_sensitivity_discovery": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"service_account_email_address": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"workload_identity_provider_id": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},

								"mdc_containers_agentless_discovery_k8s": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"service_account_email_address": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"workload_identity_provider_id": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},

								"mdc_containers_image_assessment": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"service_account_email_address": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"workload_identity_provider_id": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},

								"vm_scanner": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"exclusion_tags": commonschema.Tags(),

											"scanning_mode": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												Default:      "Default",
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
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

			if v := item.CspmMonitorAwsNativeCloudConnectionCloudRoleArn; v != "" {
				cspmMonitorAwsOffering.NativeCloudConnection = &securityconnectors.CspmMonitorAwsOfferingNativeCloudConnection{
					CloudRoleArn: pointer.To(v),
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

				defenderFoDatabasesAwsOffering.ArcAutoProvisioning = expandDatabasesAwsArcAutoProvisioning(defenderForDatabasesAws.ArcAutoProvisioningCloudRoleArn)
				defenderFoDatabasesAwsOffering.DatabasesDspm = expandDatabasesAwsDspm(defenderForDatabasesAws.DatabasesDspmCloudRoleArn)
				defenderFoDatabasesAwsOffering.Rds = expandDatabasesAwsRds(defenderForDatabasesAws.RdsCloudRoleArn)
			}

			result = append(result, defenderFoDatabasesAwsOffering)
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeDefenderForDatabasesGcp) {
			defenderForDatabasesGcpOffering := securityconnectors.DefenderForDatabasesGcpOffering{}

			if v := item.DefenderForDatabasesGcp; len(v) != 0 {
				defenderForDatabasesGcp := v[0]

				defenderForDatabasesGcpOffering.ArcAutoProvisioning = expandDatabasesGcpArcAutoProvisioning(defenderForDatabasesGcp.ArcAutoProvisioningEnabled)
				defenderForDatabasesGcpOffering.DefenderForDatabasesArcAutoProvisioning = expandDatabasesGcpDefenderForDatabasesArcAutoProvisioning(defenderForDatabasesGcp.DefenderForDatabasesArcAutoProvisioning)
			}

			result = append(result, defenderForDatabasesGcpOffering)
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeDefenderForContainersAws) {
			defenderForContainersAwsOffering := securityconnectors.DefenderForContainersAwsOffering{}

			if v := item.DefenderForContainersAws; len(v) != 0 {
				defenderForContainersAws := v[0]

				defenderForContainersAwsOffering.AutoProvisioning = pointer.To(defenderForContainersAws.AutoProvisioningEnabled)
				defenderForContainersAwsOffering.EnableContainerVulnerabilityAssessment = pointer.To(defenderForContainersAws.ContainerVulnerabilityAssessmentEnabled)
				defenderForContainersAwsOffering.MdcContainersAgentlessDiscoveryK8s = expandContainersAwsMdcContainersAgentlessDiscoveryK8s(defenderForContainersAws.MdcContainersAgentlessDiscoveryK8sCloudRoleArn)
				defenderForContainersAwsOffering.MdcContainersImageAssessment = expandContainersAwsMdcContainersImageAssessment(defenderForContainersAws.MdcContainersImageAssessmentCloudRoleArn)
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
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeDefenderForServersAws) {
			defenderForServersAwsOffering := securityconnectors.DefenderForServersAwsOffering{}

			if v := item.DefenderForServersAws; len(v) != 0 {
				defenderForServersAws := v[0]

				defenderForServersAwsOffering.ArcAutoProvisioning = expandServersAwsArcAutoProvisioning(defenderForServersAws.ArcAutoProvisioningCloudRoleArn)
				defenderForServersAwsOffering.MdeAutoProvisioning = expandServersAwsMdeAutoProvisioning(defenderForServersAws.MdeAutoProvisioningEnabled)
				defenderForServersAwsOffering.VaAutoProvisioning = expandServersAwsVaAutoProvisioning(defenderForServersAws.VaAutoProvisioningConfigurationType)
				defenderForServersAwsOffering.VMScanners = expandServersAwsVMScanners(defenderForServersAws.VmScanner)

				defenderForServersAwsOffering.DefenderForServers = &securityconnectors.DefenderForServersAwsOfferingDefenderForServers{
					CloudRoleArn: pointer.To(defenderForServersAws.DefenderForServersCloudRoleArn),
				}

				defenderForServersAwsOffering.SubPlan = &securityconnectors.DefenderForServersAwsOfferingSubPlan{
					Type: pointer.To(securityconnectors.SubPlan(defenderForServersAws.SubPlanType)),
				}
			}

			result = append(result, defenderForServersAwsOffering)
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeDefenderForServersGcp) {
			defenderForServersGcpOffering := securityconnectors.DefenderForServersGcpOffering{}

			if v := item.DefenderForServersGcp; len(v) != 0 {
				defenderForServersGcp := v[0]

				defenderForServersGcpOffering.ArcAutoProvisioning = expandServersGcpArcAutoProvisioning(defenderForServersGcp.ArcAutoProvisioningEnabled)
				defenderForServersGcpOffering.DefenderForServers = expandServersGcpDefenderForServers(defenderForServersGcp.DefenderForServers)
				defenderForServersGcpOffering.MdeAutoProvisioning = expandServersGcpMdeAutoProvisioning(defenderForServersGcp.MdeAutoProvisioningEnabled)
				defenderForServersGcpOffering.VaAutoProvisioning = expandServersGcpVaAutoProvisioning(defenderForServersGcp.VaAutoProvisioningConfigurationType)
				defenderForServersGcpOffering.VMScanners = expandServersGcpVMScanners(defenderForServersGcp.VmScanner)

				defenderForServersGcpOffering.SubPlan = &securityconnectors.DefenderForServersGcpOfferingSubPlan{
					Type: pointer.To(securityconnectors.SubPlan(defenderForServersGcp.SubPlanType)),
				}
			}

			result = append(result, defenderForServersGcpOffering)
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeDefenderCspmAws) {
			defenderCspmAwsOffering := securityconnectors.DefenderCspmAwsOffering{}

			if v := item.DefenderCspmAws; len(v) != 0 {
				defenderCspmAws := v[0]

				defenderCspmAwsOffering.Ciem = expandCspmAwsCiem(defenderCspmAws.Ciem)
				defenderCspmAwsOffering.DataSensitivityDiscovery = expandCspmAwsDataSensitivityDiscovery(defenderCspmAws.DataSensitivityDiscoveryCloudRoleArn)
				defenderCspmAwsOffering.DatabasesDspm = expandCspmAwsDatabasesDspm(defenderCspmAws.DatabasesDspmCloudRoleArn)
				defenderCspmAwsOffering.MdcContainersAgentlessDiscoveryK8s = expandCspmAwsMdcContainersAgentlessDiscoveryK8s(defenderCspmAws.MdcContainersAgentlessDiscoveryK8sCloudRoleArn)
				defenderCspmAwsOffering.MdcContainersImageAssessment = expandCspmAwsMdcContainersImageAssessment(defenderCspmAws.MdcContainersImageAssessmentCloudRoleArn)
				defenderCspmAwsOffering.VMScanners = expandCspmAwsVMScanners(defenderCspmAws.VmScanner)
			}

			result = append(result, defenderCspmAwsOffering)
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeDefenderCspmGcp) {
			defenderCspmGcpOffering := securityconnectors.DefenderCspmGcpOffering{}

			if v := item.DefenderCspmGcp; len(v) != 0 {
				defenderCspmGcp := v[0]

				defenderCspmGcpOffering.CiemDiscovery = expandCspmGcpCiemDiscovery(defenderCspmGcp.CiemDiscovery)
				defenderCspmGcpOffering.DataSensitivityDiscovery = expandCspmGcpDataSensitivityDiscovery(defenderCspmGcp.DataSensitivityDiscovery)
				defenderCspmGcpOffering.MdcContainersAgentlessDiscoveryK8s = expandCspmGcpMdcContainersAgentlessDiscoveryK8s(defenderCspmGcp.MdcContainersAgentlessDiscoveryK8s)
				defenderCspmGcpOffering.MdcContainersImageAssessment = expandCspmGcpMdcContainersImageAssessment(defenderCspmGcp.MdcContainersImageAssessment)
				defenderCspmGcpOffering.VMScanners = expandCspmGcpVMScanners(defenderCspmGcp.VmScanner)
			}

			result = append(result, defenderCspmGcpOffering)
		}
	}

	return &result
}

func expandDatabasesAwsArcAutoProvisioning(input string) *securityconnectors.DefenderFoDatabasesAwsOfferingArcAutoProvisioning {
	if input == "" {
		return &securityconnectors.DefenderFoDatabasesAwsOfferingArcAutoProvisioning{
			Enabled: pointer.To(false),
		}
	}

	result := &securityconnectors.DefenderFoDatabasesAwsOfferingArcAutoProvisioning{
		CloudRoleArn:  pointer.To(input),
		Configuration: &securityconnectors.DefenderFoDatabasesAwsOfferingArcAutoProvisioningConfiguration{},
		Enabled:       pointer.To(true),
	}

	return result
}

func expandDatabasesAwsDspm(input string) *securityconnectors.DefenderFoDatabasesAwsOfferingDatabasesDspm {
	if input == "" {
		return &securityconnectors.DefenderFoDatabasesAwsOfferingDatabasesDspm{
			Enabled: pointer.To(false),
		}
	}

	result := &securityconnectors.DefenderFoDatabasesAwsOfferingDatabasesDspm{
		CloudRoleArn: pointer.To(input),
		Enabled:      pointer.To(true),
	}

	return result
}

func expandDatabasesAwsRds(input string) *securityconnectors.DefenderFoDatabasesAwsOfferingRds {
	if input == "" {
		return &securityconnectors.DefenderFoDatabasesAwsOfferingRds{
			Enabled: pointer.To(false),
		}
	}

	result := &securityconnectors.DefenderFoDatabasesAwsOfferingRds{
		CloudRoleArn: pointer.To(input),
		Enabled:      pointer.To(true),
	}

	return result
}

func expandDatabasesGcpArcAutoProvisioning(input bool) *securityconnectors.DefenderForDatabasesGcpOfferingArcAutoProvisioning {
	result := &securityconnectors.DefenderForDatabasesGcpOfferingArcAutoProvisioning{
		Configuration: &securityconnectors.DefenderForDatabasesGcpOfferingArcAutoProvisioningConfiguration{},
		Enabled:       pointer.To(input),
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

func expandContainersAwsMdcContainersAgentlessDiscoveryK8s(input string) *securityconnectors.DefenderForContainersAwsOfferingMdcContainersAgentlessDiscoveryK8s {
	if input == "" {
		return &securityconnectors.DefenderForContainersAwsOfferingMdcContainersAgentlessDiscoveryK8s{
			Enabled: pointer.To(false),
		}
	}

	result := &securityconnectors.DefenderForContainersAwsOfferingMdcContainersAgentlessDiscoveryK8s{
		CloudRoleArn: pointer.To(input),
		Enabled:      pointer.To(true),
	}

	return result
}

func expandContainersAwsMdcContainersImageAssessment(input string) *securityconnectors.DefenderForContainersAwsOfferingMdcContainersImageAssessment {
	if input == "" {
		return &securityconnectors.DefenderForContainersAwsOfferingMdcContainersImageAssessment{
			Enabled: pointer.To(false),
		}
	}

	result := &securityconnectors.DefenderForContainersAwsOfferingMdcContainersImageAssessment{
		CloudRoleArn: pointer.To(input),
		Enabled:      pointer.To(true),
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
		return &securityconnectors.DefenderForContainersGcpOfferingMdcContainersAgentlessDiscoveryK8s{
			Enabled: pointer.To(false),
		}
	}

	gcpMdcContainersAgentlessDiscoveryK8s := input[0]

	result := &securityconnectors.DefenderForContainersGcpOfferingMdcContainersAgentlessDiscoveryK8s{
		ServiceAccountEmailAddress: pointer.To(gcpMdcContainersAgentlessDiscoveryK8s.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.To(gcpMdcContainersAgentlessDiscoveryK8s.WorkloadIdentityProviderId),
		Enabled:                    pointer.To(true),
	}

	return result
}

func expandContainersGcpMdcContainersImageAssessment(input []ContainersGcpMdcContainersImageAssessment) *securityconnectors.DefenderForContainersGcpOfferingMdcContainersImageAssessment {
	if len(input) == 0 {
		return &securityconnectors.DefenderForContainersGcpOfferingMdcContainersImageAssessment{
			Enabled: pointer.To(false),
		}
	}

	gcpMdcContainersImageAssessment := input[0]

	result := &securityconnectors.DefenderForContainersGcpOfferingMdcContainersImageAssessment{
		ServiceAccountEmailAddress: pointer.To(gcpMdcContainersImageAssessment.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.To(gcpMdcContainersImageAssessment.WorkloadIdentityProviderId),
		Enabled:                    pointer.To(true),
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

func expandServersAwsArcAutoProvisioning(input string) *securityconnectors.DefenderForServersAwsOfferingArcAutoProvisioning {
	if input == "" {
		return &securityconnectors.DefenderForServersAwsOfferingArcAutoProvisioning{
			Enabled: pointer.To(false),
		}
	}

	result := &securityconnectors.DefenderForServersAwsOfferingArcAutoProvisioning{
		CloudRoleArn:  pointer.To(input),
		Configuration: &securityconnectors.DefenderForServersAwsOfferingArcAutoProvisioningConfiguration{},
		Enabled:       pointer.To(true),
	}

	return result
}

func expandServersAwsMdeAutoProvisioning(input bool) *securityconnectors.DefenderForServersAwsOfferingMdeAutoProvisioning {
	result := &securityconnectors.DefenderForServersAwsOfferingMdeAutoProvisioning{
		Enabled: pointer.To(input),
	}

	return result
}

func expandServersAwsVaAutoProvisioning(input string) *securityconnectors.DefenderForServersAwsOfferingVaAutoProvisioning {
	if input == "" {
		return &securityconnectors.DefenderForServersAwsOfferingVaAutoProvisioning{
			Enabled: pointer.To(false),
		}
	}

	result := &securityconnectors.DefenderForServersAwsOfferingVaAutoProvisioning{
		Configuration: &securityconnectors.DefenderForServersAwsOfferingVaAutoProvisioningConfiguration{
			Type: pointer.To(securityconnectors.Type(input)),
		},
		Enabled: pointer.To(true),
	}

	return result
}

func expandServersAwsVMScanners(input []ServersAwsVmScanner) *securityconnectors.DefenderForServersAwsOfferingVMScanners {
	if len(input) == 0 {
		return &securityconnectors.DefenderForServersAwsOfferingVMScanners{
			Enabled: pointer.To(false),
		}
	}

	serversAwsVmScanner := input[0]

	result := &securityconnectors.DefenderForServersAwsOfferingVMScanners{
		Configuration: &securityconnectors.DefenderForServersAwsOfferingVMScannersConfiguration{
			CloudRoleArn:  pointer.To(serversAwsVmScanner.CloudRoleArn),
			ExclusionTags: pointer.To(serversAwsVmScanner.ExclusionTags),
			ScanningMode:  pointer.To(securityconnectors.ScanningMode(serversAwsVmScanner.ScanningMode)),
		},
		Enabled: pointer.To(true),
	}

	return result
}

func expandServersGcpArcAutoProvisioning(input bool) *securityconnectors.DefenderForServersGcpOfferingArcAutoProvisioning {
	result := &securityconnectors.DefenderForServersGcpOfferingArcAutoProvisioning{
		Configuration: &securityconnectors.DefenderForServersGcpOfferingArcAutoProvisioningConfiguration{},
		Enabled:       pointer.To(input),
	}

	return result
}

func expandServersGcpDefenderForServers(input []ServersGcpDefenderForServers) *securityconnectors.DefenderForServersGcpOfferingDefenderForServers {
	if len(input) == 0 {
		return nil
	}

	serversGcpDefenderForServers := input[0]

	result := &securityconnectors.DefenderForServersGcpOfferingDefenderForServers{
		ServiceAccountEmailAddress: pointer.To(serversGcpDefenderForServers.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.To(serversGcpDefenderForServers.WorkloadIdentityProviderId),
	}

	return result
}

func expandServersGcpMdeAutoProvisioning(input bool) *securityconnectors.DefenderForServersGcpOfferingMdeAutoProvisioning {
	result := &securityconnectors.DefenderForServersGcpOfferingMdeAutoProvisioning{
		Enabled: pointer.To(input),
	}

	return result
}

func expandServersGcpVaAutoProvisioning(input string) *securityconnectors.DefenderForServersGcpOfferingVaAutoProvisioning {
	if input == "" {
		return &securityconnectors.DefenderForServersGcpOfferingVaAutoProvisioning{
			Enabled: pointer.To(false),
		}
	}

	result := &securityconnectors.DefenderForServersGcpOfferingVaAutoProvisioning{
		Configuration: &securityconnectors.DefenderForServersGcpOfferingVaAutoProvisioningConfiguration{
			Type: pointer.To(securityconnectors.Type(input)),
		},
		Enabled: pointer.To(true),
	}

	return result
}

func expandServersGcpVMScanners(input []ServersGcpVmScanner) *securityconnectors.DefenderForServersGcpOfferingVMScanners {
	if len(input) == 0 {
		return &securityconnectors.DefenderForServersGcpOfferingVMScanners{
			Enabled: pointer.To(false),
		}
	}

	serversGcpVmScanner := input[0]

	result := &securityconnectors.DefenderForServersGcpOfferingVMScanners{
		Configuration: &securityconnectors.DefenderForServersGcpOfferingVMScannersConfiguration{
			ExclusionTags: pointer.To(serversGcpVmScanner.ExclusionTags),
			ScanningMode:  pointer.To(securityconnectors.ScanningMode(serversGcpVmScanner.ScanningMode)),
		},
		Enabled: pointer.To(true),
	}

	return result
}

func expandCspmAwsCiem(input []CspmAwsCiem) *securityconnectors.DefenderCspmAwsOfferingCiem {
	if len(input) == 0 {
		return nil
	}

	cspmAwsCiem := input[0]

	result := &securityconnectors.DefenderCspmAwsOfferingCiem{
		CiemDiscovery: &securityconnectors.DefenderCspmAwsOfferingCiemCiemDiscovery{
			CloudRoleArn: pointer.To(cspmAwsCiem.CiemDiscoveryCloudRoleArn),
		},
		CiemOidc: expandCspmAwsCiemOidc(cspmAwsCiem.CiemOidc),
	}

	return result
}

func expandCspmAwsCiemOidc(input []CspmAwsCiemOidc) *securityconnectors.DefenderCspmAwsOfferingCiemCiemOidc {
	if len(input) == 0 {
		return nil
	}

	cspmAwsCiemOidc := input[0]

	result := &securityconnectors.DefenderCspmAwsOfferingCiemCiemOidc{
		AzureActiveDirectoryAppName: pointer.To(cspmAwsCiemOidc.AzureActiveDirectoryAppName),
		CloudRoleArn:                pointer.To(cspmAwsCiemOidc.CloudRoleArn),
	}

	return result
}

func expandCspmAwsDataSensitivityDiscovery(input string) *securityconnectors.DefenderCspmAwsOfferingDataSensitivityDiscovery {
	if input == "" {
		return &securityconnectors.DefenderCspmAwsOfferingDataSensitivityDiscovery{
			Enabled: pointer.To(false),
		}
	}

	result := &securityconnectors.DefenderCspmAwsOfferingDataSensitivityDiscovery{
		CloudRoleArn: pointer.To(input),
		Enabled:      pointer.To(true),
	}

	return result
}

func expandCspmAwsDatabasesDspm(input string) *securityconnectors.DefenderCspmAwsOfferingDatabasesDspm {
	if input == "" {
		return &securityconnectors.DefenderCspmAwsOfferingDatabasesDspm{
			Enabled: pointer.To(false),
		}
	}

	result := &securityconnectors.DefenderCspmAwsOfferingDatabasesDspm{
		CloudRoleArn: pointer.To(input),
		Enabled:      pointer.To(true),
	}

	return result
}

func expandCspmAwsMdcContainersAgentlessDiscoveryK8s(input string) *securityconnectors.DefenderCspmAwsOfferingMdcContainersAgentlessDiscoveryK8s {
	if input == "" {
		return &securityconnectors.DefenderCspmAwsOfferingMdcContainersAgentlessDiscoveryK8s{
			Enabled: pointer.To(false),
		}
	}

	result := &securityconnectors.DefenderCspmAwsOfferingMdcContainersAgentlessDiscoveryK8s{
		CloudRoleArn: pointer.To(input),
		Enabled:      pointer.To(true),
	}

	return result
}

func expandCspmAwsMdcContainersImageAssessment(input string) *securityconnectors.DefenderCspmAwsOfferingMdcContainersImageAssessment {
	if input == "" {
		return &securityconnectors.DefenderCspmAwsOfferingMdcContainersImageAssessment{
			Enabled: pointer.To(false),
		}
	}

	result := &securityconnectors.DefenderCspmAwsOfferingMdcContainersImageAssessment{
		CloudRoleArn: pointer.To(input),
		Enabled:      pointer.To(true),
	}

	return result
}

func expandCspmAwsVMScanners(input []CspmAwsVmScanner) *securityconnectors.DefenderCspmAwsOfferingVMScanners {
	if len(input) == 0 {
		return &securityconnectors.DefenderCspmAwsOfferingVMScanners{
			Enabled: pointer.To(false),
		}
	}

	cspmAwsVmScanner := input[0]

	result := &securityconnectors.DefenderCspmAwsOfferingVMScanners{
		Configuration: &securityconnectors.DefenderCspmAwsOfferingVMScannersConfiguration{
			CloudRoleArn:  pointer.To(cspmAwsVmScanner.CloudRoleArn),
			ExclusionTags: pointer.To(cspmAwsVmScanner.ExclusionTags),
			ScanningMode:  pointer.To(securityconnectors.ScanningMode(cspmAwsVmScanner.ScanningMode)),
		},
		Enabled: pointer.To(true),
	}

	return result
}

func expandCspmGcpCiemDiscovery(input []CspmGcpCiemDiscovery) *securityconnectors.DefenderCspmGcpOfferingCiemDiscovery {
	if len(input) == 0 {
		return nil
	}

	cspmGcpCiemDiscovery := input[0]

	result := &securityconnectors.DefenderCspmGcpOfferingCiemDiscovery{
		AzureActiveDirectoryAppName: pointer.To(cspmGcpCiemDiscovery.AzureActiveDirectoryAppName),
		ServiceAccountEmailAddress:  pointer.To(cspmGcpCiemDiscovery.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId:  pointer.To(cspmGcpCiemDiscovery.WorkloadIdentityProviderId),
	}

	return result
}

func expandCspmGcpDataSensitivityDiscovery(input []CspmGcpDataSensitivityDiscovery) *securityconnectors.DefenderCspmGcpOfferingDataSensitivityDiscovery {
	if len(input) == 0 {
		return &securityconnectors.DefenderCspmGcpOfferingDataSensitivityDiscovery{
			Enabled: pointer.To(false),
		}
	}

	cspmGcpDataSensitivityDiscovery := input[0]

	result := &securityconnectors.DefenderCspmGcpOfferingDataSensitivityDiscovery{
		ServiceAccountEmailAddress: pointer.To(cspmGcpDataSensitivityDiscovery.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.To(cspmGcpDataSensitivityDiscovery.WorkloadIdentityProviderId),
		Enabled:                    pointer.To(true),
	}

	return result
}

func expandCspmGcpMdcContainersAgentlessDiscoveryK8s(input []CspmGcpMdcContainersAgentlessDiscoveryK8s) *securityconnectors.DefenderCspmGcpOfferingMdcContainersAgentlessDiscoveryK8s {
	if len(input) == 0 {
		return &securityconnectors.DefenderCspmGcpOfferingMdcContainersAgentlessDiscoveryK8s{
			Enabled: pointer.To(false),
		}
	}

	cspmGcpMdcContainersAgentlessDiscoveryK8s := input[0]

	result := &securityconnectors.DefenderCspmGcpOfferingMdcContainersAgentlessDiscoveryK8s{
		ServiceAccountEmailAddress: pointer.To(cspmGcpMdcContainersAgentlessDiscoveryK8s.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.To(cspmGcpMdcContainersAgentlessDiscoveryK8s.WorkloadIdentityProviderId),
		Enabled:                    pointer.To(true),
	}

	return result
}

func expandCspmGcpMdcContainersImageAssessment(input []CspmGcpMdcContainersImageAssessment) *securityconnectors.DefenderCspmGcpOfferingMdcContainersImageAssessment {
	if len(input) == 0 {
		return &securityconnectors.DefenderCspmGcpOfferingMdcContainersImageAssessment{
			Enabled: pointer.To(false),
		}
	}

	cspmGcpMdcContainersImageAssessment := input[0]

	result := &securityconnectors.DefenderCspmGcpOfferingMdcContainersImageAssessment{
		ServiceAccountEmailAddress: pointer.To(cspmGcpMdcContainersImageAssessment.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.To(cspmGcpMdcContainersImageAssessment.WorkloadIdentityProviderId),
		Enabled:                    pointer.To(true),
	}

	return result
}

func expandCspmGcpVMScanners(input []CspmGcpVmScanner) *securityconnectors.DefenderCspmGcpOfferingVMScanners {
	if len(input) == 0 {
		return &securityconnectors.DefenderCspmGcpOfferingVMScanners{
			Enabled: pointer.To(false),
		}
	}

	cspmGcpVmScanner := input[0]

	result := &securityconnectors.DefenderCspmGcpOfferingVMScanners{
		Configuration: &securityconnectors.DefenderCspmGcpOfferingVMScannersConfiguration{
			ExclusionTags: pointer.To(cspmGcpVmScanner.ExclusionTags),
			ScanningMode:  pointer.To(securityconnectors.ScanningMode(cspmGcpVmScanner.ScanningMode)),
		},
		Enabled: pointer.To(true),
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

	return append(result, awsAccount)
}

func flattenGcpProject(input securityconnectors.GcpProjectEnvironmentData) []GcpProject {
	result := make([]GcpProject, 0)

	gcpProject := GcpProject{
		ScanInterval:   int(pointer.From(input.ScanInterval)),
		ProjectDetails: flattenGcpProjectDetails(input.ProjectDetails),
	}

	return append(result, gcpProject)
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
				cspmMonitorAwsOffering.CspmMonitorAwsNativeCloudConnectionCloudRoleArn = pointer.From(nativeCloudConnection.CloudRoleArn)
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
						ArcAutoProvisioningCloudRoleArn: flattenDatabasesAwsArcAutoProvisioning(v.ArcAutoProvisioning),
						DatabasesDspmCloudRoleArn:       flattenDatabasesDspmCloudRoleArn(v.DatabasesDspm),
						RdsCloudRoleArn:                 flattenRdsCloudRoleArn(v.Rds),
					},
				},
			}

			result = append(result, defenderFoDatabasesAwsOffering)
		} else if v, ok := item.(securityconnectors.DefenderForDatabasesGcpOffering); ok {
			defenderForDatabasesGcpOffering := Offering{
				Type: string(securityconnectors.OfferingTypeDefenderForDatabasesGcp),
				DefenderForDatabasesGcp: []DefenderForDatabasesGcp{
					{
						ArcAutoProvisioningEnabled:              flattenDatabasesGcpGcpArcAutoProvisioning(v.ArcAutoProvisioning),
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
						AutoProvisioningEnabled:                        pointer.From(v.AutoProvisioning),
						ContainerVulnerabilityAssessmentEnabled:        pointer.From(v.EnableContainerVulnerabilityAssessment),
						KubeAuditRetentionTime:                         int(pointer.From(v.KubeAuditRetentionTime)),
						MdcContainersAgentlessDiscoveryK8sCloudRoleArn: flattenContainersAwsMdcContainersAgentlessDiscoveryK8s(v.MdcContainersAgentlessDiscoveryK8s),
						MdcContainersImageAssessmentCloudRoleArn:       flattenContainersAwsMdcContainersImageAssessment(v.MdcContainersImageAssessment),
						ScubaExternalId:                                pointer.From(v.ScubaExternalId),
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
		} else if v, ok := item.(securityconnectors.DefenderForServersAwsOffering); ok {
			defenderForServersAwsOffering := Offering{
				Type: string(securityconnectors.OfferingTypeDefenderForServersAws),
				DefenderForServersAws: []DefenderForServersAws{
					{
						ArcAutoProvisioningCloudRoleArn:     flattenServersAwsArcAutoProvisioning(v.ArcAutoProvisioning),
						MdeAutoProvisioningEnabled:          flattenServersAwsMdeAutoProvisioning(v.MdeAutoProvisioning),
						VaAutoProvisioningConfigurationType: flattenServersAwsVaAutoProvisioning(v.VaAutoProvisioning),
						VmScanner:                           flattenServersAwsVmScanner(v.VMScanners),
					},
				},
			}

			if defenderForServers := v.DefenderForServers; defenderForServers != nil {
				defenderForServersAwsOffering.DefenderForServersAws[0].DefenderForServersCloudRoleArn = pointer.From(defenderForServers.CloudRoleArn)
			}

			if subPlan := v.SubPlan; subPlan != nil {
				defenderForServersAwsOffering.DefenderForServersAws[0].SubPlanType = string(pointer.From(subPlan.Type))
			}

			result = append(result, defenderForServersAwsOffering)
		} else if v, ok := item.(securityconnectors.DefenderForServersGcpOffering); ok {
			defenderForServersGcpOffering := Offering{
				Type: string(securityconnectors.OfferingTypeDefenderForServersGcp),
				DefenderForServersGcp: []DefenderForServersGcp{
					{
						ArcAutoProvisioningEnabled:          flattenServersGcpArcAutoProvisioning(v.ArcAutoProvisioning),
						DefenderForServers:                  flattenServersGcpDefenderForServers(v.DefenderForServers),
						MdeAutoProvisioningEnabled:          flattenServersGcpMdeAutoProvisioning(v.MdeAutoProvisioning),
						VaAutoProvisioningConfigurationType: flattenServersGcpVaAutoProvisioning(v.VaAutoProvisioning),
						VmScanner:                           flattenServersGcpVmScanner(v.VMScanners),
					},
				},
			}

			if subPlan := v.SubPlan; subPlan != nil {
				defenderForServersGcpOffering.DefenderForServersGcp[0].SubPlanType = string(pointer.From(subPlan.Type))
			}

			result = append(result, defenderForServersGcpOffering)
		} else if v, ok := item.(securityconnectors.DefenderCspmAwsOffering); ok {
			defenderCspmAwsOffering := Offering{
				Type: string(securityconnectors.OfferingTypeDefenderCspmAws),
				DefenderCspmAws: []DefenderCspmAws{
					{
						Ciem:                                 flattenDefenderCspmAwsCiem(v.Ciem),
						DataSensitivityDiscoveryCloudRoleArn: flattenDefenderCspmAwsDataSensitivityDiscovery(v.DataSensitivityDiscovery),
						DatabasesDspmCloudRoleArn:            flattenDefenderCspmAwsDatabasesDspm(v.DatabasesDspm),
						MdcContainersAgentlessDiscoveryK8sCloudRoleArn: flattenDefenderCspmAwsMdcContainersAgentlessDiscoveryK8s(v.MdcContainersAgentlessDiscoveryK8s),
						MdcContainersImageAssessmentCloudRoleArn:       flattenDefenderCspmAwsMdcContainersImageAssessment(v.MdcContainersImageAssessment),
						VmScanner:                                      flattenDefenderCspmAwsVmScanner(v.VMScanners),
					},
				},
			}

			result = append(result, defenderCspmAwsOffering)
		} else if v, ok := item.(securityconnectors.DefenderCspmGcpOffering); ok {
			defenderCspmGcpOffering := Offering{
				Type: string(securityconnectors.OfferingTypeDefenderCspmGcp),
				DefenderCspmGcp: []DefenderCspmGcp{
					{
						CiemDiscovery:                      flattenDefenderCspmGcpCiemDiscovery(v.CiemDiscovery),
						DataSensitivityDiscovery:           flattenDefenderCspmGcpDataSensitivityDiscovery(v.DataSensitivityDiscovery),
						MdcContainersAgentlessDiscoveryK8s: flattenDefenderCspmGcpMdcContainersAgentlessDiscoveryK8s(v.MdcContainersAgentlessDiscoveryK8s),
						MdcContainersImageAssessment:       flattenDefenderCspmGcpMdcContainersImageAssessment(v.MdcContainersImageAssessment),
						VmScanner:                          flattenDefenderCspmGcpVmScanner(v.VMScanners),
					},
				},
			}

			result = append(result, defenderCspmGcpOffering)
		}
	}

	return result
}

func flattenDatabasesAwsArcAutoProvisioning(input *securityconnectors.DefenderFoDatabasesAwsOfferingArcAutoProvisioning) string {
	var arcAutoProvisioningCloudRoleArn string
	if input == nil || pointer.From(input.Enabled) == false {
		return arcAutoProvisioningCloudRoleArn
	}

	arcAutoProvisioningCloudRoleArn = pointer.From(input.CloudRoleArn)

	return arcAutoProvisioningCloudRoleArn
}

func flattenDatabasesDspmCloudRoleArn(input *securityconnectors.DefenderFoDatabasesAwsOfferingDatabasesDspm) string {
	var databasesDspmCloudRoleArn string
	if input == nil || pointer.From(input.Enabled) == false {
		return databasesDspmCloudRoleArn
	}

	databasesDspmCloudRoleArn = pointer.From(input.CloudRoleArn)

	return databasesDspmCloudRoleArn
}

func flattenRdsCloudRoleArn(input *securityconnectors.DefenderFoDatabasesAwsOfferingRds) string {
	var rdsCloudRoleArn string
	if input == nil || pointer.From(input.Enabled) == false {
		return rdsCloudRoleArn
	}

	rdsCloudRoleArn = pointer.From(input.CloudRoleArn)

	return rdsCloudRoleArn
}

func flattenDatabasesGcpGcpArcAutoProvisioning(input *securityconnectors.DefenderForDatabasesGcpOfferingArcAutoProvisioning) bool {
	if input == nil {
		return false
	}

	return pointer.From(input.Enabled)
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

func flattenContainersAwsMdcContainersAgentlessDiscoveryK8s(input *securityconnectors.DefenderForContainersAwsOfferingMdcContainersAgentlessDiscoveryK8s) string {
	var mdcContainersAgentlessDiscoveryK8sCloudRoleArn string
	if input == nil || pointer.From(input.Enabled) == false {
		return mdcContainersAgentlessDiscoveryK8sCloudRoleArn
	}

	mdcContainersAgentlessDiscoveryK8sCloudRoleArn = pointer.From(input.CloudRoleArn)

	return mdcContainersAgentlessDiscoveryK8sCloudRoleArn
}

func flattenContainersAwsMdcContainersImageAssessment(input *securityconnectors.DefenderForContainersAwsOfferingMdcContainersImageAssessment) string {
	var mdcContainersImageAssessmentCloudRoleArn string
	if input == nil || pointer.From(input.Enabled) == false {
		return mdcContainersImageAssessmentCloudRoleArn
	}

	mdcContainersImageAssessmentCloudRoleArn = pointer.From(input.CloudRoleArn)

	return mdcContainersImageAssessmentCloudRoleArn
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

func flattenServersAwsArcAutoProvisioning(input *securityconnectors.DefenderForServersAwsOfferingArcAutoProvisioning) string {
	var arcAutoProvisioningCloudRoleArn string
	if input == nil || pointer.From(input.Enabled) == false {
		return arcAutoProvisioningCloudRoleArn
	}

	arcAutoProvisioningCloudRoleArn = pointer.From(input.CloudRoleArn)

	return arcAutoProvisioningCloudRoleArn
}

func flattenServersAwsMdeAutoProvisioning(input *securityconnectors.DefenderForServersAwsOfferingMdeAutoProvisioning) bool {
	if input == nil {
		return false
	}

	return pointer.From(input.Enabled)
}

func flattenServersAwsVaAutoProvisioning(input *securityconnectors.DefenderForServersAwsOfferingVaAutoProvisioning) string {
	var vaAutoProvisioningConfigurationType string
	if input == nil || pointer.From(input.Enabled) == false {
		return vaAutoProvisioningConfigurationType
	}

	if configuration := input.Configuration; configuration != nil {
		vaAutoProvisioningConfigurationType = string(pointer.From(configuration.Type))
	}

	return vaAutoProvisioningConfigurationType
}

func flattenServersAwsVmScanner(input *securityconnectors.DefenderForServersAwsOfferingVMScanners) []ServersAwsVmScanner {
	result := make([]ServersAwsVmScanner, 0)
	if input == nil || pointer.From(input.Enabled) == false {
		return result
	}

	if v := input.Configuration; v != nil {
		serversAwsVmScanner := ServersAwsVmScanner{
			CloudRoleArn:  pointer.From(v.CloudRoleArn),
			ExclusionTags: pointer.From(v.ExclusionTags),
			ScanningMode:  string(pointer.From(v.ScanningMode)),
		}
		return append(result, serversAwsVmScanner)
	}

	return result
}

func flattenServersGcpArcAutoProvisioning(input *securityconnectors.DefenderForServersGcpOfferingArcAutoProvisioning) bool {
	if input == nil {
		return false
	}

	return pointer.From(input.Enabled)
}

func flattenServersGcpDefenderForServers(input *securityconnectors.DefenderForServersGcpOfferingDefenderForServers) []ServersGcpDefenderForServers {
	result := make([]ServersGcpDefenderForServers, 0)
	if input == nil {
		return result
	}

	serversGcpDefenderForServers := ServersGcpDefenderForServers{
		ServiceAccountEmailAddress: pointer.From(input.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.From(input.WorkloadIdentityProviderId),
	}

	return append(result, serversGcpDefenderForServers)
}

func flattenServersGcpMdeAutoProvisioning(input *securityconnectors.DefenderForServersGcpOfferingMdeAutoProvisioning) bool {
	if input == nil {
		return false
	}

	return pointer.From(input.Enabled)
}

func flattenServersGcpVaAutoProvisioning(input *securityconnectors.DefenderForServersGcpOfferingVaAutoProvisioning) string {
	var vaAutoProvisioningConfigurationType string
	if input == nil || pointer.From(input.Enabled) == false {
		return vaAutoProvisioningConfigurationType
	}

	if configuration := input.Configuration; configuration != nil {
		vaAutoProvisioningConfigurationType = string(pointer.From(configuration.Type))
	}

	return vaAutoProvisioningConfigurationType
}

func flattenServersGcpVmScanner(input *securityconnectors.DefenderForServersGcpOfferingVMScanners) []ServersGcpVmScanner {
	result := make([]ServersGcpVmScanner, 0)
	if input == nil || pointer.From(input.Enabled) == false {
		return result
	}

	if v := input.Configuration; v != nil {
		serversGcpVmScanner := ServersGcpVmScanner{
			ExclusionTags: pointer.From(v.ExclusionTags),
			ScanningMode:  string(pointer.From(v.ScanningMode)),
		}
		return append(result, serversGcpVmScanner)
	}

	return result
}

func flattenDefenderCspmAwsCiem(input *securityconnectors.DefenderCspmAwsOfferingCiem) []CspmAwsCiem {
	result := make([]CspmAwsCiem, 0)
	if input == nil {
		return result
	}

	cspmAwsCiem := CspmAwsCiem{
		CiemOidc: flattenDefenderCspmAwsCiemOidc(input.CiemOidc),
	}

	if v := input.CiemDiscovery; v != nil {
		cspmAwsCiem.CiemDiscoveryCloudRoleArn = pointer.From(v.CloudRoleArn)
	}

	return append(result, cspmAwsCiem)
}

func flattenDefenderCspmAwsCiemOidc(input *securityconnectors.DefenderCspmAwsOfferingCiemCiemOidc) []CspmAwsCiemOidc {
	result := make([]CspmAwsCiemOidc, 0)
	if input == nil {
		return result
	}

	cspmAwsCiemOidc := CspmAwsCiemOidc{
		AzureActiveDirectoryAppName: pointer.From(input.AzureActiveDirectoryAppName),
		CloudRoleArn:                pointer.From(input.CloudRoleArn),
	}

	return append(result, cspmAwsCiemOidc)
}

func flattenDefenderCspmAwsDataSensitivityDiscovery(input *securityconnectors.DefenderCspmAwsOfferingDataSensitivityDiscovery) string {
	var dataSensitivityDiscoveryCloudRoleArn string
	if input == nil || pointer.From(input.Enabled) == false {
		return dataSensitivityDiscoveryCloudRoleArn
	}

	dataSensitivityDiscoveryCloudRoleArn = pointer.From(input.CloudRoleArn)

	return dataSensitivityDiscoveryCloudRoleArn
}

func flattenDefenderCspmAwsDatabasesDspm(input *securityconnectors.DefenderCspmAwsOfferingDatabasesDspm) string {
	var databasesDspmCloudRoleArn string
	if input == nil || pointer.From(input.Enabled) == false {
		return databasesDspmCloudRoleArn
	}

	databasesDspmCloudRoleArn = pointer.From(input.CloudRoleArn)

	return databasesDspmCloudRoleArn
}

func flattenDefenderCspmAwsMdcContainersAgentlessDiscoveryK8s(input *securityconnectors.DefenderCspmAwsOfferingMdcContainersAgentlessDiscoveryK8s) string {
	var mdcContainersAgentlessDiscoveryK8s string
	if input == nil || pointer.From(input.Enabled) == false {
		return mdcContainersAgentlessDiscoveryK8s
	}

	mdcContainersAgentlessDiscoveryK8s = pointer.From(input.CloudRoleArn)

	return mdcContainersAgentlessDiscoveryK8s
}

func flattenDefenderCspmAwsMdcContainersImageAssessment(input *securityconnectors.DefenderCspmAwsOfferingMdcContainersImageAssessment) string {
	var mdcContainersImageAssessmentCloudRoleArn string
	if input == nil || pointer.From(input.Enabled) == false {
		return mdcContainersImageAssessmentCloudRoleArn
	}

	mdcContainersImageAssessmentCloudRoleArn = pointer.From(input.CloudRoleArn)

	return mdcContainersImageAssessmentCloudRoleArn
}

func flattenDefenderCspmAwsVmScanner(input *securityconnectors.DefenderCspmAwsOfferingVMScanners) []CspmAwsVmScanner {
	result := make([]CspmAwsVmScanner, 0)
	if input == nil || pointer.From(input.Enabled) == false {
		return result
	}

	if v := input.Configuration; v != nil {
		cspmAwsVmScanner := CspmAwsVmScanner{
			CloudRoleArn:  pointer.From(v.CloudRoleArn),
			ExclusionTags: pointer.From(v.ExclusionTags),
			ScanningMode:  string(pointer.From(v.ScanningMode)),
		}

		return append(result, cspmAwsVmScanner)
	}

	return result
}

func flattenDefenderCspmGcpCiemDiscovery(input *securityconnectors.DefenderCspmGcpOfferingCiemDiscovery) []CspmGcpCiemDiscovery {
	result := make([]CspmGcpCiemDiscovery, 0)
	if input == nil {
		return result
	}

	cspmGcpCiemDiscovery := CspmGcpCiemDiscovery{
		AzureActiveDirectoryAppName: pointer.From(input.AzureActiveDirectoryAppName),
		ServiceAccountEmailAddress:  pointer.From(input.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId:  pointer.From(input.WorkloadIdentityProviderId),
	}

	return append(result, cspmGcpCiemDiscovery)
}

func flattenDefenderCspmGcpDataSensitivityDiscovery(input *securityconnectors.DefenderCspmGcpOfferingDataSensitivityDiscovery) []CspmGcpDataSensitivityDiscovery {
	result := make([]CspmGcpDataSensitivityDiscovery, 0)
	if input == nil || pointer.From(input.Enabled) == false {
		return result
	}

	cspmGcpDataSensitivityDiscovery := CspmGcpDataSensitivityDiscovery{
		ServiceAccountEmailAddress: pointer.From(input.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.From(input.WorkloadIdentityProviderId),
	}

	return append(result, cspmGcpDataSensitivityDiscovery)
}

func flattenDefenderCspmGcpMdcContainersAgentlessDiscoveryK8s(input *securityconnectors.DefenderCspmGcpOfferingMdcContainersAgentlessDiscoveryK8s) []CspmGcpMdcContainersAgentlessDiscoveryK8s {
	result := make([]CspmGcpMdcContainersAgentlessDiscoveryK8s, 0)
	if input == nil || pointer.From(input.Enabled) == false {
		return result
	}

	cspmGcpMdcContainersAgentlessDiscoveryK8s := CspmGcpMdcContainersAgentlessDiscoveryK8s{
		ServiceAccountEmailAddress: pointer.From(input.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.From(input.WorkloadIdentityProviderId),
	}

	return append(result, cspmGcpMdcContainersAgentlessDiscoveryK8s)
}

func flattenDefenderCspmGcpMdcContainersImageAssessment(input *securityconnectors.DefenderCspmGcpOfferingMdcContainersImageAssessment) []CspmGcpMdcContainersImageAssessment {
	result := make([]CspmGcpMdcContainersImageAssessment, 0)
	if input == nil || pointer.From(input.Enabled) == false {
		return result
	}

	cspmGcpMdcContainersImageAssessment := CspmGcpMdcContainersImageAssessment{
		ServiceAccountEmailAddress: pointer.From(input.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.From(input.WorkloadIdentityProviderId),
	}

	return append(result, cspmGcpMdcContainersImageAssessment)
}

func flattenDefenderCspmGcpVmScanner(input *securityconnectors.DefenderCspmGcpOfferingVMScanners) []CspmGcpVmScanner {
	result := make([]CspmGcpVmScanner, 0)
	if input == nil || pointer.From(input.Enabled) == false {
		return result
	}

	if v := input.Configuration; v != nil {
		cspmGcpVmScanner := CspmGcpVmScanner{
			ExclusionTags: pointer.From(v.ExclusionTags),
			ScanningMode:  string(pointer.From(v.ScanningMode)),
		}
		return append(result, cspmGcpVmScanner)
	}

	return result
}
