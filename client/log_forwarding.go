package client

type LogForwardingType string

const (
	LogForwardingTypeSplunk             LogForwardingType = "SPLUNK"
	LogForwardingTypeNewRelic           LogForwardingType = "NEWRELIC"
	LogForwardingTypeDatadog            LogForwardingType = "DATADOG"
	LogForwardingTypeSumoLogic          LogForwardingType = "SUMOLOGIC"
	LogForwardingTypeLogzIo             LogForwardingType = "LOGZIO"
	LogForwardingTypeCoralogix          LogForwardingType = "CORALOGIX"
	LogForwardingTypeGoogleCloudLogging LogForwardingType = "GOOGLE_CLOUD_LOGGING"
	LogForwardingTypeGrafanaLoki        LogForwardingType = "GRAFANA_LOKI"
	LogForwardingTypeCloudWatch         LogForwardingType = "CLOUDWATCH"
	LogForwardingTypeS3                 LogForwardingType = "S3"
)

var LogForwardingTypes = []LogForwardingType{
	LogForwardingTypeSplunk,
	LogForwardingTypeNewRelic,
	LogForwardingTypeDatadog,
	LogForwardingTypeSumoLogic,
	LogForwardingTypeLogzIo,
	LogForwardingTypeCoralogix,
	LogForwardingTypeGoogleCloudLogging,
	LogForwardingTypeGrafanaLoki,
	LogForwardingTypeCloudWatch,
	LogForwardingTypeS3,
}

// The API never returns the configuration value, so it is not part of this model.
type LogForwarding struct {
	Id                      string            `json:"id"`
	OrganizationId          string            `json:"organizationId"`
	Name                    string            `json:"name"`
	Type                    LogForwardingType `json:"type"`
	AuditLogForwarding      bool              `json:"auditLogForwarding"`
	DeploymentLogForwarding bool              `json:"deploymentLogForwarding"`
}

type LogForwardingUpsertPayload struct {
	Type                    LogForwardingType `json:"type"`
	Value                   map[string]any    `json:"value"`
	AuditLogForwarding      bool              `json:"auditLogForwarding"`
	DeploymentLogForwarding bool              `json:"deploymentLogForwarding"`
}

type LogForwardingUpsertPayloadWith struct {
	LogForwardingUpsertPayload
	OrganizationId string `json:"organizationId"`
}

func (client *ApiClient) LogForwardings() ([]LogForwarding, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	var result []LogForwarding

	if err := client.http.Get("/log-forwarding/configurations", map[string]string{"organizationId": organizationId}, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// The API keeps one configuration per organization and type: an upsert of an existing type updates it.
func (client *ApiClient) LogForwardingUpsert(payload LogForwardingUpsertPayload) (*LogForwarding, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	payloadWithOrganizationId := LogForwardingUpsertPayloadWith{
		LogForwardingUpsertPayload: payload,
		OrganizationId:             organizationId,
	}

	var result LogForwarding

	if err := client.http.Put("/log-forwarding/configurations", payloadWithOrganizationId, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) LogForwardingDelete(id string) error {
	return client.http.Delete("/log-forwarding/configurations/"+id, nil)
}
