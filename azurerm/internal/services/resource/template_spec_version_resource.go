package resource

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2019-06-01-preview/templatespecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceTemplateSpecVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTemplateSpecVersionCreate,
		Read:   resourceTemplateSpecVersionRead,
		Update: resourceTemplateSpecVersionUpdate,
		Delete: resourceTemplateSpecVersionDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.TemplateSpecVersionID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.TemplateSpecVersionName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"template_content": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				StateFunc: utils.NormalizeJson,
			},

			"template_spec_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.TemplateSpecName,
			},

			"artifact": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"template_content": {
							Type:      schema.TypeString,
							Required:  true,
							ForceNew:  true,
							StateFunc: utils.NormalizeJson,
						},
					},
				},
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.TemplateSpecVersionDescription,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceTemplateSpecVersionCreate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Resource.TemplateSpecVersionsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	templateSpecName := d.Get("template_spec_name").(string)

	id := parse.NewTemplateSpecVersionID(subscriptionId, resourceGroup, templateSpecName, name)

	existing, err := client.Get(ctx, resourceGroup, templateSpecName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing Template Spec Version %q (Resource Group %q / Template Spec %q): %+v", name, resourceGroup, templateSpecName, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_template_spec_version", id.ID())
	}

	template, err := expandTemplateDeploymentBody(d.Get("template_content").(string))
	if err != nil {
		return fmt.Errorf("expanding `template_content`: %+v", err)
	}

	templateSpecVersion := templatespecs.VersionTemplatespecs{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		VersionProperties: &templatespecs.VersionProperties{
			Template: template,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("artifact"); ok {
		artifacts, err := expandTemplateSpecVersionTemplateArtifacts(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `artifact`: %+v", err)
		}
		templateSpecVersion.VersionProperties.Artifacts = artifacts
	}

	if v, ok := d.GetOk("description"); ok {
		templateSpecVersion.VersionProperties.Description = utils.String(v.(string))
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, templateSpecName, name, templateSpecVersion); err != nil {
		return fmt.Errorf("creating Template Spec Version %q (Resource Group %q / Template Spec %q): %+v", name, resourceGroup, templateSpecName, err)
	}

	d.SetId(id.ID())

	return resourceTemplateSpecVersionRead(d, meta)
}

func resourceTemplateSpecVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.TemplateSpecVersionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TemplateSpecVersionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.TemplateSpecName, id.VersionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Template Spec Version %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Template Spec Version %q (Resource Group %q / Template Spec %q): %+v", id.VersionName, id.ResourceGroup, id.TemplateSpecName, err)
	}

	d.Set("name", id.VersionName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("template_spec_name", id.TemplateSpecName)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.VersionProperties; props != nil {
		if props.Description != nil {
			d.Set("description", props.Description)
		}

		if props.Artifacts != nil {
			artifacts, err := flattenTemplateSpecVersionTemplateArtifacts(props.Artifacts)
			if err != nil {
				return fmt.Errorf("setting `artifact`: %+v", err)
			}
			d.Set("artifact", artifacts)
		}

		template, err := flattenTemplateDeploymentBody(props.Template)
		if err != nil {
			return fmt.Errorf("flattening `template_content`: %+v", err)
		}
		d.Set("template_content", template)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceTemplateSpecVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.TemplateSpecVersionsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TemplateSpecVersionID(d.Id())
	if err != nil {
		return err
	}

	templateSpecVersionUpdateModel := templatespecs.VersionUpdateModel{}

	if d.HasChange("tags") {
		templateSpecVersionUpdateModel.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.TemplateSpecName, id.VersionName, &templateSpecVersionUpdateModel); err != nil {
		return fmt.Errorf("updating Template Spec Version %q (Resource Group %q / Template Spec %q): %+v", id.VersionName, id.ResourceGroup, id.TemplateSpecName, err)
	}

	return resourceTemplateSpecVersionRead(d, meta)
}

func resourceTemplateSpecVersionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.TemplateSpecVersionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TemplateSpecVersionID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.TemplateSpecName, id.VersionName); err != nil {
		return fmt.Errorf("deleting Template Spec Version %q (Resource Group %q / Template Spec %q): %+v", id.VersionName, id.ResourceGroup, id.TemplateSpecName, err)
	}

	return nil
}

func expandTemplateSpecVersionTemplateArtifacts(input []interface{}) (*[]templatespecs.BasicArtifact, error) {
	results := make([]templatespecs.BasicArtifact, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		template, err := expandTemplateDeploymentBody(v["template_content"].(string))
		if err != nil {
			return nil, fmt.Errorf("expanding `template_content`: %+v", err)
		}

		results = append(results, templatespecs.TemplateArtifact{
			Path:     utils.String(v["path"].(string)),
			Template: template,
		})
	}

	return &results, nil
}

func flattenTemplateSpecVersionTemplateArtifacts(input *[]templatespecs.BasicArtifact) ([]interface{}, error) {
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

		template, err := flattenTemplateDeploymentBody(artifact.Template)
		if err != nil {
			return nil, fmt.Errorf("flattening `template_content`: %+v", err)
		}

		results = append(results, map[string]interface{}{
			"path":             path,
			"template_content": template,
		})
	}

	return results, nil
}
