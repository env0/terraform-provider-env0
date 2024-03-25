package client

type AssignTeamRoleToOrganizationPayload struct {
	TeamId string `json:"teamId"`
	Role   string `json:"role" tfschema:"role_id"`
}

type TeamRoleOrganizationAssignment struct {
	TeamId string `json:"teamId"`
	Role   string `json:"role" tfschema:"role_id"`
	Id     string `json:"id"`
}

func (client *ApiClient) AssignTeamRoleToOrganization(payload *AssignTeamRoleToOrganizationPayload) (*TeamRoleOrganizationAssignment, error) {
	var result TeamRoleOrganizationAssignment

	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	payloadWithOrganizationId := struct {
		*AssignTeamRoleToOrganizationPayload
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

func (client *ApiClient) RemoveTeamRoleFromOrganization(teamId string) error {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return err
	}

	return client.http.Delete("/roles/assignments/teams", map[string]string{"organizationId": organizationId, "teamId": teamId})
}

func (client *ApiClient) TeamRoleOrganizationAssignments() ([]TeamRoleOrganizationAssignment, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	var result []TeamRoleOrganizationAssignment

	if err := client.http.Get("/roles/assignments/teams", map[string]string{"organizationId": organizationId}, &result); err != nil {
		return nil, err
	}

	return result, nil
}
