package compute

import (
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/snapshots"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceSnapshot() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceSnapshotRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			// Computed
			"os_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"disk_size_gb": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},
			"time_created": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"creation_option": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"source_uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"source_resource_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"storage_account_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"encryption_settings": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"disk_encryption_key": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"secret_url": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"source_vault_id": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
						"key_encryption_key": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"key_url": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"source_vault_id": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"trusted_launch_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceSnapshotRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.SnapshotsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := snapshots.NewSnapshotID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("loading %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("os_type", props.OsType)
			d.Set("time_created", props.TimeCreated)

			if props.DiskSizeGB != nil {
				d.Set("disk_size_gb", int(*props.DiskSizeGB))
			}

			if err := d.Set("encryption_settings", flattenSnapshotEncryptionSettings(props.EncryptionSettingsCollection)); err != nil {
				return fmt.Errorf("setting `encryption_settings`: %+v", err)
			}

			trustedLaunchEnabled := false
			if securityProfile := props.SecurityProfile; securityProfile != nil {
				if *securityProfile.SecurityType == snapshots.DiskSecurityTypesTrustedLaunch {
					trustedLaunchEnabled = true
				}
			}
			d.Set("trusted_launch_enabled", trustedLaunchEnabled)

			data := props.CreationData
			d.Set("creation_option", string(data.CreateOption))
			d.Set("source_uri", data.SourceUri)
			d.Set("source_resource_id", data.SourceResourceId)
			d.Set("storage_account_id", data.StorageAccountId)
		}
	}

	return nil
}
