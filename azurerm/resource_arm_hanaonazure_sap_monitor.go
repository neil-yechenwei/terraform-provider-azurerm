package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/preview/hanaonazure/mgmt/2017-11-03-preview/hanaonazure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	azhana "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hanaonazure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmHanaOnAzureSapMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmHanaOnAzureSapMonitorCreate,
		Read:   resourceArmHanaOnAzureSapMonitorRead,
		Update: resourceArmHanaOnAzureSapMonitorUpdate,
		Delete: resourceArmHanaOnAzureSapMonitorDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azhana.ValidateHanaOnAzureSapMonitorName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"hana_db_username": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[\da-zA-Z]{1,32}$`),
					`The name must be between 1 and 32 characters in length and and may contain only letters and numbers.`,
				),
			},

			"hana_db_sql_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1024, 65535),
			},

			"hana_host_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.IPv4Address,
			},

			"hana_subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"hana_db_name": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[\dA-Z]{2,64}$`),
					`The name must be between 2 and 64 characters in length and and may contain only uppercase letters and numbers.`,
				),
			},

			"hana_db_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^(.*){1,64}$`),
					`The name must be between 1 and 64 characters in length.`,
				),
			},

			"key_vault_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  azure.ValidateResourceID,
				ConflictsWith: []string{"hana_db_password"},
			},

			"hana_db_password_key_vault_url": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validate.URLIsHTTPS,
				ConflictsWith: []string{"hana_db_password"},
			},

			"hana_db_credentials_msi_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  azure.ValidateResourceID,
				ConflictsWith: []string{"hana_db_password"},
			},

			"log_analytics_workspace_arm_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"managed_resource_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmHanaOnAzureSapMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).HanaOnAzure.SapMonitorClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Sap Monitor %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_hanaonazure_sap_monitor", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	hanaDbUsername := d.Get("hana_db_username").(string)
	hanaDbSqlPort := int32(d.Get("hana_db_sql_port").(int))
	hanaHostName := d.Get("hana_host_name").(string)
	hanaSubnetId := d.Get("hana_subnet_id").(string)
	hanaDbName := d.Get("hana_db_name").(string)
	hanaDbPassword := d.Get("hana_db_password").(string)
	keyVaultId := d.Get("key_vault_id").(string)
	hanaDbPasswordKeyVaultUrl := d.Get("hana_db_password_key_vault_url").(string)
	hanaDbCredentialsMsiId := d.Get("hana_db_credentials_msi_id").(string)
	t := d.Get("tags").(map[string]interface{})

	sapMonitorParameters := hanaonazure.SapMonitor{
		Location: utils.String(location),
		SapMonitorProperties: &hanaonazure.SapMonitorProperties{
			HanaDbUsername:            utils.String(hanaDbUsername),
			HanaDbSQLPort:             utils.Int32(hanaDbSqlPort),
			HanaSubnet:                utils.String(hanaSubnetId),
			HanaHostname:              utils.String(hanaHostName),
			HanaDbPassword:            utils.String(hanaDbPassword),
			HanaDbName:                utils.String(hanaDbName),
			KeyVaultID:                utils.String(keyVaultId),
			HanaDbPasswordKeyVaultURL: utils.String(hanaDbPasswordKeyVaultUrl),
			HanaDbCredentialsMsiID:    utils.String(hanaDbCredentialsMsiId),
		},
		Tags: tags.Expand(t),
	}

	future, err := client.Create(ctx, resourceGroup, name, sapMonitorParameters)
	if err != nil {
		return fmt.Errorf("Error creating Sap Monitor %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Sap Monitor %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Sap Monitor %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Sap Monitor %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmHanaOnAzureSapMonitorRead(d, meta)
}

func resourceArmHanaOnAzureSapMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).HanaOnAzure.SapMonitorClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["sapMonitors"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Sap Monitor %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Sap Monitor %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.SapMonitorProperties; props != nil {
		d.Set("hana_db_username", props.HanaDbUsername)
		d.Set("hana_db_sql_port", props.HanaDbSQLPort)
		d.Set("hana_host_name", props.HanaHostname)
		d.Set("hana_subnet_id", props.HanaSubnet)
		d.Set("hana_db_password", props.HanaDbPassword)
		d.Set("hana_db_name", props.HanaDbName)
		d.Set("key_vault_id", props.KeyVaultID)
		d.Set("hana_db_password_key_vault_url", props.HanaDbPasswordKeyVaultURL)
		d.Set("hana_db_credentials_msi_id", props.HanaDbCredentialsMsiID)
		d.Set("log_analytics_workspace_arm_id", props.LogAnalyticsWorkspaceArmID)
		d.Set("managed_resource_group_name", props.ManagedResourceGroupName)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmHanaOnAzureSapMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).HanaOnAzure.SapMonitorClient
	ctx, cancel := timeouts.ForUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	newTags := d.Get("tags").(map[string]interface{})

	props := &hanaonazure.Tags{
		Tags: tags.Expand(newTags),
	}

	_, err := client.Update(ctx, resourceGroup, name, *props)
	if err != nil {
		return fmt.Errorf("Error updating Sap Monitor %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return resourceArmHanaOnAzureSapMonitorRead(d, meta)
}

func resourceArmHanaOnAzureSapMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).HanaOnAzure.SapMonitorClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["sapMonitors"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Sap Monitor %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Sap Monitor %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}
