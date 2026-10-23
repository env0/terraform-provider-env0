package env0

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"sync"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// logForwardingStore stands in for the API's list of configurations, so a read after an upsert or delete sees it.
type logForwardingStore struct {
	mu             sync.Mutex
	configurations []client.LogForwarding
}

func (s *logForwardingStore) list() ([]client.LogForwarding, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return slices.Clone(s.configurations), nil
}

func (s *logForwardingStore) set(configurations ...client.LogForwarding) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.configurations = configurations
}

func TestUnitLogForwardingResource(t *testing.T) {
	t.Parallel()

	resourceType := "env0_log_forwarding"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	config := func(logForwardingType string, value string, deployment bool, audit bool) string {
		return fmt.Sprintf(`resource "%s" "%s" {
	type                      = "%s"
	value                     = %s
	deployment_log_forwarding = %t
	audit_log_forwarding      = %t
}`, resourceType, resourceName, logForwardingType, value, deployment, audit)
	}

	datadogValue := `jsonencode({ host = "http-intake.logs.datadoghq.com", apiKey = "$${ssm:datadog-api-key}" })`
	datadogValueJson := `{"apiKey":"${ssm:datadog-api-key}","host":"http-intake.logs.datadoghq.com"}`
	datadogPayloadValue := map[string]any{"host": "http-intake.logs.datadoghq.com", "apiKey": "${ssm:datadog-api-key}"}

	updatedDatadogValue := `jsonencode({ host = "http-intake.logs.datadoghq.eu", apiKey = "$${ssm:datadog-api-key}", port = 443 })`
	updatedDatadogPayloadValue := map[string]any{"host": "http-intake.logs.datadoghq.eu", "apiKey": "${ssm:datadog-api-key}", "port": float64(443)}

	datadog := client.LogForwarding{
		Id:                      uuid.NewString(),
		OrganizationId:          "org0",
		Name:                    "DATADOG",
		Type:                    client.LogForwardingTypeDatadog,
		DeploymentLogForwarding: true,
	}

	updatedDatadog := datadog
	updatedDatadog.AuditLogForwarding = true

	createDatadogPayload := client.LogForwardingUpsertPayload{
		Type:                    client.LogForwardingTypeDatadog,
		Value:                   datadogPayloadValue,
		DeploymentLogForwarding: true,
	}

	t.Run("Success", func(t *testing.T) {
		store := &logForwardingStore{}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: config("DATADOG", datadogValue, true, false),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", datadog.Id),
						resource.TestCheckResourceAttr(accessor, "type", "DATADOG"),
						resource.TestCheckResourceAttr(accessor, "value", datadogValueJson),
						resource.TestCheckResourceAttr(accessor, "deployment_log_forwarding", "true"),
						resource.TestCheckResourceAttr(accessor, "audit_log_forwarding", "false"),
					),
				},
				{
					Config: config("DATADOG", updatedDatadogValue, true, true),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", datadog.Id),
						resource.TestCheckResourceAttr(accessor, "deployment_log_forwarding", "true"),
						resource.TestCheckResourceAttr(accessor, "audit_log_forwarding", "true"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().LogForwardings().AnyTimes().DoAndReturn(store.list)

			mock.EXPECT().LogForwardingUpsert(createDatadogPayload).Times(1).DoAndReturn(func(client.LogForwardingUpsertPayload) (*client.LogForwarding, error) {
				store.set(datadog)

				return &datadog, nil
			})

			mock.EXPECT().LogForwardingUpsert(client.LogForwardingUpsertPayload{
				Type:                    client.LogForwardingTypeDatadog,
				Value:                   updatedDatadogPayloadValue,
				DeploymentLogForwarding: true,
				AuditLogForwarding:      true,
			}).Times(1).DoAndReturn(func(client.LogForwardingUpsertPayload) (*client.LogForwarding, error) {
				store.set(updatedDatadog)

				return &updatedDatadog, nil
			})

			mock.EXPECT().LogForwardingDelete(datadog.Id).Times(1)
		})
	})

	t.Run("Reformatted JSON value is not a change", func(t *testing.T) {
		store := &logForwardingStore{}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: config("DATADOG", datadogValue, true, false),
				},
				{
					Config: config("DATADOG", `<<EOT
{
  "host": "http-intake.logs.datadoghq.com",
  "apiKey": "$${ssm:datadog-api-key}"
}
EOT`, true, false),
					PlanOnly: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().LogForwardings().AnyTimes().DoAndReturn(store.list)
			mock.EXPECT().LogForwardingUpsert(createDatadogPayload).Times(1).DoAndReturn(func(client.LogForwardingUpsertPayload) (*client.LogForwarding, error) {
				store.set(datadog)

				return &datadog, nil
			})
			mock.EXPECT().LogForwardingDelete(datadog.Id).Times(1)
		})
	})

	t.Run("Changing the type replaces the configuration", func(t *testing.T) {
		store := &logForwardingStore{}

		splunk := client.LogForwarding{
			Id:                      uuid.NewString(),
			OrganizationId:          "org0",
			Name:                    "SPLUNK",
			Type:                    client.LogForwardingTypeSplunk,
			DeploymentLogForwarding: true,
		}

		splunkValue := `jsonencode({ url = "https://splunk.example.com", token = "$${ssm:splunk-token}", index = "env0" })`

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: config("DATADOG", datadogValue, true, false),
					Check:  resource.TestCheckResourceAttr(accessor, "id", datadog.Id),
				},
				{
					Config: config("SPLUNK", splunkValue, true, false),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", splunk.Id),
						resource.TestCheckResourceAttr(accessor, "type", "SPLUNK"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().LogForwardings().AnyTimes().DoAndReturn(store.list)
			mock.EXPECT().LogForwardingUpsert(createDatadogPayload).Times(1).DoAndReturn(func(client.LogForwardingUpsertPayload) (*client.LogForwarding, error) {
				store.set(datadog)

				return &datadog, nil
			})
			mock.EXPECT().LogForwardingDelete(datadog.Id).Times(1).DoAndReturn(func(string) error {
				store.set()

				return nil
			})
			mock.EXPECT().LogForwardingUpsert(client.LogForwardingUpsertPayload{
				Type:                    client.LogForwardingTypeSplunk,
				Value:                   map[string]any{"url": "https://splunk.example.com", "token": "${ssm:splunk-token}", "index": "env0"},
				DeploymentLogForwarding: true,
			}).Times(1).DoAndReturn(func(client.LogForwardingUpsertPayload) (*client.LogForwarding, error) {
				store.set(splunk)

				return &splunk, nil
			})
			mock.EXPECT().LogForwardingDelete(splunk.Id).Times(1)
		})
	})

	t.Run("Drift - deleted outside Terraform is created again", func(t *testing.T) {
		store := &logForwardingStore{}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: config("DATADOG", datadogValue, true, false),
				},
				{
					PreConfig: func() { store.set() },
					Config:    config("DATADOG", datadogValue, true, false),
					Check:     resource.TestCheckResourceAttr(accessor, "id", datadog.Id),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().LogForwardings().AnyTimes().DoAndReturn(store.list)
			mock.EXPECT().LogForwardingUpsert(createDatadogPayload).Times(2).DoAndReturn(func(client.LogForwardingUpsertPayload) (*client.LogForwarding, error) {
				store.set(datadog)

				return &datadog, nil
			})
			mock.EXPECT().LogForwardingDelete(datadog.Id).Times(1)
		})
	})

	t.Run("Create Failure - type already configured", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      config("DATADOG", datadogValue, true, false),
					ExpectError: regexp.MustCompile("a log forwarding configuration of type DATADOG already exists"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().LogForwardings().Times(1).Return([]client.LogForwarding{datadog}, nil)
		})
	})

	t.Run("Create Failure - list error", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      config("DATADOG", datadogValue, true, false),
					ExpectError: regexp.MustCompile("could not get log forwarding configurations: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().LogForwardings().Times(1).Return(nil, errors.New("error"))
		})
	})

	t.Run("Create Failure - upsert error", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      config("DATADOG", datadogValue, true, false),
					ExpectError: regexp.MustCompile("could not create log forwarding configuration: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().LogForwardings().Times(1).Return(nil, nil)
			mock.EXPECT().LogForwardingUpsert(createDatadogPayload).Times(1).Return(nil, errors.New("error"))
		})
	})

	t.Run("Update Failure", func(t *testing.T) {
		store := &logForwardingStore{}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: config("DATADOG", datadogValue, true, false),
				},
				{
					Config:      config("DATADOG", datadogValue, true, true),
					ExpectError: regexp.MustCompile("could not update log forwarding configuration: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().LogForwardings().AnyTimes().DoAndReturn(store.list)
			mock.EXPECT().LogForwardingUpsert(createDatadogPayload).Times(1).DoAndReturn(func(client.LogForwardingUpsertPayload) (*client.LogForwarding, error) {
				store.set(datadog)

				return &datadog, nil
			})
			mock.EXPECT().LogForwardingUpsert(client.LogForwardingUpsertPayload{
				Type:                    client.LogForwardingTypeDatadog,
				Value:                   datadogPayloadValue,
				DeploymentLogForwarding: true,
				AuditLogForwarding:      true,
			}).Times(1).Return(nil, errors.New("error"))
			mock.EXPECT().LogForwardingDelete(datadog.Id).Times(1)
		})
	})

	t.Run("Read Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      config("DATADOG", datadogValue, true, false),
					ExpectError: regexp.MustCompile("could not get log forwarding configuration: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().LogForwardings().Times(1).Return(nil, nil)
			mock.EXPECT().LogForwardingUpsert(createDatadogPayload).Times(1).Return(&datadog, nil)
			mock.EXPECT().LogForwardings().Times(1).Return(nil, errors.New("error"))
			mock.EXPECT().LogForwardingDelete(datadog.Id).Times(1)
		})
	})

	t.Run("Delete of a configuration that is already gone succeeds", func(t *testing.T) {
		store := &logForwardingStore{}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: config("DATADOG", datadogValue, true, false),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().LogForwardings().AnyTimes().DoAndReturn(store.list)
			mock.EXPECT().LogForwardingUpsert(createDatadogPayload).Times(1).DoAndReturn(func(client.LogForwardingUpsertPayload) (*client.LogForwarding, error) {
				store.set(datadog)

				return &datadog, nil
			})
			mock.EXPECT().LogForwardingDelete(datadog.Id).Times(1).Return(http.NewMockFailedResponseError(404))
		})
	})

	for name, testCase := range map[string]struct {
		logForwardingType string
		value             string
		deployment        bool
		audit             bool
		expectedError     string
	}{
		"no log type enabled": {"DATADOG", datadogValue, false, false, "at least one of deployment_log_forwarding and audit_log_forwarding must be true"},
		"invalid type":        {"DYNATRACE", datadogValue, true, false, "'DYNATRACE' must be one of: SPLUNK, NEWRELIC"},
		"value not JSON":      {"DATADOG", `"not json"`, true, false, "value must be a JSON object"},
		"value a JSON array":  {"DATADOG", `jsonencode(["a"])`, true, false, "value must be a JSON object"},
		"value JSON null":     {"DATADOG", `"null"`, true, false, "value must be a JSON object"},
	} {
		t.Run("Validation Failure - "+name, func(t *testing.T) {
			runUnitTest(t, resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config:      config(testCase.logForwardingType, testCase.value, testCase.deployment, testCase.audit),
						ExpectError: regexp.MustCompile(regexp.QuoteMeta(testCase.expectedError)),
					},
				},
			}, func(mock *client.MockApiClientInterface) {})
		})
	}

	for name, importStateId := range map[string]string{
		"Import By Id":                      datadog.Id,
		"Import By Type":                    "DATADOG",
		"Import By Type - case insensitive": "datadog",
	} {
		t.Run(name, func(t *testing.T) {
			store := &logForwardingStore{}

			testCase := resource.TestCase{
				Steps: []resource.TestStep{
					{
						Config: config("DATADOG", datadogValue, true, false),
					},
					{
						ResourceName:      resourceNameImport,
						ImportState:       true,
						ImportStateId:     importStateId,
						ImportStateVerify: true,
						// The API never returns the value.
						ImportStateVerifyIgnore: []string{"value"},
						ImportStateCheck: func(states []*terraform.InstanceState) error {
							if len(states) != 1 || states[0].ID != datadog.Id {
								return fmt.Errorf("expected one imported state with id %s, got %v", datadog.Id, states)
							}

							return nil
						},
					},
				},
			}

			runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
				mock.EXPECT().LogForwardings().AnyTimes().DoAndReturn(store.list)
				mock.EXPECT().LogForwardingUpsert(createDatadogPayload).Times(1).DoAndReturn(func(client.LogForwardingUpsertPayload) (*client.LogForwarding, error) {
					store.set(datadog)

					return &datadog, nil
				})
				mock.EXPECT().LogForwardingDelete(datadog.Id).Times(1)
			})
		})
	}

	t.Run("Import Failure - not found", func(t *testing.T) {
		store := &logForwardingStore{}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: config("DATADOG", datadogValue, true, false),
				},
				{
					ResourceName:  resourceNameImport,
					ImportState:   true,
					ImportStateId: "SPLUNK",
					ExpectError:   regexp.MustCompile("log forwarding configuration SPLUNK not found"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().LogForwardings().AnyTimes().DoAndReturn(store.list)
			mock.EXPECT().LogForwardingUpsert(createDatadogPayload).Times(1).DoAndReturn(func(client.LogForwardingUpsertPayload) (*client.LogForwarding, error) {
				store.set(datadog)

				return &datadog, nil
			})
			mock.EXPECT().LogForwardingDelete(datadog.Id).Times(1)
		})
	})
}
