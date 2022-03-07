package client_test

import (
	"errors"
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EnvironmentDriftDetection", func() {
	environmentId := "env"
	path := "/scheduling/drift-detection/environments/" + environmentId
	mockError := errors.New("I don't like milk")
	schedulingExpression := EnvironmentSchedulingExpression{Cron: "0 * * * *", Enabled: true}
	var driftResponse EnvironmentSchedulingExpression

	Describe("Get", func() {
		Describe("Success", func() {
			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get(path, nil, gomock.Any()).
					Do(func(path string, request interface{}, response *EnvironmentSchedulingExpression) {
						*response = schedulingExpression
					})
				driftResponse, _ = apiClient.EnvironmentDriftDetection(environmentId)
			})

			It("Should return the drift scheduling response", func() {
				Expect(driftResponse).To(Equal(schedulingExpression))
			})
		})

		Describe("Fail", func() {
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get(path, gomock.Any(), gomock.Any()).
					Return(mockError)

				_, err = apiClient.EnvironmentDriftDetection(environmentId)

			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(mockError))
			})
		})
	})
	Describe("Update", func() {
		Describe("Success", func() {
			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Patch(path, schedulingExpression, gomock.Any()).
					Do(func(path string, request interface{}, response *EnvironmentSchedulingExpression) {
						*response = schedulingExpression
					})
				driftResponse, _ = apiClient.EnvironmentUpdateDriftDetection(environmentId, schedulingExpression)
			})

			It("Should return the drift scheduling response", func() {
				Expect(driftResponse).To(Equal(schedulingExpression))
			})
		})

		Describe("Fail", func() {
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Patch(path, schedulingExpression, gomock.Any()).
					Return(mockError)

				_, err = apiClient.EnvironmentUpdateDriftDetection(environmentId, schedulingExpression)

			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(mockError))
			})
		})
	})
	Describe("Delete", func() {
		Describe("Success", func() {
			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Patch(path, EnvironmentSchedulingExpression{Enabled: false}, gomock.Any()).
					Do(func(path string, request interface{}, response *EnvironmentSchedulingExpression) {
						*response = schedulingExpression
					})
				_ = apiClient.EnvironmentStopDriftDetection(environmentId)
			})
		})

		Describe("Fail", func() {
			var err error

			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Patch(path, EnvironmentSchedulingExpression{Enabled: false}, gomock.Any()).
					Return(mockError)

				err = apiClient.EnvironmentStopDriftDetection(environmentId)

			})

			It("Should fail if API call fails", func() {
				Expect(err).To(Equal(mockError))
			})
		})
	})

})
