package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-06-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPrivateLinkService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPrivateLinkServiceCreateUpdate,
		Read:   resourceArmPrivateLinkServiceRead,
		Update: resourceArmPrivateLinkServiceCreateUpdate,
		Delete: resourceArmPrivateLinkServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"auto_approval_subscription_names": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"fqdns": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"ip_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"subnet_id": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     azure.ValidateResourceID,
						},

						"private_ip_address": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"private_ip_address_version": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(network.IPv4),
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.IPv4),
								string(network.IPv6),
							}, false),
						},

						"private_ip_address_allocation": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.Dynamic),
								string(network.Static),
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},
					},
				},
			},

			"load_balancer_frontend_ip_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"visibility_subscription_names": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmPrivateLinkServiceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.PrivateLinkServiceClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM private link service creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	tags := d.Get("tags").(map[string]interface{})

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Private Link Service %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_private_link_service", *existing.ID)
		}
	}

	properties := &network.PrivateLinkServiceProperties{}

	fqdnsProperties := expandArmFqdns(d)
	if fqdnsProperties != nil {
		properties.Fqdns = fqdnsProperties
	}

	visibilityProperties := expandArmVisibility(d)
	if visibilityProperties != nil {
		properties.Visibility = visibilityProperties
	}

	autoApprovalProperties := expandArmAutoApproval(d)
	if autoApprovalProperties != nil {
		properties.AutoApproval = autoApprovalProperties
	}

	frontendIPConfigurationsProperties := expandArmLoadBalancerFrontendIPConfigurations(d)
	if frontendIPConfigurationsProperties != nil {
		properties.LoadBalancerFrontendIPConfigurations = frontendIPConfigurationsProperties
	}

	ipConfigsProperties := expandArmIPConfigurations(d)
	if ipConfigsProperties != nil {
		properties.IPConfigurations = ipConfigsProperties
	}

	privateLinkService := network.PrivateLinkService{
		Name:                         &name,
		Location:                     &location,
		PrivateLinkServiceProperties: properties,
		Tags:                         expandTags(tags),
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, privateLinkService)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Private Link Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Private Link Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Private Link Service %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmPrivateLinkServiceRead(d, meta)
}

func resourceArmPrivateLinkServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.PrivateLinkServiceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["privateLinkServices"]

	resp, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Private Link Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", resp.Name)

	d.Set("resource_group_name", resGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if fqdns := resp.Fqdns; fqdns != nil {
		d.Set("fqdns", fqdns)
	}

	if autoApproval := resp.AutoApproval; autoApproval != nil {
		if autoApprovalSubs := autoApproval.Subscriptions; autoApprovalSubs != nil {
			d.Set("auto_approval_subscription_names", autoApprovalSubs)
		}
	}

	if visibility := resp.Visibility; visibility != nil {
		if visibilitySubs := visibility.Subscriptions; visibilitySubs != nil {
			d.Set("visibility_subscription_names", visibilitySubs)
		}
	}

	frontendIPConfigs := flattenArmLoadBalancerFrontendIPConfigurations(resp.LoadBalancerFrontendIPConfigurations)
	if err := d.Set("load_balancer_frontend_ip_configuration", frontendIPConfigs); err != nil {
		return fmt.Errorf("Error setting `load_balancer_frontend_ip_configuration`: %+v", err)
	}

	ipConfigs := flattenArmIPConfigurations(resp.IPConfigurations)
	if err := d.Set("ip_configuration", ipConfigs); err != nil {
		return fmt.Errorf("Error setting `ip_configuration`: %+v", err)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmPrivateLinkServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.PrivateLinkServiceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["privateLinkServices"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Private Link Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Private Link Service %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func expandArmFqdns(d *schema.ResourceData) *[]string {
	data := d.Get("fqdns").([]interface{})
	if len(data) == 0 {
		return nil
	}

	fqdnsProperties := make([]string, 0)
	for _, fqdn := range data {
		fqdnsProperties = append(fqdnsProperties, fqdn.(string))
	}

	return &fqdnsProperties
}

func expandArmVisibility(d *schema.ResourceData) *network.PrivateLinkServicePropertiesVisibility {
	data := d.Get("visibility_subscription_names").([]interface{})
	if len(data) == 0 {
		return nil
	}

	visibilitySubscriptions := make([]string, 0)
	for _, visibilitySubscription := range data {
		visibilitySubscriptions = append(visibilitySubscriptions, visibilitySubscription.(string))
	}

	visibilityProperties := &network.PrivateLinkServicePropertiesVisibility{
		Subscriptions: &visibilitySubscriptions,
	}

	return visibilityProperties
}

func expandArmAutoApproval(d *schema.ResourceData) *network.PrivateLinkServicePropertiesAutoApproval {
	data := d.Get("auto_approval_subscription_names").([]interface{})
	if len(data) == 0 {
		return nil
	}

	autoApprovalSubscriptions := make([]string, 0)
	for _, autoApprovalSubscription := range data {
		autoApprovalSubscriptions = append(autoApprovalSubscriptions, autoApprovalSubscription.(string))
	}

	autoApprovalProperties := &network.PrivateLinkServicePropertiesAutoApproval{
		Subscriptions: &autoApprovalSubscriptions,
	}

	return autoApprovalProperties
}

func expandArmLoadBalancerFrontendIPConfigurations(d *schema.ResourceData) *[]network.FrontendIPConfiguration {
	data := d.Get("load_balancer_frontend_ip_configuration").([]interface{})
	if len(data) == 0 {
		return nil
	}

	frontendIPConfigurations := make([]network.FrontendIPConfiguration, 0)
	for _, v := range data {
		frontendIPConfiguration := v.(map[string]interface{})
		id := frontendIPConfiguration["id"].(string)
		if id != "" {
			frontendIPConfigurationProperty := network.FrontendIPConfiguration{
				ID: utils.String(id),
			}

			frontendIPConfigurations = append(frontendIPConfigurations, frontendIPConfigurationProperty)
		}
	}

	return &frontendIPConfigurations
}

func expandArmIPConfigurations(d *schema.ResourceData) *[]network.PrivateLinkServiceIPConfiguration {
	data := d.Get("ip_configuration").([]interface{})
	if len(data) == 0 {
		return nil
	}

	ipConfigs := make([]network.PrivateLinkServiceIPConfiguration, 0)
	for _, v := range data {
		data := v.(map[string]interface{})

		properties := network.PrivateLinkServiceIPConfigurationProperties{}
		private_ip_allocation_method := data["private_ip_address_allocation"].(string)
		if private_ip_allocation_method != "" {
			allocationMethod := network.IPAllocationMethod(private_ip_allocation_method)
			properties.PrivateIPAllocationMethod = allocationMethod
		}

		subnet_id := data["subnet_id"].(string)
		if subnet_id != "" {
			properties.Subnet = &network.Subnet{
				ID: &subnet_id,
			}
		}

		private_ip_address := data["private_ip_address"].(string)
		if private_ip_address != "" {
			properties.PrivateIPAddress = &private_ip_address
		}

		private_ip_address_version := network.IPVersion(data["private_ip_address_version"].(string))
		if private_ip_address_version != "" {
			properties.PrivateIPAddressVersion = private_ip_address_version
		}

		ipConfig := network.PrivateLinkServiceIPConfiguration{}

		name := data["name"].(string)
		if name != "" {
			ipConfig.Name = &name
		}

		ipConfig.PrivateLinkServiceIPConfigurationProperties = &properties

		ipConfigs = append(ipConfigs, ipConfig)
	}

	return &ipConfigs
}

func flattenArmLoadBalancerFrontendIPConfigurations(input *[]network.FrontendIPConfiguration) []interface{} {
	frontendIPConfigs := make([]interface{}, 0)
	if input == nil {
		return frontendIPConfigs
	}

	for _, frontendIPConfig := range *input {
		ipConfig := make(map[string]interface{})
		if frontendIPConfig.ID != nil {
			ipConfig["id"] = *frontendIPConfig.ID
		}
		frontendIPConfigs = append(frontendIPConfigs, ipConfig)
	}

	return frontendIPConfigs
}

func flattenArmIPConfigurations(input *[]network.PrivateLinkServiceIPConfiguration) []interface{} {
	ipConfigs := make([]interface{}, 0)
	if input == nil {
		return ipConfigs
	}

	for _, ipConfig := range *input {
		data := make(map[string]interface{})

		data["private_ip_address_allocation"] = ipConfig.PrivateIPAllocationMethod

		data["private_ip_address_version"] = ipConfig.PrivateIPAddressVersion

		if ipConfig.Name != nil {
			data["name"] = *ipConfig.Name
		}

		if ipConfig.PrivateIPAddress != nil {
			data["private_ip_address"] = *ipConfig.PrivateIPAddress
		}

		if ipConfig.Subnet != nil {
			if ipConfig.Subnet.ID != nil {
				data["subnet_id"] = *ipConfig.Subnet.ID
			}
		}

		ipConfigs = append(ipConfigs, data)
	}

	return ipConfigs
}
