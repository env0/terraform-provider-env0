package client_test

import (
	"errors"
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const dummyProjectAssignmentId = "dummyId"
const dummyProjectId = "dummyProjectId"
const dummyProjectRole = "Admin"
const dummyTeamId = "dummyTeamId"

var _ = Describe("TeamProjectAssignment", func() {

	mockTeamProjectAssignment := TeamProjectAssignment{
		Id:          dummyProjectAssignmentId,
		ProjectId:   dummyProjectId,
		ProjectRole: dummyProjectRole,
		TeamId:      dummyTeamId,
	}

	Describe("Create + Update", func() {
		Describe("Success", func() {
			var teamProjectAssignment TeamProjectAssignment
			BeforeEach(func() {
				expectedPayload := TeamProjectAssignmentPayload{
					TeamId:      dummyTeamId,
					ProjectId:   dummyProjectId,
					ProjectRole: dummyProjectRole,
				}
				httpCall = mockHttpClient.EXPECT().
					Post("/teams/assignments", expectedPayload, gomock.Any()).
					Do(func(path string, request interface{}, response *TeamProjectAssignment) {
						*response = mockTeamProjectAssignment
					}).Times(1)
				teamProjectAssignment, _ = apiClient.TeamProjectAssignmentCreateOrUpdate(expectedPayload)

			})

			It("Should send POST request with params", func() {
				httpCall.Times(1)
			})

			It("Should return a new resource with id", func() {
				Expect(teamProjectAssignment).To(Equal(mockTeamProjectAssignment))
			})
		})

		Describe("Failure", func() {

			It("Should fail if assignment has no project_id", func() {
				assignmentWithoutProjectId := TeamProjectAssignmentPayload{TeamId: dummyTeamId, ProjectRole: dummyProjectRole}
				_, err := apiClient.TeamProjectAssignmentCreateOrUpdate(assignmentWithoutProjectId)
				Expect(err).To(BeEquivalentTo(errors.New("must specify project_id")))
			})

			It("Should fail if assignment has no team_id", func() {
				assignmentWithoutProjectId := TeamProjectAssignmentPayload{ProjectId: dummyProjectId, ProjectRole: dummyProjectRole}
				_, err := apiClient.TeamProjectAssignmentCreateOrUpdate(assignmentWithoutProjectId)
				Expect(err).To(BeEquivalentTo(errors.New("must specify team_id")))
			})

			It("Should fail if assignment has no project_role", func() {
				assignmentWithoutProjectId := TeamProjectAssignmentPayload{ProjectId: dummyProjectId, TeamId: dummyTeamId}
				_, err := apiClient.TeamProjectAssignmentCreateOrUpdate(assignmentWithoutProjectId)
				Expect(err).To(BeEquivalentTo(errors.New("must specify valid project_role")))
			})

			It("Should fail if assignment has invalid project_role", func() {
				assignmentWithoutProjectId := TeamProjectAssignmentPayload{ProjectId: dummyProjectId, TeamId: dummyTeamId, ProjectRole: "sdf"}
				_, err := apiClient.TeamProjectAssignmentCreateOrUpdate(assignmentWithoutProjectId)
				Expect(err).To(BeEquivalentTo(errors.New("must specify valid project_role")))
			})

			It("Should fail if assignment has empty project_role", func() {
				assignmentWithoutProjectId := TeamProjectAssignmentPayload{ProjectId: dummyProjectId, TeamId: dummyTeamId}
				_, err := apiClient.TeamProjectAssignmentCreateOrUpdate(assignmentWithoutProjectId)
				Expect(err).To(BeEquivalentTo(errors.New("must specify valid project_role")))
			})

		})

	})

	Describe("Get", func() {
		mockTeamProjectAssignments := []TeamProjectAssignment{mockTeamProjectAssignment}
		var teamProjectAssignments []TeamProjectAssignment
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/teams/assignments", map[string]string{"projectId": mockTeamProjectAssignment.ProjectId}, gomock.Any()).
				Do(func(path string, request interface{}, response *[]TeamProjectAssignment) {
					*response = mockTeamProjectAssignments
				})
			teamProjectAssignments, _ = apiClient.TeamProjectAssignments(mockTeamProjectAssignment.ProjectId)
		})

		It("Should send GET request", func() {
			httpCall.Times(1)
		})

		It("Should return the projects assignments", func() {
			Expect(teamProjectAssignments).To(Equal(mockTeamProjectAssignments))
		})

	})

	Describe("Delete", func() {

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/teams/assignments/" + mockTeamProjectAssignment.Id)
			apiClient.TeamProjectAssignmentDelete(mockTeamProjectAssignment.Id)
		})

		It("Should send DELETE request with assignment id", func() {
			httpCall.Times(1)
		})
	})

})
