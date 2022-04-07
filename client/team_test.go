package client_test

import (
	"errors"

	. "github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/copier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Teams Client", func() {
	mockTeam := Team{
		Id:   "team-id",
		Name: "team-name",
	}

	Describe("Get Single Team", func() {
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

		It("Should return team", func() {
			Expect(returnedTeam).To(Equal(mockTeam))
		})
	})

	Describe("Get All Teams", func() {
		var returnedTeams []Team
		mockTeams := []Team{mockTeam}

		BeforeEach(func() {
			mockOrganizationIdCall(organizationId)
			httpCall = mockHttpClient.EXPECT().
				Get("/teams/organizations/"+organizationId, nil, gomock.Any()).
				Do(func(path string, request interface{}, response *[]Team) {
					*response = mockTeams
				})
			returnedTeams, _ = apiClient.Teams()
		})

		It("Should get organization id", func() {
			organizationIdCall.Times(1)
		})

		It("Should send GET request", func() {
			httpCall.Times(1)
		})

		It("Should return teams", func() {
			Expect(returnedTeams).To(Equal(mockTeams))
		})
	})

	Describe("TeamCreate", func() {
		Describe("Success", func() {
			var createdTeam Team
			var err error

			BeforeEach(func() {
				mockOrganizationIdCall(organizationId)

				createTeamPayload := TeamCreatePayload{}
				copier.Copy(&createTeamPayload, &mockTeam)

				expectedCreateRequest := createTeamPayload
				expectedCreateRequest.OrganizationId = organizationId

				httpCall = mockHttpClient.EXPECT().
					Post("/teams", expectedCreateRequest, gomock.Any()).
					Do(func(path string, request interface{}, response *Team) {
						*response = mockTeam
					})

				createdTeam, err = apiClient.TeamCreate(createTeamPayload)
			})

			It("Should get organization id", func() {
				organizationIdCall.Times(1)
			})

			It("Should send POST request with params", func() {
				httpCall.Times(1)
			})

			It("Should not return error", func() {
				Expect(err).To(BeNil())
			})

			It("Should return created team", func() {
				Expect(createdTeam).To(Equal(mockTeam))
			})
		})

		Describe("Failure", func() {
			It("Should fail when team has no name", func() {
				teamWithoutNamePayload := TeamCreatePayload{Description: "team-without-name"}
				_, err := apiClient.TeamCreate(teamWithoutNamePayload)
				Expect(err).To(BeEquivalentTo(errors.New("must specify team name on creation")))
			})

			It("Should fail if request includes organizationId (should be inferred automatically)", func() {
				payloadWithOrgId := TeamCreatePayload{Name: "team-name", OrganizationId: "org-id"}
				_, err := apiClient.TeamCreate(payloadWithOrgId)
				Expect(err).To(BeEquivalentTo(errors.New("must not specify organizationId")))
			})
		})
	})

	Describe("TeamDelete", func() {
		BeforeEach(func() {
			httpCall = mockHttpClient.EXPECT().Delete("/teams/" + mockTeam.Id)
			apiClient.TeamDelete(mockTeam.Id)
		})

		It("Should send DELETE request with team id", func() {
			httpCall.Times(1)
		})
	})

	Describe("TeamUpdate", func() {
		Describe("Success", func() {
			var updatedTeam Team
			var err error

			BeforeEach(func() {
				updateTeamPayload := TeamUpdatePayload{Name: "updated-name"}

				httpCall = mockHttpClient.EXPECT().
					Put("/teams/"+mockTeam.Id, updateTeamPayload, gomock.Any()).
					Do(func(path string, request interface{}, response *Team) {
						*response = mockTeam
					})

				updatedTeam, err = apiClient.TeamUpdate(mockTeam.Id, updateTeamPayload)
			})

			It("Should send Put request with expected payload", func() {
				httpCall.Times(1)
			})

			It("Should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("Should return team received from API", func() {
				Expect(updatedTeam).To(Equal(mockTeam))
			})
		})

		Describe("Failure", func() {
			It("Should fail if team has no name", func() {
				payloadWithNoName := TeamUpdatePayload{Description: "team-without-name"}
				_, err := apiClient.TeamUpdate(mockTeam.Id, payloadWithNoName)
				Expect(err).To(BeEquivalentTo(errors.New("must specify team name on update")))
			})
		})
	})
})
