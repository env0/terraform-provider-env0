package client

import (
	"errors"
)

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
	TeamId      string `json:"teamId"`
	ProjectId   string `json:"projectId"`
	ProjectRole string `json:"projectRole" tfschema:"-"`
}

type TeamProjectAssignment struct {
	Id          string `json:"id"`
	TeamId      string `json:"teamId"`
	ProjectId   string `json:"projectId"`
	ProjectRole string `json:"projectRole" tfschema:"-"`
}

func (client *ApiClient) TeamProjectAssignmentCreateOrUpdate(payload TeamProjectAssignmentPayload) (TeamProjectAssignment, error) {
	var result TeamProjectAssignment

	var err = client.http.Post("/teams/assignments", payload, &result)

	if err != nil {
		return TeamProjectAssignment{}, err
	}
	return result, nil
}

func (client *ApiClient) TeamProjectAssignmentDelete(assignmentId string) error {
	if assignmentId == "" {
		return errors.New("empty assignmentId")
	}
	return client.http.Delete("/teams/assignments/" + assignmentId)
}

func (client *ApiClient) TeamProjectAssignments(projectId string) ([]TeamProjectAssignment, error) {

	var result []TeamProjectAssignment
	err := client.http.Get("/teams/assignments", map[string]string{"projectId": projectId}, &result)

	if err != nil {
		return []TeamProjectAssignment{}, err
	}
	return result, nil
}
