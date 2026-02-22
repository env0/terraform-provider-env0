package client

type TeamRoleAssignmentPayload struct {
	TeamId string `json:"teamId"`
	Role   string `json:"role"   tfschema:"role_id"`
	Id     string `json:"id"`
}

type TeamRoleAssignmentCreateOrUpdatePayload struct {
	TeamId         string `json:"teamId"`
	Role           string `json:"role"                     tfschema:"role_id"`
	EnvironmentId  string `json:"environmentId,omitempty"`
	OrganizationId string `json:"organizationId,omitempty"`
	ProjectId      string `json:"projectId,omitempty"`
}

type TeamRoleAssignmentDeletePayload struct {
	TeamId         string
	EnvironmentId  string
	OrganizationId string
	ProjectId      string
}

type TeamRoleAssignmentListPayload struct {
	EnvironmentId  string
	OrganizationId string
	ProjectId      string
}

func (client *ApiClient) TeamRoleAssignmentCreateOrUpdate(payload *TeamRoleAssignmentCreateOrUpdatePayload) (*TeamRoleAssignmentPayload, error) {
	var result TeamRoleAssignmentPayload

	if err := client.http.Put("/roles/assignments/teams", payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) TeamRoleAssignmentDelete(payload *TeamRoleAssignmentDeletePayload) error {
	params := make(map[string]string)

	params["teamId"] = payload.TeamId

	if payload.EnvironmentId != "" {
		params["environmentId"] = payload.EnvironmentId
	}

	if payload.OrganizationId != "" {
		params["organizationId"] = payload.OrganizationId
	}

	if payload.ProjectId != "" {
		params["projectId"] = payload.ProjectId
	}

	return client.http.Delete("/roles/assignments/teams", params)
}

func (client *ApiClient) TeamRoleAssignments(payload *TeamRoleAssignmentListPayload) ([]TeamRoleAssignmentPayload, error) {
	params := make(map[string]string)

	if payload.EnvironmentId != "" {
		params["environmentId"] = payload.EnvironmentId
	}

	if payload.OrganizationId != "" {
		params["organizationId"] = payload.OrganizationId
	}

	if payload.ProjectId != "" {
		params["projectId"] = payload.ProjectId
	}

	var result []TeamRoleAssignmentPayload

	if err := client.http.Get("/roles/assignments/teams", params, &result); err != nil {
		return nil, err
	}

	return result, nil
}
