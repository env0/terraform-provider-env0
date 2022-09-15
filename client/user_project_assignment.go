package client

type AssignUserToProjectPayload struct {
	UserId string      `json:"userId"`
	Role   ProjectRole `json:"role"`
}

type UpdateUserProjectAssignmentPayload struct {
	Role ProjectRole `json:"role"`
}

type UserProjectAssignment struct {
	UserId string      `json:"userId"`
	Role   ProjectRole `json:"role"`
	Id     string      `json:"id"`
}

func (client *ApiClient) AssignUserToProject(projectId string, payload *AssignUserToProjectPayload) (*UserProjectAssignment, error) {
	var result UserProjectAssignment

	if err := client.http.Post("/permissions/projects/"+projectId, payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) RemoveUserFromProject(projectId string, userId string) error {
	return client.http.Delete("/permissions/projects/" + projectId + "/users/" + userId)
}

func (client *ApiClient) UpdateUserProjectAssignment(projectId string, userId string, payload *UpdateUserProjectAssignmentPayload) (*UserProjectAssignment, error) {
	var result UserProjectAssignment

	if err := client.http.Put("/permissions/projects/"+projectId+"/users/"+userId, payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) UserProjectAssignments(projectId string) ([]UserProjectAssignment, error) {
	var result []UserProjectAssignment

	if err := client.http.Get("/permissions/projects/"+projectId, nil, &result); err != nil {
		return nil, err
	}

	return result, nil
}
