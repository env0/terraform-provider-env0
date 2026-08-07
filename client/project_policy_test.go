package client_test

import (
	"encoding/json"
	"errors"

	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gomock "go.uber.org/mock/gomock"
)

const (
	policyId = "policy0"
)

var _ = Describe("Policy", func() {
	mockPolicy := Policy{
		Id:                   policyId,
		ProjectId:            "project0",
		DriftDetectionCron:   "0 * * * *",
		AutoDriftRemediation: "CODE_TO_CLOUD",
	}

	Describe("Policy", func() {
		var (
			policy Policy
			err    error
		)

		path := "/policies?projectId=" + mockPolicy.ProjectId

		Describe("Success", func() {
			BeforeEach(func() {
				policiesResult := mockPolicy
				httpCall = mockHttpClient.EXPECT().
					Get(path, nil, gomock.Any()).
					Do(func(path string, request any, response *Policy) {
						*response = policiesResult
					})

				policy, err = apiClient.Policy(mockPolicy.ProjectId)
			})

			It("Should send GET request once", func() {
				httpCall.Times(1)
			})

			It("Should return policy", func() {
				Expect(policy).Should(Equal(mockPolicy))
			})

			It("Should not return an error", func() {
				Expect(err).Should(BeNil())
			})

			It("Should return policy with auto drift remediation", func() {
				Expect(policy.AutoDriftRemediation).Should(Equal("CODE_TO_CLOUD"))
			})
		})

		Describe("Failure", func() {
			It("On error from server return the error", func() {
				expectedErr := errors.New("some error")
				httpCall = mockHttpClient.EXPECT().
					Get(path, nil, gomock.Any()).
					Return(expectedErr)

				_, err = apiClient.Policy(mockPolicy.ProjectId)
				Expect(expectedErr).Should(Equal(err))
			})
		})
	})

	Describe("PolicyUpdate", func() {
		updatePolicyPayload := PolicyUpdatePayload{
			ProjectId:             "project0",
			DriftDetectionCron:    "0 * * * *",
			DriftDetectionEnabled: true,
			AutoDriftRemediation:  "CODE_TO_CLOUD",
		}

		Describe("Success", func() {
			var (
				updatedPolicy Policy
				err           error
			)

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Put("/policies", updatePolicyPayload, gomock.Any()).
					Do(func(path string, request any, response *Policy) {
						*response = mockPolicy
					})

				updatedPolicy, err = apiClient.PolicyUpdate(updatePolicyPayload)
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

			It("Should return policy with updated auto drift remediation", func() {
				Expect(updatedPolicy.AutoDriftRemediation).To(Equal("CODE_TO_CLOUD"))
			})
		})

		Describe("Failure", func() {
			It("On error from server return the error", func() {
				expectedErr := errors.New("some error")
				httpCall = mockHttpClient.EXPECT().
					Put("/policies", updatePolicyPayload, gomock.Any()).
					Return(expectedErr)

				_, err := apiClient.PolicyUpdate(updatePolicyPayload)
				Expect(expectedErr).Should(Equal(err))
			})
		})
	})

	Describe("PolicyUpdatePayload ttl serialization", func() {
		marshalToMap := func(payload PolicyUpdatePayload) map[string]any {
			serialized, err := json.Marshal(payload)
			Expect(err).To(BeNil())

			var body map[string]any
			Expect(json.Unmarshal(serialized, &body)).To(Succeed())

			return body
		}

		It("Should serialize an unset ttl as an explicit null rather than omitting it", func() {
			body := marshalToMap(PolicyUpdatePayload{ProjectId: "project0"})

			Expect(body).To(HaveKey("maxTtl"))
			Expect(body).To(HaveKey("defaultTtl"))
			Expect(body["maxTtl"]).To(BeNil())
			Expect(body["defaultTtl"]).To(BeNil())
		})

		It("Should serialize an infinite max ttl as null alongside a finite default ttl", func() {
			defaultTtl := "7-h"
			body := marshalToMap(PolicyUpdatePayload{ProjectId: "project0", DefaultTtl: &defaultTtl})

			Expect(body).To(HaveKey("maxTtl"))
			Expect(body["maxTtl"]).To(BeNil())
			Expect(body["defaultTtl"]).To(Equal("7-h"))
		})

		It("Should serialize finite ttls as their string values", func() {
			maxTtl, defaultTtl := "1-M", "7-h"
			body := marshalToMap(PolicyUpdatePayload{ProjectId: "project0", MaxTtl: &maxTtl, DefaultTtl: &defaultTtl})

			Expect(body["maxTtl"]).To(Equal("1-M"))
			Expect(body["defaultTtl"]).To(Equal("7-h"))
		})

		It("Should serialize inherited ttls as their string values", func() {
			maxTtl, defaultTtl := "inherit", "inherit"
			body := marshalToMap(PolicyUpdatePayload{ProjectId: "project0", MaxTtl: &maxTtl, DefaultTtl: &defaultTtl})

			Expect(body["maxTtl"]).To(Equal("inherit"))
			Expect(body["defaultTtl"]).To(Equal("inherit"))
		})
	})
})
