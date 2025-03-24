package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestCustomFlowDataSource(t *testing.T) {
	customFlow := client.CustomFlow{
		Id:         "id0",
		Name:       "name0",
		Repository: "rep1",
		Path:       "path",
	}

	otherCustomFlow := client.CustomFlow{
		Id:         "id1",
		Name:       "name1",
		Repository: "rep2",
		Path:       "path",
	}

	customFlowFieldsByName := map[string]any{"name": customFlow.Name}
	customFlowFieldsById := map[string]any{"id": customFlow.Id}

	resourceType := "env0_custom_flow"
	resourceName := "test_custom_flow"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(input map[string]any) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", customFlow.Id),
						resource.TestCheckResourceAttr(accessor, "name", customFlow.Name),
					),
				},
			},
		}
	}

	getErrorTestCase := func(input map[string]any, expectedError string) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      dataSourceConfigCreate(resourceType, resourceName, input),
					ExpectError: regexp.MustCompile(expectedError),
				},
			},
		}
	}

	mockListCustomFlowsCall := func(returnValue []client.CustomFlow) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CustomFlows(customFlow.Name).AnyTimes().Return(returnValue, nil)
		}
	}

	mockCustomFlowCall := func(returnValue *client.CustomFlow) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CustomFlow(customFlow.Id).AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("By ID", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(customFlowFieldsById),
			mockCustomFlowCall(&customFlow),
		)
	})

	t.Run("By Name", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(customFlowFieldsByName),
			mockListCustomFlowsCall([]client.CustomFlow{customFlow}),
		)
	})

	t.Run("Throw error when by name and more than one custom flow exists", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(customFlowFieldsByName, "found multiple custom flows with name"),
			mockListCustomFlowsCall([]client.CustomFlow{customFlow, otherCustomFlow, customFlow}),
		)
	})

	t.Run("Throw error when by id and no custom flow found with that id", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(customFlowFieldsById, "failed to get custom flow by id"),
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().CustomFlow(customFlow.Id).Times(1).Return(nil, http.NewMockFailedResponseError(404))
			},
		)
	})
}
