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

type TeamProjectAssignmentPayload struct {
	TeamId      string      `json:"teamId"`
	ProjectId   string      `json:"projectId"`
	ProjectRole ProjectRole `json:"projectRole" tfschema:"role"`
}

type TeamProjectAssignment struct {
	Id          string      `json:"id"`
	TeamId      string      `json:"teamId"`
	ProjectId   string      `json:"projectId"`
	ProjectRole ProjectRole `json:"projectRole" tfschema:"role"`
}

func (client *ApiClient) TeamProjectAssignmentCreateOrUpdate(payload TeamProjectAssignmentPayload) (TeamProjectAssignment, error) {
	if payload.ProjectId == "" {
		return TeamProjectAssignment{}, errors.New("must specify project_id")
	}
	if payload.TeamId == "" {
		return TeamProjectAssignment{}, errors.New("must specify team_id")
	}
	if payload.ProjectRole == "" ||
		payload.ProjectRole != Admin &&
			payload.ProjectRole != Deployer &&
			payload.ProjectRole != Viewer &&
			payload.ProjectRole != Planner {
		return TeamProjectAssignment{}, errors.New("must specify valid project_role")
	}
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
