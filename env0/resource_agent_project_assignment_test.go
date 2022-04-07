package env0

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitAgentProjectAssignmentResource(t *testing.T) {

	// helper functions that receives a variadic list of key value items and returns a ProjectsAgentsAssignments instance.

	GenerateProjectsAgentsAssignmentsMap := func(items ...string) map[string]interface{} {
		res := make(map[string]interface{})

		for i := 0; i < len(items)-1; i += 2 {
			res[items[i]] = items[i+1]
		}

		return res
	}

	GenerateProjectsAgentsAssignments := func(items ...string) *client.ProjectsAgentsAssignments {
		return &client.ProjectsAgentsAssignments{
			ProjectsAgents: GenerateProjectsAgentsAssignmentsMap(items...),
		}
	}

	resourceType := "env0_agent_project_assignment"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	projectId := "pid"
	agentId := "aid"

	t.Run("Create assignment", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id": projectId,
						"agent_id":   agentId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", agentId+"_"+projectId),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "agent_id", agentId),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ProjectsAgentsAssignments().Times(1).Return(GenerateProjectsAgentsAssignments("p111", "a222"), nil),
				mock.EXPECT().AssignAgentsToProjects(GenerateProjectsAgentsAssignmentsMap(
					"p111",
					"a222",
					projectId,
					agentId,
				)).Times(1).Return(nil, nil),
				mock.EXPECT().ProjectsAgentsAssignments().Times(2).Return(GenerateProjectsAgentsAssignments("p111", "a222", projectId, agentId), nil),
				mock.EXPECT().AssignAgentsToProjects(GenerateProjectsAgentsAssignmentsMap(
					"p111",
					"a222",
				)).Times(1).Return(nil, nil),
			)
		})
	})

	t.Run("Assignment already exist", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id": projectId,
						"agent_id":   agentId,
					}),
					ExpectError: regexp.MustCompile(fmt.Sprintf("assignment for project id %v and agent id %v already exist", projectId, agentId)),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ProjectsAgentsAssignments().Times(1).Return(GenerateProjectsAgentsAssignments("p111", "a222", projectId, agentId), nil),
			)
		})
	})

	t.Run("Import Assignment", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id": projectId,
						"agent_id":   agentId,
					}),
				},
				{
					ResourceName:      resourceType + "." + resourceName,
					ImportState:       true,
					ImportStateId:     agentId + "_" + projectId,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ProjectsAgentsAssignments().Times(1).Return(GenerateProjectsAgentsAssignments(), nil),
				mock.EXPECT().AssignAgentsToProjects(GenerateProjectsAgentsAssignmentsMap(
					projectId,
					agentId,
				)).Times(1).Return(nil, nil),
				mock.EXPECT().ProjectsAgentsAssignments().Times(4).Return(GenerateProjectsAgentsAssignments(projectId, agentId), nil),
				mock.EXPECT().AssignAgentsToProjects(GenerateProjectsAgentsAssignmentsMap()).Times(1).Return(nil, nil),
			)
		})
	})

	t.Run("Import Assignment with invalid id", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id": projectId,
						"agent_id":   agentId,
					}),
				},
				{
					ResourceName:      resourceType + "." + resourceName,
					ImportState:       true,
					ImportStateId:     "invalid",
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile("the id invalid is invalid must be <agent_id>_<project_id>"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ProjectsAgentsAssignments().Times(1).Return(GenerateProjectsAgentsAssignments(), nil),
				mock.EXPECT().AssignAgentsToProjects(GenerateProjectsAgentsAssignmentsMap(
					projectId,
					agentId,
				)).Times(1).Return(nil, nil),
				mock.EXPECT().ProjectsAgentsAssignments().Times(2).Return(GenerateProjectsAgentsAssignments(projectId, agentId), nil),
				mock.EXPECT().AssignAgentsToProjects(GenerateProjectsAgentsAssignmentsMap()).Times(1).Return(nil, nil),
			)
		})
	})

	t.Run("Import Assignment id not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id": projectId,
						"agent_id":   agentId,
					}),
				},
				{
					ResourceName:      resourceType + "." + resourceName,
					ImportState:       true,
					ImportStateId:     "pid22_aid22",
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile("assignment with id pid22_aid22 not found"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ProjectsAgentsAssignments().Times(1).Return(GenerateProjectsAgentsAssignments(), nil),
				mock.EXPECT().AssignAgentsToProjects(GenerateProjectsAgentsAssignmentsMap(
					projectId,
					agentId,
				)).Times(1).Return(nil, nil),
				mock.EXPECT().ProjectsAgentsAssignments().Times(3).Return(GenerateProjectsAgentsAssignments(projectId, agentId), nil),
				mock.EXPECT().AssignAgentsToProjects(GenerateProjectsAgentsAssignmentsMap()).Times(1).Return(nil, nil),
			)
		})
	})
}
