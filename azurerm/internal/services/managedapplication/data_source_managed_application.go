package managedapplication

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmManagedApplication() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmManagedApplicationRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateManagedAppName,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"managed_resource_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"application_definition_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"parameters": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"plan": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"publisher": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"promotion_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmManagedApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Managed Application (Managed Application Name %q / Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading Managed Application (Managed Application Name %q / Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Error retrieving Managed Application %q (Resource Group %q): ID was nil or empty", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	d.Set("resource_group_name", resourceGroup)
	d.Set("name", name)
	d.Set("kind", resp.Kind)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if err := d.Set("plan", flattenArmManagedApplicationPlan(resp.Plan)); err != nil {
		return fmt.Errorf("Error setting `plan`: %+v", err)
	}
	if props := resp.ApplicationProperties; props != nil {
		d.Set("managed_resource_group_id", props.ManagedResourceGroupID)
		d.Set("application_definition_id", props.ApplicationDefinitionID)

		if v, ok := d.GetOk("parameters"); ok {
			json, err := structure.ExpandJsonFromString(v.(string))
			if err != nil {
				return fmt.Errorf("Error expanding JSON from Parameters %q: %+v", v, err)
			}

			d.Set("parameters", json)
		}
	}

	return nil
}
