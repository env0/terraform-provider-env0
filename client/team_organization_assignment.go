package client

type AssignOrganizationRoleToTeamPayload struct {
	TeamId string `json:"teamId"`
	Role   string `json:"role" tfschema:"role_id"`
}

type OrganizationRoleTeamAssignment struct {
	TeamId string `json:"teamId"`
	Role   string `json:"role" tfschema:"role_id"`
	Id     string `json:"id"`
}

func (client *ApiClient) AssignOrganizationRoleToTeam(payload *AssignOrganizationRoleToTeamPayload) (*OrganizationRoleTeamAssignment, error) {
	var result OrganizationRoleTeamAssignment

	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	payloadWithOrganizationId := struct {
		*AssignOrganizationRoleToTeamPayload
		OrganizationId string `json:"organizationId"`
	}{
		payload,
		organizationId,
	}

	if err := client.http.Put("/roles/assignments/teams", payloadWithOrganizationId, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) RemoveOrganizationRoleFromTeam(teamId string) error {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return err
	}

	return client.http.Delete("/roles/assignments/teams", map[string]string{"organizationId": organizationId, "teamId": teamId})
}

func (client *ApiClient) OrganizationRoleTeamAssignments() ([]OrganizationRoleTeamAssignment, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	var result []OrganizationRoleTeamAssignment

	if err := client.http.Get("/roles/assignments/teams", map[string]string{"organizationId": organizationId}, &result); err != nil {
		return nil, err
	}

	return result, nil
}
