package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TeamProjectAssignment", func() {
	const dummyProjectAssignmentId = "dummyId"
	const dummyProjectId = "dummyProjectId"
	const dummyProjectRole = "Admin"
	const dummyTeamId = "dummyTeamId"

	mockTeamProjectAssignment := TeamProjectAssignment{
		Id:          dummyProjectAssignmentId,
		ProjectId:   dummyProjectId,
		ProjectRole: dummyProjectRole,
		TeamId:      dummyTeamId,
	}

	Describe("CreateOrUpdate", func() {
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
			httpCall = mockHttpClient.EXPECT().Delete("/teams/assignments/"+mockTeamProjectAssignment.Id, nil)
			apiClient.TeamProjectAssignmentDelete(mockTeamProjectAssignment.Id)
		})

		It("Should send DELETE request with assignment id", func() {
			httpCall.Times(1)
		})
	})

})
