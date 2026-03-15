package client

type Agent struct {
	AgentKey string `json:"agentKey"`
}

func (client *ApiClient) Agents() ([]Agent, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	var result []Agent

	if err := client.http.Get("/agents", map[string]string{"organizationId": organizationId}, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (client *ApiClient) AgentValues(id string) (string, error) {
	var result string
	if err := client.http.Get("/agents/"+id+"/values", nil, &result); err != nil {
		return "", err
	}

	return result, nil
}

// Agent Pool types and methods (self-service agents)

type AgentPoolSelfHostedLogs struct {
	AccountId  string `json:"accountId"`
	Region     string `json:"region"`
	ExternalId string `json:"externalId,omitempty"`
}

type AgentPoolDynamoLogs struct {
	SelfHosted *AgentPoolSelfHostedLogs `json:"selfHosted,omitempty"`
}

type AgentPoolLogsConfig struct {
	Dynamo *AgentPoolDynamoLogs `json:"dynamo,omitempty"`
}

type AgentPool struct {
	Id          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	AgentKey    string               `json:"agentKey"`
	CreatedBy   string               `json:"createdBy,omitempty"`
	Logs        *AgentPoolLogsConfig `json:"logs,omitempty"`
}

type AgentPoolCreatePayload struct {
	OrganizationId string `json:"organizationId"`
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
}

type AgentPoolUpdatePayload struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Logs        *AgentPoolLogsConfig `json:"logs"`
}

type AgentSecret struct {
	Id          string `json:"id"`
	Secret      string `json:"secret"`
	AgentId     string `json:"agentId"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
	CreatedBy   User   `json:"createdBy,omitempty"`
}

type AgentSecretCreatePayload struct {
	Description string `json:"description,omitempty"`
}

func (client *ApiClient) AgentPools() ([]AgentPool, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	var result []AgentPool

	if err := client.http.Get("/agents", map[string]string{"organizationId": organizationId}, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (client *ApiClient) AgentPoolCreate(payload AgentPoolCreatePayload) (*AgentPool, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	payload.OrganizationId = organizationId

	var result AgentPool
	if err := client.http.Post("/agents", payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) AgentPool(id string) (*AgentPool, error) {
	var result AgentPool
	if err := client.http.Get("/agents/"+id, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) AgentPoolUpdate(id string, payload AgentPoolUpdatePayload) (*AgentPool, error) {
	var result AgentPool
	if err := client.http.Patch("/agents/"+id, payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) AgentPoolDelete(id string) error {
	return client.http.Delete("/agents/"+id, nil)
}

func (client *ApiClient) AgentSecretCreate(agentId string, payload AgentSecretCreatePayload) (*AgentSecret, error) {
	var result AgentSecret
	if err := client.http.Post("/agents/"+agentId+"/secrets", payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) AgentSecrets(agentId string) ([]AgentSecret, error) {
	var result []AgentSecret
	if err := client.http.Get("/agents/"+agentId+"/secrets", nil, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (client *ApiClient) AgentSecretDelete(agentId string, secretId string) error {
	return client.http.Delete("/agents/"+agentId+"/secrets/"+secretId, nil)
}
