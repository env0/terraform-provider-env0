package client

type CloudCredentialIdsInProjectResponse struct {
	CredentialIds []string `json:"credentialIds"`
}

type CloudCredentialsProjectAssignment struct {
	Id           string `json:"id"`
	CredentialId string `json:"credentialId"`
	ProjectId    string `json:"projectId"`
}

func (client *ApiClient) AssignCloudCredentialsToProject(projectId string, credentialId string) (CloudCredentialsProjectAssignment, error) {
	var result CloudCredentialsProjectAssignment

	err := client.http.Put("/credentials/deployment/"+credentialId+"/project/"+projectId, nil, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (client *ApiClient) RemoveCloudCredentialsFromProject(projectId string, credentialId string) error {
	return client.http.Delete("/credentials/deployment/" + credentialId + "/project/" + projectId)
}

func (client *ApiClient) CloudCredentialIdsInProject(projectId string) ([]string, error) {
	var result CloudCredentialIdsInProjectResponse
	err := client.http.Get("/credentials/deployment/project/"+projectId, nil, &result)

	if err != nil {
		return nil, err
	}
	return result.CredentialIds, nil
}
