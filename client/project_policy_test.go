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

		Describe("TTL serialization", func() {
			It("serializes empty TTL values as explicit nulls", func() {
				empty := ""
				payload := PolicyUpdatePayload{
					ProjectId:  "project0",
					MaxTtl:     &empty,
					DefaultTtl: &empty,
				}

				serialized, err := json.Marshal(payload)
				Expect(err).NotTo(HaveOccurred())

				var body map[string]any
				Expect(json.Unmarshal(serialized, &body)).To(Succeed())
				Expect(body).To(HaveKey("maxTtl"))
				Expect(body["maxTtl"]).To(BeNil())
				Expect(body).To(HaveKey("defaultTtl"))
				Expect(body["defaultTtl"]).To(BeNil())
			})

			It("omits unset TTL values", func() {
				serialized, err := json.Marshal(PolicyUpdatePayload{ProjectId: "project0"})
				Expect(err).NotTo(HaveOccurred())

				var body map[string]any
				Expect(json.Unmarshal(serialized, &body)).To(Succeed())
				Expect(body).NotTo(HaveKey("maxTtl"))
				Expect(body).NotTo(HaveKey("defaultTtl"))
			})

			It("preserves concrete TTL values", func() {
				maxTtl := "1-M"
				defaultTtl := "1-w"
				payload := PolicyUpdatePayload{
					ProjectId:  "project0",
					MaxTtl:     &maxTtl,
					DefaultTtl: &defaultTtl,
				}

				serialized, err := json.Marshal(payload)
				Expect(err).NotTo(HaveOccurred())

				var body map[string]any
				Expect(json.Unmarshal(serialized, &body)).To(Succeed())
				Expect(body).To(HaveKeyWithValue("maxTtl", maxTtl))
				Expect(body).To(HaveKeyWithValue("defaultTtl", defaultTtl))
			})
		})
	})
})
