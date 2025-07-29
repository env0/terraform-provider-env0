package client

import "fmt"

type LogForwardingConfigurationType string

const (
	LogForwardingConfigurationTypeSplunk            LogForwardingConfigurationType = "SPLUNK"
	LogForwardingConfigurationTypeNewRelic          LogForwardingConfigurationType = "NEWRELIC"
	LogForwardingConfigurationTypeDatadog           LogForwardingConfigurationType = "DATADOG"
	LogForwardingConfigurationTypeSumoLogic         LogForwardingConfigurationType = "SUMOLOGIC"
	LogForwardingConfigurationTypeLogzIo            LogForwardingConfigurationType = "LOGZIO"
	LogForwardingConfigurationTypeCoralogix         LogForwardingConfigurationType = "CORALOGIX"
	LogForwardingConfigurationTypeGoogleCloudLogging LogForwardingConfigurationType = "GOOGLE_CLOUD_LOGGING"
	LogForwardingConfigurationTypeGrafanaLoki       LogForwardingConfigurationType = "GRAFANA_LOKI"
	LogForwardingConfigurationTypeCloudWatch        LogForwardingConfigurationType = "CLOUDWATCH"
	LogForwardingConfigurationTypeS3                LogForwardingConfigurationType = "S3"
)

type LogForwardingConfiguration struct {
	Id                      string                         `json:"id"`
	OrganizationId          string                         `json:"organizationId"`
	Type                    LogForwardingConfigurationType `json:"type"`
	Value                   map[string]interface{}         `json:"value"`
	AuditLogForwarding      bool                           `json:"auditLogForwarding"`
	DeploymentLogForwarding bool                           `json:"deploymentLogForwarding"`
}

type LogForwardingConfigurationCreatePayload struct {
	Type                    LogForwardingConfigurationType `json:"type"`
	Value                   map[string]interface{}         `json:"value"`
	AuditLogForwarding      bool                           `json:"auditLogForwarding"`
	DeploymentLogForwarding bool                           `json:"deploymentLogForwarding"`
	OrganizationId          string                         `json:"organizationId,omitempty"`
}

type LogForwardingConfigurationUpdatePayload struct {
	Value                   map[string]interface{} `json:"value"`
	AuditLogForwarding      bool                   `json:"auditLogForwarding"`
	DeploymentLogForwarding bool                   `json:"deploymentLogForwarding"`
}

func (client *ApiClient) LogForwardingConfigurationCreate(payload *LogForwardingConfigurationCreatePayload) (*LogForwardingConfiguration, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, fmt.Errorf("failed to get organization id: %w", err)
	}

	payloadWithOrganizationId := *payload
	payloadWithOrganizationId.OrganizationId = organizationId

	var logForwardingConfig LogForwardingConfiguration
	if err := client.http.Put("/log-forwarding/configurations", &payloadWithOrganizationId, &logForwardingConfig); err != nil {
		return nil, err
	}

	return &logForwardingConfig, nil
}

func (client *ApiClient) LogForwardingConfigurationUpdate(id string, payload *LogForwardingConfigurationUpdatePayload) (*LogForwardingConfiguration, error) {
	var logForwardingConfig LogForwardingConfiguration
	if err := client.http.Put("/log-forwarding/configurations/"+id, payload, &logForwardingConfig); err != nil {
		return nil, err
	}

	return &logForwardingConfig, nil
}

func (client *ApiClient) LogForwardingConfigurationDelete(id string) error {
	return client.http.Delete("/log-forwarding/configurations/"+id, nil)
}

func (client *ApiClient) LogForwardingConfiguration(id string) (*LogForwardingConfiguration, error) {
	var logForwardingConfig LogForwardingConfiguration
	if err := client.http.Get("/log-forwarding/configurations/"+id, nil, &logForwardingConfig); err != nil {
		return nil, err
	}

	return &logForwardingConfig, nil
}

func (client *ApiClient) LogForwardingConfigurations() ([]LogForwardingConfiguration, error) {
	var logForwardingConfigs []LogForwardingConfiguration
	if err := client.http.Get("/log-forwarding/configurations", nil, &logForwardingConfigs); err != nil {
		return nil, err
	}

	return logForwardingConfigs, nil
}
