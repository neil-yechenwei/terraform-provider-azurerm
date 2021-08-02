package domainservices

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/domainservices/mgmt/2020-01-01/aad"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/domainservices/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceActiveDirectoryDomainServiceOuContainer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceActiveDirectoryDomainServiceOuContainerCreate,
		Read:   resourceActiveDirectoryDomainServiceOuContainerRead,
		Update: resourceActiveDirectoryDomainServiceOuContainerUpdate,
		Delete: resourceActiveDirectoryDomainServiceOuContainerDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DomainServiceOuContainerID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"domain_service_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"account_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"password": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"spn": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceActiveDirectoryDomainServiceOuContainerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DomainServices.OuContainerClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	domainServiceName := d.Get("domain_service_name").(string)

	id := parse.NewDomainServiceOuContainerID(subscriptionId, resourceGroup, domainServiceName, name).ID()

	existing, err := client.Get(ctx, resourceGroup, domainServiceName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_active_directory_domain_service_ou_container", id)
	}

	containerAccount := aad.ContainerAccount{}

	if v, ok := d.GetOk("account_name"); ok {
		containerAccount.AccountName = utils.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		containerAccount.Password = utils.String(v.(string))
	}

	if v, ok := d.GetOk("spn"); ok {
		containerAccount.Spn = utils.String(v.(string))
	}

	future, err := client.Create(ctx, resourceGroup, domainServiceName, name, containerAccount)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id)
	return resourceActiveDirectoryDomainServiceOuContainerRead(d, meta)
}

func resourceActiveDirectoryDomainServiceOuContainerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DomainServices.OuContainerClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DomainServiceOuContainerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.DomainServiceName, id.OuContainerName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] aad %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.OuContainerName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("domain_service_name", id.DomainServiceName)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.OuContainerProperties; props != nil {
		d.Set("domain_name", props.DomainName)

		if accounts := props.Accounts; accounts != nil {
			for _, item := range *accounts {
				var accountName string
				if item.AccountName != nil {
					accountName = *item.AccountName
				}

				var password string
				if item.Password != nil {
					password = *item.Password
				}

				var spn string
				if item.Spn != nil {
					spn = *item.Spn
				}

				d.Set("account_name", accountName)
				d.Set("password", password)
				d.Set("spn", spn)
			}
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceActiveDirectoryDomainServiceOuContainerUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DomainServices.OuContainerClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DomainServiceOuContainerID(d.Id())
	if err != nil {
		return err
	}

	containerAccount := aad.ContainerAccount{}

	if d.HasChange("account_name") {
		containerAccount.AccountName = utils.String(d.Get("account_name").(string))
	}

	if d.HasChange("password") {
		containerAccount.Password = utils.String(d.Get("password").(string))
	}

	if d.HasChange("spn") {
		containerAccount.Spn = utils.String(d.Get("spn").(string))
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.DomainServiceName, id.OuContainerName, containerAccount)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", *id, err)
	}

	return resourceActiveDirectoryDomainServiceOuContainerRead(d, meta)
}

func resourceActiveDirectoryDomainServiceOuContainerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DomainServices.OuContainerClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DomainServiceOuContainerID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.DomainServiceName, id.OuContainerName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}
