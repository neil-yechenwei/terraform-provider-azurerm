package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2019-06-01/netapp"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmNetAppSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetAppSnapshotCreate,
		Read:   resourceArmNetAppSnapshotRead,
		Update: resourceArmNetAppSnapshotUpdate,
		Delete: resourceArmNetAppSnapshotDelete,

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

			"pool_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"volume_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"file_system_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"snapshot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmNetAppSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).netapp.SnapshotClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("account_name").(string)
	poolName := d.Get("pool_name").(string)
	volumeName := d.Get("volume_name").(string)
	fileSystemId := d.Get("file_system_id").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, accountName, poolName, volumeName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for present of existing NetApp Snapshot %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_netapp_snapshot", *resp.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))

	snapshotParameters := netapp.Snapshot{
		Location:           utils.String(location),
		SnapshotProperties: &netapp.SnapshotProperties{},
	}

	if fileSystemId != "" {
		snapshotParameters.SnapshotProperties.FileSystemID = utils.String(fileSystemId)
	}

	future, err := client.Create(ctx, snapshotParameters, resourceGroup, accountName, poolName, volumeName, name)
	if err != nil {
		return fmt.Errorf("Error creating NetApp Snapshot %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of NetApp Snapshot %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, accountName, poolName, volumeName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving NetApp Snapshot %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read NetApp Snapshot %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmNetAppSnapshotRead(d, meta)
}

func resourceArmNetAppSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).netapp.SnapshotClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for NetApp Snapshot updating")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("account_name").(string)
	poolName := d.Get("pool_name").(string)
	volumeName := d.Get("volume_name").(string)

	snapshotPatchParameters := netapp.SnapshotPatch{}

	result, err := client.Update(ctx, snapshotPatchParameters, resourceGroup, accountName, poolName, volumeName, name)
	if err != nil {
		return fmt.Errorf("Error updating NetApp Snapshot %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if result.ID == nil {
		return fmt.Errorf("Cannot read NetApp Snapshot %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*result.ID)

	return resourceArmNetAppSnapshotRead(d, meta)
}

func resourceArmNetAppSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).netapp.SnapshotClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	accountName := id.Path["netAppAccounts"]
	poolName := id.Path["capacityPools"]
	volumeName := id.Path["volumes"]
	name := id.Path["snapshots"]

	resp, err := client.Get(ctx, resourceGroup, accountName, poolName, volumeName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] NetApp Snapshots %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading NetApp Snapshots %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	snapshotName := strings.Split(*resp.Name, "/")[3]
	d.Set("name", snapshotName)
	d.Set("resource_group_name", resourceGroup)
	d.Set("account_name", accountName)
	d.Set("pool_name", poolName)
	d.Set("volume_name", volumeName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if snapshotProperties := resp.SnapshotProperties; snapshotProperties != nil {
		if err := d.Set("file_system_id", snapshotProperties.FileSystemID); err != nil {
			return fmt.Errorf("Error setting `file_system_id`: %+v", err)
		}

		if err := d.Set("snapshot_id", snapshotProperties.SnapshotID); err != nil {
			return fmt.Errorf("Error setting `snapshot_id`: %+v", err)
		}
	}

	if t := resp.Tags; t != nil {
		return tags.FlattenAndSetTags(d, t)
	}

	return nil
}

func resourceArmNetAppSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).netapp.SnapshotClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	accountName := id.Path["netAppAccounts"]
	poolName := id.Path["capacityPools"]
	volumeName := id.Path["volumes"]
	name := id.Path["snapshots"]

	future, err := client.Delete(ctx, resourceGroup, accountName, poolName, volumeName, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting NetApp Snapshot %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting NetApp Snapshot %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}
