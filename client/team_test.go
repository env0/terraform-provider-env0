package client_test

import (
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Teams Client", func() {
	mockTeam := Team{
		Id:   "team-id",
		Name: "team-name",
	}

	Describe("Team", func() {
		var returnedTeam Team

		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().
				Get("/teams/"+mockTeam.Id, gomock.Nil(), gomock.Any()).
				Do(func(path string, request interface{}, response *Team) {
					*response = mockTeam
				})
			returnedTeam, _ = apiClient.Team(mockTeam.Id)
		})

		It("Should send GET request", func() {
			httpCall.Times(1)
		})

		It("Should return template", func() {
			Expect(returnedTeam).To(Equal(mockTeam))
		})
	})
})
