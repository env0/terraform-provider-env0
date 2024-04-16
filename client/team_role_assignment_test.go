package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("TeamRoleAssignment", func() {
	const dummyProjectAssignmentId = "dummyId"
	const dummyProjectId = "dummyProjectId"
	const dummyEnvironmentId = "dummyEnvironmentId"
	const dummyOrganizationId = "dummyOrganizationId"
	const dummyProjectRole = "Admin"
	const dummyTeamId = "dummyTeamId"

	mockTeamRoleAssignment := TeamRoleAssignmentPayload{
		Id:     dummyProjectAssignmentId,
		Role:   dummyProjectRole,
		TeamId: dummyTeamId,
	}

	mockTeamRoleAssignments := []TeamRoleAssignmentPayload{mockTeamRoleAssignment}

	Describe("CreateOrUpdate", func() {
		Describe("ProjectId", func() {
			var assignment *TeamRoleAssignmentPayload
			BeforeEach(func() {
				createPayload := TeamRoleAssignmentCreateOrUpdatePayload{
					TeamId:    dummyTeamId,
					ProjectId: dummyProjectId,
					Role:      dummyProjectRole,
				}
				httpCall = mockHttpClient.EXPECT().
					Put("/roles/assignments/teams", &createPayload, gomock.Any()).
					Do(func(path string, request interface{}, response *TeamRoleAssignmentPayload) {
						*response = mockTeamRoleAssignment
					}).Times(1)
				assignment, _ = apiClient.TeamRoleAssignmentCreateOrUpdate(&createPayload)

			})

			It("Should send PUT request with params", func() {})

			It("Should return a new resource with id", func() {
				Expect(*assignment).To(Equal(mockTeamRoleAssignment))
			})
		})

		Describe("EnvironmentId", func() {
			var assignment *TeamRoleAssignmentPayload
			BeforeEach(func() {
				createPayload := TeamRoleAssignmentCreateOrUpdatePayload{
					TeamId:        dummyTeamId,
					EnvironmentId: dummyEnvironmentId,
					Role:          dummyProjectRole,
				}
				httpCall = mockHttpClient.EXPECT().
					Put("/roles/assignments/teams", &createPayload, gomock.Any()).
					Do(func(path string, request interface{}, response *TeamRoleAssignmentPayload) {
						*response = mockTeamRoleAssignment
					}).Times(1)
				assignment, _ = apiClient.TeamRoleAssignmentCreateOrUpdate(&createPayload)

			})

			It("Should send PUT request with params", func() {})

			It("Should return a new resource with id", func() {
				Expect(*assignment).To(Equal(mockTeamRoleAssignment))
			})
		})

		Describe("OrganizationId", func() {
			var assignment *TeamRoleAssignmentPayload
			BeforeEach(func() {
				createPayload := TeamRoleAssignmentCreateOrUpdatePayload{
					TeamId:         dummyTeamId,
					OrganizationId: dummyOrganizationId,
					Role:           dummyProjectRole,
				}
				httpCall = mockHttpClient.EXPECT().
					Put("/roles/assignments/teams", &createPayload, gomock.Any()).
					Do(func(path string, request interface{}, response *TeamRoleAssignmentPayload) {
						*response = mockTeamRoleAssignment
					}).Times(1)
				assignment, _ = apiClient.TeamRoleAssignmentCreateOrUpdate(&createPayload)
			})

			It("Should send PUT request with params", func() {})

			It("Should return a new resource with id", func() {
				Expect(*assignment).To(Equal(mockTeamRoleAssignment))
			})
		})
	})

	Describe("Get Project Assignments", func() {
		var assignments []TeamRoleAssignmentPayload
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/roles/assignments/teams", map[string]string{"projectId": dummyProjectId}, gomock.Any()).
				Do(func(path string, request interface{}, response *[]TeamRoleAssignmentPayload) {
					*response = mockTeamRoleAssignments
				}).Times(1)
			assignments, _ = apiClient.TeamRoleAssignments(&TeamRoleAssignmentListPayload{ProjectId: dummyProjectId})
		})

		It("Should send GET request", func() {})

		It("Should return the projects assignments", func() {
			Expect(assignments).To(Equal(mockTeamRoleAssignments))
		})
	})

	Describe("Get Environment Assignments", func() {
		var assignments []TeamRoleAssignmentPayload
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/roles/assignments/teams", map[string]string{"environmentId": dummyEnvironmentId}, gomock.Any()).
				Do(func(path string, request interface{}, response *[]TeamRoleAssignmentPayload) {
					*response = mockTeamRoleAssignments
				}).Times(1)
			assignments, _ = apiClient.TeamRoleAssignments(&TeamRoleAssignmentListPayload{EnvironmentId: dummyEnvironmentId})
		})

		It("Should send GET request", func() {})

		It("Should return the environment assignments", func() {
			Expect(assignments).To(Equal(mockTeamRoleAssignments))
		})
	})

	Describe("Get Organization Assignments", func() {
		var assignments []TeamRoleAssignmentPayload
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/roles/assignments/teams", map[string]string{"organizationId": dummyOrganizationId}, gomock.Any()).
				Do(func(path string, request interface{}, response *[]TeamRoleAssignmentPayload) {
					*response = mockTeamRoleAssignments
				}).Times(1)
			assignments, _ = apiClient.TeamRoleAssignments(&TeamRoleAssignmentListPayload{OrganizationId: dummyOrganizationId})
		})

		It("Should send GET request", func() {})

		It("Should return the organization assignments", func() {
			Expect(assignments).To(Equal(mockTeamRoleAssignments))
		})
	})

	Describe("Delete Project Assignment", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/roles/assignments/teams", map[string]string{"projectId": dummyProjectId, "teamId": dummyTeamId}).Times(1)
			apiClient.TeamRoleAssignmentDelete(&TeamRoleAssignmentDeletePayload{TeamId: dummyTeamId, ProjectId: dummyProjectId})
		})

		It("Should send Delete request with correct params", func() {})
	})

	Describe("Delete Environment Assignment", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/roles/assignments/teams", map[string]string{"environmentId": dummyEnvironmentId, "teamId": dummyTeamId}).Times(1)
			apiClient.TeamRoleAssignmentDelete(&TeamRoleAssignmentDeletePayload{TeamId: dummyTeamId, EnvironmentId: dummyEnvironmentId})
		})

		It("Should send Delete request with correct params", func() {})
	})

	Describe("Delete Organization Assignment", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/roles/assignments/teams", map[string]string{"organizationId": dummyOrganizationId, "teamId": dummyTeamId}).Times(1)
			apiClient.TeamRoleAssignmentDelete(&TeamRoleAssignmentDeletePayload{TeamId: dummyTeamId, OrganizationId: dummyOrganizationId})
		})

		It("Should send Delete request with correct params", func() {})
	})

})
