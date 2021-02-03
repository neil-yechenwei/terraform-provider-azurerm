package resource

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2019-06-01-preview/templatespecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceTemplateSpecVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTemplateSpecVersionRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.TemplateSpecVersionName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"template_spec_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.TemplateSpecName,
			},

			"artifact": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTemplateSpecVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.TemplateSpecVersionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	templateSpecName := d.Get("template_spec_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, templateSpecName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Template Spec Version %s (Resource Group %s) was not found: %+v", name, resourceGroup, err)
		}

		return fmt.Errorf("making Read request on Template Spec Version '%s': %+v", name, err)
	}

	d.SetId(*resp.ID)
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("template_spec_name", templateSpecName)

	if props := resp.VersionProperties; props != nil {
		if props.Artifacts != nil {
			artifacts, err := flattenTemplateSpecVersionTemplateArtifactPaths(props.Artifacts)
			if err != nil {
				return fmt.Errorf("setting `artifact`: %+v", err)
			}
			d.Set("artifact", artifacts)
		}
	}

	return nil
}

func flattenTemplateSpecVersionTemplateArtifactPaths(input *[]templatespecs.BasicArtifact) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	for _, item := range *input {
		artifact := item.(templatespecs.TemplateArtifact)

		var path string
		if artifact.Path != nil {
			path = *artifact.Path
		}

		results = append(results, map[string]interface{}{
			"path": path,
		})
	}

	return results, nil
}
