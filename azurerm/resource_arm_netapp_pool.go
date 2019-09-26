package azurerm

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2019-06-01/netapp"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"log"
	"strconv"
	"strings"
	"time"
)

func resourceArmNetAppPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetAppPoolCreateUpdate,
		Read:   resourceArmNetAppPoolRead,
		Update: resourceArmNetAppPoolCreateUpdate,
		Delete: resourceArmNetAppPoolDelete,

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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"service_level": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(netapp.Premium),
					string(netapp.Standard),
					string(netapp.Ultra),
				}, false),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(4398046511104, 549755813888000),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmNetAppPoolCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).netapp.PoolClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("account_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, accountName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for present of existing NetApp Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_netapp_pool", *resp.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	serviceLevel := d.Get("service_level").(string)
	size := int64(d.Get("size").(int))
	t := d.Get("tags").(map[string]interface{})

	capacityPoolParameters := netapp.CapacityPool{
		Location: utils.String(location),
		PoolProperties: &netapp.PoolProperties{
			ServiceLevel: netapp.ServiceLevel(serviceLevel),
			Size:         utils.Int64(size),
		},
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, capacityPoolParameters, resourceGroup, accountName, name)
	if err != nil {
		return fmt.Errorf("Error creating NetApp Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of NetApp Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving NetApp Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read NetApp Pool %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmNetAppPoolRead(d, meta)
}

func resourceArmNetAppPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).netapp.PoolClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	accountName := id.Path["netAppAccounts"]
	name := id.Path["capacityPools"]

	resp, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] NetApp Pools %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading NetApp Pools %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	poolName := strings.Split(*resp.Name, "/")[1]
	d.Set("name", poolName)
	d.Set("resource_group_name", resourceGroup)
	d.Set("account_name", accountName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if poolProperties := resp.PoolProperties; poolProperties != nil {
		if err := d.Set("service_level", poolProperties.ServiceLevel); err != nil {
			return fmt.Errorf("Error setting `service_level`: %+v", err)
		}

		if err := d.Set("size", poolProperties.Size); err != nil {
			return fmt.Errorf("Error setting `size`: %+v", err)
		}
	}

	return tags.FlattenAndSetTags(d, resp.Tags)
}

func resourceArmNetAppPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).netapp.PoolClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	accountName := id.Path["netAppAccounts"]
	name := id.Path["capacityPools"]

	future, err := client.Delete(ctx, resourceGroup, accountName, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting NetApp Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return waitForNetAppPoolToBeDeleted(ctx, client, resourceGroup, accountName, name)
}

func waitForNetAppPoolToBeDeleted(ctx context.Context, client *netapp.PoolsClient, resourceGroup, accountName, name string) error {
	log.Printf("[DEBUG] Waiting for NetApp Pool Provisioning Service %q (Resource Group %q) to be deleted", name, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"200", "202"},
		Target:  []string{"404"},
		Refresh: netappPoolDeleteStateRefreshFunc(ctx, client, resourceGroup, accountName, name),
		Timeout: 20 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for NetApp Pool Provisioning Service %q (Resource Group %q) to be deleted: %+v", name, resourceGroup, err)
	}

	return nil
}

func netappPoolDeleteStateRefreshFunc(ctx context.Context, client *netapp.PoolsClient, resourceGroupName string, accountName string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, accountName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(res.Response) {
				return nil, "", fmt.Errorf("Error retrieving NetApp Pool %q (Resource Group %q): %s", name, resourceGroupName, err)
			}
		}

		if _, err := client.Delete(ctx, resourceGroupName, accountName, name); err != nil {
			log.Printf("Error reissuing NetApp Pool %q delete request (Resource Group %q): %+v", name, resourceGroupName, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}
