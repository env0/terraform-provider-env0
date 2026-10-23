package client_test

import (
	"errors"

	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gomock "go.uber.org/mock/gomock"
)

var _ = Describe("Log Forwarding Client", func() {
	mockLogForwarding := LogForwarding{
		Id:                      "log-forwarding-id",
		OrganizationId:          organizationId,
		Name:                    "DATADOG",
		Type:                    LogForwardingTypeDatadog,
		AuditLogForwarding:      false,
		DeploymentLogForwarding: true,
	}

	Describe("LogForwardings", func() {
		var (
			returnedLogForwardings []LogForwarding
			err                    error
		)

		mockLogForwardings := []LogForwarding{mockLogForwarding}

		BeforeEach(func() {
			mockOrganizationIdCall()

			httpCall = mockHttpClient.EXPECT().
				Get("/log-forwarding/configurations", map[string]string{"organizationId": organizationId}, gomock.Any()).
				Do(func(path string, request any, response *[]LogForwarding) {
					*response = mockLogForwardings
				})

			returnedLogForwardings, err = apiClient.LogForwardings()
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send GET request", func() {
			httpCall.Times(1)
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})

		It("Should return log forwarding configurations", func() {
			Expect(returnedLogForwardings).To(Equal(mockLogForwardings))
		})
	})

	Describe("LogForwardingUpsert", func() {
		payload := LogForwardingUpsertPayload{
			Type: LogForwardingTypeDatadog,
			Value: map[string]any{
				"host":   "http-intake.logs.datadoghq.com",
				"apiKey": "${ssm:datadog-api-key}",
			},
			AuditLogForwarding:      false,
			DeploymentLogForwarding: true,
		}

		Describe("Success", func() {
			var (
				upserted *LogForwarding
				err      error
			)

			BeforeEach(func() {
				mockOrganizationIdCall()

				expectedRequest := LogForwardingUpsertPayloadWith{
					LogForwardingUpsertPayload: payload,
					OrganizationId:             organizationId,
				}

				httpCall = mockHttpClient.EXPECT().
					Put("/log-forwarding/configurations", expectedRequest, gomock.Any()).
					Do(func(path string, request any, response *LogForwarding) {
						*response = mockLogForwarding
					})

				upserted, err = apiClient.LogForwardingUpsert(payload)
			})

			It("Should get organization id", func() {
				organizationIdCall.Times(1)
			})

			It("Should send PUT request with the organization id", func() {
				httpCall.Times(1)
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})

			It("Should return the upserted configuration", func() {
				Expect(*upserted).To(Equal(mockLogForwarding))
			})
		})

		Describe("Failure", func() {
			var err error

			BeforeEach(func() {
				mockOrganizationIdCall()

				httpCall = mockHttpClient.EXPECT().
					Put("/log-forwarding/configurations", gomock.Any(), gomock.Any()).
					Return(errors.New("error"))

				_, err = apiClient.LogForwardingUpsert(payload)
			})

			It("Should return the error", func() {
				Expect(err).To(MatchError("error"))
			})
		})
	})

	Describe("LogForwardingDelete", func() {
		var err error

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/log-forwarding/configurations/"+mockLogForwarding.Id, nil)
			err = apiClient.LogForwardingDelete(mockLogForwarding.Id)
		})

		It("Should send DELETE request with the configuration id", func() {
			httpCall.Times(1)
		})

		It("Should not return error", func() {
			Expect(err).To(BeNil())
		})
	})
})
