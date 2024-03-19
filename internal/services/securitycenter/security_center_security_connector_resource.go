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
	Name                string               `tfschema:"name"`
	ResourceGroupName   string               `tfschema:"resource_group_name"`
	Location            string               `tfschema:"location"`
	EnvironmentName     string               `tfschema:"environment_name"`
	HierarchyIdentifier string               `tfschema:"hierarchy_identifier"`
	AwsEnvironmentData  []AwsEnvironmentData `tfschema:"aws_environment_data"`
	Tags                map[string]string    `tfschema:"tags"`
}

type AwsEnvironmentData struct {
	OrganizationalDataMaster                  []OrganizationalDataMaster `tfschema:"organizational_data_master"`
	OrganizationalDataMemberParentHierarchyId string                     `tfschema:"organizational_data_member_parent_hierarchy_id"`
	Regions                                   []string                   `tfschema:"regions"`
	ScanInterval                              int                        `tfschema:"scan_interval"`
}

type OrganizationalDataMaster struct {
	StacksetName       string   `tfschema:"stackset_name"`
	ExcludedAccountIds []string `tfschema:"excluded_account_ids"`
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
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"stackset_name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"excluded_account_ids": {
									Type:     pluginsdk.TypeList,
									Optional: true,
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
						ValidateFunc: validation.IntBetween(1, 24),
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
				},
				Tags: pointer.To(model.Tags),
			}

			if model.EnvironmentName == string(securityconnectors.CloudNameAzureDevOps) {
				parameters.Properties.EnvironmentData = securityconnectors.AzureDevOpsScopeEnvironmentData{}
			} else if model.EnvironmentName == string(securityconnectors.CloudNameGithub) {
				parameters.Properties.EnvironmentData = securityconnectors.GithubScopeEnvironmentData{}
			} else if model.EnvironmentName == string(securityconnectors.CloudNameGitLab) {
				parameters.Properties.EnvironmentData = securityconnectors.GitlabScopeEnvironmentData{}
			} else if model.EnvironmentName == string(securityconnectors.CloudNameAWS) {
				parameters.Properties.EnvironmentData = expandAwsEnvironmentData(model.AwsEnvironmentData)
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

					if v, ok := props.EnvironmentData.(securityconnectors.AwsEnvironmentData); ok {
						state.AwsEnvironmentData = flattenAwsEnvironmentData(v)
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

			if metadata.ResourceData.HasChange("environment_name") {
				parameters.Properties.EnvironmentName = pointer.To(securityconnectors.CloudName(model.EnvironmentName))
			}

			if metadata.ResourceData.HasChange("hierarchy_identifier") {
				parameters.Properties.HierarchyIdentifier = pointer.To(model.HierarchyIdentifier)
			}

			if model.EnvironmentName == string(securityconnectors.CloudNameAzureDevOps) {
				parameters.Properties.EnvironmentData = securityconnectors.AzureDevOpsScopeEnvironmentData{}
			} else if model.EnvironmentName == string(securityconnectors.CloudNameGithub) {
				parameters.Properties.EnvironmentData = securityconnectors.GithubScopeEnvironmentData{}
			} else if model.EnvironmentName == string(securityconnectors.CloudNameGitLab) {
				parameters.Properties.EnvironmentData = securityconnectors.GitlabScopeEnvironmentData{}
			} else if model.EnvironmentName == string(securityconnectors.CloudNameAWS) {
				parameters.Properties.EnvironmentData = expandAwsEnvironmentData(model.AwsEnvironmentData)
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
		Regions: pointer.To(awsEnvironmentData.Regions),
	}

	if v := awsEnvironmentData.OrganizationalDataMaster; v != nil {
		result.OrganizationalData = expandAwsOrganizationalDataMaster(v)
	} else if v := awsEnvironmentData.OrganizationalDataMemberParentHierarchyId; v != "" {
		result.OrganizationalData = expandAwsOrganizationalDataMember(v)
	}

	if v := awsEnvironmentData.ScanInterval; v != 0 {
		result.ScanInterval = pointer.To(int64(awsEnvironmentData.ScanInterval))
	}

	return result
}

func expandAwsOrganizationalDataMaster(input []OrganizationalDataMaster) *securityconnectors.AwsOrganizationalDataMaster {
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

func flattenAwsOrganizationalDataMaster(input securityconnectors.AwsOrganizationalDataMaster) []OrganizationalDataMaster {
	result := make([]OrganizationalDataMaster, 0)

	awsOrganizationalDataMaster := OrganizationalDataMaster{
		ExcludedAccountIds: pointer.From(input.ExcludedAccountIds),
		StacksetName:       pointer.From(input.StacksetName),
	}

	return append(result, awsOrganizationalDataMaster)
}

func flattenAwsOrganizationalDataMember(input securityconnectors.AwsOrganizationalDataMember) string {
	organizationalDataMemberParentHierarchyId := pointer.From(input.ParentHierarchyId)

	return organizationalDataMemberParentHierarchyId
}
