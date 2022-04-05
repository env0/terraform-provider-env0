package client

type WorkflowTrigger struct {
	Id string `json:"id"`
}

type WorkflowTriggerUpsertPayload struct {
	DownstreamEnvironmentIds []string `json:"downstreamEnvironmentIds"`
}

func (client *ApiClient) WorkflowTrigger(environmentId string) ([]WorkflowTrigger, error) {
	var result []WorkflowTrigger
	err := client.http.Get("environments/"+environmentId+"/downstream", nil, &result)
	if err != nil {
		return []WorkflowTrigger{}, err
	}

	return result, nil
}

func (client *ApiClient) WorkflowTriggerUpsert(environmentId string, request WorkflowTriggerUpsertPayload) ([]WorkflowTrigger, error) {
	var result []WorkflowTrigger

	err := client.http.Put("/environments/"+environmentId+"/downstream", request, &result)
	if err != nil {
		return []WorkflowTrigger{}, err
	}
	return result, nil
}
