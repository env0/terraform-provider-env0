package client

import (
	"errors"
)

func (self *ApiClient) TeamProjectAssignmentCreateOrUpdate(payload TeamProjectAssignmentPayload) (TeamProjectAssignment, error) {
	if payload.ProjectId == "" {
		return TeamProjectAssignment{}, errors.New("must specify project_id")
	}
	if payload.TeamId == "" {
		return TeamProjectAssignment{}, errors.New("must specify team_id")
	}
	if payload.ProjectRole == "" {
		return TeamProjectAssignment{}, errors.New("must specify project_role")
	}
	var result TeamProjectAssignment

	var err = self.http.Post("/teams/assignments/", payload, &result)

	if err != nil {
		return TeamProjectAssignment{}, err
	}
	return result, nil
}

func (self *ApiClient) TeamProjectAssignmentDelete(assignmentId string) error {
	return self.http.Delete("/teams/assignments/" + assignmentId)
}
