package env0

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAgentPoolDataSource(t *testing.T) {
	agentPool := client.AgentPool{
		Id:          "id0",
		Name:        "pool0",
		Description: "desc0",
		AgentKey:    "key0",
	}

	otherAgentPool := client.AgentPool{
		Id:          "id1",
		Name:        "pool1",
		Description: "desc1",
		AgentKey:    "key1",
	}

	fieldsByName := map[string]any{"name": agentPool.Name}
	fieldsById := map[string]any{"id": agentPool.Id}

	resourceType := "env0_agent_pool"
	resourceName := "test_agent_pool"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(input map[string]any) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", agentPool.Id),
						resource.TestCheckResourceAttr(accessor, "name", agentPool.Name),
						resource.TestCheckResourceAttr(accessor, "description", agentPool.Description),
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

	t.Run("By ID", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(fieldsById),
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().AgentPool(agentPool.Id).AnyTimes().Return(&agentPool, nil)
			},
		)
	})

	t.Run("By Name", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(fieldsByName),
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().AgentPools().AnyTimes().Return([]client.AgentPool{agentPool, otherAgentPool}, nil)
			},
		)
	})

	t.Run("By ID With Logs", func(t *testing.T) {
		agentPoolWithLogs := client.AgentPool{
			Id:          "id-logs",
			Name:        "pool-logs",
			Description: "with logs",
			AgentKey:    "key-logs",
			Logs: &client.AgentPoolLogsConfig{
				Dynamo: &client.AgentPoolDynamoLogs{
					SelfHosted: &client.AgentPoolSelfHostedLogs{
						AccountId: "123456789",
						Region:    "us-east-1",
					},
				},
			},
		}

		runUnitTest(t,
			resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: dataSourceConfigCreate(resourceType, resourceName, map[string]any{"id": agentPoolWithLogs.Id}),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(accessor, "id", agentPoolWithLogs.Id),
							resource.TestCheckResourceAttr(accessor, "logs.0.account_id", "123456789"),
							resource.TestCheckResourceAttr(accessor, "logs.0.region", "us-east-1"),
						),
					},
				},
			},
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().AgentPool(agentPoolWithLogs.Id).AnyTimes().Return(&agentPoolWithLogs, nil)
			},
		)
	})

	t.Run("Throw error when no name or id is supplied", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(map[string]any{}, "one of `id,name` must be specified"),
			func(mock *client.MockApiClientInterface) {},
		)
	})

	t.Run("Throw error when by name and more than one agent pool exists", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(fieldsByName, "found multiple agent pools"),
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().AgentPools().AnyTimes().Return([]client.AgentPool{agentPool, otherAgentPool, agentPool}, nil)
			},
		)
	})

	t.Run("Throw error when by name and no agent pool found", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(fieldsByName, "not found"),
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().AgentPools().AnyTimes().Return([]client.AgentPool{otherAgentPool}, nil)
			},
		)
	})

	t.Run("Throw error when by id and API error", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(fieldsById, "could not read agent pool"),
			func(mock *client.MockApiClientInterface) {
				mock.EXPECT().AgentPool(agentPool.Id).AnyTimes().Return(nil, fmt.Errorf("error"))
			},
		)
	})
}
