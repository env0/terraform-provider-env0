package client_test

import (
	"errors"

	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// mockOrganizationIdCall(organizationId)

var _ = Describe("Agent Project Assignment", func() {
	psas := map[string]interface{}{
		"pid1": "aid1",
		"pid2": "aid2",
	}

	expectedResponse := ProjectsAgentsAssignments{
		ProjectsAgents: psas,
	}

	errorMock := errors.New("error")

	Describe("AssignAgentsToProjects", func() {

		Describe("Successful", func() {
			var actualResult *ProjectsAgentsAssignments
			var err error

			BeforeEach(func() {
				mockOrganizationIdCall(organizationId)

				httpCall = mockHttpClient.EXPECT().
					Post("/agents/projects-assignments?organizationId="+organizationId, gomock.Any(), gomock.Any()).
					Do(func(path string, request interface{}, response *ProjectsAgentsAssignments) {
						*response = expectedResponse
					}).Times(1)
				actualResult, err = apiClient.AssignAgentsToProjects(psas)

			})

			It("Should get organization id", func() {
				organizationIdCall.Times(1)
			})

			It("Should send POST request with params", func() {
				httpCall.Times(1)
			})

			It("should return the POST result", func() {
				Expect(*actualResult).To(Equal(expectedResponse))
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("Failure", func() {
			var actualResult *ProjectsAgentsAssignments
			var err error

			BeforeEach(func() {
				mockOrganizationIdCall(organizationId)

				httpCall = mockHttpClient.EXPECT().
					Post("/agents/projects-assignments?organizationId="+organizationId, gomock.Any(), gomock.Any()).
					Return(errorMock)

				actualResult, err = apiClient.AssignAgentsToProjects(psas)
			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(errorMock))
			})

			It("Should not return results", func() {
				Expect(actualResult).To(BeNil())
			})
		})
	})

	Describe("ProjectsAgentsAssignments", func() {
		Describe("Successful", func() {
			var actualResult *ProjectsAgentsAssignments
			var err error

			BeforeEach(func() {
				mockOrganizationIdCall(organizationId)

				httpCall = mockHttpClient.EXPECT().
					Get("/agents/projects-assignments", gomock.Any(), gomock.Any()).
					Do(func(path string, request interface{}, response *ProjectsAgentsAssignments) {
						*response = expectedResponse
					})
				actualResult, err = apiClient.ProjectsAgentsAssignments()
			})

			It("Should get organization id", func() {
				organizationIdCall.Times(1)
			})

			It("Should send GET request with params", func() {
				httpCall.Times(1)
			})

			It("Should return the GET result", func() {
				Expect(*actualResult).To(Equal(expectedResponse))
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("Failure", func() {
			var actualResult *ProjectsAgentsAssignments
			var err error

			BeforeEach(func() {
				mockOrganizationIdCall(organizationId)

				httpCall = mockHttpClient.EXPECT().
					Get("/agents/projects-assignments", gomock.Any(), gomock.Any()).
					Return(errorMock)

				actualResult, err = apiClient.ProjectsAgentsAssignments()
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
