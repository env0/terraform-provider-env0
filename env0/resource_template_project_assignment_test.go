package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitTemplateProjectAssignmentResource(t *testing.T) {

	resourceType := "env0_template_project_assignment"

	resourceName := "test"

	resourceTemplateAssignment := map[string]interface{}{
		"template_id": "id",
		"project_id":  "pid",
	}

	resourceTemplateAssignmentUpdate := map[string]interface{}{
		"template_id": "update",
		"project_id":  "updatepid",
	}

	accessor := resourceAccessor(resourceType, resourceName)

	payLoad := client.TemplateAssignmentToProjectPayload{
		ProjectId: resourceTemplateAssignment["project_id"].(string),
	}

	updatePayload := client.TemplateAssignmentToProjectPayload{
		ProjectId: resourceTemplateAssignmentUpdate["project_id"].(string),
	}

	returnValues := client.Template{
		Id:         "tid",
		ProjectIds: []string{"pid", "other-id"},
	}
	driftReturnValues := client.Template{
		Id:         "tid",
		ProjectIds: []string{"other-id"},
	}
	updateReturnValues := client.Template{
		Id:         "updatetid",
		ProjectIds: []string{"updatepid"},
	}

	testCaseforCreate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, resourceTemplateAssignment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "template_id", resourceTemplateAssignment["template_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "project_id", resourceTemplateAssignment["project_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", "tid|pid"),
				),
			},
			{
				Config: resourceConfigCreate(resourceType, resourceName, resourceTemplateAssignmentUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "template_id", resourceTemplateAssignmentUpdate["template_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "project_id", resourceTemplateAssignmentUpdate["project_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", "updatetid|updatepid"),
				),
			},
		},
	}

	testCaseForError := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config:      resourceConfigCreate(resourceType, resourceName, map[string]interface{}{"template_id": "id"}),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	}

	testCaseForApiclientError := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config:      resourceConfigCreate(resourceType, resourceName, resourceTemplateAssignment),
				ExpectError: regexp.MustCompile("could not assign template to project: error"),
			},
		},
	}

	t.Run("create", func(t *testing.T) {
		runUnitTest(t, testCaseforCreate, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AssignTemplateToProject(resourceTemplateAssignment["template_id"].(string), payLoad).
				Times(1).Return(returnValues, nil)

			mock.EXPECT().RemoveTemplateFromProject(resourceTemplateAssignment["template_id"].(string),
				resourceTemplateAssignment["project_id"].(string)).Times(1).Return(nil)

			mock.EXPECT().AssignTemplateToProject(resourceTemplateAssignmentUpdate["template_id"].(string), updatePayload).
				Times(1).Return(updateReturnValues, nil)

			mock.EXPECT().RemoveTemplateFromProject(resourceTemplateAssignmentUpdate["template_id"].(string),
				resourceTemplateAssignmentUpdate["project_id"].(string)).Times(1).Return(nil)

			gomock.InOrder(
				mock.EXPECT().Template(resourceTemplateAssignment["template_id"].(string)).Times(2).
					Return(returnValues, nil),
				mock.EXPECT().Template(resourceTemplateAssignmentUpdate["template_id"].(string)).Times(1).
					Return(updateReturnValues, nil),
			)
		})

	})

	t.Run("throw error when missing values", func(t *testing.T) {
		runUnitTest(t, testCaseForError, func(mock *client.MockApiClientInterface) {

		})
	})

	t.Run("detect error when apiclient.AssignTemplateToProject throw error", func(t *testing.T) {
		runUnitTest(t, testCaseForApiclientError, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AssignTemplateToProject(resourceTemplateAssignment["template_id"].(string), payLoad).
				Times(1).Return(client.Template{}, errors.New("error"))
		})
	})

	t.Run("detect drift", func(t *testing.T) {

		runUnitTest(t, testCaseforCreate, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AssignTemplateToProject(resourceTemplateAssignment["template_id"].(string), payLoad).
				Times(1).Return(returnValues, nil)

			mock.EXPECT().AssignTemplateToProject(resourceTemplateAssignmentUpdate["template_id"].(string), updatePayload).
				Times(1).Return(updateReturnValues, nil)

			mock.EXPECT().RemoveTemplateFromProject(resourceTemplateAssignmentUpdate["template_id"].(string),
				resourceTemplateAssignmentUpdate["project_id"].(string)).Times(1).Return(nil)

			gomock.InOrder(
				mock.EXPECT().Template(resourceTemplateAssignment["template_id"].(string)).Times(1).
					Return(returnValues, nil),
				mock.EXPECT().Template(resourceTemplateAssignment["template_id"].(string)).Times(1).
					Return(driftReturnValues, nil),
				mock.EXPECT().Template(resourceTemplateAssignmentUpdate["template_id"].(string)).Times(1).
					Return(updateReturnValues, nil),
			)
		})
	})
}
