package client_test

import (
	"errors"

	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Team Envrionment Assignment", func() {
	environmentId := "environmentId"
	teamId := "teamId"

	assignPayload := &AssignTeamRoleToEnvironmentPayload{
		TeamId:        teamId,
		EnvironmentId: environmentId,
		Role:          "role1",
	}

	expectedResponse := &TeamRoleEnvironmentAssignment{
		Id: "id",
	}

	errorMock := errors.New("error")

	Describe("AssignTeamToEnvironment", func() {

		Describe("Successful", func() {
			var actualResult *TeamRoleEnvironmentAssignment
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Put("/roles/assignments/teams", assignPayload, gomock.Any()).
					Do(func(path string, request interface{}, response *TeamRoleEnvironmentAssignment) {
						*response = *expectedResponse
					}).Times(1)
				actualResult, err = apiClient.AssignTeamRoleToEnvironment(assignPayload)

			})

			It("Should send POST request with params", func() {
				httpCall.Times(1)
			})

			It("should return the PUT result", func() {
				Expect(*actualResult).To(Equal(*expectedResponse))
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("Failure", func() {
			var actualResult *TeamRoleEnvironmentAssignment
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Put("/roles/assignments/teams", gomock.Any(), gomock.Any()).Return(errorMock)
				actualResult, err = apiClient.AssignTeamRoleToEnvironment(assignPayload)
			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(errorMock))
			})

			It("Should not return results", func() {
				Expect(actualResult).To(BeNil())
			})
		})
	})

	Describe("RemoveTeamFromEnvironment", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/roles/assignments/teams", map[string]string{"environmentId": environmentId, "teamId": teamId})
			apiClient.RemoveTeamRoleFromEnvironment(environmentId, teamId)
		})

		It("Should send DELETE request with assignment id", func() {
			httpCall.Times(1)
		})
	})

	Describe("TeamEnvironmentAssignments", func() {

		Describe("Successful", func() {
			var actualResult []TeamRoleEnvironmentAssignment
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/roles/assignments/teams", map[string]string{"environmentId": environmentId}, gomock.Any()).
					Do(func(path string, request interface{}, response *[]TeamRoleEnvironmentAssignment) {
						*response = []TeamRoleEnvironmentAssignment{*expectedResponse}
					}).Times(1)
				actualResult, err = apiClient.TeamRoleEnvironmentAssignments(environmentId)

			})

			It("Should send GET request with params", func() {
				httpCall.Times(1)
			})

			It("should return the GET result", func() {
				Expect(actualResult).To(Equal([]TeamRoleEnvironmentAssignment{*expectedResponse}))
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("Failure", func() {
			var actualResult []TeamRoleEnvironmentAssignment
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/roles/assignments/teams", map[string]string{"environmentId": environmentId}, gomock.Any()).Return(errorMock)
				actualResult, err = apiClient.TeamRoleEnvironmentAssignments(environmentId)
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
