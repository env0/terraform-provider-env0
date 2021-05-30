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
		Name:           sshKeyName,
		Value:          sshKeyValue,
		OrganizationId: organizationId,
	}

	Describe("SshKeyCreate", func() {
		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)
			expectedPayload := SshKeyCreatePayloadExtended{SshKeyCreatePayload: SshKeyCreatePayload{
				Name:  sshKeyName,
				Value: sshKeyValue,
			}, OrganizationId: organizationId}
			httpCall = mockHttpClient.EXPECT().
				Post("/ssh-keys", expectedPayload,
					gomock.Any()).
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
			Expect(sshKey).To(Equal(mockSshKey))
		})
	})
})
