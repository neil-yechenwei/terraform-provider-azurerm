package databox

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databox/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmDataBoxCredential() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmDataBoxCredentialRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"databox_job_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.DataBoxJobName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),
		},
	}
}

func dataSourceArmDataBoxCredentialRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBox.JobClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("databox_job_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.ListCredentials(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("DataBox Credential (DataBox Credential Name %q / Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("reading DataBox Job (DataBox Credential Name %q / Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.Request.URL.Path == "" {
		return fmt.Errorf("API returns a nil/empty id on DataBox Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	d.SetId(resp.Request.URL.Path)

	d.Set("databox_job_name", name)
	d.Set("resource_group_name", resourceGroup)

	return nil
}
