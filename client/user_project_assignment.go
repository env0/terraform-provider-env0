package client

type AssignUserToProjectPayload struct {
	UserId string `json:"userId"`
	Role   Role   `json:"role"`
}

type UserProjectAssignment struct {
	Id string `json:"id"`
}

func (client *ApiClient) AssignUserToProject(projectId string, payload *AssignUserToProjectPayload) (*UserProjectAssignment, error) {
	var result UserProjectAssignment

	err := client.http.Post("/permissions/projects/"+projectId, payload, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (client *ApiClient) RemoveUserFromProject(projectId string, userId string) error {
	return client.http.Delete("/permissions/projects/" + projectId + "/users/" + userId)
}
