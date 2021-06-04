package eventhub

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventhub/mgmt/2018-01-01-preview/eventhub"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/validate"
	privateEndpointValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceEventHubNamespacePrivateEndpointConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceEventHubNamespacePrivateEndpointConnectionCreateUpdate,
		Read:   resourceEventHubNamespacePrivateEndpointConnectionRead,
		Update: resourceEventHubNamespacePrivateEndpointConnectionCreateUpdate,
		Delete: resourceEventHubNamespacePrivateEndpointConnectionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.EventHubNamespacePrivateEndpointConnectionID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"namespace_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateEventHubNamespaceName(),
			},

			"private_endpoint_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: privateEndpointValidate.PrivateEndpointID,
			},

			"private_link_service_connection_state": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"description": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"status": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(eventhub.Approved),
								string(eventhub.Disconnected),
								string(eventhub.Pending),
								string(eventhub.Rejected),
							}, true),
						},
					},
				},
			},
		},
	}
}

func resourceEventHubNamespacePrivateEndpointConnectionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Eventhub.PrivateEndpointConnectionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	namespaceName := d.Get("namespace_name").(string)

	id := parse.NewEventHubNamespacePrivateEndpointConnectionID(subscriptionId, resourceGroup, namespaceName, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, namespaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_eventhub_namespace_private_endpoint_connection", id.ID())
		}
	}

	parameters := eventhub.PrivateEndpointConnection{
		PrivateEndpointConnectionProperties: &eventhub.PrivateEndpointConnectionProperties{},
	}

	if v, ok := d.GetOk("private_endpoint_id"); ok {
		parameters.PrivateEndpointConnectionProperties.PrivateEndpoint = &eventhub.PrivateEndpoint{
			ID: utils.String(v.(string)),
		}
	}

	if v, ok := d.GetOk("private_link_service_connection_state"); ok {
		parameters.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState = expandPrivateEndpointConnectionConnectionState(v.([]interface{}))
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, namespaceName, name, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceEventHubNamespacePrivateEndpointConnectionRead(d, meta)
}

func resourceEventHubNamespacePrivateEndpointConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.PrivateEndpointConnectionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EventHubNamespacePrivateEndpointConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.NamespaceName, id.PrivateEndpointConnectionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.PrivateEndpointConnectionName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("namespace_name", id.NamespaceName)

	if props := resp.PrivateEndpointConnectionProperties; props != nil {
		if privateEndpoint := props.PrivateEndpoint; privateEndpoint != nil {
			d.Set("private_endpoint_id", privateEndpoint.ID)
		}

		if err := d.Set("private_link_service_connection_state", flattenPrivateEndpointConnectionConnectionState(props.PrivateLinkServiceConnectionState)); err != nil {
			return fmt.Errorf("setting `private_link_service_connection_state`: %+v", err)
		}
	}

	return nil
}

func resourceEventHubNamespacePrivateEndpointConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.PrivateEndpointConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EventHubNamespacePrivateEndpointConnectionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.NamespaceName, id.PrivateEndpointConnectionName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandPrivateEndpointConnectionConnectionState(input []interface{}) *eventhub.ConnectionState {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	result := eventhub.ConnectionState{}

	if description := v["description"].(string); description != "" {
		result.Description = &description
	}

	if status := v["status"].(string); status != "" {
		result.Status = eventhub.PrivateLinkConnectionStatus(status)
	}

	return &result
}

func flattenPrivateEndpointConnectionConnectionState(input *eventhub.ConnectionState) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var description string
	if input.Description != nil {
		description = *input.Description
	}

	var status string
	if input.Status != "" {
		status = string(input.Status)
	}

	return []interface{}{
		map[string]interface{}{
			"description": description,
			"status":      status,
		},
	}
}
