package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVirtualHub() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualHubCreateUpdate,
		Read:   resourceArmVirtualHubRead,
		Update: resourceArmVirtualHubCreateUpdate,
		Delete: resourceArmVirtualHubDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"address_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"virtual_wan_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmVirtualHubCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.VirtualHubClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Virtual Hub creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	addressPrefix := d.Get("address_prefix").(string)
	virtualWanId := d.Get("virtual_wan_id").(string)
	tags := d.Get("tags").(map[string]interface{})

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_virtual_hub", *existing.ID)
		}
	}

	hub := network.VirtualHub{
		Location: utils.String(location),
		Tags:     expandTags(tags),
		VirtualHubProperties: &network.VirtualHubProperties{
			AddressPrefix: &addressPrefix,
			VirtualWan: &network.SubResource{
				ID: utils.String(virtualWanId),
			},
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, hub)
	if err != nil {
		return fmt.Errorf("Error creating/updating Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation/update of Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Virtual Hub %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmVirtualHubRead(d, meta)
}

func resourceArmVirtualHubRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.VirtualHubClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["virtualHubs"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Virtual Hub %q (Resource Group %q) was not found - removing from state", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.VirtualHubProperties; props != nil {
		d.Set("address_prefix", props.AddressPrefix)
		d.Set("virtual_wan_id", props.VirtualWan.ID)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmVirtualHubDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.VirtualHubClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["virtualHubs"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		// deleted outside of Terraform
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error deleting Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for the deletion of Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}
