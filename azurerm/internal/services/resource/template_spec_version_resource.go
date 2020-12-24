package resource

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2019-06-01-preview/templatespecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
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

func resourceArmTemplateSpecVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmTemplateSpecVersionCreate,
		Read:   resourceArmTemplateSpecVersionRead,
		Update: resourceArmTemplateSpecVersionUpdate,
		Delete: resourceArmTemplateSpecVersionDelete,

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

			"template_spec_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.TemplateSpecName,
			},

			"artifact": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  templatespecs.KindTemplate,
							ValidateFunc: validation.StringInSlice([]string{
								string(templatespecs.KindTemplate),
							}, false),
						},

						"path": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"template_content": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							ValidateFunc:     validation.StringIsJSON,
							DiffSuppressFunc: structure.SuppressJsonDiff,
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

			"template_content": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmTemplateSpecVersionCreate(d *schema.ResourceData, meta interface{}) error {
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

	template, err := structure.ExpandJsonFromString(d.Get("template_content").(string))
	if err != nil {
		return fmt.Errorf("failed to parse JSON from `template`: %+v", err)
	}

	templateSpecVersion := templatespecs.VersionTemplatespecs{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		VersionProperties: &templatespecs.VersionProperties{
			Artifacts:   expandTemplateSpecVersionTemplateArtifacts(d.Get("artifact").(*schema.Set).List()),
			Description: utils.String(d.Get("description").(string)),
			Template:    template,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, templateSpecName, name, templateSpecVersion); err != nil {
		return fmt.Errorf("creating Template Spec Version %q (Resource Group %q / Template Spec %q): %+v", name, resourceGroup, templateSpecName, err)
	}

	d.SetId(id.ID())

	return resourceArmTemplateSpecVersionRead(d, meta)
}

func resourceArmTemplateSpecVersionRead(d *schema.ResourceData, meta interface{}) error {
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
		d.Set("description", props.Description)

		if err := d.Set("artifact", flattenTemplateSpecVersionTemplateArtifacts(props.Artifacts)); err != nil {
			return fmt.Errorf("setting `artifact`: %+v", err)
		}

		if props.Template != nil {
			tValue := props.Template.(map[string]interface{})
			tStr, _ := structure.FlattenJsonToString(tValue)
			d.Set("template_content", tStr)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmTemplateSpecVersionUpdate(d *schema.ResourceData, meta interface{}) error {
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

	return resourceArmTemplateSpecVersionRead(d, meta)
}

func resourceArmTemplateSpecVersionDelete(d *schema.ResourceData, meta interface{}) error {
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

func expandTemplateSpecVersionTemplateArtifacts(input []interface{}) *[]templatespecs.BasicArtifact {
	results := make([]templatespecs.BasicArtifact, 0)

	for _, item := range input {
		v := item.(map[string]interface{})
		template, _ := structure.ExpandJsonFromString(v["template_content"].(string))

		results = append(results, templatespecs.TemplateArtifact{
			Path:     utils.String(v["path"].(string)),
			Kind:     templatespecs.Kind(v["kind"].(string)),
			Template: template,
		})
	}

	return &results
}

func flattenTemplateSpecVersionTemplateArtifacts(input *[]templatespecs.BasicArtifact) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		artifact := item.(templatespecs.TemplateArtifact)

		var kind templatespecs.Kind
		if artifact.Kind != "" {
			kind = artifact.Kind
		}

		var p string
		if artifact.Path != nil {
			p = *artifact.Path
		}

		var t string
		if artifact.Template != nil {
			t, _ = structure.FlattenJsonToString(artifact.Template.(map[string]interface{}))
		}

		results = append(results, map[string]interface{}{
			"kind":             kind,
			"path":             p,
			"template_content": t,
		})
	}

	return results
}
