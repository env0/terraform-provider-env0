package client

type CostCredentialProjectAssignment struct {
	ProjectId       string `json:"projectId"`
	CredentialsId   string `json:"credentialsId"`
	CredentialsType string `json:"credentialsType"`
}

func (self *ApiClient) AssignCostCredentialsToProject(projectId string, credentialId string) (CostCredentialProjectAssignment, error) {
	var result CostCredentialProjectAssignment

	err := self.http.Put("/costs/project/"+projectId+"/credentials", credentialId, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (self *ApiClient) RemoveCostCredentialsFromProject(projectId string, credentialId string) error {
	return self.http.Delete("/costs/project/" + projectId + "/credentials/" + credentialId)
}

func (self *ApiClient) CostCredentialIdsInProject(projectId string) ([]CostCredentialProjectAssignment, error) {
	var result []CostCredentialProjectAssignment
	err := self.http.Get("/costs/project/"+projectId, nil, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}
