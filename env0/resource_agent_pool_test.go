package env0

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitAgentPoolResource(t *testing.T) {
	resourceType := "env0_agent_pool"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	agentPool := client.AgentPool{
		Id:          "id0",
		Name:        "my-pool",
		Description: "my pool description",
		AgentKey:    "agent-key-0",
	}

	updatedAgentPool := client.AgentPool{
		Id:          "id0",
		Name:        "updated-pool",
		Description: "updated description",
		AgentKey:    "agent-key-0",
	}

	t.Run("Create and Update", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":        agentPool.Name,
						"description": agentPool.Description,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", agentPool.Id),
						resource.TestCheckResourceAttr(accessor, "name", agentPool.Name),
						resource.TestCheckResourceAttr(accessor, "description", agentPool.Description),
						resource.TestCheckResourceAttr(accessor, "agent_key", agentPool.AgentKey),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":        updatedAgentPool.Name,
						"description": updatedAgentPool.Description,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedAgentPool.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedAgentPool.Name),
						resource.TestCheckResourceAttr(accessor, "description", updatedAgentPool.Description),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AgentPoolCreate(client.AgentPoolCreatePayload{
					Name:        agentPool.Name,
					Description: agentPool.Description,
				}).Times(1).Return(&agentPool, nil),
				// Read after create + refresh reads before update
				mock.EXPECT().AgentPool(agentPool.Id).Times(3).Return(&agentPool, nil),
				mock.EXPECT().AgentPoolUpdate(agentPool.Id, client.AgentPoolUpdatePayload{
					Name:        updatedAgentPool.Name,
					Description: updatedAgentPool.Description,
				}).Times(1).Return(&updatedAgentPool, nil),
				// Read after update + refresh before destroy
				mock.EXPECT().AgentPool(updatedAgentPool.Id).Times(2).Return(&updatedAgentPool, nil),
				mock.EXPECT().AgentPoolDelete(agentPool.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Create Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name": agentPool.Name,
					}),
					ExpectError: regexp.MustCompile("could not create agent pool: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AgentPoolCreate(gomock.Any()).Times(1).Return(nil, errors.New("error"))
		})
	})

	t.Run("Drift Detection", func(t *testing.T) {
		driftPool := client.AgentPool{
			Id:       "id-drift",
			Name:     "drift-pool",
			AgentKey: "key-drift",
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name": driftPool.Name,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", driftPool.Id),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name": driftPool.Name,
					}),
					ExpectNonEmptyPlan: true,
					PlanOnly:           true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AgentPoolCreate(gomock.Any()).Times(1).Return(&driftPool, nil),
				// Read after create + post-apply refresh
				mock.EXPECT().AgentPool(driftPool.Id).Times(2).Return(&driftPool, nil),
				// Step 2 refresh: resource gone → drift detected
				mock.EXPECT().AgentPool(driftPool.Id).Times(1).Return(nil, &client.NotFoundError{}),
			)
		})
	})

	t.Run("Create With Logs", func(t *testing.T) {
		logs := &client.AgentPoolLogsConfig{
			Dynamo: &client.AgentPoolDynamoLogs{
				SelfHosted: &client.AgentPoolSelfHostedLogs{
					AccountId:  "123456789",
					Region:     "us-east-1",
					ExternalId: "ext-id",
				},
			},
		}

		poolWithLogs := client.AgentPool{
			Id:       "id-logs",
			Name:     "pool-with-logs",
			AgentKey: "key-logs",
			Logs:     logs,
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
						resource "%s" "%s" {
							name = "%s"

							logs {
								account_id  = "%s"
								region      = "%s"
								external_id = "%s"
							}
						}
					`, resourceType, resourceName, poolWithLogs.Name,
						logs.Dynamo.SelfHosted.AccountId,
						logs.Dynamo.SelfHosted.Region,
						logs.Dynamo.SelfHosted.ExternalId),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", poolWithLogs.Id),
						resource.TestCheckResourceAttr(accessor, "logs.0.account_id", logs.Dynamo.SelfHosted.AccountId),
						resource.TestCheckResourceAttr(accessor, "logs.0.region", logs.Dynamo.SelfHosted.Region),
						resource.TestCheckResourceAttr(accessor, "logs.0.external_id", logs.Dynamo.SelfHosted.ExternalId),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				// Create (POST) — no logs in create payload
				mock.EXPECT().AgentPoolCreate(client.AgentPoolCreatePayload{
					Name: poolWithLogs.Name,
				}).Times(1).Return(&client.AgentPool{
					Id:       poolWithLogs.Id,
					Name:     poolWithLogs.Name,
					AgentKey: poolWithLogs.AgentKey,
				}, nil),
				// PATCH to set logs after create
				mock.EXPECT().AgentPoolUpdate(poolWithLogs.Id, client.AgentPoolUpdatePayload{
					Logs: logs,
				}).Times(1).Return(&poolWithLogs, nil),
				// Read after create + pre-destroy refresh
				mock.EXPECT().AgentPool(poolWithLogs.Id).Times(2).Return(&poolWithLogs, nil),
				mock.EXPECT().AgentPoolDelete(poolWithLogs.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Update Add and Remove Logs", func(t *testing.T) {
		logs := &client.AgentPoolLogsConfig{
			Dynamo: &client.AgentPoolDynamoLogs{
				SelfHosted: &client.AgentPoolSelfHostedLogs{
					AccountId: "111111111",
					Region:    "eu-west-1",
				},
			},
		}

		poolNoLogs := client.AgentPool{
			Id:       "id-logs-update",
			Name:     "pool-logs-update",
			AgentKey: "key-logs-update",
		}

		poolWithLogs := client.AgentPool{
			Id:       poolNoLogs.Id,
			Name:     poolNoLogs.Name,
			AgentKey: poolNoLogs.AgentKey,
			Logs:     logs,
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					// Step 1: create without logs
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name": poolNoLogs.Name,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", poolNoLogs.Id),
						resource.TestCheckResourceAttr(accessor, "logs.#", "0"),
					),
				},
				{
					// Step 2: add logs
					Config: fmt.Sprintf(`
						resource "%s" "%s" {
							name = "%s"

							logs {
								account_id = "%s"
								region     = "%s"
							}
						}
					`, resourceType, resourceName, poolNoLogs.Name,
						logs.Dynamo.SelfHosted.AccountId,
						logs.Dynamo.SelfHosted.Region),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "logs.0.account_id", logs.Dynamo.SelfHosted.AccountId),
						resource.TestCheckResourceAttr(accessor, "logs.0.region", logs.Dynamo.SelfHosted.Region),
					),
				},
				{
					// Step 3: remove logs
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name": poolNoLogs.Name,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "logs.#", "0"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				// Step 1: create without logs
				mock.EXPECT().AgentPoolCreate(client.AgentPoolCreatePayload{
					Name: poolNoLogs.Name,
				}).Times(1).Return(&poolNoLogs, nil),
				// Read after create + pre-step-2 refreshes
				mock.EXPECT().AgentPool(poolNoLogs.Id).Times(3).Return(&poolNoLogs, nil),
				// Step 2: update to add logs
				mock.EXPECT().AgentPoolUpdate(poolNoLogs.Id, client.AgentPoolUpdatePayload{
					Name: poolNoLogs.Name,
					Logs: logs,
				}).Times(1).Return(&poolWithLogs, nil),
				// Read after update + pre-step-3 refreshes
				mock.EXPECT().AgentPool(poolNoLogs.Id).Times(3).Return(&poolWithLogs, nil),
				// Step 3: update to remove logs (nil logs sent explicitly)
				mock.EXPECT().AgentPoolUpdate(poolNoLogs.Id, client.AgentPoolUpdatePayload{
					Name: poolNoLogs.Name,
				}).Times(1).Return(&poolNoLogs, nil),
				// Read after update + pre-destroy refresh
				mock.EXPECT().AgentPool(poolNoLogs.Id).Times(2).Return(&poolNoLogs, nil),
				mock.EXPECT().AgentPoolDelete(poolNoLogs.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Create With Logs PATCH Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
						resource "%s" "%s" {
							name = "pool-logs-fail"

							logs {
								account_id = "123456789"
								region     = "us-east-1"
							}
						}
					`, resourceType, resourceName),
					ExpectError: regexp.MustCompile("failed to set logs configuration"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AgentPoolCreate(gomock.Any()).Times(1).Return(&client.AgentPool{
				Id:   "id-logs-fail",
				Name: "pool-logs-fail",
			}, nil)
			mock.EXPECT().AgentPoolUpdate(gomock.Any(), gomock.Any()).Times(1).Return(nil, errors.New("error"))
			// Resource is in state (SetId called before PATCH), so Terraform cleans up on failure
			mock.EXPECT().AgentPoolDelete("id-logs-fail").Times(1).Return(nil)
		})
	})

	t.Run("Import By Id", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":        agentPool.Name,
						"description": agentPool.Description,
					}),
				},
				{
					ResourceName:            resourceNameImport,
					ImportState:             true,
					ImportStateId:           agentPool.Id,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"agent_key"},
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AgentPoolCreate(gomock.Any()).Times(1).Return(&agentPool, nil),
				// Read after create (1) + pre-import refresh (1) + import handler (1) + post-import read (1)
				mock.EXPECT().AgentPool(agentPool.Id).Times(4).Return(&agentPool, nil),
				mock.EXPECT().AgentPoolDelete(agentPool.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Import By Name", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":        agentPool.Name,
						"description": agentPool.Description,
					}),
				},
				{
					ResourceName:            resourceNameImport,
					ImportState:             true,
					ImportStateId:           agentPool.Name,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"agent_key"},
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AgentPoolCreate(gomock.Any()).Times(1).Return(&agentPool, nil),
				// Read after create + pre-import refresh
				mock.EXPECT().AgentPool(agentPool.Id).Times(2).Return(&agentPool, nil),
				// Import: try by id first (fails with name string), then by name via list
				mock.EXPECT().AgentPool(agentPool.Name).Times(1).Return(nil, &client.NotFoundError{}),
				mock.EXPECT().AgentPools().Times(1).Return([]client.AgentPool{agentPool}, nil),
				// Post-import read
				mock.EXPECT().AgentPool(agentPool.Id).Times(1).Return(&agentPool, nil),
				mock.EXPECT().AgentPoolDelete(agentPool.Id).Times(1).Return(nil),
			)
		})
	})
}
