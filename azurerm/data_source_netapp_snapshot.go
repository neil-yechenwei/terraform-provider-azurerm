package azurerm

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmNetAppSnapshot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmNetAppSnapshotRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocationForDataSource(),

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
				Computed: true,
			},

			"snapshot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmNetAppSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).netapp.SnapshotClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	accountName := d.Get("account_name").(string)
	poolName := d.Get("pool_name").(string)
	volumeName := d.Get("volume_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, accountName, poolName, volumeName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: NetApp Snapshot %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading NetApp Snapshot %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

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
