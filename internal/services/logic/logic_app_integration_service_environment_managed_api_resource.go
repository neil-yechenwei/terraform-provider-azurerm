package logic

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2019-05-01/logic"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogicAppIntegrationServiceEnvironmentManagedApi() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogicAppIntegrationServiceEnvironmentManagedApiCreateUpdate,
		Read:   resourceLogicAppIntegrationServiceEnvironmentManagedApiRead,
		Update: resourceLogicAppIntegrationServiceEnvironmentManagedApiCreateUpdate,
		Delete: resourceLogicAppIntegrationServiceEnvironmentManagedApiDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.IntegrationServiceEnvironmentManagedApiID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"as2",
					"azureautomation",
					"azureblob",
					"azureeventgrid",
					"azureeventgridpublish",
					"azurefile",
					"azurequeues",
					"azuretables",
					"db2",
					"documentdb",
					"edifact",
					"eventhubs",
					"ftp",
					"isefilesystem",
					"keyvault",
					"mq",
					"sap",
					"servicebus",
					"sftp",
					"sftpwithssh",
					"si3270",
					"smtp",
					"sql",
					"sqldw",
					"x12",
				}, false),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"integration_service_environment_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationServiceEnvironmentID,
			},

			"content_link_uri": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceLogicAppIntegrationServiceEnvironmentManagedApiCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Logic.IntegrationServiceEnvironmentManagedApiClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	iseId, err := parse.IntegrationServiceEnvironmentID(d.Get("integration_service_environment_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewIntegrationServiceEnvironmentManagedApiID(subscriptionId, resourceGroup, iseId.Name, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, id.IntegrationServiceEnvironmentName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_logic_app_integration_service_environment_managed_api", *existing.ID)
		}
	}

	parameters := logic.IntegrationServiceEnvironmentManagedAPI{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		IntegrationServiceEnvironmentManagedAPIProperties: &logic.IntegrationServiceEnvironmentManagedAPIProperties{
			IntegrationServiceEnvironment: &logic.ResourceReference{
				ID: utils.String(iseId.ID()),
			},
		},
	}

	if v, ok := d.GetOk("content_link_uri"); ok {
		parameters.IntegrationServiceEnvironmentManagedAPIProperties.DeploymentParameters = &logic.IntegrationServiceEnvironmentManagedAPIDeploymentParameters{
			ContentLinkDefinition: &logic.ContentLink{
				URI: utils.String(v.(string)),
			},
		}
	}

	future, err := client.Put(ctx, resourceGroup, id.IntegrationServiceEnvironmentName, name, parameters)
	if err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %q: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogicAppIntegrationServiceEnvironmentManagedApiRead(d, meta)
}

func resourceLogicAppIntegrationServiceEnvironmentManagedApiRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationServiceEnvironmentManagedApiClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationServiceEnvironmentManagedApiID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.IntegrationServiceEnvironmentName, id.ManagedApiName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ManagedApiName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.IntegrationServiceEnvironmentManagedAPIProperties; props != nil {
		if props.IntegrationServiceEnvironment != nil && props.IntegrationServiceEnvironment.ID != nil {
			d.Set("integration_service_environment_id", props.IntegrationServiceEnvironment.ID)
		}

		if props.DeploymentParameters != nil && props.DeploymentParameters.ContentLinkDefinition != nil && props.DeploymentParameters.ContentLinkDefinition.URI != nil {
			d.Set("content_link_uri", props.DeploymentParameters.ContentLinkDefinition.URI)
		}
	}

	return nil
}

func resourceLogicAppIntegrationServiceEnvironmentManagedApiDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationServiceEnvironmentManagedApiClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationServiceEnvironmentManagedApiID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.IntegrationServiceEnvironmentName, id.ManagedApiName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}
