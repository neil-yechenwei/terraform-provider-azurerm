package compute

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/sdk/2022-03-02/disks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceManagedDisk() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceManagedDiskRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"create_option": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"disk_encryption_set_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"disk_iops_read_write": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"disk_mbps_read_write": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"disk_size_gb": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"image_reference_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"os_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"source_resource_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"source_uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"storage_account_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"storage_account_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),

			"zones": commonschema.ZonesMultipleComputed(),
		},
	}
}

func dataSourceManagedDiskRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DisksClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := disks.NewDiskID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("making Read request on %s: %s", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.DiskName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("zones", zones.Flatten(model.Zones))

		storageAccountType := ""
		if sku := model.Sku; sku != nil {
			storageAccountType = string(*sku.Name)
		}
		d.Set("storage_account_type", storageAccountType)

		if props := model.Properties; props != nil {
			d.Set("create_option", string(props.CreationData.CreateOption))

			imageReferenceID := ""
			if props.CreationData.ImageReference != nil && props.CreationData.ImageReference.Id != nil {
				imageReferenceID = *props.CreationData.ImageReference.Id
			}
			d.Set("image_reference_id", imageReferenceID)

			d.Set("source_resource_id", props.CreationData.SourceResourceId)
			d.Set("source_uri", props.CreationData.SourceUri)
			d.Set("storage_account_id", props.CreationData.StorageAccountId)

			d.Set("disk_size_gb", props.DiskSizeGB)
			d.Set("disk_iops_read_write", props.DiskIOPSReadWrite)
			d.Set("disk_mbps_read_write", props.DiskMBpsReadWrite)
			d.Set("os_type", props.OsType)

			diskEncryptionSetId := ""
			if props.Encryption != nil && props.Encryption.DiskEncryptionSetId != nil {
				diskEncryptionSetId = *props.Encryption.DiskEncryptionSetId
			}
			d.Set("disk_encryption_set_id", diskEncryptionSetId)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}
