package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("VCSToken", func() {
	vcsType := "github"
	repository := "http://myrepo.com/"

	mockVcsToken := VscToken{
		Token: 12345,
	}

	Describe("get", func() {
		var ret *VscToken
		var err error

		BeforeEach(func() {
			mockOrganizationIdCall()

			httpCall = mockHttpClient.EXPECT().
				Get("/vcs-token/"+vcsType, map[string]string{
					"organizationId": organizationId,
					"repository":     repository,
				}, gomock.Any()).
				Do(func(path string, request any, response *VscToken) {
					*response = mockVcsToken
				}).Times(1)

			ret, err = apiClient.VcsToken(vcsType, repository)
		})

		It("should return vcs token", func() {
			Expect(*ret).To(Equal(mockVcsToken))
		})

		It("should not return error", func() {
			Expect(err).To(BeNil())
		})
	})

})
