package workloads

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/monitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/providerinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkloadsSAPMonitorProviderInstanceModel struct {
	Name                                string                                `tfschema:"name"`
	MonitorId                           string                                `tfschema:"monitor_id"`
	Identity                            []identity.ModelUserAssigned          `tfschema:"identity"`
	PrometheusOSProviderSettings        []PrometheusOSProviderSettings        `tfschema:"prometheus_os_provider_settings"`
	PrometheusHAClusterProviderSettings []PrometheusHAClusterProviderSettings `tfschema:"prometheus_ha_cluster_provider_settings"`
	MssqlServerProviderSettings         []MssqlServerProviderSettings         `tfschema:"mssql_server_provider_settings"`
	Db2ProviderSettings                 []Db2ProviderSettings                 `tfschema:"db2_provider_settings"`
	SAPHanaProviderSettings             []SAPHanaProviderSettings             `tfschema:"sap_hana_provider_settings"`
	SAPNetWeaverProviderSettings        []SAPNetWeaverProviderSettings        `tfschema:"sap_net_weaver_provider_settings"`
}

type PrometheusOSProviderSettings struct {
	PrometheusUrl     string `tfschema:"prometheus_url"`
	SID               string `tfschema:"sid"`
	SslCertificateUri string `tfschema:"ssl_certificate_uri"`
	SslPreference     string `tfschema:"ssl_preference"`
}

type PrometheusHAClusterProviderSettings struct {
	ClusterName       string `tfschema:"cluster_name"`
	HostName          string `tfschema:"host_name"`
	PrometheusUrl     string `tfschema:"prometheus_url"`
	SID               string `tfschema:"sid"`
	SslCertificateUri string `tfschema:"ssl_certificate_uri"`
	SslPreference     string `tfschema:"ssl_preference"`
}

type MssqlServerProviderSettings struct {
	DbPassword        string `tfschema:"db_password"`
	DbPasswordUri     string `tfschema:"db_password_uri"`
	DbPort            string `tfschema:"db_port"`
	DbUsername        string `tfschema:"db_username"`
	HostName          string `tfschema:"host_name"`
	SID               string `tfschema:"sid"`
	SslCertificateUri string `tfschema:"ssl_certificate_uri"`
	SslPreference     string `tfschema:"ssl_preference"`
}

type Db2ProviderSettings struct {
	DbName            string `tfschema:"db_name"`
	DbPassword        string `tfschema:"db_password"`
	DbPasswordUri     string `tfschema:"db_password_uri"`
	DbPort            string `tfschema:"db_port"`
	DbUsername        string `tfschema:"db_username"`
	HostName          string `tfschema:"host_name"`
	SID               string `tfschema:"sid"`
	SslCertificateUri string `tfschema:"ssl_certificate_uri"`
	SslPreference     string `tfschema:"ssl_preference"`
}

type SAPHanaProviderSettings struct {
	DbName                   string `tfschema:"db_name"`
	DbPassword               string `tfschema:"db_password"`
	DbPasswordUri            string `tfschema:"db_password_uri"`
	DbUsername               string `tfschema:"db_username"`
	HostName                 string `tfschema:"host_name"`
	InstanceNumber           string `tfschema:"instance_number"`
	SID                      string `tfschema:"sid"`
	SqlPort                  string `tfschema:"sql_port"`
	SslCertificateUri        string `tfschema:"ssl_certificate_uri"`
	SslHostNameInCertificate string `tfschema:"ssl_host_name_in_certificate"`
	SslPreference            string `tfschema:"ssl_preference"`
}

type SAPNetWeaverProviderSettings struct {
	ClientId          string   `tfschema:"client_id"`
	HostFileEntries   []string `tfschema:"host_file_entries"`
	HostName          string   `tfschema:"host_name"`
	InstanceNr        string   `tfschema:"instance_nr"`
	Password          string   `tfschema:"password"`
	PasswordUri       string   `tfschema:"password_uri"`
	PortNumber        string   `tfschema:"port_number"`
	SID               string   `tfschema:"sid"`
	Username          string   `tfschema:"username"`
	SslCertificateUri string   `tfschema:"ssl_certificate_uri"`
	SslPreference     string   `tfschema:"ssl_preference"`
}

type WorkloadsSAPMonitorProviderInstanceResource struct{}

func (r WorkloadsSAPMonitorProviderInstanceResource) ResourceType() string {
	return "azurerm_workloads_sap_monitor_provider_instance"
}

func (r WorkloadsSAPMonitorProviderInstanceResource) ModelObject() interface{} {
	return &WorkloadsSAPMonitorProviderInstanceModel{}
}

func (r WorkloadsSAPMonitorProviderInstanceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return providerinstances.ValidateProviderInstanceID
}

func (r WorkloadsSAPMonitorProviderInstanceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"monitor_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: monitors.ValidateMonitorID,
		},

		"identity": commonschema.UserAssignedIdentityRequired(),

		"prometheus_os_provider_settings": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"prometheus_url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"sid": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"ssl_certificate_uri": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"ssl_preference": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
			ExactlyOneOf: []string{"prometheus_os_provider_settings", "prometheus_ha_cluster_provider_settings", "mssql_server_provider_settings", "db2_provider_settings", "sap_hana_provider_settings", "sap_net_weaver_provider_settings"},
		},

		"prometheus_ha_cluster_provider_settings": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cluster_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"host_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"prometheus_url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"sid": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"ssl_certificate_uri": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"ssl_preference": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
			ExactlyOneOf: []string{"prometheus_os_provider_settings", "prometheus_ha_cluster_provider_settings", "mssql_server_provider_settings", "db2_provider_settings", "sap_hana_provider_settings", "sap_net_weaver_provider_settings"},
		},

		"mssql_server_provider_settings": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"db_password": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"db_password_uri": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"db_port": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"db_username": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"host_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"sid": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"ssl_certificate_uri": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"ssl_preference": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
			ExactlyOneOf: []string{"prometheus_os_provider_settings", "prometheus_ha_cluster_provider_settings", "mssql_server_provider_settings", "db2_provider_settings", "sap_hana_provider_settings", "sap_net_weaver_provider_settings"},
		},

		"db2_provider_settings": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"db_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"db_password": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"db_password_uri": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"db_port": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"db_username": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"host_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"sid": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"ssl_certificate_uri": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"ssl_preference": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
			ExactlyOneOf: []string{"prometheus_os_provider_settings", "prometheus_ha_cluster_provider_settings", "mssql_server_provider_settings", "db2_provider_settings", "sap_hana_provider_settings", "sap_net_weaver_provider_settings"},
		},

		"sap_hana_provider_settings": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"db_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"db_password": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"db_password_uri": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"db_username": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"host_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"instance_number": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"sid": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"sql_port": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"ssl_certificate_uri": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"ssl_host_name_in_certificate": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"ssl_preference": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
			ExactlyOneOf: []string{"prometheus_os_provider_settings", "prometheus_ha_cluster_provider_settings", "mssql_server_provider_settings", "db2_provider_settings", "sap_hana_provider_settings", "sap_net_weaver_provider_settings"},
		},

		"sap_net_weaver_provider_settings": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"host_file_entries": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"host_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"instance_nr": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"password": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"password_uri": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"port_number": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"sid": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"username": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"ssl_certificate_uri": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"ssl_preference": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
			ExactlyOneOf: []string{"prometheus_os_provider_settings", "prometheus_ha_cluster_provider_settings", "mssql_server_provider_settings", "db2_provider_settings", "sap_hana_provider_settings", "sap_net_weaver_provider_settings"},
		},
	}
}

func (r WorkloadsSAPMonitorProviderInstanceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r WorkloadsSAPMonitorProviderInstanceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkloadsSAPMonitorProviderInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Workloads.ProviderInstances
			monitorId, err := monitors.ParseMonitorID(model.MonitorId)
			if err != nil {
				return err
			}

			id := providerinstances.NewProviderInstanceID(monitorId.SubscriptionId, monitorId.ResourceGroupName, monitorId.MonitorName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			identity, err := identity.ExpandUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			properties := &providerinstances.ProviderInstance{
				Identity:   identity,
				Properties: &providerinstances.ProviderInstanceProperties{},
			}

			if model.PrometheusOSProviderSettings != nil {
				properties.Properties.ProviderSettings = expandPrometheusOSProviderSettings(model.PrometheusOSProviderSettings)
			}

			if model.PrometheusHAClusterProviderSettings != nil {
				properties.Properties.ProviderSettings = expandPrometheusHAClusterProviderSettings(model.PrometheusHAClusterProviderSettings)
			}

			if model.MssqlServerProviderSettings != nil {
				properties.Properties.ProviderSettings = expandMssqlServerProviderSettings(model.MssqlServerProviderSettings)
			}

			if model.Db2ProviderSettings != nil {
				properties.Properties.ProviderSettings = expandDb2ProviderSettings(model.Db2ProviderSettings)
			}

			if model.SAPHanaProviderSettings != nil {
				properties.Properties.ProviderSettings = expandSAPHanaProviderSettings(model.SAPHanaProviderSettings)
			}

			if model.SAPNetWeaverProviderSettings != nil {
				properties.Properties.ProviderSettings = expandSAPNetWeaverProviderSettings(model.SAPNetWeaverProviderSettings)
			}

			if err := client.CreateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r WorkloadsSAPMonitorProviderInstanceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.ProviderInstances

			id, err := providerinstances.ParseProviderInstanceID(metadata.ResourceData.Id())
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

			state := WorkloadsSAPMonitorProviderInstanceModel{
				Name:      id.ProviderInstanceName,
				MonitorId: monitors.NewMonitorID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName).ID(),
			}

			identity, err := identity.FlattenUserAssignedMapToModel(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}
			state.Identity = *identity

			if properties := model.Properties; properties != nil {
				if v := properties.ProviderSettings; v != nil {
					if v, ok := v.(providerinstances.PrometheusOSProviderInstanceProperties); ok {
						state.PrometheusOSProviderSettings = flattenPrometheusOSProviderSettings(v)
					}

					if v, ok := v.(providerinstances.PrometheusHaClusterProviderInstanceProperties); ok {
						state.PrometheusHAClusterProviderSettings = flattenPrometheusHAClusterProviderSettings(v)
					}

					if v, ok := v.(providerinstances.MsSqlServerProviderInstanceProperties); ok {
						state.MssqlServerProviderSettings = flattenMssqlServerProviderSettings(v)
					}

					if v, ok := v.(providerinstances.DB2ProviderInstanceProperties); ok {
						state.Db2ProviderSettings = flattenDb2ProviderSettings(v)
					}

					if v, ok := v.(providerinstances.HanaDbProviderInstanceProperties); ok {
						state.SAPHanaProviderSettings = flattenSAPHanaProviderSettings(v)
					}

					if v, ok := v.(providerinstances.SapNetWeaverProviderInstanceProperties); ok {
						state.SAPNetWeaverProviderSettings = flattenSAPNetWeaverProviderSettings(v)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r WorkloadsSAPMonitorProviderInstanceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.ProviderInstances

			id, err := providerinstances.ParseProviderInstanceID(metadata.ResourceData.Id())
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

func expandPrometheusOSProviderSettings(input []PrometheusOSProviderSettings) *providerinstances.PrometheusOSProviderInstanceProperties {
	prometheusOSProviderSettings := &input[0]

	result := providerinstances.PrometheusOSProviderInstanceProperties{}

	if v := prometheusOSProviderSettings.PrometheusUrl; v != "" {
		result.PrometheusUrl = utils.String(v)
	}

	if v := prometheusOSProviderSettings.SID; v != "" {
		result.SapSid = utils.String(v)
	}

	if v := prometheusOSProviderSettings.SslCertificateUri; v != "" {
		result.SslCertificateUri = utils.String(v)
	}

	if v := prometheusOSProviderSettings.SslPreference; v != "" {
		sslPreference := providerinstances.SslPreference(v)
		result.SslPreference = &sslPreference
	}

	return &result
}

func flattenPrometheusOSProviderSettings(input providerinstances.PrometheusOSProviderInstanceProperties) []PrometheusOSProviderSettings {
	var result []PrometheusOSProviderSettings

	prometheusOSProviderSettings := PrometheusOSProviderSettings{}

	if v := input.PrometheusUrl; v != nil {
		prometheusOSProviderSettings.PrometheusUrl = *v
	}

	if v := input.SapSid; v != nil {
		prometheusOSProviderSettings.SID = *v
	}

	if v := input.SslCertificateUri; v != nil {
		prometheusOSProviderSettings.SslCertificateUri = *v
	}

	if v := input.SslPreference; v != nil {
		prometheusOSProviderSettings.SslPreference = string(*v)
	}

	return append(result, prometheusOSProviderSettings)
}

func expandPrometheusHAClusterProviderSettings(input []PrometheusHAClusterProviderSettings) *providerinstances.PrometheusHaClusterProviderInstanceProperties {
	prometheusHAClusterProviderSettings := &input[0]

	result := providerinstances.PrometheusHaClusterProviderInstanceProperties{}

	if v := prometheusHAClusterProviderSettings.ClusterName; v != "" {
		result.ClusterName = utils.String(v)
	}

	if v := prometheusHAClusterProviderSettings.HostName; v != "" {
		result.Hostname = utils.String(v)
	}

	if v := prometheusHAClusterProviderSettings.PrometheusUrl; v != "" {
		result.PrometheusUrl = utils.String(v)
	}

	if v := prometheusHAClusterProviderSettings.SID; v != "" {
		result.Sid = utils.String(v)
	}

	if v := prometheusHAClusterProviderSettings.SslCertificateUri; v != "" {
		result.SslCertificateUri = utils.String(v)
	}

	if v := prometheusHAClusterProviderSettings.SslPreference; v != "" {
		sslPreference := providerinstances.SslPreference(v)
		result.SslPreference = &sslPreference
	}

	return &result
}

func flattenPrometheusHAClusterProviderSettings(input providerinstances.PrometheusHaClusterProviderInstanceProperties) []PrometheusHAClusterProviderSettings {
	var result []PrometheusHAClusterProviderSettings

	prometheusHAClusterProviderSettings := PrometheusHAClusterProviderSettings{}

	if v := input.ClusterName; v != nil {
		prometheusHAClusterProviderSettings.ClusterName = *v
	}

	if v := input.Hostname; v != nil {
		prometheusHAClusterProviderSettings.HostName = *v
	}

	if v := input.PrometheusUrl; v != nil {
		prometheusHAClusterProviderSettings.PrometheusUrl = *v
	}

	if v := input.Sid; v != nil {
		prometheusHAClusterProviderSettings.SID = *v
	}

	if v := input.SslCertificateUri; v != nil {
		prometheusHAClusterProviderSettings.SslCertificateUri = *v
	}

	if v := input.SslPreference; v != nil {
		prometheusHAClusterProviderSettings.SslPreference = string(*v)
	}

	return append(result, prometheusHAClusterProviderSettings)
}

func expandMssqlServerProviderSettings(input []MssqlServerProviderSettings) *providerinstances.MsSqlServerProviderInstanceProperties {
	mssqlServerProviderSettings := &input[0]

	result := providerinstances.MsSqlServerProviderInstanceProperties{}

	if v := mssqlServerProviderSettings.DbPassword; v != "" {
		result.DbPassword = utils.String(v)
	}

	if v := mssqlServerProviderSettings.DbPasswordUri; v != "" {
		result.DbPasswordUri = utils.String(v)
	}

	if v := mssqlServerProviderSettings.DbPort; v != "" {
		result.DbPort = utils.String(v)
	}

	if v := mssqlServerProviderSettings.DbUsername; v != "" {
		result.DbUsername = utils.String(v)
	}

	if v := mssqlServerProviderSettings.HostName; v != "" {
		result.Hostname = utils.String(v)
	}

	if v := mssqlServerProviderSettings.SID; v != "" {
		result.SapSid = utils.String(v)
	}

	if v := mssqlServerProviderSettings.SslCertificateUri; v != "" {
		result.SslCertificateUri = utils.String(v)
	}

	if v := mssqlServerProviderSettings.SslPreference; v != "" {
		sslPreference := providerinstances.SslPreference(v)
		result.SslPreference = &sslPreference
	}

	return &result
}

func flattenMssqlServerProviderSettings(input providerinstances.MsSqlServerProviderInstanceProperties) []MssqlServerProviderSettings {
	var result []MssqlServerProviderSettings

	mssqlServerProviderSettings := MssqlServerProviderSettings{}

	if v := input.DbPassword; v != nil {
		mssqlServerProviderSettings.DbPassword = *v
	}

	if v := input.DbPasswordUri; v != nil {
		mssqlServerProviderSettings.DbPasswordUri = *v
	}

	if v := input.DbPort; v != nil {
		mssqlServerProviderSettings.DbPort = *v
	}

	if v := input.DbUsername; v != nil {
		mssqlServerProviderSettings.DbUsername = *v
	}

	if v := input.Hostname; v != nil {
		mssqlServerProviderSettings.HostName = *v
	}

	if v := input.SapSid; v != nil {
		mssqlServerProviderSettings.SID = *v
	}

	if v := input.SslCertificateUri; v != nil {
		mssqlServerProviderSettings.SslCertificateUri = *v
	}

	if v := input.SslPreference; v != nil {
		mssqlServerProviderSettings.SslPreference = string(*v)
	}

	return append(result, mssqlServerProviderSettings)
}

func expandDb2ProviderSettings(input []Db2ProviderSettings) *providerinstances.DB2ProviderInstanceProperties {
	db2ProviderSettings := &input[0]

	result := providerinstances.DB2ProviderInstanceProperties{}

	if v := db2ProviderSettings.DbName; v != "" {
		result.DbName = utils.String(v)
	}

	if v := db2ProviderSettings.DbPassword; v != "" {
		result.DbPassword = utils.String(v)
	}

	if v := db2ProviderSettings.DbPasswordUri; v != "" {
		result.DbPasswordUri = utils.String(v)
	}

	if v := db2ProviderSettings.DbPort; v != "" {
		result.DbPort = utils.String(v)
	}

	if v := db2ProviderSettings.DbUsername; v != "" {
		result.DbUsername = utils.String(v)
	}

	if v := db2ProviderSettings.HostName; v != "" {
		result.Hostname = utils.String(v)
	}

	if v := db2ProviderSettings.SID; v != "" {
		result.SapSid = utils.String(v)
	}

	if v := db2ProviderSettings.SslCertificateUri; v != "" {
		result.SslCertificateUri = utils.String(v)
	}

	if v := db2ProviderSettings.SslPreference; v != "" {
		sslPreference := providerinstances.SslPreference(v)
		result.SslPreference = &sslPreference
	}

	return &result
}

func flattenDb2ProviderSettings(input providerinstances.DB2ProviderInstanceProperties) []Db2ProviderSettings {
	var result []Db2ProviderSettings

	db2ProviderSettings := Db2ProviderSettings{}

	if v := input.DbName; v != nil {
		db2ProviderSettings.DbName = *v
	}

	if v := input.DbPassword; v != nil {
		db2ProviderSettings.DbPassword = *v
	}

	if v := input.DbPasswordUri; v != nil {
		db2ProviderSettings.DbPasswordUri = *v
	}

	if v := input.DbPort; v != nil {
		db2ProviderSettings.DbPort = *v
	}

	if v := input.DbUsername; v != nil {
		db2ProviderSettings.DbUsername = *v
	}

	if v := input.Hostname; v != nil {
		db2ProviderSettings.HostName = *v
	}

	if v := input.SapSid; v != nil {
		db2ProviderSettings.SID = *v
	}

	if v := input.SslCertificateUri; v != nil {
		db2ProviderSettings.SslCertificateUri = *v
	}

	if v := input.SslPreference; v != nil {
		db2ProviderSettings.SslPreference = string(*v)
	}

	return append(result, db2ProviderSettings)
}

func expandSAPHanaProviderSettings(input []SAPHanaProviderSettings) *providerinstances.HanaDbProviderInstanceProperties {
	sapHanaProviderSettings := &input[0]

	result := providerinstances.HanaDbProviderInstanceProperties{}

	if v := sapHanaProviderSettings.DbName; v != "" {
		result.DbName = utils.String(v)
	}

	if v := sapHanaProviderSettings.DbPassword; v != "" {
		result.DbPassword = utils.String(v)
	}

	if v := sapHanaProviderSettings.DbPasswordUri; v != "" {
		result.DbPasswordUri = utils.String(v)
	}

	if v := sapHanaProviderSettings.DbUsername; v != "" {
		result.DbUsername = utils.String(v)
	}

	if v := sapHanaProviderSettings.HostName; v != "" {
		result.Hostname = utils.String(v)
	}

	if v := sapHanaProviderSettings.InstanceNumber; v != "" {
		result.InstanceNumber = utils.String(v)
	}

	if v := sapHanaProviderSettings.SID; v != "" {
		result.SapSid = utils.String(v)
	}

	if v := sapHanaProviderSettings.SqlPort; v != "" {
		result.SqlPort = utils.String(v)
	}

	if v := sapHanaProviderSettings.SslCertificateUri; v != "" {
		result.SslCertificateUri = utils.String(v)
	}

	if v := sapHanaProviderSettings.SslHostNameInCertificate; v != "" {
		result.SslHostNameInCertificate = utils.String(v)
	}

	if v := sapHanaProviderSettings.SslPreference; v != "" {
		sslPreference := providerinstances.SslPreference(v)
		result.SslPreference = &sslPreference
	}

	return &result
}

func flattenSAPHanaProviderSettings(input providerinstances.HanaDbProviderInstanceProperties) []SAPHanaProviderSettings {
	var result []SAPHanaProviderSettings

	sapHanaProviderSettings := SAPHanaProviderSettings{}

	if v := input.DbName; v != nil {
		sapHanaProviderSettings.DbName = *v
	}

	if v := input.DbPassword; v != nil {
		sapHanaProviderSettings.DbPassword = *v
	}

	if v := input.DbPasswordUri; v != nil {
		sapHanaProviderSettings.DbPasswordUri = *v
	}

	if v := input.DbUsername; v != nil {
		sapHanaProviderSettings.DbUsername = *v
	}

	if v := input.Hostname; v != nil {
		sapHanaProviderSettings.HostName = *v
	}

	if v := input.InstanceNumber; v != nil {
		sapHanaProviderSettings.InstanceNumber = *v
	}

	if v := input.SapSid; v != nil {
		sapHanaProviderSettings.SID = *v
	}

	if v := input.SqlPort; v != nil {
		sapHanaProviderSettings.SqlPort = *v
	}

	if v := input.SslCertificateUri; v != nil {
		sapHanaProviderSettings.SslCertificateUri = *v
	}

	if v := input.SslHostNameInCertificate; v != nil {
		sapHanaProviderSettings.SslHostNameInCertificate = *v
	}

	if v := input.SslPreference; v != nil {
		sapHanaProviderSettings.SslPreference = string(*v)
	}

	return append(result, sapHanaProviderSettings)
}

func expandSAPNetWeaverProviderSettings(input []SAPNetWeaverProviderSettings) *providerinstances.SapNetWeaverProviderInstanceProperties {
	sapNetWeaverProviderSettings := &input[0]

	result := providerinstances.SapNetWeaverProviderInstanceProperties{}

	if v := sapNetWeaverProviderSettings.ClientId; v != "" {
		result.SapClientId = utils.String(v)
	}

	if v := sapNetWeaverProviderSettings.HostFileEntries; len(v) != 0 {
		result.SapHostFileEntries = &v
	}

	if v := sapNetWeaverProviderSettings.HostName; v != "" {
		result.SapHostname = utils.String(v)
	}

	if v := sapNetWeaverProviderSettings.InstanceNr; v != "" {
		result.SapInstanceNr = utils.String(v)
	}

	if v := sapNetWeaverProviderSettings.Password; v != "" {
		result.SapPassword = utils.String(v)
	}

	if v := sapNetWeaverProviderSettings.PasswordUri; v != "" {
		result.SapPasswordUri = utils.String(v)
	}

	if v := sapNetWeaverProviderSettings.PortNumber; v != "" {
		result.SapPortNumber = utils.String(v)
	}

	if v := sapNetWeaverProviderSettings.SID; v != "" {
		result.SapSid = utils.String(v)
	}

	if v := sapNetWeaverProviderSettings.Username; v != "" {
		result.SapUsername = utils.String(v)
	}

	if v := sapNetWeaverProviderSettings.SslCertificateUri; v != "" {
		result.SslCertificateUri = utils.String(v)
	}

	if v := sapNetWeaverProviderSettings.SslPreference; v != "" {
		sslPreference := providerinstances.SslPreference(v)
		result.SslPreference = &sslPreference
	}

	return &result
}

func flattenSAPNetWeaverProviderSettings(input providerinstances.SapNetWeaverProviderInstanceProperties) []SAPNetWeaverProviderSettings {
	var result []SAPNetWeaverProviderSettings

	sapNetWeaverProviderSettings := SAPNetWeaverProviderSettings{}

	if v := input.SapClientId; v != nil {
		sapNetWeaverProviderSettings.ClientId = *v
	}

	if v := input.SapHostFileEntries; v != nil {
		sapNetWeaverProviderSettings.HostFileEntries = *v
	}

	if v := input.SapHostname; v != nil {
		sapNetWeaverProviderSettings.HostName = *v
	}

	if v := input.SapInstanceNr; v != nil {
		sapNetWeaverProviderSettings.InstanceNr = *v
	}

	if v := input.SapPassword; v != nil {
		sapNetWeaverProviderSettings.Password = *v
	}

	if v := input.SapPasswordUri; v != nil {
		sapNetWeaverProviderSettings.PasswordUri = *v
	}

	if v := input.SapPortNumber; v != nil {
		sapNetWeaverProviderSettings.PortNumber = *v
	}

	if v := input.SapSid; v != nil {
		sapNetWeaverProviderSettings.SID = *v
	}

	if v := input.SapUsername; v != nil {
		sapNetWeaverProviderSettings.Username = *v
	}

	if v := input.SslCertificateUri; v != nil {
		sapNetWeaverProviderSettings.SslCertificateUri = *v
	}

	if v := input.SslPreference; v != nil {
		sapNetWeaverProviderSettings.SslPreference = string(*v)
	}

	return append(result, sapNetWeaverProviderSettings)
}
