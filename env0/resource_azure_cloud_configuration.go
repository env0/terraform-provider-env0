package env0

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAzureCloudConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureCloudConfigurationCreate,
		UpdateContext: resourceAzureCloudConfigurationUpdate,
		ReadContext:   readCloudConfiguration,
		DeleteContext: deleteCloudConfiguration,

		Importer: &schema.ResourceImporter{StateContext: importCloudConfiguration},

		Description: "configure an Azure cloud account (Cloud Compass)",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name for the cloud configuration for insights",
				Required:    true,
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Description: "the Azure tenant ID",
				Required:    true,
			},
			"client_id": {
				Type:        schema.TypeString,
				Description: "the Azure client ID",
				Required:    true,
			},
			"log_analytics_workspace_id": {
				Type:        schema.TypeString,
				Description: "the Azure Log Analytics Workspace ID",
				Required:    true,
			},
			"health": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "an indicator if the configuration is valid",
			},
		},
	}
}

func resourceAzureCloudConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return createCloudConfiguration(d, meta, "AzureLAW")
}

func resourceAzureCloudConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return updateCloudConfiguration(d, meta, "AzureLAW")
}
