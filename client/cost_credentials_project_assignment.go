package client

type CostCredentialProjectAssignment struct {
	ProjectId       string `json:"projectId"`
	CredentialsId   string `json:"credentialsId"`
	CredentialsType string `json:"credentialsType"`
}

func (client *ApiClient) AssignCostCredentialsToProject(projectId string, credentialId string) (CostCredentialProjectAssignment, error) {
	var result CostCredentialProjectAssignment

	err := client.http.Put("/costs/project/"+projectId+"/credentials", map[string]string{
		"credentialsId": credentialId,
	}, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (client *ApiClient) RemoveCostCredentialsFromProject(projectId string, credentialId string) error {
	return client.http.Delete("/costs/project/"+projectId+"/credentials/"+credentialId, nil)
}

func (client *ApiClient) CostCredentialIdsInProject(projectId string) ([]CostCredentialProjectAssignment, error) {
	var result []CostCredentialProjectAssignment
	err := client.http.Get("/costs/project/"+projectId+"/credentials", nil, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}
