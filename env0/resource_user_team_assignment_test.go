package env0

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitUserTeamAssignmentResource(t *testing.T) {
	userId := "uid"
	teamId := "tid"

	GenerateTeam := func(teamId string, userIds []string) client.Team {
		team := client.Team{
			Id:          teamId,
			Name:        "name",
			Description: "description",
		}

		for _, userId := range userIds {
			team.Users = append(team.Users, client.User{UserId: userId})
		}

		return team
	}

	GenerateUpdateTeamPayload := func(userIds []string) client.TeamUpdatePayload {
		return client.TeamUpdatePayload{
			Name:        "name",
			Description: "description",
			UserIds:     userIds,
		}
	}

	resourceType := "env0_user_team_assignment"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	t.Run("Create assignment", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"user_id": userId,
						"team_id": teamId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", userId+"_"+teamId),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "team_id", teamId),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().Team(teamId).Times(1).Return(GenerateTeam(teamId, []string{"otherId"}), nil),
				mock.EXPECT().TeamUpdate(teamId, GenerateUpdateTeamPayload([]string{userId, "otherId"})).Times(1).Return(client.Team{}, nil),
				mock.EXPECT().Team(teamId).Times(2).Return(GenerateTeam(teamId, []string{userId, "otherId"}), nil),
				mock.EXPECT().TeamUpdate(teamId, GenerateUpdateTeamPayload([]string{"otherId"})).Times(1).Return(client.Team{}, nil),
			)
		})
	})

	t.Run("Create Assignment - already exist", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"user_id": userId,
						"team_id": teamId,
					}),
					ExpectError: regexp.MustCompile(fmt.Sprintf("assignment for user id %v and team id %v already exist", userId, teamId)),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().Team(teamId).Times(1).Return(GenerateTeam(teamId, []string{userId}), nil),
			)
		})
	})

	t.Run("Create Assignment - failed to get team", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"user_id": userId,
						"team_id": teamId,
					}),
					ExpectError: regexp.MustCompile("could not get team: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().Team(teamId).Times(1).Return(client.Team{}, errors.New("error")),
			)
		})
	})

	t.Run("Create Assignment - failed to update team", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"user_id": userId,
						"team_id": teamId,
					}),
					ExpectError: regexp.MustCompile("could not update team with new assignment: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().Team(teamId).Times(1).Return(GenerateTeam(teamId, []string{"otherId"}), nil),
				mock.EXPECT().TeamUpdate(teamId, GenerateUpdateTeamPayload([]string{userId, "otherId"})).Times(1).Return(client.Team{}, errors.New("error")),
			)
		})
	})

	t.Run("Create Assignment - drift detected", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"user_id": userId,
						"team_id": teamId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", userId+"_"+teamId),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "team_id", teamId),
					),
					ExpectNonEmptyPlan: true,
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"user_id": userId,
						"team_id": teamId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", userId+"_"+teamId),
						resource.TestCheckResourceAttr(accessor, "user_id", userId),
						resource.TestCheckResourceAttr(accessor, "team_id", teamId),
					),
					PlanOnly:           true,
					ExpectNonEmptyPlan: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().Team(teamId).Times(1).Return(GenerateTeam(teamId, []string{"otherId"}), nil),
				mock.EXPECT().TeamUpdate(teamId, GenerateUpdateTeamPayload([]string{userId, "otherId"})).Times(1).Return(client.Team{}, nil),
				mock.EXPECT().Team(teamId).Times(1).Return(GenerateTeam(teamId, []string{"otherId"}), nil),
			)
		})
	})

	t.Run("Import Assignment", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"user_id": userId,
						"team_id": teamId,
					}),
				},
				{
					ResourceName:      resourceType + "." + resourceName,
					ImportState:       true,
					ImportStateId:     userId + "_" + teamId,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().Team(teamId).Times(1).Return(GenerateTeam(teamId, []string{"otherId"}), nil),
				mock.EXPECT().TeamUpdate(teamId, GenerateUpdateTeamPayload([]string{userId, "otherId"})).Times(1).Return(client.Team{}, nil),
				mock.EXPECT().Team(teamId).Times(4).Return(GenerateTeam(teamId, []string{userId, "otherId"}), nil),
				mock.EXPECT().TeamUpdate(teamId, GenerateUpdateTeamPayload([]string{"otherId"})).Times(1).Return(client.Team{}, nil),
			)
		})
	})

	t.Run("Import Assignment - invalid id", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"user_id": userId,
						"team_id": teamId,
					}),
				},
				{
					ResourceName:      resourceType + "." + resourceName,
					ImportState:       true,
					ImportStateId:     "invalid",
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile("the id invalid is invalid must be <user_id>_<team_id>"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().Team(teamId).Times(1).Return(GenerateTeam(teamId, []string{"otherId"}), nil),
				mock.EXPECT().TeamUpdate(teamId, GenerateUpdateTeamPayload([]string{userId, "otherId"})).Times(1).Return(client.Team{}, nil),
				mock.EXPECT().Team(teamId).Times(2).Return(GenerateTeam(teamId, []string{userId, "otherId"}), nil),
				mock.EXPECT().TeamUpdate(teamId, GenerateUpdateTeamPayload([]string{"otherId"})).Times(1).Return(client.Team{}, nil),
			)
		})
	})

	t.Run("Import Assignment - team not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"user_id": userId,
						"team_id": teamId,
					}),
				},
				{
					ResourceName:      resourceType + "." + resourceName,
					ImportState:       true,
					ImportStateId:     "uid22_tid22",
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile("team tid22 not found"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().Team(teamId).Times(1).Return(GenerateTeam(teamId, []string{"otherId"}), nil),
				mock.EXPECT().TeamUpdate(teamId, GenerateUpdateTeamPayload([]string{userId, "otherId"})).Times(1).Return(client.Team{}, nil),
				mock.EXPECT().Team(teamId).Times(1).Return(GenerateTeam(teamId, []string{userId, "otherId"}), nil),
				mock.EXPECT().Team("tid22").Times(1).Return(client.Team{}, http.NewMockFailedResponseError(404)),
				mock.EXPECT().Team(teamId).Times(1).Return(GenerateTeam(teamId, []string{userId, "otherId"}), nil),
				mock.EXPECT().TeamUpdate(teamId, GenerateUpdateTeamPayload([]string{"otherId"})).Times(1).Return(client.Team{}, nil),
			)
		})
	})

	t.Run("Import Assignment - user not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"user_id": userId,
						"team_id": teamId,
					}),
				},
				{
					ResourceName:      resourceType + "." + resourceName,
					ImportState:       true,
					ImportStateId:     "uid22_tid22",
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile("user uid22 not assigned to team tid22"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().Team(teamId).Times(1).Return(GenerateTeam(teamId, []string{"otherId"}), nil),
				mock.EXPECT().TeamUpdate(teamId, GenerateUpdateTeamPayload([]string{userId, "otherId"})).Times(1).Return(client.Team{}, nil),
				mock.EXPECT().Team(teamId).Times(1).Return(GenerateTeam(teamId, []string{userId, "otherId"}), nil),
				mock.EXPECT().Team("tid22").Times(1).Return(GenerateTeam(teamId, []string{userId, "otherId"}), nil),
				mock.EXPECT().Team(teamId).Times(1).Return(GenerateTeam(teamId, []string{userId, "otherId"}), nil),
				mock.EXPECT().TeamUpdate(teamId, GenerateUpdateTeamPayload([]string{"otherId"})).Times(1).Return(client.Team{}, nil),
			)
		})
	})
}
