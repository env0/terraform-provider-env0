package env0

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/time/rate"
)

// Limit to a burst of 5 and accumulate 1 request every 20 seconds
var awsCloudConfigRateLimiter = rate.NewLimiter(rate.Limit(0.05), 5)

func resourceAwsCloudConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAwsCloudConfigurationCreate,
		UpdateContext: resourceAwsCloudConfigurationUpdate,
		ReadContext:   readCloudConfiguration,
		DeleteContext: deleteCloudConfiguration,

		Importer: &schema.ResourceImporter{StateContext: importCloudConfiguration},

		Description: "configure an AWS cloud account (Cloud Compass)",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name for the cloud configuration for insights",
				Required:    true,
			},
			"account_id": {
				Type:        schema.TypeString,
				Description: "the AWS account id",
				Required:    true,
			},
			"bucket_name": {
				Type:        schema.TypeString,
				Description: "the CloudTrail bucket name",
				Required:    true,
			},
			"regions": {
				Type:        schema.TypeList,
				Description: "a list of AWS regions",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "an optional bucket prefix (folder)",
			},
			"health": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "an indicator if the configuration is valid",
			},
			"should_prefix_under_logs_folder": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If 'true' than the prefix will be under 'AWSLogs' folder (default: false)",
			},
		},
	}
}

func resourceAwsCloudConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if err := awsCloudConfigRateLimiter.Wait(ctx); err != nil {
		return diag.Errorf("rate limit wait error: %v", err)
	}

	return createCloudConfiguration(d, meta, "AWS")
}

func resourceAwsCloudConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if err := awsCloudConfigRateLimiter.Wait(ctx); err != nil {
		return diag.Errorf("rate limit wait error: %v", err)
	}

	return updateCloudConfiguration(d, meta, "AWS")
}
