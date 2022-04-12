package client

// Key is project id.
// Value is agent id.
type AssignProjectsAgentsAssignmentsPayload map[string]interface{}

type ProjectsAgentsAssignments struct {
	DefaultAgent   string                 `json:"defaultAgent"`
	ProjectsAgents map[string]interface{} `json:"ProjectsAgents"`
}

type Agent struct {
	AgentKey string `json:"agentKey"`
}

func (client *ApiClient) AssignAgentsToProjects(payload AssignProjectsAgentsAssignmentsPayload) (*ProjectsAgentsAssignments, error) {
	organizationId, err := client.organizationId()
	if err != nil {
		return nil, err
	}

	var result ProjectsAgentsAssignments
	if err := client.http.Post("/agents/projects-assignments?organizationId="+organizationId, payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) ProjectsAgentsAssignments() (*ProjectsAgentsAssignments, error) {
	organizationId, err := client.organizationId()
	if err != nil {
		return nil, err
	}

	var result ProjectsAgentsAssignments
	err = client.http.Get("/agents/projects-assignments", map[string]string{"organizationId": organizationId}, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) Agents() ([]Agent, error) {
	organizationId, err := client.organizationId()
	if err != nil {
		return nil, err
	}

	var result []Agent
	err = client.http.Get("/agents", map[string]string{"organizationId": organizationId}, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
