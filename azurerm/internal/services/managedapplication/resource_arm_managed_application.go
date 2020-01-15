package managedapplication

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-06-01/managedapplications"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmManagedApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmManagedApplicationCreateOrUpdate,
		Read:   resourceArmManagedApplicationRead,
		Update: resourceArmManagedApplicationCreateOrUpdate,
		Delete: resourceArmManagedApplicationDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateManagedAppName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"kind": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"MarketPlace",
					"ServiceCatalog",
				}, false),
			},

			"managed_resource_group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"application_definition_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"parameters": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},

			"plan": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"publisher": {
							Type:     schema.TypeString,
							Required: true,
						},
						"product": {
							Type:     schema.TypeString,
							Required: true,
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
						},
						"promotion_code": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmManagedApplicationCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroupName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Managed Application (Managed Application Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_managed_application", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	kind := d.Get("kind").(string)
	managedResourceGroupId := d.Get("managed_resource_group_id").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := managedapplications.Application{
		Location: utils.String(location),
		Kind:     utils.String(kind),
		ApplicationProperties: &managedapplications.ApplicationProperties{
			ManagedResourceGroupID: utils.String(managedResourceGroupId),
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("application_definition_id"); ok {
		parameters.ApplicationDefinitionID = utils.String(v.(string))
	}
	if _, ok := d.GetOk("plan"); ok {
		parameters.Plan = expandArmManagedApplicationPlan(d)
	}
	if v, ok := d.GetOk("parameters"); ok {
		expandedParams, err := structure.ExpandJsonFromString(v.(string))
		if err != nil {
			return fmt.Errorf("Error expanding JSON from Parameters %q: %+v", v, err)
		}

		parameters.Parameters = &expandedParams
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroupName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Managed Application (Managed Application Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Managed Application (Managed Application Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}

	resp, err := client.Get(ctx, resourceGroupName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Managed Application (Managed Application Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read Managed Application (Managed Application Name %q / Resource Group %q) ID", name, resourceGroupName)
	}
	d.SetId(*resp.ID)

	return resourceArmManagedApplicationRead(d, meta)
}

func resourceArmManagedApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroupName := id.ResourceGroup
	name := id.Path["applications"]

	resp, err := client.Get(ctx, resourceGroupName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Managed Application %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Managed Application (Managed Application Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroupName)
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

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmManagedApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroupName := id.ResourceGroup
	name := id.Path["applications"]

	future, err := client.Delete(ctx, resourceGroupName, name)
	if err != nil {
		return fmt.Errorf("Error deleting Managed Application (Managed Application Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Managed Application (Managed Application Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
		}
	}

	return nil
}

func expandArmManagedApplicationPlan(d *schema.ResourceData) *managedapplications.Plan {
	plans := d.Get("plan").(*schema.Set).List()
	plan := plans[0].(map[string]interface{})

	return &managedapplications.Plan{
		Name:          utils.String(plan["name"].(string)),
		Publisher:     utils.String(plan["publisher"].(string)),
		Product:       utils.String(plan["product"].(string)),
		Version:       utils.String(plan["version"].(string)),
		PromotionCode: utils.String(plan["promotion_code"].(string)),
	}
}

func flattenArmManagedApplicationPlan(input *managedapplications.Plan) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	name := ""
	if input.Name != nil {
		name = *input.Name
	}
	publisher := ""
	if input.Publisher != nil {
		publisher = *input.Publisher
	}
	product := ""
	if input.Product != nil {
		product = *input.Product
	}
	version := ""
	if input.Version != nil {
		version = *input.Version
	}
	promotionCode := ""
	if input.PromotionCode != nil {
		promotionCode = *input.PromotionCode
	}

	results = append(results, map[string]interface{}{
		"name":           name,
		"publisher":      publisher,
		"product":        product,
		"version":        version,
		"promotion_code": promotionCode,
	})

	return results
}
