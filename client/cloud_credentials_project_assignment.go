package client

func (self *ApiClient) AssignCloudCredentialsToProject(projectId string, credentialId string) (CloudCredentialsProjectAssignment, error) {
	var result CloudCredentialsProjectAssignment

	err := self.http.Put("/credentials/deployment/"+credentialId+"/project/"+projectId, nil, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (self *ApiClient) RemoveCloudCredentialsFromProject(projectId string, credentialId string) error {
	return self.http.Delete("/credentials/deployment/" + credentialId + "/project/" + projectId)
}

func (self *ApiClient) CloudCredentialIdsInProject(projectId string) ([]string, error) {
	var result CloudCredentialIdsInProjectResponse
	err := self.http.Get("/credentials/deployment/project/"+projectId, nil, &result)

	if err != nil {
		return nil, err
	}
	return result.CredentialIds, nil
}
