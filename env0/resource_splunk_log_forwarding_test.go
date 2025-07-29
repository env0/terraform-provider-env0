package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitSplunkLogForwardingResource(t *testing.T) {
	resourceType := "env0_splunk_log_forwarding"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	logForwardingConfig := client.LogForwardingConfiguration{
		Id:                      "id0",
		OrganizationId:          "org0",
		Type:                    client.LogForwardingConfigurationTypeSplunk,
		Value:                   map[string]interface{}{
			"url":   "https://splunk.example.com",
			"token": "splunk-token-123",
			"index": "main",
		},
		AuditLogForwarding:      true,
		DeploymentLogForwarding: true,
	}

	updatedLogForwardingConfig := client.LogForwardingConfiguration{
		Id:                      logForwardingConfig.Id,
		OrganizationId:          logForwardingConfig.OrganizationId,
		Type:                    client.LogForwardingConfigurationTypeSplunk,
		Value:                   map[string]interface{}{
			"url":   "https://updated-splunk.example.com",
			"token": "updated-splunk-token-456",
			"index": "updated-index",
		},
		AuditLogForwarding:      false,
		DeploymentLogForwarding: false,
	}

	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"url":                       "https://splunk.example.com",
						"token":                     "splunk-token-123",
						"index":                     "main",
						"audit_log_forwarding":      true,
						"deployment_log_forwarding": true,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", logForwardingConfig.Id),
						resource.TestCheckResourceAttr(accessor, "url", "https://splunk.example.com"),
						resource.TestCheckResourceAttr(accessor, "token", "splunk-token-123"),
						resource.TestCheckResourceAttr(accessor, "index", "main"),
						resource.TestCheckResourceAttr(accessor, "audit_log_forwarding", "true"),
						resource.TestCheckResourceAttr(accessor, "deployment_log_forwarding", "true"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"url":                       "https://updated-splunk.example.com",
						"token":                     "updated-splunk-token-456",
						"index":                     "updated-index",
						"audit_log_forwarding":      false,
						"deployment_log_forwarding": false,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedLogForwardingConfig.Id),
						resource.TestCheckResourceAttr(accessor, "url", "https://updated-splunk.example.com"),
						resource.TestCheckResourceAttr(accessor, "index", "updated-index"),
						resource.TestCheckResourceAttr(accessor, "audit_log_forwarding", "false"),
						resource.TestCheckResourceAttr(accessor, "deployment_log_forwarding", "false"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().LogForwardingConfigurationCreate(&client.LogForwardingConfigurationCreatePayload{
				Type: client.LogForwardingConfigurationTypeSplunk,
				Value: map[string]interface{}{
					"url":   "https://splunk.example.com",
					"token": "splunk-token-123",
					"index": "main",
				},
				AuditLogForwarding:      true,
				DeploymentLogForwarding: true,
			}).Times(1).Return(&logForwardingConfig, nil)

			mock.EXPECT().LogForwardingConfigurationUpdate(updatedLogForwardingConfig.Id, &client.LogForwardingConfigurationUpdatePayload{
				Value: map[string]interface{}{
					"url":   "https://updated-splunk.example.com",
					"token": "updated-splunk-token-456",
					"index": "updated-index",
				},
				AuditLogForwarding:      false,
				DeploymentLogForwarding: false,
			}).Times(1).Return(&updatedLogForwardingConfig, nil)

			gomock.InOrder(
				mock.EXPECT().LogForwardingConfiguration(gomock.Any()).Times(2).Return(&logForwardingConfig, nil),        // 2 after create (shared + specific)
				mock.EXPECT().LogForwardingConfiguration(gomock.Any()).Times(2).Return(&updatedLogForwardingConfig, nil), // 2 after update (shared + specific)
			)

			mock.EXPECT().LogForwardingConfigurationDelete(logForwardingConfig.Id).Times(1)
		})
	})

	t.Run("Create Failure - Missing URL", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"token": "splunk-token-123",
						"index": "main",
					}),
					ExpectError: regexp.MustCompile(`The argument "url" is required`),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("Create Failure - Missing Token", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"url":   "https://splunk.example.com",
						"index": "main",
					}),
					ExpectError: regexp.MustCompile(`The argument "token" is required`),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("Create Failure - Missing Index", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"url":   "https://splunk.example.com",
						"token": "splunk-token-123",
					}),
					ExpectError: regexp.MustCompile(`The argument "index" is required`),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("Create Failure - API Error", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"url":   "https://splunk.example.com",
						"token": "splunk-token-123",
						"index": "main",
					}),
					ExpectError: regexp.MustCompile("failed to create log forwarding configuration: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().LogForwardingConfigurationCreate(&client.LogForwardingConfigurationCreatePayload{
				Type: client.LogForwardingConfigurationTypeSplunk,
				Value: map[string]interface{}{
					"url":   "https://splunk.example.com",
					"token": "splunk-token-123",
					"index": "main",
				},
				AuditLogForwarding:      true,
				DeploymentLogForwarding: true,
			}).Times(1).Return(nil, errors.New("error"))
		})
	})

	t.Run("Update Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"url":   "https://splunk.example.com",
						"token": "splunk-token-123",
						"index": "main",
					}),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"url":   "https://updated-splunk.example.com",
						"token": "updated-splunk-token-456",
						"index": "updated-index",
					}),
					ExpectError: regexp.MustCompile("failed to update log forwarding configuration: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().LogForwardingConfigurationCreate(&client.LogForwardingConfigurationCreatePayload{
				Type: client.LogForwardingConfigurationTypeSplunk,
				Value: map[string]interface{}{
					"url":   "https://splunk.example.com",
					"token": "splunk-token-123",
					"index": "main",
				},
				AuditLogForwarding:      true,
				DeploymentLogForwarding: true,
			}).Times(1).Return(&logForwardingConfig, nil)

			mock.EXPECT().LogForwardingConfigurationUpdate(updatedLogForwardingConfig.Id, &client.LogForwardingConfigurationUpdatePayload{
				Value: map[string]interface{}{
					"url":   "https://updated-splunk.example.com",
					"token": "updated-splunk-token-456",
					"index": "updated-index",
				},
				AuditLogForwarding:      true,
				DeploymentLogForwarding: true,
			}).Times(1).Return(nil, errors.New("error"))

			mock.EXPECT().LogForwardingConfiguration(gomock.Any()).Times(2).Return(&logForwardingConfig, nil) // 1 after create, 1 before update
			mock.EXPECT().LogForwardingConfigurationDelete(logForwardingConfig.Id).Times(1)
		})
	})

	t.Run("Read Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"url":   "https://splunk.example.com",
						"token": "splunk-token-123",
						"index": "main",
					}),
					ExpectError: regexp.MustCompile("error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().LogForwardingConfigurationCreate(&client.LogForwardingConfigurationCreatePayload{
				Type: client.LogForwardingConfigurationTypeSplunk,
				Value: map[string]interface{}{
					"url":   "https://splunk.example.com",
					"token": "splunk-token-123",
					"index": "main",
				},
				AuditLogForwarding:      true,
				DeploymentLogForwarding: true,
			}).Times(1).Return(&logForwardingConfig, nil)

			mock.EXPECT().LogForwardingConfiguration(gomock.Any()).Return(nil, errors.New("error"))
			mock.EXPECT().LogForwardingConfigurationDelete(logForwardingConfig.Id).Times(1)
		})
	})

	t.Run("Resource removed in UI", func(t *testing.T) {
		stepConfig := resourceConfigCreate(resourceType, resourceName, map[string]any{
			"url":   "https://splunk.example.com",
			"token": "splunk-token-123",
			"index": "main",
		})

		createPayload := &client.LogForwardingConfigurationCreatePayload{
			Type: client.LogForwardingConfigurationTypeSplunk,
			Value: map[string]interface{}{
				"url":   "https://splunk.example.com",
				"token": "splunk-token-123",
				"index": "main",
			},
			AuditLogForwarding:      true,
			DeploymentLogForwarding: true,
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfig,
				},
				{
					Config: stepConfig,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().LogForwardingConfigurationCreate(createPayload).Times(1).Return(&logForwardingConfig, nil),
				mock.EXPECT().LogForwardingConfiguration(logForwardingConfig.Id).Times(1).Return(&logForwardingConfig, nil),
				mock.EXPECT().LogForwardingConfiguration(logForwardingConfig.Id).Times(1).Return(nil, http.NewMockFailedResponseError(404)),
				mock.EXPECT().LogForwardingConfigurationCreate(createPayload).Times(1).Return(&logForwardingConfig, nil),
				mock.EXPECT().LogForwardingConfiguration(logForwardingConfig.Id).Times(1).Return(&logForwardingConfig, nil),
				mock.EXPECT().LogForwardingConfigurationDelete(logForwardingConfig.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Import By Id", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"url":   "https://splunk.example.com",
						"token": "splunk-token-123",
						"index": "main",
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     logForwardingConfig.Id,
					ImportStateVerify: true,
					ImportStateVerifyIgnore: []string{"token"}, // token is sensitive and not returned by API
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().LogForwardingConfigurationCreate(&client.LogForwardingConfigurationCreatePayload{
				Type: client.LogForwardingConfigurationTypeSplunk,
				Value: map[string]interface{}{
					"url":   "https://splunk.example.com",
					"token": "splunk-token-123",
					"index": "main",
				},
				AuditLogForwarding:      true,
				DeploymentLogForwarding: true,
			}).Times(1).Return(&logForwardingConfig, nil)

			mock.EXPECT().LogForwardingConfiguration(logForwardingConfig.Id).Times(2).Return(&logForwardingConfig, nil) // 1 after create, 1 for import
			mock.EXPECT().LogForwardingConfigurationDelete(logForwardingConfig.Id).Times(1)
		})
	})
}
