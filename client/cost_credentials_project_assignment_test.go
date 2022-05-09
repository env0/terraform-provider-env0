package client_test

import (
	"errors"

	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe(" Cost Credentials Project Assignment", func() {
	projectId := "projectId"
	credentialId := "credentialId"

	Describe("AssignCostCredentialsToProject", func() {
		expectedResponse := CostCredentialProjectAssignment{
			ProjectId:       "assigment id",
			CredentialsId:   "credentialId",
			CredentialsType: "GCP_CREDENTIALS",
		}
		Describe("Successful", func() {
			var actualResult CostCredentialProjectAssignment
			BeforeEach(func() {

				httpCall = mockHttpClient.EXPECT().
					Put("/costs/project/"+projectId+"/credentials", map[string]string{"credentialsId": credentialId}, gomock.Any()).
					Do(func(path string, request interface{}, response *CostCredentialProjectAssignment) {
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
			var actualError error
			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Put("/costs/project/"+projectId+"/credentials", map[string]string{"credentialsId": credentialId}, gomock.Any()).
					Return(errors.New(errorInfo)).
					Times(1)
				_, actualError = apiClient.AssignCostCredentialsToProject(projectId, credentialId)

			})

			It("should return the error from the api call", func() {
				Expect(actualError).ShouldNot(BeNil())
				Expect(actualError.Error()).Should(Equal(errorInfo))
			})
		})
	})

	Describe("RemoveCostCredentialsFromProject", func() {
		errorInfo := "error"
		var actualError error
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Delete("/costs/project/" + projectId + "/credentials/" + credentialId).
				Return(errors.New(errorInfo)).
				Times(1)
			actualError = apiClient.RemoveCostCredentialsFromProject(projectId, credentialId)

		})

		It("should return the error from the api call", func() {
			Expect(actualError).ShouldNot(BeNil())
			Expect(actualError.Error()).Should(Equal(errorInfo))
		})
	})

	Describe("CostCredentialIdsInProject", func() {
		Describe("Successful", func() {

			var actualResult []CostCredentialProjectAssignment

			firstResulteResponse := CostCredentialProjectAssignment{
				ProjectId:       "assigment id",
				CredentialsId:   "credentialId",
				CredentialsType: "GCP_CREDENTIALS",
			}

			secondResulteResponse := CostCredentialProjectAssignment{
				ProjectId:       "assigment id",
				CredentialsId:   "credentialId-2",
				CredentialsType: "GCP_CREDENTIALS",
			}

			expectedResponse := []CostCredentialProjectAssignment{firstResulteResponse, secondResulteResponse}
			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/costs/project/"+projectId+"/credentials", nil, gomock.Any()).
					Do(func(path string, request interface{}, response *[]CostCredentialProjectAssignment) {
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
			var actualError error
			BeforeEach(func() {
				httpCall = mockHttpClient.EXPECT().
					Get("/costs/project/"+projectId+"/credentials", nil, gomock.Any()).
					Return(errors.New(errorInfo)).
					Times(1)
				_, actualError = apiClient.CostCredentialIdsInProject(projectId)

			})

			It("should return the error from the api call", func() {
				Expect(actualError).ShouldNot(BeNil())
				Expect(actualError.Error()).Should(Equal(errorInfo))
			})
		})

	})
})
