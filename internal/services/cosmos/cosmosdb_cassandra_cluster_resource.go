package cosmos

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCassandraCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCassandraClusterCreateUpdate,
		Read:   resourceCassandraClusterRead,
		Update: resourceCassandraClusterCreateUpdate,
		Delete: resourceCassandraClusterDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.CassandraClusterID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": location.Schema(),

			"delegated_management_subnet_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.SubnetID,
			},

			"default_admin_password": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"audit_logging_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"authentication_method": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(documentdb.AuthenticationMethodNone),
					string(documentdb.AuthenticationMethodCassandra),
				}, false),
			},

			"client_certificate": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"pem": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
					},
				},
			},

			"cluster_name_override": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"deallocated": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"external_gossip_certificate": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"pem": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
					},
				},
			},

			"external_seed_node": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"ip_address": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
					},
				},
			},

			"hours_between_backups": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
			},

			"prometheus_endpoint": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"ip_address": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
					},
				},
			},

			"repair_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"restore_from_backup_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"identity": commonschema.SystemAssignedIdentity(),

			"tags": tags.Schema(),
		},
	}
}

func resourceCassandraClusterCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClustersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroupName := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	id := parse.NewCassandraClusterID(subscriptionId, resourceGroupName, name)

	existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_cosmosdb_cassandra_cluster", id.ID())
	}

	identity, err := expandCassandraClusterIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	body := documentdb.ClusterResource{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Properties: &documentdb.ClusterResourceProperties{
			DelegatedManagementSubnetID:   utils.String(d.Get("delegated_management_subnet_id").(string)),
			InitialCassandraAdminPassword: utils.String(d.Get("default_admin_password").(string)),
		},
		Identity: identity,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("authentication_method"); ok {
		body.Properties.AuthenticationMethod = documentdb.AuthenticationMethod(v.(string))
	}

	if v, ok := d.GetOk("audit_logging_enabled"); ok {
		body.Properties.CassandraAuditLoggingEnabled = utils.Bool(v.(bool))
	}

	if v, ok := d.GetOk("client_certificate"); ok {
		body.Properties.ClientCertificates = expandCassandraClusterCertificate(v.(*pluginsdk.Set).List())
	}

	if v, ok := d.GetOk("version"); ok {
		body.Properties.CassandraVersion = utils.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_name_override"); ok {
		body.Properties.ClusterNameOverride = utils.String(v.(string))
	}

	if v, ok := d.GetOk("deallocated"); ok {
		body.Properties.Deallocated = utils.Bool(v.(bool))
	}

	if v, ok := d.GetOk("external_gossip_certificate"); ok {
		body.Properties.ExternalGossipCertificates = expandCassandraClusterCertificate(v.(*pluginsdk.Set).List())
	}

	if v, ok := d.GetOk("external_seed_node"); ok {
		body.Properties.ExternalSeedNodes = expandCassandraClusterSeedNodes(v.(*pluginsdk.Set).List())
	}

	if v, ok := d.GetOk("hours_between_backups"); ok {
		body.Properties.HoursBetweenBackups = utils.Int32(int32(v.(int)))
	}

	if v, ok := d.GetOk("prometheus_endpoint"); ok {
		body.Properties.PrometheusEndpoint = expandCassandraClusterPrometheusEndpoint(v.([]interface{}))
	}

	if v, ok := d.GetOk("repair_enabled"); ok {
		body.Properties.RepairEnabled = utils.Bool(v.(bool))
	}

	if v, ok := d.GetOk("restore_from_backup_id"); ok {
		body.Properties.RestoreFromBackupID = utils.String(v.(string))
	}

	future, err := client.CreateUpdate(ctx, id.ResourceGroup, id.Name, body)
	if err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on create/update for %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCassandraClusterRead(d, meta)
}

func resourceCassandraClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClustersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CassandraClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading %q - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading %q: %+v", id, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("name", id.Name)

	identity, err := flattenCassandraClusterIdentity(resp.Identity)
	if err != nil {
		return err
	}
	d.Set("identity", identity)

	if props := resp.Properties; props != nil {
		if res := props; res != nil {
			d.Set("delegated_management_subnet_id", props.DelegatedManagementSubnetID)
			d.Set("authentication_method", string(props.AuthenticationMethod))
			d.Set("audit_logging_enabled", props.CassandraAuditLoggingEnabled)
			d.Set("version", props.CassandraVersion)
			d.Set("cluster_name_override", props.ClusterNameOverride)
			d.Set("deallocated", props.Deallocated)
			d.Set("hours_between_backups", props.HoursBetweenBackups)
			d.Set("repair_enabled", props.RepairEnabled)
			d.Set("restore_from_backup_id", props.RestoreFromBackupID)

			if err := d.Set("client_certificate", flattenCassandraClusterCertificate(props.ClientCertificates)); err != nil {
				return fmt.Errorf("setting `client_certificate`: %+v", err)
			}

			if err := d.Set("external_gossip_certificate", flattenCassandraClusterCertificate(props.ExternalGossipCertificates)); err != nil {
				return fmt.Errorf("setting `external_gossip_certificate`: %+v", err)
			}

			if err := d.Set("external_seed_node", flattenCassandraClusterSeedNodes(props.ExternalSeedNodes)); err != nil {
				return fmt.Errorf("setting `external_seed_node`: %+v", err)
			}

			if err := d.Set("prometheus_endpoint", flattenCassandraClusterPrometheusEndpoint(props.PrometheusEndpoint)); err != nil {
				return fmt.Errorf("setting `prometheus_endpoint`: %+v", err)
			}
		}
	}
	// The "default_admin_password" is not returned in GET response, hence setting it from config.
	d.Set("default_admin_password", d.Get("default_admin_password").(string))
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceCassandraClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClustersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CassandraClusterID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("deleting %q: %+v", id, err)
		}
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("waiting on delete future for %q: %+v", id, err)
	}

	return nil
}

func expandCassandraClusterIdentity(input []interface{}) (*documentdb.ManagedCassandraManagedServiceIdentity, error) {
	config, err := identity.ExpandSystemAssigned(input)
	if err != nil {
		return nil, err
	}

	return &documentdb.ManagedCassandraManagedServiceIdentity{
		Type: documentdb.ManagedCassandraResourceIdentityType(config.Type),
	}, nil
}

func expandCassandraClusterCertificate(input []interface{}) *[]documentdb.Certificate {
	results := make([]documentdb.Certificate, 0)

	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, documentdb.Certificate{
			Pem: utils.String(v["pem"].(string)),
		})
	}

	return &results
}

func expandCassandraClusterSeedNodes(input []interface{}) *[]documentdb.SeedNode {
	results := make([]documentdb.SeedNode, 0)

	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, documentdb.SeedNode{
			IPAddress: utils.String(v["ip_address"].(string)),
		})
	}

	return &results
}

func expandCassandraClusterPrometheusEndpoint(input []interface{}) *documentdb.SeedNode {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &documentdb.SeedNode{
		IPAddress: utils.String(v["ip_address"].(string)),
	}
}

func flattenCassandraClusterIdentity(input *documentdb.ManagedCassandraManagedServiceIdentity) ([]interface{}, error) {
	var config *identity.SystemAssigned

	if input != nil {
		principalId := ""
		if input.PrincipalID != nil {
			principalId = *input.PrincipalID
		}

		tenantId := ""
		if input.TenantID != nil {
			tenantId = *input.TenantID
		}

		config = &identity.SystemAssigned{
			Type:        identity.Type(string(input.Type)),
			PrincipalId: principalId,
			TenantId:    tenantId,
		}
	}
	return identity.FlattenSystemAssigned(config), nil
}

func flattenCassandraClusterCertificate(input *[]documentdb.Certificate) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var p string
		if item.Pem != nil {
			p = *item.Pem
		}
		results = append(results, map[string]interface{}{
			"pem": p,
		})
	}

	return results
}

func flattenCassandraClusterSeedNodes(input *[]documentdb.SeedNode) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var ipAddress string
		if item.IPAddress != nil {
			ipAddress = *item.IPAddress
		}
		results = append(results, map[string]interface{}{
			"ip_address": ipAddress,
		})
	}

	return results
}

func flattenCassandraClusterPrometheusEndpoint(input *documentdb.SeedNode) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var ipAddress string
	if input.IPAddress != nil {
		ipAddress = *input.IPAddress
	}

	return []interface{}{
		map[string]interface{}{
			"ip_address": ipAddress,
		},
	}
}
