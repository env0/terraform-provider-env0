package client_test

import (
	"errors"

	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Agent Client", func() {
	errorMock := errors.New("error")

	Describe("Agents", func() {
		mockAgent := Agent{
			AgentKey: "key",
		}

		expectedResponse := []Agent{mockAgent}

		Describe("Successful", func() {
			var actualResult []Agent
			var err error

			BeforeEach(func() {
				mockOrganizationIdCall(organizationId)

				httpCall = mockHttpClient.EXPECT().
					Get("/agents", gomock.Any(), gomock.Any()).
					Do(func(path string, request interface{}, response *[]Agent) {
						*response = expectedResponse
					})
				actualResult, err = apiClient.Agents()
			})

			It("Should get organization id", func() {
				organizationIdCall.Times(1)
			})

			It("Should send GET request with params", func() {
				httpCall.Times(1)
			})

			It("Should return the GET result", func() {
				Expect(actualResult).To(Equal(expectedResponse))
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("Failure", func() {
			var actualResult []Agent
			var err error

			BeforeEach(func() {
				mockOrganizationIdCall(organizationId)

				httpCall = mockHttpClient.EXPECT().
					Get("/agents", gomock.Any(), gomock.Any()).
					Return(errorMock)

				actualResult, err = apiClient.Agents()
			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(errorMock))
			})

			It("Should not return results", func() {
				Expect(actualResult).To(BeNil())
			})
		})
	})

	Describe("AgentValues", func() {
		expectedResponse := "response"
		agentId := "id"

		Describe("Successful", func() {
			var actualResult string
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/agents/"+agentId+"/values", nil, gomock.Any()).
					Do(func(path string, request interface{}, response *string) {
						*response = expectedResponse
					})
				actualResult, err = apiClient.AgentValues(agentId)
			})

			It("Should send GET request with params", func() {
				httpCall.Times(1)
			})

			It("Should return the GET result", func() {
				Expect(actualResult).To(Equal(expectedResponse))
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})
		})

		Describe("Failure", func() {
			var err error
			var actualResult string

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/agents/"+agentId+"/values", nil, gomock.Any()).
					Return(errorMock)

				actualResult, err = apiClient.AgentValues(agentId)
			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(errorMock))
				Expect(actualResult).To(Equal(""))
			})
		})
	})
})
