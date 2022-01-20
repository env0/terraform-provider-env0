package client_test

import (
	"errors"
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EnvironmentScheduling", func() {
	mockEnvironmentId := "env0"
	mockError := errors.New("very cool error")

	mockDeployCronPayload := EnvironmentSchedulingExpression{Cron: "1 * * * *", Enabled: true}
	mockDestroyCronPayload := EnvironmentSchedulingExpression{Cron: "0 * * * *", Enabled: true}

	mockEnvironmentSchedulingPayload := EnvironmentScheduling{
		Deploy:  mockDeployCronPayload,
		Destroy: mockDestroyCronPayload,
	}

	var environmentSchedulingResponse EnvironmentScheduling

	Describe("Update", func() {
		Describe("Success", func() {
			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Put("/scheduling/environments/"+mockEnvironmentId, mockEnvironmentSchedulingPayload, gomock.Any()).
					Do(func(path string, request interface{}, response *EnvironmentScheduling) {
						*response = mockEnvironmentSchedulingPayload
					}).Times(1)
				environmentSchedulingResponse, _ = apiClient.EnvironmentSchedulingUpdate(mockEnvironmentId, mockEnvironmentSchedulingPayload)

			})

			It("Should send PUT request with params", func() {
				httpCall.Times(1)
			})

			It("Should return environment scheduling response", func() {
				Expect(environmentSchedulingResponse).To(Equal(mockEnvironmentSchedulingPayload))
			})
		})

		Describe("Failure", func() {
			It("Should fail if cron expressions are the same", func() {
				mockFailedEnvironmentSchedulingPayload := EnvironmentScheduling{
					Deploy:  mockDeployCronPayload,
					Destroy: mockDeployCronPayload,
				}
				_, err := apiClient.EnvironmentSchedulingUpdate(mockEnvironmentId, mockFailedEnvironmentSchedulingPayload)
				Expect(err).To(BeEquivalentTo(errors.New("deploy and destroy cron expressions must not be the same")))
			})

			It("Should fail if API call fails", func() {
				httpCall = mockHttpClient.EXPECT().
					Put("/scheduling/environments/"+mockEnvironmentId, gomock.Any(), gomock.Any()).
					Return(mockError).
					Times(1)

				_, err := apiClient.EnvironmentSchedulingUpdate(mockEnvironmentId, mockEnvironmentSchedulingPayload)

				Expect(err).To(BeEquivalentTo(mockError))
			})

		})

	})

	Describe("Get", func() {
		Describe("Success", func() {
			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/scheduling/environments/"+mockEnvironmentId, nil, gomock.Any()).
					Do(func(path string, request interface{}, response *EnvironmentScheduling) {
						*response = mockEnvironmentSchedulingPayload
					})
				environmentSchedulingResponse, _ = apiClient.EnvironmentScheduling(mockEnvironmentId)
			})

			It("Should send GET request", func() {
				httpCall.Times(1)
			})

			It("Should return the environment scheduling response", func() {
				Expect(environmentSchedulingResponse).To(Equal(mockEnvironmentSchedulingPayload))
			})
		})

		Describe("Fail", func() {
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/scheduling/environments/"+mockEnvironmentId, gomock.Any(), gomock.Any()).
					Return(mockError).
					Times(1)

				_, err = apiClient.EnvironmentScheduling(mockEnvironmentId)

			})

			It("Should send GET request", func() {
				httpCall.Times(1)
			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(mockError))
			})
		})

	})

	Describe("Delete", func() {
		Describe("Success", func() {
			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().Put("/scheduling/environments/"+mockEnvironmentId, EnvironmentScheduling{}, nil)
				apiClient.EnvironmentSchedulingDelete(mockEnvironmentId)
			})

			It("Should send PUT request with empty environment scheduling object", func() {
				httpCall.Times(1)
			})
		})

		Describe("Fail", func() {
			var err error
			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Put("/scheduling/environments/"+mockEnvironmentId, gomock.Any(), gomock.Any()).
					Return(mockError).
					Times(1)
				err = apiClient.EnvironmentSchedulingDelete(mockEnvironmentId)
			})

			It("Should return error", func() {
				Expect(err).To(BeEquivalentTo(mockError))
			})
		})
	})

})
