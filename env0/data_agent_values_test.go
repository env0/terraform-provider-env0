package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAgentValues(t *testing.T) {
	resourceType := "env0_agent_values"
	resourceName := "test"
	accessor := dataSourceAccessor(resourceType, resourceName)

	agentKey := "key"
	values := "values"

	getTestCase := func() resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"agent_key": agentKey,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "agent_key", agentKey),
						resource.TestCheckResourceAttr(accessor, "values", values),
					),
				},
			},
		}
	}

	mockAgentValues := func(returnValue string) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AgentValues(agentKey).AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("Success", func(t *testing.T) {
		runUnitTest(t,
			getTestCase(),
			mockAgentValues(values),
		)
	})

	t.Run("API Call Error", func(t *testing.T) {
		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{
							"agent_key": agentKey,
						}),
						ExpectError: regexp.MustCompile("could not get agent values: error"),
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().AgentValues(agentKey).AnyTimes().Return("", errors.New("error"))
			},
		)
	})

}
