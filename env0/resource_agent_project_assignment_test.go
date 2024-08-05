package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitAgentProjectAssignmentResource(t *testing.T) {
	resourceType := "env0_agent_project_assignment"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	assignment := AgentProjectAssignment{
		ProjectId: "pid",
		AgentId:   "aid1",
	}

	updatedAssignment := AgentProjectAssignment{
		ProjectId: "pid",
		AgentId:   "aid2",
	}

	otherAssignment := AgentProjectAssignment{
		ProjectId: "pid_other",
		AgentId:   "aid_other",
	}

	defaultAssignment := AgentProjectAssignment{
		ProjectId: "pid",
		AgentId:   "default_aid",
	}

	getConfig := func(a *AgentProjectAssignment) string {
		return resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
			"project_id": a.ProjectId,
			"agent_id":   a.AgentId,
		})
	}

	getCheck := func(a *AgentProjectAssignment) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(accessor, "project_id", a.ProjectId),
			resource.TestCheckResourceAttr(accessor, "agent_id", a.AgentId),
		)
	}

	getAssignmentPayload := func(a *AgentProjectAssignment) client.AssignProjectsAgentsAssignmentsPayload {
		return client.AssignProjectsAgentsAssignmentsPayload{
			a.ProjectId: a.AgentId,
		}
	}

	getAssignmentsResponse := func(as []AgentProjectAssignment) *client.ProjectsAgentsAssignments {
		assignments := client.ProjectsAgentsAssignments{
			DefaultAgent: defaultAssignment.AgentId,
		}

		if len(as) > 0 {
			assignments.ProjectsAgents = map[string]interface{}{}

			for _, a := range as {
				assignments.ProjectsAgents[a.ProjectId] = a.AgentId
			}
		}

		return &assignments
	}

	t.Run("create and update", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: getConfig(&assignment),
					Check:  getCheck(&assignment),
				},
				{
					Config: getConfig(&updatedAssignment),
					Check:  getCheck(&updatedAssignment),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignAgentsToProjects(getAssignmentPayload(&assignment)).Times(1).Return(nil, nil),
				mock.EXPECT().ProjectsAgentsAssignments().Times(2).Return(getAssignmentsResponse([]AgentProjectAssignment{otherAssignment, assignment}), nil),
				mock.EXPECT().AssignAgentsToProjects(getAssignmentPayload(&updatedAssignment)).Times(1).Return(nil, nil),
				mock.EXPECT().ProjectsAgentsAssignments().Times(2).Return(getAssignmentsResponse([]AgentProjectAssignment{otherAssignment, updatedAssignment}), nil),
				mock.EXPECT().AssignAgentsToProjects(getAssignmentPayload(&defaultAssignment)).Times(1).Return(nil, nil),
			)
		})
	})

	t.Run("default agent when assignment not found - validate no drift", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: getConfig(&defaultAssignment),
					Check:  getCheck(&defaultAssignment),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignAgentsToProjects(getAssignmentPayload(&defaultAssignment)).Times(1).Return(nil, nil),
				mock.EXPECT().ProjectsAgentsAssignments().Times(2).Return(getAssignmentsResponse([]AgentProjectAssignment{otherAssignment}), nil),
				mock.EXPECT().AssignAgentsToProjects(getAssignmentPayload(&defaultAssignment)).Times(1).Return(nil, nil),
			)
		})
	})

	t.Run("drift", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: getConfig(&assignment),
					Check:  getCheck(&assignment),
				},
				{
					Config:             getConfig(&assignment),
					ExpectNonEmptyPlan: true,
					PlanOnly:           true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignAgentsToProjects(getAssignmentPayload(&assignment)).Times(1).Return(nil, nil),
				mock.EXPECT().ProjectsAgentsAssignments().Times(1).Return(getAssignmentsResponse([]AgentProjectAssignment{assignment}), nil),
				mock.EXPECT().ProjectsAgentsAssignments().Times(3).Return(getAssignmentsResponse([]AgentProjectAssignment{otherAssignment}), nil),
				mock.EXPECT().AssignAgentsToProjects(getAssignmentPayload(&defaultAssignment)).Times(1).Return(nil, nil),
			)
		})
	})

	t.Run("import", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: getConfig(&assignment),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     assignment.ProjectId,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignAgentsToProjects(getAssignmentPayload(&assignment)).Times(1).Return(nil, nil),
				mock.EXPECT().ProjectsAgentsAssignments().Times(1).Return(getAssignmentsResponse([]AgentProjectAssignment{otherAssignment, assignment}), nil),
				mock.EXPECT().Project(assignment.ProjectId).Times(1).Return(client.Project{}, nil),
				mock.EXPECT().ProjectsAgentsAssignments().Times(3).Return(getAssignmentsResponse([]AgentProjectAssignment{otherAssignment, assignment}), nil),
				mock.EXPECT().AssignAgentsToProjects(getAssignmentPayload(&defaultAssignment)).Times(1).Return(nil, nil),
			)
		})
	})

	t.Run("import for unassigned project is default", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: getConfig(&defaultAssignment),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     defaultAssignment.ProjectId,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignAgentsToProjects(getAssignmentPayload(&defaultAssignment)).Times(1).Return(nil, nil),
				mock.EXPECT().ProjectsAgentsAssignments().Times(1).Return(getAssignmentsResponse([]AgentProjectAssignment{otherAssignment, defaultAssignment}), nil),
				mock.EXPECT().Project(assignment.ProjectId).Times(1).Return(client.Project{}, nil),
				mock.EXPECT().ProjectsAgentsAssignments().Times(3).Return(getAssignmentsResponse([]AgentProjectAssignment{otherAssignment}), nil),
				mock.EXPECT().AssignAgentsToProjects(getAssignmentPayload(&defaultAssignment)).Times(1).Return(nil, nil),
			)
		})
	})

	t.Run("import error - project not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: getConfig(&assignment),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     defaultAssignment.ProjectId,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile("unable to get or find a project with id 'pid': error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignAgentsToProjects(getAssignmentPayload(&assignment)).Times(1).Return(nil, nil),
				mock.EXPECT().ProjectsAgentsAssignments().Times(1).Return(getAssignmentsResponse([]AgentProjectAssignment{otherAssignment, assignment}), nil),
				mock.EXPECT().Project(assignment.ProjectId).Times(1).Return(client.Project{}, errors.New("error")),
				mock.EXPECT().ProjectsAgentsAssignments().Times(1).Return(getAssignmentsResponse([]AgentProjectAssignment{otherAssignment, assignment}), nil),
				mock.EXPECT().AssignAgentsToProjects(getAssignmentPayload(&defaultAssignment)).Times(1).Return(nil, nil),
			)
		})
	})
}
