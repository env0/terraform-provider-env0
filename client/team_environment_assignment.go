package client

type AssignTeamRoleToEnvironmentPayload struct {
	TeamId        string `json:"teamId"`
	Role          string `json:"role" tfschema:"role_id"`
	EnvironmentId string `json:"environmentId"`
}

type TeamRoleEnvironmentAssignment struct {
	TeamId string `json:"teamId"`
	Role   string `json:"role" tfschema:"role_id"`
	Id     string `json:"id"`
}

func (client *ApiClient) AssignTeamRoleToEnvironment(payload *AssignTeamRoleToEnvironmentPayload) (*TeamRoleEnvironmentAssignment, error) {
	var result TeamRoleEnvironmentAssignment

	if err := client.http.Put("/roles/assignments/teams", payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) RemoveTeamRoleFromEnvironment(environmentId string, teamId string) error {
	return client.http.Delete("/roles/assignments/teams", map[string]string{"environmentId": environmentId, "teamId": teamId})
}

func (client *ApiClient) TeamRoleEnvironmentAssignments(environmentId string) ([]TeamRoleEnvironmentAssignment, error) {
	var result []TeamRoleEnvironmentAssignment

	if err := client.http.Get("/roles/assignments/teams", map[string]string{"environmentId": environmentId}, &result); err != nil {
		return nil, err
	}

	return result, nil
}
