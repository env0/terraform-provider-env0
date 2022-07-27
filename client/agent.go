package client

type Agent struct {
	AgentKey string `json:"agentKey"`
}

func (client *ApiClient) Agents() ([]Agent, error) {
	organizationId, err := client.organizationId()
	if err != nil {
		return nil, err
	}

	var result []Agent
	err = client.http.Get("/agents", map[string]string{"organizationId": organizationId}, &result)
	if err != nil {
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
