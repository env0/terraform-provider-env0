package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("TeamProjectAssignment", func() {
	const dummyProjectAssignmentId = "dummyId"
	const dummyProjectId = "dummyProjectId"
	const dummyProjectRole = "Admin"
	const dummyTeamId = "dummyTeamId"

	mockTeamProjectAssignment := TeamProjectAssignment{
		Id:        dummyProjectAssignmentId,
		ProjectId: dummyProjectId,
		Role:      dummyProjectRole,
		TeamId:    dummyTeamId,
	}

	Describe("CreateOrUpdate", func() {
		Describe("Success", func() {
			var teamProjectAssignment *TeamProjectAssignment
			BeforeEach(func() {
				expectedPayload := TeamProjectAssignmentPayload{
					TeamId:    dummyTeamId,
					ProjectId: dummyProjectId,
					Role:      dummyProjectRole,
				}
				httpCall = mockHttpClient.EXPECT().
					Post("/roles/assignments/teams", &expectedPayload, gomock.Any()).
					Do(func(path string, request interface{}, response *TeamProjectAssignment) {
						*response = mockTeamProjectAssignment
					}).Times(1)
				teamProjectAssignment, _ = apiClient.TeamProjectAssignmentCreateOrUpdate(&expectedPayload)

			})

			It("Should send POST request with params", func() {})

			It("Should return a new resource with id", func() {
				Expect(*teamProjectAssignment).To(Equal(mockTeamProjectAssignment))
			})
		})
	})

	Describe("Get", func() {
		mockTeamProjectAssignments := []TeamProjectAssignment{mockTeamProjectAssignment}
		var teamProjectAssignments []TeamProjectAssignment
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/roles/assignments/teams", map[string]string{"projectId": mockTeamProjectAssignment.ProjectId}, gomock.Any()).
				Do(func(path string, request interface{}, response *[]TeamProjectAssignment) {
					*response = mockTeamProjectAssignments
				}).Times(1)
			teamProjectAssignments, _ = apiClient.TeamProjectAssignments(mockTeamProjectAssignment.ProjectId)
		})

		It("Should send GET request", func() {})

		It("Should return the projects assignments", func() {
			Expect(teamProjectAssignments).To(Equal(mockTeamProjectAssignments))
		})

	})

	Describe("Delete", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/roles/assignments/teams", map[string]string{"projectId": mockTeamProjectAssignment.ProjectId, "teamId": mockTeamProjectAssignment.TeamId}).Times(1)
			apiClient.TeamProjectAssignmentDelete(mockTeamProjectAssignment.ProjectId, mockTeamProjectAssignment.TeamId)
		})

		It("Should send DELETE request with assignment id", func() {})
	})

})
