package env0

import (
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitVariableSetAssignmentResource(t *testing.T) {
	resourceType := "env0_variable_set_assignment"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	scope := "environment"
	scopeId := "environment_id"

	t.Run("add assignment and then delete and add", func(t *testing.T) {
		setIds := []string{"a1", "a2"}
		configurationSetIds := []client.ConfigurationSet{
			{
				Id: "a1",
			},
			{
				Id: "a2",
			},
		}
		// Validates that drifts do not occur due to ordering.
		flippedConfigurationSetIds := []client.ConfigurationSet{
			{
				Id: "a2",
			},
			{
				Id: "a1",
			},
		}

		updatedSetIds := []string{"a1", "a3"}
		updatedConfigurationSetIds := []client.ConfigurationSet{
			{
				Id: "a3",
			},
			{
				Id: "a1",
			},
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"scope":    scope,
						"scope_id": scopeId,
						"set_ids":  setIds,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "scope", scope),
						resource.TestCheckResourceAttr(accessor, "scope_id", scopeId),
						resource.TestCheckResourceAttr(accessor, "set_ids.0", setIds[0]),
						resource.TestCheckResourceAttr(accessor, "set_ids.1", setIds[1]),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"scope":    scope,
						"scope_id": scopeId,
						"set_ids":  updatedSetIds,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "scope", scope),
						resource.TestCheckResourceAttr(accessor, "scope_id", scopeId),
						resource.TestCheckResourceAttr(accessor, "set_ids.0", updatedSetIds[0]),
						resource.TestCheckResourceAttr(accessor, "set_ids.1", updatedSetIds[1]),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignConfigurationSets(scope, scopeId, setIds).Times(1).Return(nil),
				mock.EXPECT().ConfigurationSetsAssignments(scope, scopeId).Times(1).Return(configurationSetIds, nil),
				// order doesn't matter - tests that there is no drift due to flipped order.
				mock.EXPECT().ConfigurationSetsAssignments(scope, scopeId).Times(2).Return(flippedConfigurationSetIds, nil),
				// a2 was removed - check that it is unassigned.
				mock.EXPECT().UnassignConfigurationSets(scope, scopeId, []string{"a2"}).Times(1).Return(nil),
				// a3 was added - check that it is assigned.
				mock.EXPECT().AssignConfigurationSets(scope, scopeId, []string{"a3"}).Times(1).Return(nil),
				mock.EXPECT().ConfigurationSetsAssignments(scope, scopeId).Times(1).Return(updatedConfigurationSetIds, nil),
				mock.EXPECT().UnassignConfigurationSets(scope, scopeId, updatedSetIds).Times(1).Return(nil),
			)
		})
	})

	t.Run("add assignment and cause a drift", func(t *testing.T) {
		setIds := []string{"a1"}
		configurationSetIds := []client.ConfigurationSet{
			{
				Id: "a1",
			},
			{
				Id: "a2",
			},
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"scope":    scope,
						"scope_id": scopeId,
						"set_ids":  setIds,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "scope", scope),
						resource.TestCheckResourceAttr(accessor, "scope_id", scopeId),
						resource.TestCheckResourceAttr(accessor, "set_ids.0", setIds[0]),
					),
					ExpectNonEmptyPlan: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AssignConfigurationSets(scope, scopeId, setIds).Times(1).Return(nil),
				mock.EXPECT().ConfigurationSetsAssignments(scope, scopeId).Times(1).Return(configurationSetIds, nil),
				mock.EXPECT().UnassignConfigurationSets(scope, scopeId, []string{"a1", "a2"}).Times(1).Return(nil),
			)
		})
	})
}
