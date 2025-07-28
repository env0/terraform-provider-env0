package env0

import (
	"context"
	"encoding/json"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLogForwardingConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLogForwardingConfigurationCreate,
		ReadContext:   resourceLogForwardingConfigurationRead,
		UpdateContext: resourceLogForwardingConfigurationUpdate,
		DeleteContext: resourceLogForwardingConfigurationDelete,
		Description:   "Manages log forwarding configuration for audit and deployment logs",

		Importer: &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Description: "type of log forwarding configuration (e.g., 'SPLUNK', 'NEWRELIC', 'DATADOG', 'SUMOLOGIC', 'LOGZIO', 'CORALOGIX', 'GOOGLE_CLOUD_LOGGING', 'GRAFANA_LOKI', 'CLOUDWATCH', 'S3')",
				Required:    true,
			},
			"value": {
				Type:        schema.TypeString,
				Description: "configuration value as JSON string containing connection details",
				Required:    true,
				Sensitive:   true,
			},
			"audit_log_forwarding": {
				Type:        schema.TypeBool,
				Description: "whether to forward audit logs",
				Optional:    true,
				Default:     false,
			},
			"deployment_log_forwarding": {
				Type:        schema.TypeBool,
				Description: "whether to forward deployment logs",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceLogForwardingConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var value any
	valueStr := d.Get("value").(string)
	if err := json.Unmarshal([]byte(valueStr), &value); err != nil {
		return diag.Errorf("invalid JSON in value field: %v", err)
	}

	payload := client.LogForwardingConfigurationCreatePayload{
		Type:                    d.Get("type").(string),
		Value:                   value,
		AuditLogForwarding:      boolPtr(d.Get("audit_log_forwarding").(bool)),
		DeploymentLogForwarding: boolPtr(d.Get("deployment_log_forwarding").(bool)),
	}

	configuration, err := apiClient.LogForwardingConfigurationCreate(payload)
	if err != nil {
		return diag.Errorf("could not create log forwarding configuration: %v", err)
	}

	d.SetId(configuration.Id)

	return nil
}

func resourceLogForwardingConfigurationRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	configuration, err := apiClient.LogForwardingConfiguration(id)
	if err != nil {
		return ResourceGetFailure(ctx, "log forwarding configuration", d, err)
	}

	d.Set("type", configuration.Type)

	if configuration.Value != nil {
		valueBytes, err := json.Marshal(configuration.Value)
		if err != nil {
			return diag.Errorf("could not marshal configuration value: %v", err)
		}
		d.Set("value", string(valueBytes))
	}

	if configuration.AuditLogForwarding != nil {
		d.Set("audit_log_forwarding", *configuration.AuditLogForwarding)
	}
	if configuration.DeploymentLogForwarding != nil {
		d.Set("deployment_log_forwarding", *configuration.DeploymentLogForwarding)
	}

	return nil
}

func resourceLogForwardingConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var value any
	valueStr := d.Get("value").(string)
	if err := json.Unmarshal([]byte(valueStr), &value); err != nil {
		return diag.Errorf("invalid JSON in value field: %v", err)
	}

	payload := client.LogForwardingConfigurationUpdatePayload{
		Id:                      d.Id(),
		Type:                    d.Get("type").(string),
		Value:                   value,
		AuditLogForwarding:      boolPtr(d.Get("audit_log_forwarding").(bool)),
		DeploymentLogForwarding: boolPtr(d.Get("deployment_log_forwarding").(bool)),
	}

	_, err := apiClient.LogForwardingConfigurationUpdate(payload)
	if err != nil {
		return diag.Errorf("could not update log forwarding configuration: %v", err)
	}

	return nil
}

func resourceLogForwardingConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	err := apiClient.LogForwardingConfigurationDelete(id)
	if err != nil {
		return diag.Errorf("could not delete log forwarding configuration: %v", err)
	}

	return nil
}
