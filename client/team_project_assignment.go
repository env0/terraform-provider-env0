package client

import (
	"errors"
)

type Role string

const (
	Admin    Role = "Admin"
	Deployer Role = "Deployer"
	Planner  Role = "Planner"
	Viewer   Role = "Viewer"
)

type TeamProjectAssignmentPayload struct {
	TeamId      string `json:"teamId"`
	ProjectId   string `json:"projectId"`
	ProjectRole Role   `json:"projectRole"`
}

type TeamProjectAssignment struct {
	Id          string `json:"id"`
	TeamId      string `json:"teamId"`
	ProjectId   string `json:"projectId"`
	ProjectRole Role   `json:"projectRole"`
}

func (self *ApiClient) TeamProjectAssignmentCreateOrUpdate(payload TeamProjectAssignmentPayload) (TeamProjectAssignment, error) {
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

	var err = self.http.Post("/teams/assignments", payload, &result)

	if err != nil {
		return TeamProjectAssignment{}, err
	}
	return result, nil
}

func (self *ApiClient) TeamProjectAssignmentDelete(assignmentId string) error {
	if assignmentId == "" {
		return errors.New("empty assignmentId")
	}
	return self.http.Delete("/teams/assignments/" + assignmentId)
}

func (self *ApiClient) TeamProjectAssignments(projectId string) ([]TeamProjectAssignment, error) {

	var result []TeamProjectAssignment
	err := self.http.Get("/teams/assignments", map[string]string{"projectId": projectId}, &result)

	if err != nil {
		return []TeamProjectAssignment{}, err
	}
	return result, nil
}
