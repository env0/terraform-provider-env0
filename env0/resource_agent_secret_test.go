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

func TestUnitAgentSecretResource(t *testing.T) {
	resourceType := "env0_agent_secret"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	agentId := "agent-id-0"

	agentSecret := client.AgentSecret{
		Id:          "secret-id-0",
		Secret:      "secret-value-0",
		AgentId:     agentId,
		Description: "my secret",
	}

	t.Run("Create and Read", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"agent_id":    agentId,
						"description": agentSecret.Description,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", agentSecret.Id),
						resource.TestCheckResourceAttr(accessor, "agent_id", agentId),
						resource.TestCheckResourceAttr(accessor, "description", agentSecret.Description),
						resource.TestCheckResourceAttr(accessor, "secret", agentSecret.Secret),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AgentSecretCreate(agentId, client.AgentSecretCreatePayload{
					Description: agentSecret.Description,
				}).Times(1).Return(&agentSecret, nil),
				mock.EXPECT().AgentSecrets(agentId).Times(1).Return([]client.AgentSecret{agentSecret}, nil),
				mock.EXPECT().AgentSecretDelete(agentId, agentSecret.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Create Without Description", func(t *testing.T) {
		secretNoDesc := client.AgentSecret{
			Id:      "secret-no-desc",
			Secret:  "secret-val",
			AgentId: agentId,
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"agent_id": agentId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", secretNoDesc.Id),
						resource.TestCheckResourceAttr(accessor, "agent_id", agentId),
						resource.TestCheckResourceAttr(accessor, "secret", secretNoDesc.Secret),
						resource.TestCheckNoResourceAttr(accessor, "description"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AgentSecretCreate(agentId, client.AgentSecretCreatePayload{}).Times(1).Return(&secretNoDesc, nil),
				mock.EXPECT().AgentSecrets(agentId).Times(1).Return([]client.AgentSecret{secretNoDesc}, nil),
				mock.EXPECT().AgentSecretDelete(agentId, secretNoDesc.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Create Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"agent_id": agentId,
					}),
					ExpectError: regexp.MustCompile("could not create agent secret: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AgentSecretCreate(agentId, gomock.Any()).Times(1).Return(nil, errors.New("error"))
		})
	})

	t.Run("Drift Detection", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"agent_id":    agentId,
						"description": agentSecret.Description,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", agentSecret.Id),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"agent_id":    agentId,
						"description": agentSecret.Description,
					}),
					ExpectNonEmptyPlan: true,
					PlanOnly:           true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AgentSecretCreate(agentId, gomock.Any()).Times(1).Return(&agentSecret, nil),
				// Read after create
				mock.EXPECT().AgentSecrets(agentId).Times(1).Return([]client.AgentSecret{agentSecret}, nil),
				// Step 2 refresh: secret was revoked externally, drift detected
				mock.EXPECT().AgentSecrets(agentId).Times(1).Return([]client.AgentSecret{}, nil),
			)
		})
	})

	t.Run("Import", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"agent_id":    agentId,
						"description": agentSecret.Description,
					}),
				},
				{
					ResourceName:            resourceNameImport,
					ImportState:             true,
					ImportStateId:           fmt.Sprintf("%s/%s", agentId, agentSecret.Id),
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"secret"},
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().AgentSecretCreate(agentId, gomock.Any()).Times(1).Return(&agentSecret, nil),
				// Read after create (1) + pre-import refresh (1) + post-import read (1)
				mock.EXPECT().AgentSecrets(agentId).Times(3).Return([]client.AgentSecret{agentSecret}, nil),
				mock.EXPECT().AgentSecretDelete(agentId, agentSecret.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Import Invalid Format", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"agent_id":    agentId,
						"description": agentSecret.Description,
					}),
				},
				{
					ResourceName:  resourceNameImport,
					ImportState:   true,
					ImportStateId: "invalid-no-slash",
					ExpectError:   regexp.MustCompile("invalid import format"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().AgentSecretCreate(agentId, gomock.Any()).Times(1).Return(&agentSecret, nil)
			mock.EXPECT().AgentSecrets(agentId).AnyTimes().Return([]client.AgentSecret{agentSecret}, nil)
			mock.EXPECT().AgentSecretDelete(agentId, agentSecret.Id).Times(1).Return(nil)
		})
	})
}
