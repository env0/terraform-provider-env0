package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const sshKeyName = "new_ssh_key"
const sshKeyValue = "fake key"

var _ = Describe("SshKey", func() {
	var sshKey SshKey
	mockSshKey := SshKey{
		Id:             "idX",
		Name:           sshKeyName,
		Value:          sshKeyValue,
		OrganizationId: organizationId,
	}

	Describe("SshKeyCreate", func() {
		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)

			httpCall = mockHttpClient.EXPECT().
				Post("/ssh-key", map[string]interface{}{
					"name":           sshKeyName,
					"value":          "***************",
					"organizationId": organizationId,
					"id":             "idX",
				},
					gomock.Any()).
				Do(func(path string, request interface{}, response *SshKey) {
					*response = mockSshKey
				})

			sshKey, _ = apiClient.SshKeyCreate(SshKeyCreatePayload{Name: sshKeyName, Value: sshKeyValue, OrganizationId: organizationId})
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send POST request with params", func() {
			httpCall.Times(1)
		})

		It("Should return project", func() {
			Expect(sshKey).To(Equal(mockSshKey))
		})
	})

	//Describe("ProjectDelete", func() {
	//	BeforeEach(func() {
	//		httpCall = mockHttpClient.EXPECT().Delete("/projects/" + mockProject.Id)
	//		apiClient.ProjectDelete(mockProject.Id)
	//	})
	//
	//	It("Should send DELETE request with project id", func() {
	//		httpCall.Times(1)
	//	})
	//})
	//
	//Describe("Project", func() {
	//	BeforeEach(func() {
	//		httpCall = mockHttpClient.EXPECT().
	//			Get("/projects/"+mockProject.Id, nil, gomock.Any()).
	//			Do(func(path string, request interface{}, response *Project) {
	//				*response = mockProject
	//			})
	//		project, _ = apiClient.Project(mockProject.Id)
	//	})
	//
	//	It("Should send GET request with project id", func() {
	//		httpCall.Times(1)
	//	})
	//
	//	It("Should return project", func() {
	//		Expect(project).To(Equal(mockProject))
	//	})
	//})

})
