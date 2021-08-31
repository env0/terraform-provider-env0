package env0

import (
	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestUnitTeamResource(t *testing.T) {
	resourceType := "env0_team"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	t.Run("Success", func(t *testing.T) {
		team := client.Team{
			Id:          "id0",
			Name:        "my-team",
			Description: "team description",
		}

		updatedTeam := client.Team{
			Id:          team.Id,
			Name:        "my-updated-team",
			Description: "updated team description",
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        team.Name,
						"description": team.Description,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", team.Id),
						resource.TestCheckResourceAttr(accessor, "name", team.Name),
						resource.TestCheckResourceAttr(accessor, "description", team.Description),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        updatedTeam.Name,
						"description": updatedTeam.Description,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedTeam.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedTeam.Name),
						resource.TestCheckResourceAttr(accessor, "description", updatedTeam.Description),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().TeamCreate(client.TeamCreatePayload{
				Name:        team.Name,
				Description: team.Description,
			}).Times(1).Return(team, nil)
			mock.EXPECT().TeamUpdate(updatedTeam.Id, client.TeamUpdatePayload{
				Name:        updatedTeam.Name,
				Description: updatedTeam.Description,
			}).Times(1).Return(updatedTeam, nil)

			gomock.InOrder(
				mock.EXPECT().Team(gomock.Any()).Times(2).Return(team, nil),        // 1 after create, 1 before update
				mock.EXPECT().Team(gomock.Any()).Times(1).Return(updatedTeam, nil), // 1 after update
			)

			mock.EXPECT().TeamDelete(team.Id).Times(1)
		})
	})
}
