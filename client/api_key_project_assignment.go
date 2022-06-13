package client

type ApiKeyRole string

const (
	ApiKeyRoleDeployer ApiKeyRole = "Deployer"
	ApiKeyRoleViewer   ApiKeyRole = "Viewer"
	ApiKeyRolePlanner  ApiKeyRole = "Planner"
	ApiKeyRoleAdmin    ApiKeyRole = "Admin"
)

type AssignAPIKeyToProjectPayload struct {
	UserId string     `json:"userId"`
	Role   ApiKeyRole `json:"role"`
}

type ApiKeyProjectAssignment struct {
	Id string `json:"id"`
}

func (client *ApiClient) AssignApiKeyToProject(projectId string, payload *AssignAPIKeyToProjectPayload) (*ApiKeyProjectAssignment, error) {
	var result ApiKeyProjectAssignment

	err := client.http.Post("/permissions/projects/"+projectId, payload, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *ApiClient) RemoveApiKeyFromProject(projectId string, apiKeyId string) error {
	return client.http.Delete("/permissions/projects/" + projectId + "/users/" + apiKeyId)
}
