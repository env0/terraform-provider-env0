package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAgentsDataSource(t *testing.T) {
	agent1 := client.Agent{
		AgentKey: "akey1",
	}

	agent2 := client.Agent{
		AgentKey: "akey2",
	}

	resourceType := "env0_agents"
	resourceName := "test_agents"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getTestCase := func() resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "agents.0.agent_key", agent1.AgentKey),
						resource.TestCheckResourceAttr(accessor, "agents.1.agent_key", agent2.AgentKey),
					),
				},
			},
		}
	}

	mockAgents := func(returnValue []client.Agent) func(mockFunc *client.MockApiClientInterface) {
		return func(mock *client.MockApiClientInterface) {
			mock.EXPECT().Agents().AnyTimes().Return(returnValue, nil)
		}
	}

	t.Run("Success", func(t *testing.T) {
		runUnitTest(t,
			getTestCase(),
			mockAgents([]client.Agent{agent1, agent2}),
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
				mock.EXPECT().Agents().AnyTimes().Return(nil, errors.New("error"))
			},
		)
	})
}
