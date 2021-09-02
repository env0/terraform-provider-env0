package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestTeamDataSource(t *testing.T) {
	team := client.Team{
		Id:          "id0",
		Name:        "my-team-1",
		Description: "A team's description",
	}

	otherTeam := client.Team{
		Id:          "other-id",
		Name:        "other-name",
		Description: team.Description,
	}

	teamFieldsByName := map[string]interface{}{"name": team.Name}
	teamFieldsById := map[string]interface{}{"id": team.Id}

	resourceType := "env0_team"
	resourceName := "test_team"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(input map[string]interface{}) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", team.Id),
						resource.TestCheckResourceAttr(accessor, "name", team.Name),
						resource.TestCheckResourceAttr(accessor, "description", team.Description),
					),
				},
			},
		}
	}

	getErrorTestCase := func(input map[string]interface{}, expectedError string) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      dataSourceConfigCreate(resourceType, resourceName, input),
					ExpectError: regexp.MustCompile(expectedError),
				},
			},
		}
	}

	mockGetTeamCall := func(returnValue client.Team) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Team(team.Id).AnyTimes().Return(returnValue, nil)
		}
	}

	mockListTeamsCall := func(returnValue []client.Team) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Teams().AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("By ID", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(teamFieldsById),
			mockGetTeamCall(team),
		)
	})

	t.Run("By Name", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(teamFieldsByName),
			mockListTeamsCall([]client.Team{team, otherTeam}),
		)
	})

	t.Run("Throw error when no name or id is supplied", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(map[string]interface{}{}, "one of `id,name` must be specified"),
			func(mock *client.MockApiClientInterface) {},
		)
	})

	t.Run("Throw error when by name and more than one team exists", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(teamFieldsByName, "Found multiple teams for name"),
			mockListTeamsCall([]client.Team{team, team, otherTeam}),
		)
	})

	t.Run("Throw error when by name and no projects found at all", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(teamFieldsByName, "Could not find an env0 team with name"),
			mockListTeamsCall([]client.Team{}),
		)
	})

	t.Run("Throw error when by name and no teams found with that name", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(teamFieldsByName, "Could not find an env0 team with name"),
			mockListTeamsCall([]client.Team{otherTeam}),
		)
	})
}
