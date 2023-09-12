package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("SshKey", func() {
	const sshKeyName = "new_ssh_key"
	const sshKeyValue = "fake key"
	mockSshKey := SshKey{
		Id:             "123",
		Name:           sshKeyName,
		Value:          sshKeyValue,
		OrganizationId: organizationId,
	}

	Describe("SshKeyCreate", func() {
		var sshKey *SshKey

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)
			expectedPayload := SshKeyCreatePayload{Name: sshKeyName, Value: sshKeyValue, OrganizationId: organizationId}
			httpCall = mockHttpClient.EXPECT().
				Post("/ssh-keys", expectedPayload, gomock.Any()).
				Do(func(path string, request interface{}, response *SshKey) {
					*response = mockSshKey
				})

			sshKey, _ = apiClient.SshKeyCreate(SshKeyCreatePayload{Name: sshKeyName, Value: sshKeyValue})
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send POST request with params", func() {
			httpCall.Times(1)
		})

		It("Should return project", func() {
			Expect(*sshKey).To(Equal(mockSshKey))
		})
	})

	Describe("SshKeyDelete", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/ssh-keys/"+mockSshKey.Id, nil)
			_ = apiClient.SshKeyDelete(mockSshKey.Id)
		})

		It("Should send DELETE request once for correct id", func() {
			httpCall.Times(1)
		})
	})

	Describe("SshKeys", func() {
		var sshKeys []SshKey
		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)
			httpCall = mockHttpClient.EXPECT().
				Get("/ssh-keys",
					map[string]string{"organizationId": organizationId},
					gomock.Any()).
				Do(func(path string, request interface{}, response *[]SshKey) {
					*response = []SshKey{mockSshKey}
				})

			sshKeys, _ = apiClient.SshKeys()
		})

		It("Should send GET request once", func() {
			httpCall.Times(1)
		})

		It("Should return ssh keys", func() {
			Expect(sshKeys).Should(HaveLen(1))
			Expect(sshKeys).Should(ContainElement(mockSshKey))
		})
	})
})
