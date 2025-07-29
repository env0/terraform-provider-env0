package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSplunkLogForwarding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSplunkLogForwardingCreate,
		ReadContext:   resourceSplunkLogForwardingRead,
		UpdateContext: resourceSplunkLogForwardingUpdate,
		DeleteContext: resourceSplunkLogForwardingDelete,

		Importer: &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},

		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Splunk URL",
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The Splunk token",
			},
			"index": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Splunk index",
			},
			"audit_log_forwarding": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to forward audit logs",
			},
			"deployment_log_forwarding": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to forward deployment logs",
			},
		},
	}
}

func resourceSplunkLogForwardingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return createLogForwardingConfiguration(d, meta, client.LogForwardingConfigurationTypeSplunk)
}

func resourceSplunkLogForwardingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	logForwardingConfig, err := apiClient.LogForwardingConfiguration(d.Id())
	if err != nil {
		return ResourceGetFailure(ctx, "log forwarding configuration", d, err)
	}

	if err := d.Set("audit_log_forwarding", logForwardingConfig.AuditLogForwarding); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("deployment_log_forwarding", logForwardingConfig.DeploymentLogForwarding); err != nil {
		return diag.FromErr(err)
	}

	if value, ok := logForwardingConfig.Value["url"]; ok {
		d.Set("url", value)
	}

	if value, ok := logForwardingConfig.Value["index"]; ok {
		d.Set("index", value)
	}

	return nil
}

func resourceSplunkLogForwardingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return updateLogForwardingConfiguration(d, meta, client.LogForwardingConfigurationTypeSplunk)
}

func resourceSplunkLogForwardingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return deleteLogForwardingConfiguration(d, meta)
}
