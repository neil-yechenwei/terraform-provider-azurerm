package containers

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/hybridkubernetes/mgmt/2021-03-01/hybridkubernetes"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceConnectedKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceConnectedKubernetesClusterCreateUpdate,
		Read:   resourceConnectedKubernetesClusterRead,
		Update: resourceConnectedKubernetesClusterCreateUpdate,
		Delete: resourceConnectedKubernetesClusterDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ConnectedClusterID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"agent_public_key_certificate": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"identity_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(hybridkubernetes.None),
					string(hybridkubernetes.SystemAssigned),
				}, false),
				Default: hybridkubernetes.SystemAssigned,
			},

			"distribution": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"aks",
					"aks_engine",
					"aks_hci",
					"auto",
					"capz",
					"eks",
					"generic",
					"gke",
					"k3s",
					"kind",
					"minikube",
					"openshift",
					"rancher_rke",
					"tkg",
				}, false),
			},

			"infrastructure": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"auto",
					"aws",
					"azure",
					"azure_stack_edge",
					"azure_stack_hci",
					"azure_stack_hub",
					"gcp",
					"generic",
					"vsphere",
				}, false),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceConnectedKubernetesClusterCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Containers.ConnectedKubernetesClustersClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewConnectedClusterID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_connected_kubernetes_cluster", id.ID())
		}
	}

	parameters := hybridkubernetes.ConnectedCluster{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Identity: &hybridkubernetes.ConnectedClusterIdentity{
			Type: hybridkubernetes.ResourceIdentityType(d.Get("identity_type").(string)),
		},
		ConnectedClusterProperties: &hybridkubernetes.ConnectedClusterProperties{
			AgentPublicKeyCertificate: utils.String(d.Get("agent_public_key_certificate").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("distribution"); ok {
		parameters.ConnectedClusterProperties.Distribution = utils.String(v.(string))
	}

	if v, ok := d.GetOk("infrastructure"); ok {
		parameters.ConnectedClusterProperties.Infrastructure = utils.String(v.(string))
	}

	future, err := client.Create(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceConnectedKubernetesClusterRead(d, meta)
}

func resourceConnectedKubernetesClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ConnectedKubernetesClustersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ConnectedClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s does not exist - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("identity_type", resp.Identity.Type)

	if props := resp.ConnectedClusterProperties; props != nil {
		d.Set("agent_public_key_certificate", props.AgentPublicKeyCertificate)
		d.Set("distribution", props.Distribution)
		d.Set("infrastructure", props.Infrastructure)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceConnectedKubernetesClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ConnectedKubernetesClustersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ConnectedClusterID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}

	return nil
}
