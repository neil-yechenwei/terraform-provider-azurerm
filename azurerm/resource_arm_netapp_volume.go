package azurerm

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2019-06-01/netapp"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmNetAppVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetAppVolumeCreateUpdate,
		Read:   resourceArmNetAppVolumeRead,
		Update: resourceArmNetAppVolumeCreateUpdate,
		Delete: resourceArmNetAppVolumeDelete,

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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"export_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rules": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed_clients": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"cifs": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"nfsv3": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"nfsv4": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"rule_index": {
										Type:         schema.TypeInt,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.IntBetween(1, 6),
									},
									"unix_read_only": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"unix_read_write": {
										Type:     schema.TypeBool,
										Optional: true,
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
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(netapp.Premium),
					string(netapp.Standard),
					string(netapp.Ultra),
				}, false),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"usage_threshold": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(107374182400, 109951162780000),
			},

			"file_system_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmNetAppVolumeCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).netapp.VolumeClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("account_name").(string)
	poolName := d.Get("pool_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, accountName, poolName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for present of existing NetApp Volume %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_netapp_volume", *resp.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	creationToken := d.Get("creation_token").(string)
	exportPolicy := d.Get("export_policy").([]interface{})
	serviceLevel := d.Get("service_level").(string)
	subnetId := d.Get("subnet_id").(string)
	usageThreshold := int64(d.Get("usage_threshold").(int))
	t := d.Get("tags").(map[string]interface{})

	volumeParameters := netapp.Volume{
		Location: utils.String(location),
		VolumeProperties: &netapp.VolumeProperties{
			CreationToken:  utils.String(creationToken),
			ExportPolicy:   expandArmExportPolicy(exportPolicy),
			ServiceLevel:   netapp.ServiceLevel(serviceLevel),
			SubnetID:       utils.String(subnetId),
			UsageThreshold: utils.Int64(usageThreshold),
		},
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, volumeParameters, resourceGroup, accountName, poolName, name)
	if err != nil {
		return fmt.Errorf("Error creating NetApp Volume %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of NetApp Volume %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, accountName, poolName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving NetApp Volume %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read NetApp Volume %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmNetAppVolumeRead(d, meta)
}

func resourceArmNetAppVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).netapp.VolumeClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	accountName := id.Path["netAppAccounts"]
	poolName := id.Path["capacityPools"]
	name := id.Path["volumes"]

	resp, err := client.Get(ctx, resourceGroup, accountName, poolName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] NetApp Volumes %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading NetApp Volumes %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

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

		if err := d.Set("file_system_id", volumeProperties.FileSystemID); err != nil {
			return fmt.Errorf("Error setting `file_system_id`: %+v", err)
		}
	}

	return tags.FlattenAndSetTags(d, resp.Tags)
}

func resourceArmNetAppVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).netapp.VolumeClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	accountName := id.Path["netAppAccounts"]
	poolName := id.Path["capacityPools"]
	name := id.Path["volumes"]

	future, err := client.Delete(ctx, resourceGroup, accountName, poolName, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting NetApp Volume %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return waitForNetAppVolumeToBeDeleted(ctx, client, resourceGroup, accountName, poolName, name)
}

func waitForNetAppVolumeToBeDeleted(ctx context.Context, client *netapp.VolumesClient, resourceGroup, accountName, poolName, name string) error {
	log.Printf("[DEBUG] Waiting for NetApp Volume Provisioning Service %q (Resource Group %q) to be deleted", name, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"200", "202"},
		Target:  []string{"404"},
		Refresh: netappVolumeDeleteStateRefreshFunc(ctx, client, resourceGroup, accountName, poolName, name),
		Timeout: 20 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for NetApp Volume Provisioning Service %q (Resource Group %q) to be deleted: %+v", name, resourceGroup, err)
	}

	return nil
}

func netappVolumeDeleteStateRefreshFunc(ctx context.Context, client *netapp.VolumesClient, resourceGroupName string, accountName string, poolName string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, accountName, poolName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(res.Response) {
				return nil, "", fmt.Errorf("Error retrieving NetApp Volume %q (Resource Group %q): %s", name, resourceGroupName, err)
			}
		}

		if _, err := client.Delete(ctx, resourceGroupName, accountName, poolName, name); err != nil {
			log.Printf("Error reissuing NetApp Volume %q delete request (Resource Group %q): %+v", name, resourceGroupName, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}

func expandArmExportPolicy(input []interface{}) *netapp.VolumePropertiesExportPolicy {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	rules := v["rules"].([]interface{})

	result := netapp.VolumePropertiesExportPolicy{
		Rules: expandArmExportPolicyRule(rules),
	}
	return &result
}

func expandArmExportPolicyRule(input []interface{}) *[]netapp.ExportPolicyRule {
	results := make([]netapp.ExportPolicyRule, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		allowedClients := v["allowed_clients"].(string)
		cifs := v["cifs"].(bool)
		nfsv3 := v["nfsv3"].(bool)
		nfsv4 := v["nfsv4"].(bool)
		ruleIndex := int32(v["rule_index"].(int))
		unixReadOnly := v["unix_read_only"].(bool)
		unixReadWrite := v["unix_read_write"].(bool)

		result := netapp.ExportPolicyRule{
			AllowedClients: utils.String(allowedClients),
			Cifs:           utils.Bool(cifs),
			Nfsv3:          utils.Bool(nfsv3),
			Nfsv4:          utils.Bool(nfsv4),
			RuleIndex:      utils.Int32(ruleIndex),
			UnixReadOnly:   utils.Bool(unixReadOnly),
			UnixReadWrite:  utils.Bool(unixReadWrite),
		}

		results = append(results, result)
	}
	return &results
}

func flattenArmExportPolicy(input *netapp.VolumePropertiesExportPolicy) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	result["rules"] = flattenArmExportPolicyRule(input.Rules)

	return []interface{}{result}
}

func flattenArmExportPolicyRule(input *[]netapp.ExportPolicyRule) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		if allowedClients := item.AllowedClients; allowedClients != nil {
			v["allowed_clients"] = *allowedClients
		}

		if cifs := item.Cifs; cifs != nil {
			v["cifs"] = *cifs
		}

		if nfsv3 := item.Nfsv3; nfsv3 != nil {
			v["nfsv3"] = *nfsv3
		}

		if nfsv4 := item.Nfsv4; nfsv4 != nil {
			v["nfsv4"] = *nfsv4
		}

		if ruleIndex := item.RuleIndex; ruleIndex != nil {
			v["rule_index"] = *ruleIndex
		}

		if unixReadOnly := item.UnixReadOnly; unixReadOnly != nil {
			v["unix_read_only"] = *unixReadOnly
		}

		if unixReadWrite := item.UnixReadWrite; unixReadWrite != nil {
			v["unix_read_write"] = *unixReadWrite
		}

		results = append(results, v)
	}

	return results
}
