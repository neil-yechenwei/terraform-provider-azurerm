package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	azhana "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/hanaonazure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmHanaOnAzureSapMonitor() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmHanaOnAzureSapMonitorRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azhana.ValidateHanaOnAzureSapMonitorName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"hana_db_username": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"hana_db_sql_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"hana_host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"hana_subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"hana_db_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"hana_db_password": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"key_vault_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"hana_db_password_key_vault_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"hana_db_credentials_msi_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"log_analytics_workspace_arm_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"managed_resource_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmHanaOnAzureSapMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).HanaOnAzure.SapMonitorClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Sap Monitor %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading Sap Monitor %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

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
		d.Set("hana_db_name", props.HanaDbName)
		d.Set("hana_db_password", props.HanaDbPassword)
		d.Set("key_vault_id", props.KeyVaultID)
		d.Set("hana_db_password_key_vault_url", props.HanaDbPasswordKeyVaultURL)
		d.Set("hana_db_credentials_msi_id", props.HanaDbCredentialsMsiID)
		d.Set("log_analytics_workspace_arm_id", props.LogAnalyticsWorkspaceArmID)
		d.Set("managed_resource_group_name", props.ManagedResourceGroupName)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
