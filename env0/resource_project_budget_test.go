package env0

import (
	"errors"
	"regexp"
	"strconv"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitProjectBudgetResource(t *testing.T) {
	resourceType := "env0_project_budget"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	projectBudget := &client.ProjectBudget{
		Id:         "id",
		ProjectId:  "pid",
		Amount:     10,
		Timeframe:  "MONTHLY",
		Thresholds: []int{1},
	}

	updatedProjectBudget := &client.ProjectBudget{
		Id:         "id",
		ProjectId:  "pid",
		Amount:     20,
		Timeframe:  "WEEKLY",
		Thresholds: []int{2},
	}

	t.Run("create and update", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id": projectBudget.ProjectId,
						"amount":     strconv.Itoa(projectBudget.Amount),
						"timeframe":  projectBudget.Timeframe,
						"thresholds": projectBudget.Thresholds,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "project_id", projectBudget.ProjectId),
						resource.TestCheckResourceAttr(accessor, "amount", strconv.Itoa(projectBudget.Amount)),
						resource.TestCheckResourceAttr(accessor, "timeframe", projectBudget.Timeframe),
						resource.TestCheckResourceAttr(accessor, "thresholds.0", strconv.Itoa(projectBudget.Thresholds[0])),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id": updatedProjectBudget.ProjectId,
						"amount":     strconv.Itoa(updatedProjectBudget.Amount),
						"timeframe":  updatedProjectBudget.Timeframe,
						"thresholds": updatedProjectBudget.Thresholds,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "project_id", updatedProjectBudget.ProjectId),
						resource.TestCheckResourceAttr(accessor, "amount", strconv.Itoa(updatedProjectBudget.Amount)),
						resource.TestCheckResourceAttr(accessor, "timeframe", updatedProjectBudget.Timeframe),
						resource.TestCheckResourceAttr(accessor, "thresholds.0", strconv.Itoa(updatedProjectBudget.Thresholds[0])),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ProjectBudgetUpdate(projectBudget.ProjectId, &client.ProjectBudgetUpdatePayload{
					Amount:     projectBudget.Amount,
					Timeframe:  projectBudget.Timeframe,
					Thresholds: projectBudget.Thresholds,
				}).Times(1).Return(projectBudget, nil),
				mock.EXPECT().ProjectBudget(projectBudget.ProjectId).Times(2).Return(projectBudget, nil),
				mock.EXPECT().ProjectBudgetUpdate(updatedProjectBudget.ProjectId, &client.ProjectBudgetUpdatePayload{
					Amount:     updatedProjectBudget.Amount,
					Timeframe:  updatedProjectBudget.Timeframe,
					Thresholds: updatedProjectBudget.Thresholds,
				}).Times(1).Return(updatedProjectBudget, nil),
				mock.EXPECT().ProjectBudget(updatedProjectBudget.ProjectId).Times(1).Return(updatedProjectBudget, nil),
				mock.EXPECT().ProjectBudgetDelete(projectBudget.ProjectId).Times(1).Return(nil),
			)
		})
	})

	t.Run("Create Failure - invalid timeframe", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id": projectBudget.ProjectId,
						"amount":     strconv.Itoa(projectBudget.Amount),
						"timeframe":  "invalid",
						"thresholds": projectBudget.Thresholds,
					}),
					ExpectError: regexp.MustCompile("must be one of: WEEKLY, MONTHLY, QUARTERLY, YEARLY"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("Detect drift", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id": projectBudget.ProjectId,
						"amount":     strconv.Itoa(projectBudget.Amount),
						"timeframe":  projectBudget.Timeframe,
						"thresholds": projectBudget.Thresholds,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "project_id", projectBudget.ProjectId),
						resource.TestCheckResourceAttr(accessor, "amount", strconv.Itoa(projectBudget.Amount)),
						resource.TestCheckResourceAttr(accessor, "timeframe", projectBudget.Timeframe),
						resource.TestCheckResourceAttr(accessor, "thresholds.0", strconv.Itoa(projectBudget.Thresholds[0])),
					),
					ExpectNonEmptyPlan: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ProjectBudgetUpdate(projectBudget.ProjectId, &client.ProjectBudgetUpdatePayload{
					Amount:     projectBudget.Amount,
					Timeframe:  projectBudget.Timeframe,
					Thresholds: projectBudget.Thresholds,
				}).Times(1).Return(projectBudget, nil),
				mock.EXPECT().ProjectBudget(projectBudget.ProjectId).Times(1).Return(nil, &client.NotFoundError{}),
				mock.EXPECT().ProjectBudgetDelete(projectBudget.ProjectId).Times(1).Return(nil),
			)
		})
	})

	t.Run("Failure in create", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id": projectBudget.ProjectId,
						"amount":     strconv.Itoa(projectBudget.Amount),
						"timeframe":  projectBudget.Timeframe,
						"thresholds": projectBudget.Thresholds,
					}),
					ExpectError: regexp.MustCompile("could not create or update budget: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ProjectBudgetUpdate(projectBudget.ProjectId, &client.ProjectBudgetUpdatePayload{
				Amount:     projectBudget.Amount,
				Timeframe:  projectBudget.Timeframe,
				Thresholds: projectBudget.Thresholds,
			}).Times(1).Return(nil, errors.New("error"))
		})
	})
}
