package client_test

import (
	"errors"

	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe(" Cost Credentials Project Assignment", func() {
	projectId := "projectId"
	credentialId := "credentialId"

	Describe("AssignCostCredentialsToProject", func() {
		expectedResponse := CostCredentialProjectAssignment{
			ProjectId:       "assignment id",
			CredentialsId:   "credentialId",
			CredentialsType: "GCP_CREDENTIALS",
		}
		Describe("Successful", func() {
			var actualResult CostCredentialProjectAssignment
			BeforeEach(func() {

				httpCall = mockHttpClient.EXPECT().
					Put("/costs/project/"+projectId+"/credentials", map[string]string{"credentialsId": credentialId}, gomock.Any()).
					Do(func(path string, request any, response *CostCredentialProjectAssignment) {
						*response = expectedResponse
					}).Times(1)
				actualResult, _ = apiClient.AssignCostCredentialsToProject(projectId, credentialId)

			})

			It("should return the PUT result", func() {
				Expect(actualResult).To(Equal(expectedResponse))
			})
		})
		Describe("On Error", func() {
			errorInfo := "error"
			var err error
			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Put("/costs/project/"+projectId+"/credentials", map[string]string{"credentialsId": credentialId}, gomock.Any()).
					Return(errors.New(errorInfo)).
					Times(1)
				_, err = apiClient.AssignCostCredentialsToProject(projectId, credentialId)

			})

			It("should return the error from the api call", func() {
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).Should(Equal(errorInfo))
			})
		})
	})

	Describe("RemoveCostCredentialsFromProject", func() {
		errorInfo := "error"
		var err error
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Delete("/costs/project/"+projectId+"/credentials/"+credentialId, nil).
				Return(errors.New(errorInfo)).
				Times(1)
			err = apiClient.RemoveCostCredentialsFromProject(projectId, credentialId)

		})

		It("should return the error from the api call", func() {
			Expect(err).ShouldNot(BeNil())
			Expect(err.Error()).Should(Equal(errorInfo))
		})
	})

	Describe("CostCredentialIdsInProject", func() {
		Describe("Successful", func() {

			var actualResult []CostCredentialProjectAssignment

			firstResulteResponse := CostCredentialProjectAssignment{
				ProjectId:       "assignment id",
				CredentialsId:   "credentialId",
				CredentialsType: "GCP_CREDENTIALS",
			}

			secondResulteResponse := CostCredentialProjectAssignment{
				ProjectId:       "assignment id",
				CredentialsId:   "credentialId-2",
				CredentialsType: "GCP_CREDENTIALS",
			}

			expectedResponse := []CostCredentialProjectAssignment{firstResulteResponse, secondResulteResponse}
			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/costs/project/"+projectId+"/credentials", nil, gomock.Any()).
					Do(func(path string, request any, response *[]CostCredentialProjectAssignment) {
						*response = expectedResponse
					}).Times(1)
				actualResult, _ = apiClient.CostCredentialIdsInProject(projectId)

			})

			It("should return the GET result", func() {
				Expect(actualResult).To(Equal(expectedResponse))
			})
		})
		Describe("On Error", func() {
			errorInfo := "error"
			var err error
			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/costs/project/"+projectId+"/credentials", nil, gomock.Any()).
					Return(errors.New(errorInfo)).
					Times(1)
				_, err = apiClient.CostCredentialIdsInProject(projectId)

			})

			It("should return the error from the api call", func() {
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).Should(Equal(errorInfo))
			})
		})

	})
})
