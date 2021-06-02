package client

import "errors"

func (self *ApiClient) AssignCloudCredentialsToProject(projectId string, payload CloudCredentialsProjectAssignmentPatchPayload) (CloudCredentialsProjectAssignment, error) {
	var result CloudCredentialsProjectAssignment
	if payload.CredentialIds == nil || len(payload.CredentialIds) == 0 {
		return result, errors.New("Must specify cloud credentials to assign to be assigned to project")
	}
	err := self.http.Patch("/credentials/deployment/project/"+projectId, payload, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (self *ApiClient) RemoveCloudCredentialsFromProject(credentialId string, projectId string) error {
	return self.http.Delete("/credentials/deployment/" + credentialId + "/project/" + projectId)
}

func (self *ApiClient) CloudCredentialProjectAssginments(projectId string) ([]CloudCredentialsProjectAssignment, error) {
	var result []CloudCredentialsProjectAssignment
	err := self.http.Get("/credentials/deployment/project/"+projectId, nil, &result)

	if err != nil {
		return []CloudCredentialsProjectAssignment{}, err
	}
	return result, nil
}
