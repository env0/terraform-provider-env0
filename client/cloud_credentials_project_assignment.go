package client

import (
	"fmt"
)

func (self *ApiClient) AssignCloudCredentialsToProject(projectId string, credentialId string) (CloudCredentialsProjectAssignment, error) {
	var result CloudCredentialsProjectAssignment

	//err := self.http.Patch("/credentials/deployment/project/"+projectId, payload, &result)
	sprintf := fmt.Sprintf("/credentials/deployment/%s/project/%s", projectId, credentialId)
	err := self.http.Put(sprintf, nil, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (self *ApiClient) RemoveCloudCredentialsFromProject(credentialId string, projectId string) error {
	return self.http.Delete("/credentials/deployment/" + credentialId + "/project/" + projectId)
}

func (self *ApiClient) CloudCredentialProjectAssignments(projectId string) ([]CloudCredentialsProjectAssignment, error) {
	var result []CloudCredentialsProjectAssignment
	err := self.http.Get("/credentials/deployment/project/"+projectId, nil, &result)

	if err != nil {
		return []CloudCredentialsProjectAssignment{}, err
	}
	return result, nil
}
