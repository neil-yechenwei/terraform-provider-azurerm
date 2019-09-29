package azurerm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"strings"
)

func dataSourceArmNetAppVolume() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmNetAppVolumeRead,

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

			"creation_token": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"export_policy": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed_clients": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cifs": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"nfsv3": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"nfsv4": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"rule_index": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"unix_read_only": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"unix_read_write": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"service_level": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"usage_threshold": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmNetAppVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).netapp.VolumeClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	accountName := d.Get("account_name").(string)
	poolName := d.Get("pool_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, accountName, poolName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: NetApp Volume %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading NetApp Volume %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	volumeName := strings.Split(*resp.Name, "/")[2]
	d.Set("name", volumeName)
	d.Set("resource_group_name", resourceGroup)
	d.Set("account_name", accountName)
	d.Set("pool_name", poolName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if volumeProperties := resp.VolumeProperties; volumeProperties != nil {
		if err := d.Set("creation_token", volumeProperties.CreationToken); err != nil {
			return fmt.Errorf("Error setting `creation_token`: %+v", err)
		}

		if err := d.Set("export_policy", flattenArmExportPolicy(volumeProperties.ExportPolicy)); err != nil {
			return fmt.Errorf("Error setting `export_policy`: %+v", err)
		}

		if err := d.Set("service_level", volumeProperties.ServiceLevel); err != nil {
			return fmt.Errorf("Error setting `service_level`: %+v", err)
		}

		if err := d.Set("subnet_id", volumeProperties.SubnetID); err != nil {
			return fmt.Errorf("Error setting `subnet_id`: %+v", err)
		}

		if err := d.Set("usage_threshold", volumeProperties.UsageThreshold); err != nil {
			return fmt.Errorf("Error setting `usage_threshold`: %+v", err)
		}
	}

	return tags.FlattenAndSetTags(d, resp.Tags)
}
