package client

import (
	"errors"
)

func (self *ApiClient) TeamProjectAssignmentCreateOrUpdate(payload TeamProjectAssignmentPayload) (TeamProjectAssignmentResponse, error) {
	if payload.ProjectId == "" {
		return TeamProjectAssignmentResponse{}, errors.New("must specify project_id")
	}
	if payload.TeamId == "" {
		return TeamProjectAssignmentResponse{}, errors.New("must specify team_id")
	}
	if payload.ProjectRole == "" {
		return TeamProjectAssignmentResponse{}, errors.New("must specify project_role")
	}
	var result TeamProjectAssignmentResponse

	var err = self.http.Post("/teams/assignments/", payload, &result)

	if err != nil {
		return TeamProjectAssignmentResponse{}, err
	}
	return result, nil
}

func (self *ApiClient) TeamProjectAssignmentDelete(assignmentId string) error {
	return self.http.Delete("/teams/assignments/" + assignmentId)
}
