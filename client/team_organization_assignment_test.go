package client_test

import (
	"errors"

	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Team Orgnization Assignment", func() {
	organizationId := "organizationId"
	teamId := "teamId"

	assignPayload := &AssignTeamRoleToOrganizationPayload{
		TeamId: teamId,
		Role:   "role1",
	}

	assignPayloadWithOrganizationId := struct {
		*AssignTeamRoleToOrganizationPayload
		OrganizationId string `json:"organizationId"`
	}{
		assignPayload,
		organizationId,
	}

	expectedResponse := &TeamRoleOrganizationAssignment{
		Id: "id",
	}

	errorMock := errors.New("error")

	Describe("AssignTeamToOrganization", func() {

		Describe("Successful", func() {
			var actualResult *TeamRoleOrganizationAssignment
			var err error

			BeforeEach(func() {
				mockOrganizationIdCall(organizationId).Times(1)
				httpCall = mockHttpClient.EXPECT().
					Put("/roles/assignments/teams", assignPayloadWithOrganizationId, gomock.Any()).
					Do(func(path string, request interface{}, response *TeamRoleOrganizationAssignment) {
						*response = *expectedResponse
					}).Times(1)
				actualResult, err = apiClient.AssignTeamRoleToOrganization(assignPayload)

			})

			It("Should send POST request with params", func() {})

			It("should return the PUT result", func() {
				Expect(*actualResult).To(Equal(*expectedResponse))
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("Failure", func() {
			var actualResult *TeamRoleOrganizationAssignment
			var err error

			BeforeEach(func() {
				mockOrganizationIdCall(organizationId).Times(1)
				httpCall = mockHttpClient.EXPECT().
					Put("/roles/assignments/teams", gomock.Any(), gomock.Any()).Return(errorMock).Times(1)
				actualResult, err = apiClient.AssignTeamRoleToOrganization(assignPayload)
			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(errorMock))
			})

			It("Should not return results", func() {
				Expect(actualResult).To(BeNil())
			})
		})
	})

	Describe("RemoveTeamFromOrganization", func() {
		BeforeEach(func() {
			mockOrganizationIdCall(organizationId).Times(1)
			httpCall = mockHttpClient.EXPECT().Delete("/roles/assignments/teams", map[string]string{"organizationId": organizationId, "teamId": teamId}).Times(1)
			apiClient.RemoveTeamRoleFromOrganization(teamId)
		})

		It("Should send DELETE request with assignment id", func() {})
	})

	Describe("TeamOrganizationAssignments", func() {

		Describe("Successful", func() {
			var actualResult []TeamRoleOrganizationAssignment
			var err error

			BeforeEach(func() {
				mockOrganizationIdCall(organizationId).Times(1)
				httpCall = mockHttpClient.EXPECT().
					Get("/roles/assignments/teams", map[string]string{"organizationId": organizationId}, gomock.Any()).
					Do(func(path string, request interface{}, response *[]TeamRoleOrganizationAssignment) {
						*response = []TeamRoleOrganizationAssignment{*expectedResponse}
					}).Times(1)
				actualResult, err = apiClient.TeamRoleOrganizationAssignments()

			})

			It("Should send GET request with params", func() {})

			It("should return the GET result", func() {
				Expect(actualResult).To(Equal([]TeamRoleOrganizationAssignment{*expectedResponse}))
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("Failure", func() {
			var actualResult []TeamRoleOrganizationAssignment
			var err error

			BeforeEach(func() {
				mockOrganizationIdCall(organizationId).Times(1)
				httpCall = mockHttpClient.EXPECT().
					Get("/roles/assignments/teams", map[string]string{"organizationId": organizationId}, gomock.Any()).Return(errorMock).Times(1)
				actualResult, err = apiClient.TeamRoleOrganizationAssignments()
			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(errorMock))
			})

			It("Should not return results", func() {
				Expect(actualResult).To(BeNil())
			})
		})
	})
})
