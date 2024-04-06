package client

type ProjectRole string

const (
	Admin    ProjectRole = "Admin"
	Deployer ProjectRole = "Deployer"
	Planner  ProjectRole = "Planner"
	Viewer   ProjectRole = "Viewer"
)

func IsBuiltinProjectRole(role string) bool {
	return role == string(Admin) || role == string(Deployer) || role == string(Planner) || role == string(Viewer)
}

type TeamProjectAssignmentPayload struct {
	TeamId    string `json:"teamId"`
	ProjectId string `json:"projectId"`
	Role      string `json:"role" tfschema:"-"`
}

type TeamProjectAssignment struct {
	Id        string `json:"id"`
	TeamId    string `json:"teamId"`
	ProjectId string `json:"projectId"`
	Role      string `json:"role" tfschema:"-"`
}

func (client *ApiClient) TeamProjectAssignmentCreateOrUpdate(payload *TeamProjectAssignmentPayload) (*TeamProjectAssignment, error) {
	var result TeamProjectAssignment

	if err := client.http.Post("/roles/assignments/teams", payload, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) TeamProjectAssignmentDelete(projectId string, teamId string) error {
	return client.http.Delete("/roles/assignments/teams", map[string]string{"projectId": projectId, "teamId": teamId})
}

func (client *ApiClient) TeamProjectAssignments(projectId string) ([]TeamProjectAssignment, error) {
	var result []TeamProjectAssignment
	err := client.http.Get("/roles/assignments/teams", map[string]string{"projectId": projectId}, &result)

	if err != nil {
		return []TeamProjectAssignment{}, err
	}
	return result, nil
}
