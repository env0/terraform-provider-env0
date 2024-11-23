package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestTeamsDataSource(t *testing.T) {
	team1 := client.Team{
		Id:          "id0",
		Name:        "name1",
		Description: "A team's description",
	}

	team2 := client.Team{
		Id:          "id1",
		Name:        "name2",
		Description: "A team's description",
	}

	team3 := client.Team{
		Id:          "id2",
		Name:        "name3",
		Description: "A team's description",
	}

	resourceType := "env0_teams"
	resourceName := "test_teams"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getTestCase := func(params map[string]interface{}) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, params),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "names.0", team1.Name),
						resource.TestCheckResourceAttr(accessor, "names.1", team2.Name),
						resource.TestCheckNoResourceAttr(accessor, "names.2"),
					),
				},
			},
		}
	}

	mockTeams := func(returnValue []client.Team) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Teams().AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("Success", func(t *testing.T) {
		runUnitTest(t,
			getTestCase(map[string]interface{}{}),
			mockTeams([]client.Team{team1, team2}),
		)
	})

	t.Run("Success with regex filter", func(t *testing.T) {
		runUnitTest(t,
			getTestCase(map[string]interface{}{
				"filter": "name(?:1|2)",
			}),
			mockTeams([]client.Team{team1, team2, team3}),
		)
	})

	t.Run("API Call Error", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config:      dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{}),
						ExpectError: regexp.MustCompile("error"),
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Teams().AnyTimes().Return(nil, errors.New("error"))
			},
		)
	})

	t.Run("invalid regex filter", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{
							"filter": "(ab",
						}),
						ExpectError: regexp.MustCompile("Invalid filter:.+"),
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().Teams().AnyTimes().Return(nil, errors.New("error"))
			},
		)
	})
}
