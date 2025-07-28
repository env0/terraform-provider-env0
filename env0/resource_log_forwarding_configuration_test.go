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

func TestUnitLogForwardingConfigurationResource(t *testing.T) {
	resourceType := "env0_log_forwarding_configuration"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	logConfig := client.LogForwardingConfiguration{
		Id:                      "config-id-123",
		Type:                    "SPLUNK",
		Value:                   map[string]any{"url": "https://example.com/webhook", "token": "secret-token"},
		AuditLogForwarding:      boolPtr(true),
		DeploymentLogForwarding: boolPtr(false),
	}

	valueJson := `{"token":"secret-token","url":"https://example.com/webhook"}`

	stepConfig := resourceConfigCreate(resourceType, resourceName, map[string]any{
		"type":                      logConfig.Type,
		"value":                     valueJson,
		"audit_log_forwarding":      *logConfig.AuditLogForwarding,
		"deployment_log_forwarding": *logConfig.DeploymentLogForwarding,
	})

	createPayload := client.LogForwardingConfigurationCreatePayload{
		Type:                    logConfig.Type,
		Value:                   logConfig.Value,
		AuditLogForwarding:      logConfig.AuditLogForwarding,
		DeploymentLogForwarding: logConfig.DeploymentLogForwarding,
	}

	t.Run("Create", func(t *testing.T) {
		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", logConfig.Id),
						resource.TestCheckResourceAttr(accessor, "type", logConfig.Type),
						resource.TestCheckResourceAttr(accessor, "value", valueJson),
						resource.TestCheckResourceAttr(accessor, "audit_log_forwarding", strconv.FormatBool(*logConfig.AuditLogForwarding)),
						resource.TestCheckResourceAttr(accessor, "deployment_log_forwarding", strconv.FormatBool(*logConfig.DeploymentLogForwarding)),
					),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().LogForwardingConfigurationCreate(createPayload).Times(1).Return(logConfig, nil)
			mock.EXPECT().LogForwardingConfiguration(logConfig.Id).Times(1).Return(logConfig, nil)
			mock.EXPECT().LogForwardingConfigurationDelete(logConfig.Id).Times(1).Return(nil)
		})
	})

	t.Run("Create with API error", func(t *testing.T) {
		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      stepConfig,
					ExpectError: regexp.MustCompile("could not create log forwarding configuration: error"),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().LogForwardingConfigurationCreate(createPayload).Times(1).Return(client.LogForwardingConfiguration{}, errors.New("error"))
		})
	})

	t.Run("Update", func(t *testing.T) {
		updatedConfig := logConfig
		updatedConfig.Type = "DATADOG"
		updatedValueJson := `{"api_key":"datadog-key","host":"datadog.example.com"}`
		updatedValue := map[string]any{"api_key": "datadog-key", "host": "datadog.example.com"}
		updatedConfig.Value = updatedValue

		updatedStepConfig := resourceConfigCreate(resourceType, resourceName, map[string]any{
			"type":                      updatedConfig.Type,
			"value":                     updatedValueJson,
			"audit_log_forwarding":      *updatedConfig.AuditLogForwarding,
			"deployment_log_forwarding": *updatedConfig.DeploymentLogForwarding,
		})

		updatePayload := client.LogForwardingConfigurationUpdatePayload{
			Id:                      updatedConfig.Id,
			Type:                    updatedConfig.Type,
			Value:                   updatedValue,
			AuditLogForwarding:      updatedConfig.AuditLogForwarding,
			DeploymentLogForwarding: updatedConfig.DeploymentLogForwarding,
		}

		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "type", logConfig.Type),
					),
				},
				{
					Config: updatedStepConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "type", updatedConfig.Type),
						resource.TestCheckResourceAttr(accessor, "value", updatedValueJson),
					),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().LogForwardingConfigurationCreate(createPayload).Times(1).Return(logConfig, nil),
				mock.EXPECT().LogForwardingConfiguration(logConfig.Id).Times(2).Return(logConfig, nil),
				mock.EXPECT().LogForwardingConfigurationUpdate(updatePayload).Times(1).Return(updatedConfig, nil),
				mock.EXPECT().LogForwardingConfiguration(logConfig.Id).Times(1).Return(updatedConfig, nil),
				mock.EXPECT().LogForwardingConfigurationDelete(logConfig.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Read with API error", func(t *testing.T) {
		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      stepConfig,
					ExpectError: regexp.MustCompile("could not get log forwarding configuration: error"),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().LogForwardingConfigurationCreate(createPayload).Times(1).Return(logConfig, nil)
			mock.EXPECT().LogForwardingConfiguration(logConfig.Id).Times(1).Return(client.LogForwardingConfiguration{}, errors.New("error"))
			mock.EXPECT().LogForwardingConfigurationDelete(logConfig.Id).Times(1).Return(nil)
		})
	})

	t.Run("Invalid JSON value", func(t *testing.T) {
		invalidStepConfig := resourceConfigCreate(resourceType, resourceName, map[string]any{
			"type":                      logConfig.Type,
			"value":                     "invalid-json",
			"audit_log_forwarding":      *logConfig.AuditLogForwarding,
			"deployment_log_forwarding": *logConfig.DeploymentLogForwarding,
		})

		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      invalidStepConfig,
					ExpectError: regexp.MustCompile("invalid JSON in value field"),
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {})
	})
}
