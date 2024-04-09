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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/validate"
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
	ProjectId    string `tfschema:"project_id"`
	ScanInterval int    `tfschema:"scan_interval"`
}

type Offering struct {
	Type                                            string                                       `tfschema:"type"`
	CspmMonitorAwsNativeCloudConnectionCloudRoleArn string                                       `tfschema:"cspm_monitor_aws_native_cloud_connection_cloud_role_arn"`
	CspmMonitorGcp                                  []CspmMonitorGcp                             `tfschema:"cspm_monitor_gcp"`
	DefenderForDatabasesAws                         []DefenderForDatabasesAws                    `tfschema:"defender_for_databases_aws"`
	DefenderForDatabasesGcpArcAutoProvisioning      []DefenderForDatabasesGcpArcAutoProvisioning `tfschema:"defender_for_databases_gcp_arc_auto_provisioning"`
	DefenderForContainersAws                        []DefenderForContainersAws                   `tfschema:"defender_for_containers_aws"`
	DefenderForContainersGcp                        []DefenderForContainersGcp                   `tfschema:"defender_for_containers_gcp"`
	DefenderCspmAws                                 []DefenderCspmAws                            `tfschema:"defender_cspm_aws"`
	DefenderCspmGcp                                 []DefenderCspmGcp                            `tfschema:"defender_cspm_gcp"`
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

type DefenderForDatabasesGcpArcAutoProvisioning struct {
	ServiceAccountEmailAddress string `tfschema:"service_account_email_address"`
	WorkloadIdentityProviderId string `tfschema:"workload_identity_provider_id"`
}

type DefenderForContainersAws struct {
	AutoProvisioningEnabled                          bool   `tfschema:"auto_provisioning_enabled"`
	CloudWatchToKinesisCloudRoleArn                  string `tfschema:"cloud_watch_to_kinesis_cloud_role_arn"`
	ContainerVulnerabilityAssessmentCloudRoleArn     string `tfschema:"container_vulnerability_assessment_cloud_role_arn"`
	ContainerVulnerabilityAssessmentTaskCloudRoleArn string `tfschema:"container_vulnerability_assessment_task_cloud_role_arn"`
	KinesisToS3CloudRoleArn                          string `tfschema:"kinesis_to_s3_cloud_role_arn"`
	KubeAuditRetentionTime                           int    `tfschema:"kube_audit_retention_time"`
	KubernetesScubaReaderCloudRoleArn                string `tfschema:"kubernetes_scuba_reader_cloud_role_arn"`
	KubernetesServiceCloudRoleArn                    string `tfschema:"kubernetes_service_cloud_role_arn"`
	MdcContainersAgentlessDiscoveryK8sCloudRoleArn   string `tfschema:"mdc_containers_agentless_discovery_k8s_cloud_role_arn"`
	MdcContainersImageAssessmentCloudRoleArn         string `tfschema:"mdc_containers_image_assessment_cloud_role_arn"`
	ScubaExternalId                                  string `tfschema:"scuba_external_id"`
}

type DefenderForContainersGcp struct {
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
}

type DefenderCspmGcp struct {
	CiemDiscovery                      []CspmGcpCiemDiscovery                      `tfschema:"ciem_discovery"`
	DataSensitivityDiscovery           []CspmGcpDataSensitivityDiscovery           `tfschema:"data_sensitivity_discovery"`
	MdcContainersAgentlessDiscoveryK8s []CspmGcpMdcContainersAgentlessDiscoveryK8s `tfschema:"mdc_containers_agentless_discovery_k8s"`
	MdcContainersImageAssessment       []CspmGcpMdcContainersImageAssessment       `tfschema:"mdc_containers_image_assessment"`
	VmScannerExclusionTags             map[string]string                           `tfschema:"vm_scanner_exclusion_tags"`
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
			ValidateFunc: validation.StringLenBetween(1, 260),
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
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"project_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
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
									ValidateFunc: validate.SecurityConnectorServiceAccountEmailAddress,
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

					"defender_for_databases_gcp_arc_auto_provisioning": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"service_account_email_address": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validate.SecurityConnectorServiceAccountEmailAddress,
								},

								"workload_identity_provider_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
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

								"auto_provisioning_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  true,
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
												ValidateFunc: validate.SecurityConnectorServiceAccountEmailAddress,
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
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"service_account_email_address": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validate.SecurityConnectorServiceAccountEmailAddress,
											},

											"workload_identity_provider_id": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
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
												ValidateFunc: validate.SecurityConnectorServiceAccountEmailAddress,
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
												ValidateFunc: validate.SecurityConnectorServiceAccountEmailAddress,
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

								"databases_dspm_cloud_role_arn": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"data_sensitivity_discovery_cloud_role_arn": {
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
												ValidateFunc: validate.SecurityConnectorServiceAccountEmailAddress,
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
												ValidateFunc: validate.SecurityConnectorServiceAccountEmailAddress,
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
												ValidateFunc: validate.SecurityConnectorServiceAccountEmailAddress,
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
												ValidateFunc: validate.SecurityConnectorServiceAccountEmailAddress,
											},

											"workload_identity_provider_id": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
										},
									},
								},

								"vm_scanner_exclusion_tags": commonschema.Tags(),
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
					EnvironmentData:     expandEnvironmentData(model.EnvironmentData),
				},
				Tags: pointer.To(model.Tags),
			}

			offerings, err := expandOfferings(model.Offerings)
			if err != nil {
				return err
			}
			parameters.Properties.Offerings = offerings

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
				offerings, err := expandOfferings(model.Offerings)
				if err != nil {
					return err
				}
				parameters.Properties.Offerings = offerings
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
		ProjectDetails: &securityconnectors.GcpProjectDetails{
			ProjectId: pointer.To(gcpProject.ProjectId),
		},
		ScanInterval: pointer.To(int64(gcpProject.ScanInterval)),
	}

	return result
}

func expandOfferings(input []Offering) (*[]securityconnectors.CloudOffering, error) {
	result := make([]securityconnectors.CloudOffering, 0)
	if len(input) == 0 {
		return &result, nil
	}

	for _, item := range input {
		if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeCspmMonitorAzureDevOps) {
			if item.CspmMonitorAwsNativeCloudConnectionCloudRoleArn != "" || len(item.CspmMonitorGcp) != 0 || len(item.DefenderForDatabasesAws) != 0 || len(item.DefenderForDatabasesGcpArcAutoProvisioning) != 0 || len(item.DefenderForContainersAws) != 0 || len(item.DefenderForContainersGcp) != 0 || len(item.DefenderCspmAws) != 0 || len(item.DefenderCspmGcp) != 0 {
				return nil, fmt.Errorf("`cspm_monitor_aws_native_cloud_connection_cloud_role_arn`, `cspm_monitor_gcp`, `defender_for_databases_aws`, `defender_for_databases_gcp_arc_auto_provisioning`, `defender_for_containers_aws`, `defender_for_containers_gcp`, `defender_cspm_aws` and `defender_cspm_gcp` cannot be set for the offering type `AzureDevOps`")
			}

			result = append(result, securityconnectors.CspmMonitorAzureDevOpsOffering{})
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeCspmMonitorGithub) {
			if item.CspmMonitorAwsNativeCloudConnectionCloudRoleArn != "" || len(item.CspmMonitorGcp) != 0 || len(item.DefenderForDatabasesAws) != 0 || len(item.DefenderForDatabasesGcpArcAutoProvisioning) != 0 || len(item.DefenderForContainersAws) != 0 || len(item.DefenderForContainersGcp) != 0 || len(item.DefenderCspmAws) != 0 || len(item.DefenderCspmGcp) != 0 {
				return nil, fmt.Errorf("`cspm_monitor_aws_native_cloud_connection_cloud_role_arn`, `cspm_monitor_gcp`, `defender_for_databases_aws`, `defender_for_databases_gcp_arc_auto_provisioning`, `defender_for_containers_aws`, `defender_for_containers_gcp`, `defender_cspm_aws` and `defender_cspm_gcp` cannot be set for the offering type `Github`")
			}

			result = append(result, securityconnectors.CspmMonitorGithubOffering{})
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeCspmMonitorGitLab) {
			if item.CspmMonitorAwsNativeCloudConnectionCloudRoleArn != "" || len(item.CspmMonitorGcp) != 0 || len(item.DefenderForDatabasesAws) != 0 || len(item.DefenderForDatabasesGcpArcAutoProvisioning) != 0 || len(item.DefenderForContainersAws) != 0 || len(item.DefenderForContainersGcp) != 0 || len(item.DefenderCspmAws) != 0 || len(item.DefenderCspmGcp) != 0 {
				return nil, fmt.Errorf("`cspm_monitor_aws_native_cloud_connection_cloud_role_arn`, `cspm_monitor_gcp`, `defender_for_databases_aws`, `defender_for_databases_gcp_arc_auto_provisioning`, `defender_for_containers_aws`, `defender_for_containers_gcp`, `defender_cspm_aws` and `defender_cspm_gcp` cannot be set for the offering type `GitLab`")
			}

			result = append(result, securityconnectors.CspmMonitorGitLabOffering{})
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeCspmMonitorAws) {
			if len(item.CspmMonitorGcp) != 0 || len(item.DefenderForDatabasesAws) != 0 || len(item.DefenderForDatabasesGcpArcAutoProvisioning) != 0 || len(item.DefenderForContainersAws) != 0 || len(item.DefenderForContainersGcp) != 0 || len(item.DefenderCspmAws) != 0 || len(item.DefenderCspmGcp) != 0 {
				return nil, fmt.Errorf("`cspm_monitor_gcp`, `defender_for_databases_aws`, `defender_for_databases_gcp_arc_auto_provisioning`, `defender_for_containers_aws`, `defender_for_containers_gcp`, `defender_cspm_aws` and `defender_cspm_gcp` cannot be set for the offering type `CspmMonitorAws`")
			}

			result = append(result, securityconnectors.CspmMonitorAwsOffering{
				NativeCloudConnection: expandCspmMonitorAwsNativeCloudConnection(item.CspmMonitorAwsNativeCloudConnectionCloudRoleArn),
			})
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeCspmMonitorGcp) {
			if item.CspmMonitorAwsNativeCloudConnectionCloudRoleArn != "" || len(item.DefenderForDatabasesAws) != 0 || len(item.DefenderForDatabasesGcpArcAutoProvisioning) != 0 || len(item.DefenderForContainersAws) != 0 || len(item.DefenderForContainersGcp) != 0 || len(item.DefenderCspmAws) != 0 || len(item.DefenderCspmGcp) != 0 {
				return nil, fmt.Errorf("`cspm_monitor_aws_native_cloud_connection_cloud_role_arn`, `defender_for_databases_aws`, `defender_for_databases_gcp_arc_auto_provisioning`, `defender_for_containers_aws`, `defender_for_containers_gcp`, `defender_cspm_aws` and `defender_cspm_gcp` cannot be set for the offering type `CspmMonitorGcp`")
			}

			result = append(result, securityconnectors.CspmMonitorGcpOffering{
				NativeCloudConnection: expandCspmMonitorGcpNativeCloudConnection(item.CspmMonitorGcp),
			})
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeDefenderForDatabasesAws) {
			if item.CspmMonitorAwsNativeCloudConnectionCloudRoleArn != "" || len(item.CspmMonitorGcp) != 0 || len(item.DefenderForDatabasesGcpArcAutoProvisioning) != 0 || len(item.DefenderForContainersAws) != 0 || len(item.DefenderForContainersGcp) != 0 || len(item.DefenderCspmAws) != 0 || len(item.DefenderCspmGcp) != 0 {
				return nil, fmt.Errorf("`cspm_monitor_aws_native_cloud_connection_cloud_role_arn`, `cspm_monitor_gcp`, `defender_for_databases_gcp_arc_auto_provisioning`, `defender_for_containers_aws`, `defender_for_containers_gcp`, `defender_cspm_aws` and `defender_cspm_gcp` cannot be set for the offering type `DefenderForDatabasesAws`")
			}

			defenderFoDatabasesAwsOffering := securityconnectors.DefenderFoDatabasesAwsOffering{}

			if v := item.DefenderForDatabasesAws; len(v) != 0 {
				defenderForDatabasesAws := v[0]

				defenderFoDatabasesAwsOffering.ArcAutoProvisioning = expandDatabasesAwsArcAutoProvisioning(defenderForDatabasesAws.ArcAutoProvisioningCloudRoleArn)
				defenderFoDatabasesAwsOffering.DatabasesDspm = expandDatabasesAwsDspm(defenderForDatabasesAws.DatabasesDspmCloudRoleArn)
				defenderFoDatabasesAwsOffering.Rds = expandDatabasesAwsRds(defenderForDatabasesAws.RdsCloudRoleArn)
			}

			result = append(result, defenderFoDatabasesAwsOffering)
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeDefenderForDatabasesGcp) {
			if item.CspmMonitorAwsNativeCloudConnectionCloudRoleArn != "" || len(item.CspmMonitorGcp) != 0 || len(item.DefenderForDatabasesAws) != 0 || len(item.DefenderForContainersAws) != 0 || len(item.DefenderForContainersGcp) != 0 || len(item.DefenderCspmAws) != 0 || len(item.DefenderCspmGcp) != 0 {
				return nil, fmt.Errorf("`cspm_monitor_aws_native_cloud_connection_cloud_role_arn`, `cspm_monitor_gcp`, `defender_for_databases_aws`, `defender_for_containers_aws`, `defender_for_containers_gcp`, `defender_cspm_aws` and `defender_cspm_gcp` cannot be set for the offering type `DefenderForDatabasesGcp`")
			}

			defenderForDatabasesGcpOffering := securityconnectors.DefenderForDatabasesGcpOffering{}

			if v := item.DefenderForDatabasesGcpArcAutoProvisioning; len(v) != 0 {
				defenderForDatabasesGcp := v[0]

				defenderForDatabasesGcpOffering.ArcAutoProvisioning = expandDefenderForDatabasesGcpArcAutoProvisioning()
				defenderForDatabasesGcpOffering.DefenderForDatabasesArcAutoProvisioning = expandDefenderForDatabasesGcpDefenderForDatabasesArcAutoProvisioning(defenderForDatabasesGcp)
			}

			result = append(result, defenderForDatabasesGcpOffering)
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeDefenderForContainersAws) {
			if item.CspmMonitorAwsNativeCloudConnectionCloudRoleArn != "" || len(item.CspmMonitorGcp) != 0 || len(item.DefenderForDatabasesAws) != 0 || len(item.DefenderForDatabasesGcpArcAutoProvisioning) != 0 || len(item.DefenderForContainersGcp) != 0 || len(item.DefenderCspmAws) != 0 || len(item.DefenderCspmGcp) != 0 {
				return nil, fmt.Errorf("`cspm_monitor_aws_native_cloud_connection_cloud_role_arn`, `cspm_monitor_gcp`, `defender_for_databases_aws`, `defender_for_databases_gcp_arc_auto_provisioning`, `defender_for_containers_gcp`, `defender_cspm_aws` and `defender_cspm_gcp` cannot be set for the offering type `DefenderForContainersAws`")
			}

			defenderForContainersAwsOffering := securityconnectors.DefenderForContainersAwsOffering{}

			if v := item.DefenderForContainersAws; len(v) != 0 {
				defenderForContainersAws := v[0]

				defenderForContainersAwsOffering.AutoProvisioning = pointer.To(defenderForContainersAws.AutoProvisioningEnabled)
				defenderForContainersAwsOffering.MdcContainersAgentlessDiscoveryK8s = expandContainersAwsMdcContainersAgentlessDiscoveryK8s(defenderForContainersAws.MdcContainersAgentlessDiscoveryK8sCloudRoleArn)
				defenderForContainersAwsOffering.MdcContainersImageAssessment = expandContainersAwsMdcContainersImageAssessment(defenderForContainersAws.MdcContainersImageAssessmentCloudRoleArn)
				defenderForContainersAwsOffering.ScubaExternalId = pointer.To(defenderForContainersAws.ScubaExternalId)
				defenderForContainersAwsOffering.CloudWatchToKinesis = expandDefenderForContainersAwsCloudWatchToKinesis(defenderForContainersAws.CloudWatchToKinesisCloudRoleArn)
				defenderForContainersAwsOffering.EnableContainerVulnerabilityAssessment = expandDefenderForContainersAwsEnableContainerVulnerabilityAssessment(defenderForContainersAws.ContainerVulnerabilityAssessmentCloudRoleArn, defenderForContainersAws.ContainerVulnerabilityAssessmentTaskCloudRoleArn)
				defenderForContainersAwsOffering.ContainerVulnerabilityAssessment = expandDefenderForContainersAwsContainerVulnerabilityAssessment(defenderForContainersAws.ContainerVulnerabilityAssessmentCloudRoleArn)
				defenderForContainersAwsOffering.ContainerVulnerabilityAssessmentTask = expandDefenderForContainersAwsContainerVulnerabilityAssessmentTask(defenderForContainersAws.ContainerVulnerabilityAssessmentTaskCloudRoleArn)
				defenderForContainersAwsOffering.KinesisToS3 = expandDefenderForContainersAwsKinesisToS3(defenderForContainersAws.KinesisToS3CloudRoleArn)
				defenderForContainersAwsOffering.KubeAuditRetentionTime = pointer.To(int64(defenderForContainersAws.KubeAuditRetentionTime))
				defenderForContainersAwsOffering.KubernetesScubaReader = expandDefenderForContainersAwsKubernetesScubaReader(defenderForContainersAws.KubernetesScubaReaderCloudRoleArn)
				defenderForContainersAwsOffering.KubernetesService = expanddefenderForContainersAwsKubernetesService(defenderForContainersAws.KubernetesServiceCloudRoleArn)
			}

			result = append(result, defenderForContainersAwsOffering)
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeDefenderForContainersGcp) {
			if item.CspmMonitorAwsNativeCloudConnectionCloudRoleArn != "" || len(item.CspmMonitorGcp) != 0 || len(item.DefenderForDatabasesAws) != 0 || len(item.DefenderForDatabasesGcpArcAutoProvisioning) != 0 || len(item.DefenderForContainersAws) != 0 || len(item.DefenderCspmAws) != 0 || len(item.DefenderCspmGcp) != 0 {
				return nil, fmt.Errorf("`cspm_monitor_aws_native_cloud_connection_cloud_role_arn`, `cspm_monitor_gcp`, `defender_for_databases_aws`, `defender_for_databases_gcp_arc_auto_provisioning`, `defender_for_containers_aws`, `defender_cspm_aws` and `defender_cspm_gcp` cannot be set for the offering type `DefenderForContainersGcp`")
			}

			defenderForContainersGcpOffering := securityconnectors.DefenderForContainersGcpOffering{}

			if v := item.DefenderForContainersGcp; len(v) != 0 {
				defenderForContainersGcp := v[0]

				if len(defenderForContainersGcp.DataPipelineNativeCloudConnection) != 0 {
					defenderForContainersGcpOffering.AuditLogsAutoProvisioningFlag = pointer.To(true)
				}
				defenderForContainersGcpOffering.DataPipelineNativeCloudConnection = expandContainersGcpDataPipelineNativeCloudConnection(defenderForContainersGcp.DataPipelineNativeCloudConnection)
				defenderForContainersGcpOffering.DefenderAgentAutoProvisioningFlag = pointer.To(defenderForContainersGcp.DefenderAgentAutoProvisioningFlagEnabled)
				defenderForContainersGcpOffering.MdcContainersAgentlessDiscoveryK8s = expandContainersGcpMdcContainersAgentlessDiscoveryK8s(defenderForContainersGcp.MdcContainersAgentlessDiscoveryK8s)
				defenderForContainersGcpOffering.MdcContainersImageAssessment = expandContainersGcpMdcContainersImageAssessment(defenderForContainersGcp.MdcContainersImageAssessment)
				defenderForContainersGcpOffering.NativeCloudConnection = expandContainersGcpNativeCloudConnection(defenderForContainersGcp.NativeCloudConnection)
				defenderForContainersGcpOffering.PolicyAgentAutoProvisioningFlag = pointer.To(defenderForContainersGcp.PolicyAgentAutoProvisioningFlagEnabled)
			}

			result = append(result, defenderForContainersGcpOffering)
		} else if offeringType := item.Type; offeringType == string(securityconnectors.OfferingTypeDefenderCspmAws) {
			if item.CspmMonitorAwsNativeCloudConnectionCloudRoleArn != "" || len(item.CspmMonitorGcp) != 0 || len(item.DefenderForDatabasesAws) != 0 || len(item.DefenderForDatabasesGcpArcAutoProvisioning) != 0 || len(item.DefenderForContainersAws) != 0 || len(item.DefenderForContainersGcp) != 0 || len(item.DefenderCspmGcp) != 0 {
				return nil, fmt.Errorf("`cspm_monitor_aws_native_cloud_connection_cloud_role_arn`, `cspm_monitor_gcp`, `defender_for_databases_aws`, `defender_for_databases_gcp_arc_auto_provisioning`, `defender_for_containers_aws`, `defender_for_containers_gcp` and `defender_cspm_gcp` cannot be set for the offering type `DefenderCspmAws`")
			}

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
			if item.CspmMonitorAwsNativeCloudConnectionCloudRoleArn != "" || len(item.CspmMonitorGcp) != 0 || len(item.DefenderForDatabasesAws) != 0 || len(item.DefenderForDatabasesGcpArcAutoProvisioning) != 0 || len(item.DefenderForContainersAws) != 0 || len(item.DefenderForContainersGcp) != 0 || len(item.DefenderCspmAws) != 0 {
				return nil, fmt.Errorf("`cspm_monitor_aws_native_cloud_connection_cloud_role_arn`, `cspm_monitor_gcp`, `defender_for_databases_aws`, `defender_for_databases_gcp_arc_auto_provisioning`, `defender_for_containers_aws`, `defender_for_containers_gcp` and `defender_cspm_aws` cannot be set for the offering type `DefenderCspmGcp`")
			}

			defenderCspmGcpOffering := securityconnectors.DefenderCspmGcpOffering{}

			if v := item.DefenderCspmGcp; len(v) != 0 {
				defenderCspmGcp := v[0]

				defenderCspmGcpOffering.CiemDiscovery = expandCspmGcpCiemDiscovery(defenderCspmGcp.CiemDiscovery)
				defenderCspmGcpOffering.DataSensitivityDiscovery = expandCspmGcpDataSensitivityDiscovery(defenderCspmGcp.DataSensitivityDiscovery)
				defenderCspmGcpOffering.MdcContainersAgentlessDiscoveryK8s = expandCspmGcpMdcContainersAgentlessDiscoveryK8s(defenderCspmGcp.MdcContainersAgentlessDiscoveryK8s)
				defenderCspmGcpOffering.MdcContainersImageAssessment = expandCspmGcpMdcContainersImageAssessment(defenderCspmGcp.MdcContainersImageAssessment)
				defenderCspmGcpOffering.VMScanners = expandCspmGcpVMScanners(defenderCspmGcp.VmScannerExclusionTags)
			}

			result = append(result, defenderCspmGcpOffering)
		}
	}

	return &result, nil
}

func expandCspmMonitorAwsNativeCloudConnection(input string) *securityconnectors.CspmMonitorAwsOfferingNativeCloudConnection {
	if input == "" {
		return nil
	}

	return &securityconnectors.CspmMonitorAwsOfferingNativeCloudConnection{
		CloudRoleArn: pointer.To(input),
	}
}

func expandCspmMonitorGcpNativeCloudConnection(input []CspmMonitorGcp) *securityconnectors.CspmMonitorGcpOfferingNativeCloudConnection {
	if len(input) == 0 {
		return nil
	}

	cspmMonitorGcp := input[0]

	return &securityconnectors.CspmMonitorGcpOfferingNativeCloudConnection{
		ServiceAccountEmailAddress: pointer.To(cspmMonitorGcp.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.To(cspmMonitorGcp.WorkloadIdentityProviderId),
	}
}

func expandDefenderForDatabasesGcpArcAutoProvisioning() *securityconnectors.DefenderForDatabasesGcpOfferingArcAutoProvisioning {
	return &securityconnectors.DefenderForDatabasesGcpOfferingArcAutoProvisioning{
		Configuration: &securityconnectors.DefenderForDatabasesGcpOfferingArcAutoProvisioningConfiguration{},
		Enabled:       pointer.To(true),
	}
}

func expandDefenderForDatabasesGcpDefenderForDatabasesArcAutoProvisioning(input DefenderForDatabasesGcpArcAutoProvisioning) *securityconnectors.DefenderForDatabasesGcpOfferingDefenderForDatabasesArcAutoProvisioning {
	return &securityconnectors.DefenderForDatabasesGcpOfferingDefenderForDatabasesArcAutoProvisioning{
		ServiceAccountEmailAddress: pointer.To(input.ServiceAccountEmailAddress),
		WorkloadIdentityProviderId: pointer.To(input.WorkloadIdentityProviderId),
	}
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

func expandDefenderForContainersAwsCloudWatchToKinesis(input string) *securityconnectors.DefenderForContainersAwsOfferingCloudWatchToKinesis {
	return &securityconnectors.DefenderForContainersAwsOfferingCloudWatchToKinesis{
		CloudRoleArn: pointer.To(input),
	}
}

func expandDefenderForContainersAwsEnableContainerVulnerabilityAssessment(containerVulnerabilityAssessmentCloudRoleArn, containerVulnerabilityAssessmentTaskCloudRoleArn string) *bool {
	if containerVulnerabilityAssessmentCloudRoleArn != "" || containerVulnerabilityAssessmentTaskCloudRoleArn != "" {
		return pointer.To(true)
	}

	return pointer.To(false)
}

func expandDefenderForContainersAwsContainerVulnerabilityAssessment(input string) *securityconnectors.DefenderForContainersAwsOfferingContainerVulnerabilityAssessment {
	return &securityconnectors.DefenderForContainersAwsOfferingContainerVulnerabilityAssessment{
		CloudRoleArn: pointer.To(input),
	}
}

func expandDefenderForContainersAwsContainerVulnerabilityAssessmentTask(input string) *securityconnectors.DefenderForContainersAwsOfferingContainerVulnerabilityAssessmentTask {
	return &securityconnectors.DefenderForContainersAwsOfferingContainerVulnerabilityAssessmentTask{
		CloudRoleArn: pointer.To(input),
	}
}

func expandDefenderForContainersAwsKinesisToS3(input string) *securityconnectors.DefenderForContainersAwsOfferingKinesisToS3 {
	return &securityconnectors.DefenderForContainersAwsOfferingKinesisToS3{
		CloudRoleArn: pointer.To(input),
	}
}

func expandDefenderForContainersAwsKubernetesScubaReader(input string) *securityconnectors.DefenderForContainersAwsOfferingKubernetesScubaReader {
	return &securityconnectors.DefenderForContainersAwsOfferingKubernetesScubaReader{
		CloudRoleArn: pointer.To(input),
	}
}

func expanddefenderForContainersAwsKubernetesService(input string) *securityconnectors.DefenderForContainersAwsOfferingKubernetesService {
	return &securityconnectors.DefenderForContainersAwsOfferingKubernetesService{
		CloudRoleArn: pointer.To(input),
	}
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
			ScanningMode:  pointer.To(securityconnectors.ScanningModeDefault),
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

func expandCspmGcpVMScanners(input map[string]string) *securityconnectors.DefenderCspmGcpOfferingVMScanners {
	if len(input) == 0 {
		return &securityconnectors.DefenderCspmGcpOfferingVMScanners{
			Enabled: pointer.To(false),
		}
	}

	result := &securityconnectors.DefenderCspmGcpOfferingVMScanners{
		Configuration: &securityconnectors.DefenderCspmGcpOfferingVMScannersConfiguration{
			ExclusionTags: pointer.To(input),
			ScanningMode:  pointer.To(securityconnectors.ScanningModeDefault),
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
		ScanInterval: int(pointer.From(input.ScanInterval)),
		ProjectId:    flattenGcpProjectDetails(input.ProjectDetails),
	}

	return append(result, gcpProject)
}

func flattenGcpProjectDetails(input *securityconnectors.GcpProjectDetails) string {
	var projectId string
	if input == nil {
		return projectId
	}

	projectId = pointer.From(input.ProjectId)

	return projectId
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
				CspmMonitorAwsNativeCloudConnectionCloudRoleArn: flattenCspmMonitorAwsNativeCloudConnectionCloudRoleArn(v.NativeCloudConnection),
			}

			result = append(result, cspmMonitorAwsOffering)
		} else if v, ok := item.(securityconnectors.CspmMonitorGcpOffering); ok {
			cspmMonitorGcpOffering := Offering{
				Type:           string(securityconnectors.OfferingTypeCspmMonitorGcp),
				CspmMonitorGcp: flattenCspmMonitorGcp(v),
			}

			result = append(result, cspmMonitorGcpOffering)
		} else if v, ok := item.(securityconnectors.DefenderFoDatabasesAwsOffering); ok {
			defenderFoDatabasesAwsOffering := Offering{
				Type:                    string(securityconnectors.OfferingTypeDefenderForDatabasesAws),
				DefenderForDatabasesAws: flattenDefenerForDatabasesAws(v),
			}

			result = append(result, defenderFoDatabasesAwsOffering)
		} else if v, ok := item.(securityconnectors.DefenderForDatabasesGcpOffering); ok {
			defenderForDatabasesGcpOffering := Offering{
				Type: string(securityconnectors.OfferingTypeDefenderForDatabasesGcp),
				DefenderForDatabasesGcpArcAutoProvisioning: flattenDefenderForDatabasesGcpArcAutoProvisioning(v.DefenderForDatabasesArcAutoProvisioning),
			}

			result = append(result, defenderForDatabasesGcpOffering)
		} else if v, ok := item.(securityconnectors.DefenderForContainersAwsOffering); ok {
			defenderForContainersAwsOffering := Offering{
				Type:                     string(securityconnectors.OfferingTypeDefenderForContainersAws),
				DefenderForContainersAws: flattenDefenderForContainersAws(v),
			}

			result = append(result, defenderForContainersAwsOffering)
		} else if v, ok := item.(securityconnectors.DefenderForContainersGcpOffering); ok {
			defenderForContainersGcpOffering := Offering{
				Type:                     string(securityconnectors.OfferingTypeDefenderForContainersGcp),
				DefenderForContainersGcp: flattenDefenderForContainersGcp(v),
			}

			result = append(result, defenderForContainersGcpOffering)
		} else if v, ok := item.(securityconnectors.DefenderCspmAwsOffering); ok {
			defenderCspmAwsOffering := Offering{
				Type:            string(securityconnectors.OfferingTypeDefenderCspmAws),
				DefenderCspmAws: flattenDefenderCspmAws(v),
			}

			result = append(result, defenderCspmAwsOffering)
		} else if v, ok := item.(securityconnectors.DefenderCspmGcpOffering); ok {
			defenderCspmGcpOffering := Offering{
				Type:            string(securityconnectors.OfferingTypeDefenderCspmGcp),
				DefenderCspmGcp: flattenDefenderCspmGcp(v),
			}

			result = append(result, defenderCspmGcpOffering)
		}
	}

	return result
}

func flattenCspmMonitorGcp(input securityconnectors.CspmMonitorGcpOffering) []CspmMonitorGcp {
	result := make([]CspmMonitorGcp, 0)
	if input.NativeCloudConnection == nil {
		return result
	}

	cspmMonitorGcp := CspmMonitorGcp{}

	if v := input.NativeCloudConnection; v != nil {
		if serviceAccountEmailAddress := v.ServiceAccountEmailAddress; serviceAccountEmailAddress != nil {
			cspmMonitorGcp.ServiceAccountEmailAddress = pointer.From(serviceAccountEmailAddress)
		}

		if workloadIdentityProviderId := v.WorkloadIdentityProviderId; workloadIdentityProviderId != nil {
			cspmMonitorGcp.WorkloadIdentityProviderId = pointer.From(workloadIdentityProviderId)
		}
	}

	return append(result, cspmMonitorGcp)
}

func flattenCspmMonitorAwsNativeCloudConnectionCloudRoleArn(input *securityconnectors.CspmMonitorAwsOfferingNativeCloudConnection) string {
	var result string
	if input == nil {
		return result
	}

	return pointer.From(input.CloudRoleArn)
}

func flattenDefenerForDatabasesAws(input securityconnectors.DefenderFoDatabasesAwsOffering) []DefenderForDatabasesAws {
	result := make([]DefenderForDatabasesAws, 0)
	if (input.ArcAutoProvisioning == nil || pointer.From(input.ArcAutoProvisioning.Enabled) == false) && (input.DatabasesDspm == nil || pointer.From(input.DatabasesDspm.Enabled) == false) && (input.Rds == nil || pointer.From(input.Rds.Enabled) == false) {
		return result
	}

	defenderForDatabasesAws := DefenderForDatabasesAws{}

	if v := input.ArcAutoProvisioning; v != nil && pointer.From(v.Enabled) == true {
		defenderForDatabasesAws.ArcAutoProvisioningCloudRoleArn = pointer.From(v.CloudRoleArn)
	}

	if v := input.DatabasesDspm; v != nil && pointer.From(v.Enabled) == true {
		defenderForDatabasesAws.DatabasesDspmCloudRoleArn = pointer.From(v.CloudRoleArn)
	}

	if v := input.Rds; v != nil && pointer.From(v.Enabled) == true {
		defenderForDatabasesAws.RdsCloudRoleArn = pointer.From(v.CloudRoleArn)
	}

	return append(result, defenderForDatabasesAws)
}

func flattenDefenderForDatabasesGcpArcAutoProvisioning(input *securityconnectors.DefenderForDatabasesGcpOfferingDefenderForDatabasesArcAutoProvisioning) []DefenderForDatabasesGcpArcAutoProvisioning {
	result := make([]DefenderForDatabasesGcpArcAutoProvisioning, 0)
	if input == nil {
		return result
	}

	defenderForDatabasesGcpArcAutoProvisioning := DefenderForDatabasesGcpArcAutoProvisioning{}

	if v := input.ServiceAccountEmailAddress; v != nil {
		defenderForDatabasesGcpArcAutoProvisioning.ServiceAccountEmailAddress = pointer.From(v)
	}

	if v := input.WorkloadIdentityProviderId; v != nil {
		defenderForDatabasesGcpArcAutoProvisioning.WorkloadIdentityProviderId = pointer.From(v)
	}

	return append(result, defenderForDatabasesGcpArcAutoProvisioning)
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
	if input == nil || pointer.From(input.Enabled) == false {
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
	if input == nil || pointer.From(input.Enabled) == false {
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

func flattenDefenderForContainersAws(input securityconnectors.DefenderForContainersAwsOffering) []DefenderForContainersAws {
	result := make([]DefenderForContainersAws, 0)

	defenderForContainersAws := DefenderForContainersAws{
		AutoProvisioningEnabled:                        pointer.From(input.AutoProvisioning),
		KubeAuditRetentionTime:                         int(pointer.From(input.KubeAuditRetentionTime)),
		MdcContainersAgentlessDiscoveryK8sCloudRoleArn: flattenContainersAwsMdcContainersAgentlessDiscoveryK8s(input.MdcContainersAgentlessDiscoveryK8s),
		MdcContainersImageAssessmentCloudRoleArn:       flattenContainersAwsMdcContainersImageAssessment(input.MdcContainersImageAssessment),
		ScubaExternalId:                                pointer.From(input.ScubaExternalId),
	}

	if v := input.CloudWatchToKinesis; v != nil {
		defenderForContainersAws.CloudWatchToKinesisCloudRoleArn = pointer.From(v.CloudRoleArn)
	}

	if v := input.ContainerVulnerabilityAssessment; v != nil {
		defenderForContainersAws.ContainerVulnerabilityAssessmentCloudRoleArn = pointer.From(v.CloudRoleArn)
	}

	if v := input.ContainerVulnerabilityAssessmentTask; v != nil {
		defenderForContainersAws.ContainerVulnerabilityAssessmentTaskCloudRoleArn = pointer.From(v.CloudRoleArn)
	}

	if v := input.KinesisToS3; v != nil {
		defenderForContainersAws.KinesisToS3CloudRoleArn = pointer.From(v.CloudRoleArn)
	}

	if v := input.KubernetesScubaReader; v != nil {
		defenderForContainersAws.KubernetesScubaReaderCloudRoleArn = pointer.From(v.CloudRoleArn)
	}

	if v := input.KubernetesService; v != nil {
		defenderForContainersAws.KubernetesServiceCloudRoleArn = pointer.From(v.CloudRoleArn)
	}

	return append(result, defenderForContainersAws)
}

func flattenDefenderCspmAws(input securityconnectors.DefenderCspmAwsOffering) []DefenderCspmAws {
	result := make([]DefenderCspmAws, 0)
	if input.Ciem == nil && (input.DataSensitivityDiscovery == nil || pointer.From(input.DataSensitivityDiscovery.Enabled) == false) && (input.DatabasesDspm == nil || pointer.From(input.DatabasesDspm.Enabled) == false) && (input.MdcContainersAgentlessDiscoveryK8s == nil || pointer.From(input.MdcContainersAgentlessDiscoveryK8s.Enabled) == false) && (input.MdcContainersImageAssessment == nil || pointer.From(input.MdcContainersImageAssessment.Enabled) == false) && (input.VMScanners == nil || pointer.From(input.VMScanners.Enabled) == false) {
		return result
	}

	defenderCspmAws := DefenderCspmAws{
		Ciem:                                 flattenDefenderCspmAwsCiem(input.Ciem),
		DataSensitivityDiscoveryCloudRoleArn: flattenDefenderCspmAwsDataSensitivityDiscovery(input.DataSensitivityDiscovery),
		DatabasesDspmCloudRoleArn:            flattenDefenderCspmAwsDatabasesDspm(input.DatabasesDspm),
		MdcContainersAgentlessDiscoveryK8sCloudRoleArn: flattenDefenderCspmAwsMdcContainersAgentlessDiscoveryK8s(input.MdcContainersAgentlessDiscoveryK8s),
		MdcContainersImageAssessmentCloudRoleArn:       flattenDefenderCspmAwsMdcContainersImageAssessment(input.MdcContainersImageAssessment),
		VmScanner:                                      flattenDefenderCspmAwsVmScanner(input.VMScanners),
	}

	return append(result, defenderCspmAws)
}

func flattenDefenderForContainersGcp(input securityconnectors.DefenderForContainersGcpOffering) []DefenderForContainersGcp {
	result := make([]DefenderForContainersGcp, 0)

	defenderForContainersGcp := DefenderForContainersGcp{
		DataPipelineNativeCloudConnection:        flattenContainersGcpDataPipelineNativeCloudConnection(input.DataPipelineNativeCloudConnection),
		DefenderAgentAutoProvisioningFlagEnabled: pointer.From(input.DefenderAgentAutoProvisioningFlag),
		MdcContainersAgentlessDiscoveryK8s:       flattenContainersGcpMdcContainersAgentlessDiscoveryK8s(input.MdcContainersAgentlessDiscoveryK8s),
		MdcContainersImageAssessment:             flattenContainersGcpMdcContainersImageAssessment(input.MdcContainersImageAssessment),
		NativeCloudConnection:                    flattenContainersGcpNativeCloudConnection(input.NativeCloudConnection),
		PolicyAgentAutoProvisioningFlagEnabled:   pointer.From(input.PolicyAgentAutoProvisioningFlag),
	}

	return append(result, defenderForContainersGcp)
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
		}

		return append(result, cspmAwsVmScanner)
	}

	return result
}

func flattenDefenderCspmGcp(input securityconnectors.DefenderCspmGcpOffering) []DefenderCspmGcp {
	result := make([]DefenderCspmGcp, 0)
	if input.CiemDiscovery == nil && (input.DataSensitivityDiscovery == nil || pointer.From(input.DataSensitivityDiscovery.Enabled) == false) && (input.MdcContainersAgentlessDiscoveryK8s == nil || pointer.From(input.MdcContainersAgentlessDiscoveryK8s.Enabled) == false) && (input.MdcContainersImageAssessment == nil || pointer.From(input.MdcContainersImageAssessment.Enabled) == false) && (input.VMScanners == nil || pointer.From(input.VMScanners.Enabled) == false) {
		return result
	}

	defenderCspmGcp := DefenderCspmGcp{
		CiemDiscovery:                      flattenDefenderCspmGcpCiemDiscovery(input.CiemDiscovery),
		DataSensitivityDiscovery:           flattenDefenderCspmGcpDataSensitivityDiscovery(input.DataSensitivityDiscovery),
		MdcContainersAgentlessDiscoveryK8s: flattenDefenderCspmGcpMdcContainersAgentlessDiscoveryK8s(input.MdcContainersAgentlessDiscoveryK8s),
		MdcContainersImageAssessment:       flattenDefenderCspmGcpMdcContainersImageAssessment(input.MdcContainersImageAssessment),
		VmScannerExclusionTags:             flattenDefenderCspmGcpVmScanner(input.VMScanners),
	}

	return append(result, defenderCspmGcp)
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

func flattenDefenderCspmGcpVmScanner(input *securityconnectors.DefenderCspmGcpOfferingVMScanners) map[string]string {
	var result map[string]string
	if input == nil || pointer.From(input.Enabled) == false {
		return result
	}

	if v := input.Configuration; v != nil {
		return pointer.From(v.ExclusionTags)
	}

	return result
}
