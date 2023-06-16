package client

type AssignUserRoleToEnvironmentPayload struct {
	UserId        string `json:"userId"`
	Role          string `json:"role"`
	EnvironmentId string `json:"environmentId"`
}

type UserRoleEnvironmentAssignment struct {
	UserId string `json:"userId"`
	Role   string `json:"role"`
	Id     string `json:"id"`
}

func (client *ApiClient) AssignUserRoleToEnvironment(payload *AssignUserRoleToEnvironmentPayload) (*UserRoleEnvironmentAssignment, error) {
	var result UserRoleEnvironmentAssignment

	if err := client.http.Put("/roles/assignments/users", payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) RemoveUserRoleFromEnvironment(environmentId string, userId string) error {
	return client.http.Delete("/roles/assignments/users", map[string]string{"environmentId": environmentId, "userId": userId})
}

func (client *ApiClient) UserRoleEnvironmentAssignments(environmentId string) ([]UserRoleEnvironmentAssignment, error) {
	var result []UserRoleEnvironmentAssignment

	if err := client.http.Get("/roles/assignments/users", map[string]string{"environmentId": environmentId}, &result); err != nil {
		return nil, err
	}

	return result, nil
}
