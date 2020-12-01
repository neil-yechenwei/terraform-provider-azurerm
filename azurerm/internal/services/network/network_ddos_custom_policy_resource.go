package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmNetworkDDoSCustomPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetworkDDoSCustomPolicyCreateUpdate,
		Read:   resourceArmNetworkDDoSCustomPolicyRead,
		Update: resourceArmNetworkDDoSCustomPolicyCreateUpdate,
		Delete: resourceArmNetworkDDoSCustomPolicyDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.NetworkDDoSCustomPolicyID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"protocol_custom_setting": {
				Type:     schema.TypeSet,
				Optional: true,
				//ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							//ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.DdosCustomPolicyProtocolSyn),
								string(network.DdosCustomPolicyProtocolTCP),
								string(network.DdosCustomPolicyProtocolUDP),
							}, false),
						},

						"source_rate_override": {
							Type:     schema.TypeString,
							Optional: true,
							//ForceNew: true,
						},

						"trigger_rate_override": {
							Type:     schema.TypeString,
							Optional: true,
							//ForceNew: true,
						},

						"trigger_sensitivity_override": {
							Type:     schema.TypeString,
							Optional: true,
							//ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.Default),
								string(network.High),
								string(network.Low),
								string(network.Relaxed),
							}, false),
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmNetworkDDoSCustomPolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.DDOSCustomPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing Network DDoSCustomPolicy %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_network_ddos_custom_policy", *existing.ID)
	}

	parameters := network.DdosCustomPolicy{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("protocol_custom_setting"); ok {
		parameters.DdosCustomPolicyPropertiesFormat = &network.DdosCustomPolicyPropertiesFormat{
			ProtocolCustomSettings: expandArmDDoSCustomPolicyProtocolCustomSettings(v.(*schema.Set).List()),
		}
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("creating Network DDoSCustomPolicy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating future for Network DDoSCustomPolicy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Network DDoSCustomPolicy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Network DDoSCustomPolicy %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmNetworkDDoSCustomPolicyRead(d, meta)
}

func resourceArmNetworkDDoSCustomPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.DDOSCustomPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NetworkDDoSCustomPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] network DDoSCustomPolicy %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Network DDoSCustomPolicy %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.DdosCustomPolicyPropertiesFormat; props != nil && props.ProtocolCustomSettings != nil {
		if err := d.Set("protocol_custom_setting", flattenArmDDoSCustomPolicyProtocolCustomSettings(props.ProtocolCustomSettings)); err != nil {
			return fmt.Errorf("setting `protocol_custom_setting`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmNetworkDDoSCustomPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.DDOSCustomPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NetworkDDoSCustomPolicyID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Network DDoSCustomPolicy %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on deleting future for Network DDoSCustomPolicy %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandArmDDoSCustomPolicyProtocolCustomSettings(input []interface{}) *[]network.ProtocolCustomSettingsFormat {
	results := make([]network.ProtocolCustomSettingsFormat, 0)

	for _, item := range input {
		block := item.(map[string]interface{})
		result := network.ProtocolCustomSettingsFormat{}

		if v, ok := block["protocol"]; ok && v.(string) != "" {
			result.Protocol = network.DdosCustomPolicyProtocol(v.(string))
		}

		if v, ok := block["trigger_rate_override"]; ok && v.(string) != "" {
			result.TriggerRateOverride = utils.String(v.(string))
		}

		if v, ok := block["source_rate_override"]; ok && v.(string) != "" {
			result.SourceRateOverride = utils.String(v.(string))
		}

		if v, ok := block["trigger_sensitivity_override"]; ok && v.(string) != "" {
			result.TriggerSensitivityOverride = network.DdosCustomPolicyTriggerSensitivityOverride(v.(string))
		}

		results = append(results, result)
	}

	return &results
}

func flattenArmDDoSCustomPolicyProtocolCustomSettings(input *[]network.ProtocolCustomSettingsFormat) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var protocol network.DdosCustomPolicyProtocol
		if item.Protocol != "" {
			protocol = item.Protocol
		}

		var sourceRateOverride string
		if item.SourceRateOverride != nil {
			sourceRateOverride = *item.SourceRateOverride
		}

		var triggerRateOverride string
		if item.TriggerRateOverride != nil {
			triggerRateOverride = *item.TriggerRateOverride
		}

		var triggerSensitivityOverride network.DdosCustomPolicyTriggerSensitivityOverride
		if item.TriggerSensitivityOverride != "" {
			triggerSensitivityOverride = item.TriggerSensitivityOverride
		}

		results = append(results, map[string]interface{}{
			"protocol":                     protocol,
			"source_rate_override":         sourceRateOverride,
			"trigger_rate_override":        triggerRateOverride,
			"trigger_sensitivity_override": triggerSensitivityOverride,
		})
	}

	return results
}
