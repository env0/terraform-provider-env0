package env0

import (
	"context"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getLogForwardingConfigurationFromSchema(d *schema.ResourceData, configType client.LogForwardingConfigurationType) (map[string]any, error) {
	value := make(map[string]any)

	switch configType {
	case client.LogForwardingConfigurationTypeSplunk:
		value["url"] = d.Get("url").(string)
		value["token"] = d.Get("token").(string)
		value["index"] = d.Get("index").(string)
	default:
		return nil, fmt.Errorf("unsupported log forwarding configuration type: %s", configType)
	}

	return value, nil
}

func createLogForwardingConfiguration(d *schema.ResourceData, meta any, configType client.LogForwardingConfigurationType) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	value, err := getLogForwardingConfigurationFromSchema(d, configType)
	if err != nil {
		return diag.FromErr(err)
	}

	createPayload := &client.LogForwardingConfigurationCreatePayload{
		Type:                    configType,
		Value:                   value,
		AuditLogForwarding:      d.Get("audit_log_forwarding").(bool),
		DeploymentLogForwarding: d.Get("deployment_log_forwarding").(bool),
	}

	logForwardingConfig, err := apiClient.LogForwardingConfigurationCreate(createPayload)
	if err != nil {
		return diag.Errorf("failed to create log forwarding configuration: %v", err)
	}

	d.SetId(logForwardingConfig.Id)

	return nil
}

func readLogForwardingConfiguration(d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	logForwardingConfig, err := apiClient.LogForwardingConfiguration(d.Id())
	if err != nil {
		return ResourceGetFailure(context.Background(), "log forwarding configuration", d, err)
	}

	if err := d.Set("audit_log_forwarding", logForwardingConfig.AuditLogForwarding); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("deployment_log_forwarding", logForwardingConfig.DeploymentLogForwarding); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func updateLogForwardingConfiguration(d *schema.ResourceData, meta any, configType client.LogForwardingConfigurationType) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	value, err := getLogForwardingConfigurationFromSchema(d, configType)
	if err != nil {
		return diag.FromErr(err)
	}

	updatePayload := &client.LogForwardingConfigurationUpdatePayload{
		Value:                   value,
		AuditLogForwarding:      d.Get("audit_log_forwarding").(bool),
		DeploymentLogForwarding: d.Get("deployment_log_forwarding").(bool),
	}

	_, err = apiClient.LogForwardingConfigurationUpdate(d.Id(), updatePayload)
	if err != nil {
		return diag.Errorf("failed to update log forwarding configuration: %v", err)
	}

	return readLogForwardingConfiguration(d, meta)
}

func deleteLogForwardingConfiguration(d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.LogForwardingConfigurationDelete(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
