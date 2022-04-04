package client

type WorkflowTrigger struct {
	Id string `json:"id"`
}

type WorkflowTriggerUpsertPayload struct {
	DownstreamEnvironmentIds []string `json:"downstreamEnvironmentIds"`
}

func (self *ApiClient) WorkflowTrigger(environmentId string) ([]WorkflowTrigger, error) {
	var result []WorkflowTrigger
	err := self.http.Get("environments/"+environmentId+"/downstream", nil, &result)
	if err != nil {
		return []WorkflowTrigger{}, err
	}

	return result, nil
}

func (self *ApiClient) WorkflowTriggerUpsert(environmentId string, request WorkflowTriggerUpsertPayload) ([]WorkflowTrigger, error) {
	var result []WorkflowTrigger

	err := self.http.Put("/environments/"+environmentId+"/downstream", request, &result)
	if err != nil {
		return []WorkflowTrigger{}, err
	}
	return result, nil
}
