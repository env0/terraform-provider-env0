package client

type LogForwardingConfiguration struct {
	Id                      string `json:"id,omitempty"`
	Type                    string `json:"type"`
	Value                   any    `json:"value"`
	AuditLogForwarding      *bool  `json:"auditLogForwarding,omitempty"`
	DeploymentLogForwarding *bool  `json:"deploymentLogForwarding,omitempty"`
}

type LogForwardingConfigurationCreatePayload struct {
	Type                    string `json:"type"`
	Value                   any    `json:"value"`
	AuditLogForwarding      *bool  `json:"auditLogForwarding,omitempty"`
	DeploymentLogForwarding *bool  `json:"deploymentLogForwarding,omitempty"`
}

type LogForwardingConfigurationUpdatePayload struct {
	Id                      string `json:"id"`
	Type                    string `json:"type"`
	Value                   any    `json:"value"`
	AuditLogForwarding      *bool  `json:"auditLogForwarding,omitempty"`
	DeploymentLogForwarding *bool  `json:"deploymentLogForwarding,omitempty"`
}

func (client *ApiClient) LogForwardingConfigurations() ([]LogForwardingConfiguration, error) {
	var result []LogForwardingConfiguration
	err := client.http.Get("/log-forwarding/configurations", nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (client *ApiClient) LogForwardingConfiguration(id string) (LogForwardingConfiguration, error) {
	configurations, err := client.LogForwardingConfigurations()
	if err != nil {
		return LogForwardingConfiguration{}, err
	}

	for _, config := range configurations {
		if config.Id == id {
			return config, nil
		}
	}

	return LogForwardingConfiguration{}, &NotFoundError{}
}

func (client *ApiClient) LogForwardingConfigurationCreate(payload LogForwardingConfigurationCreatePayload) (LogForwardingConfiguration, error) {
	var result LogForwardingConfiguration
	err := client.http.Post("/log-forwarding/configurations", payload, &result)
	if err != nil {
		return LogForwardingConfiguration{}, err
	}
	return result, nil
}

func (client *ApiClient) LogForwardingConfigurationUpdate(payload LogForwardingConfigurationUpdatePayload) (LogForwardingConfiguration, error) {
	var result LogForwardingConfiguration
	err := client.http.Post("/log-forwarding/configurations", payload, &result)
	if err != nil {
		return LogForwardingConfiguration{}, err
	}
	return result, nil
}

func (client *ApiClient) LogForwardingConfigurationDelete(id string) error {
	return client.http.Delete("/log-forwarding/configurations/"+id, nil)
}
