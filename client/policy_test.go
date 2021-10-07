package client_test

import (
	"errors"

	. "github.com/env0/terraform-provider-env0/client"
	gomock "github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Policy", func() {
	mockPolicy := Policy{
		Id: policyId,
	}

	Describe("Policy", func() {
		var policy Policy
		var err error

		Describe("Success", func() {
			BeforeEach(func() {
				policiesResult := []Policy{mockPolicy}
				httpCall = mockHttpClient.EXPECT().
					Get("/policies", nil, gomock.Any()).
					Do(func(path string, request interface{}, response *[]Policy) {
						*response = policiesResult
					})

				policy, err = apiClient.Policy()
			})

			It("Should send GET request once", func() {
				httpCall.Times(1)
			})

			It("Should return policy", func() {
				Expect(policy).Should(Equal(mockPolicy))
				Expect(err).Should(BeNil())
			})
		})

		Describe("Failure", func() {
			It("On error from server return the error", func() {
				expectedErr := errors.New("some error")
				httpCall = mockHttpClient.EXPECT().
					Get("/policies", nil, gomock.Any()).
					Return(expectedErr)

				_, err = apiClient.Policy()
				Expect(expectedErr).Should(Equal(err))
			})

			It("On too many policies return error", func() {
				policiesResult := []Policy{mockPolicy, mockPolicy}
				httpCall = mockHttpClient.EXPECT().
					Get("/policies", nil, gomock.Any()).
					Do(func(path string, request interface{}, response *[]Policy) {
						*response = policiesResult
					})

				_, err = apiClient.Policy()
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).Should(Equal("Server responded with too many policies"))
			})
		})
	})

	Describe("PolicyUpdate", func() {
		Describe("Success", func() {
			var updatedPolicy Policy
			var err error

			BeforeEach(func() {
				updatePolicyPayload := PolicyUpdatePayload{ProjectId: "project0"}

				httpCall = mockHttpClient.EXPECT().
					Put("/policies/"+mockPolicy.Id, updatePolicyPayload, gomock.Any()).
					Do(func(path string, request interface{}, response *Policy) {
						*response = mockPolicy
					})

				updatedPolicy, err = apiClient.PolicyUpdate(mockPolicy.Id, updatePolicyPayload)
			})

			It("Should send Put request with expected payload", func() {
				httpCall.Times(1)
			})

			It("Should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("Should return team received from API", func() {
				Expect(updatedPolicy).To(Equal(mockPolicy))
			})
		})

		Describe("Failure", func() {
			It("Should fail if project ID is empty", func() {
				payloadWithNoProject := PolicyUpdatePayload{NumberOfEnvironments: 1}
				_, err := apiClient.PolicyUpdate(mockPolicy.Id, payloadWithNoProject)
				Expect(err).To(BeEquivalentTo(errors.New("Must specify project ID on update")))
			})
		})
	})
})
