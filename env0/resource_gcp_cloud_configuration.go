package env0

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGcpCloudConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGcpCloudConfigurationCreate,
		UpdateContext: resourceGcpCloudConfigurationUpdate,
		ReadContext:   readCloudConfiguration,
		DeleteContext: deleteCloudConfiguration,

		Importer: &schema.ResourceImporter{StateContext: importCloudConfiguration},

		Description: "configure a GCP cloud account (Cloud Compass)",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name for the cloud configuration for insights",
				Required:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "the GCP project ID",
				Required:    true,
			},
			"credential_configuration_file_content": {
				Type:        schema.TypeString,
				Description: "the JSON configuration file content containing the service account key",
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

func resourceGcpCloudConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return createCloudConfiguration(d, meta, "GCP")
}

func resourceGcpCloudConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return updateCloudConfiguration(d, meta, "GCP")
}
