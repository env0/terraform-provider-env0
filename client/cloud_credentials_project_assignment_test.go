package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Credentials Project Assignment", func() {
	projectId := "projectId"
	credentialId := "credentialId"

	expectedResponse := CloudCredentialsProjectAssignment{
		Id:           "assigment id",
		CredentialId: "credentialId",
		ProjectId:    projectId,
	}

	Describe("AssignCloudCredentialsToProject", func() {
		var actualResult CloudCredentialsProjectAssignment
		BeforeEach(func() {

			httpCall = mockHttpClient.EXPECT().
				Put("/credentials/deployment/"+projectId+"/project"+credentialId, nil, gomock.Any()).
				Do(func(path string, request interface{}, response *CloudCredentialsProjectAssignment) {
					*response = expectedResponse
				}).Times(1)
			actualResult, _ = apiClient.AssignCloudCredentialsToProject(projectId, credentialId)

		})

		It("should return the PATCH result", func() {
			Expect(actualResult).To(Equal(expectedResponse))
		})
	})
})
