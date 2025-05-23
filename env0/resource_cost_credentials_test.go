package env0

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/env0/terraform-provider-env0/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitAwsCostCredentialsResource(t *testing.T) {
	resourceType := "env0_aws_cost_credentials"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	awsCredentialResource := map[string]any{
		"name":     "test",
		"arn":      "11111",
		"duration": 7200,
	}

	updatedAwsCredentialResource := map[string]any{
		"name":     "update",
		"arn":      "33333",
		"duration": 3600,
	}

	invalidDurationAwsCredentialResource := map[string]any{
		"name":     "update",
		"arn":      "33333",
		"duration": 1234,
	}

	awsCredCreatePayload := client.AwsCredentialsCreatePayload{
		Name: awsCredentialResource["name"].(string),
		Value: client.AwsCredentialsValuePayload{
			RoleArn:  awsCredentialResource["arn"].(string),
			Duration: awsCredentialResource["duration"].(int),
		},
		Type: client.AwsCostCredentialsType,
	}

	updateAwsCredCreatePayload := client.AwsCredentialsCreatePayload{
		Name: updatedAwsCredentialResource["name"].(string),
		Value: client.AwsCredentialsValuePayload{
			RoleArn:  updatedAwsCredentialResource["arn"].(string),
			Duration: updatedAwsCredentialResource["duration"].(int),
		},
		Type: client.AwsCostCredentialsType,
	}

	returnValues := client.Credentials{
		Id:             "id",
		Name:           "test",
		OrganizationId: "id",
		Type:           "AWS_ASSUMED_ROLE",
	}

	updateReturnValues := client.Credentials{
		Id:             "id",
		Name:           "update",
		OrganizationId: "id",
		Type:           "AWS_ASSUMED_ROLE",
	}

	testCaseForCreate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, awsCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", awsCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "arn", awsCredentialResource["arn"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", "id"),
					resource.TestCheckResourceAttr(accessor, "duration", strconv.Itoa(awsCredentialResource["duration"].(int))),
				),
			},
		},
	}

	testCaseForUpdate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, awsCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", awsCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "arn", awsCredentialResource["arn"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", returnValues.Id),
					resource.TestCheckResourceAttr(accessor, "duration", strconv.Itoa(awsCredentialResource["duration"].(int))),
				),
			},
			{
				Config: resourceConfigCreate(resourceType, resourceName, updatedAwsCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", updatedAwsCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "arn", updatedAwsCredentialResource["arn"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", updateReturnValues.Id),
					resource.TestCheckResourceAttr(accessor, "duration", strconv.Itoa(updatedAwsCredentialResource["duration"].(int))),
				),
			},
		},
	}

	missingValidInputErrorResource := map[string]any{
		"name": "update",
	}
	testCaseFormMissingValidInputError := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config:      resourceConfigCreate(resourceType, resourceName, missingValidInputErrorResource),
				ExpectError: regexp.MustCompile("Error: Missing required argument"),
			},
		},
	}

	t.Run("create", func(t *testing.T) {
		runUnitTest(t, testCaseForCreate, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CredentialsCreate(&awsCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("any update cause a destroy before a new create", func(t *testing.T) {
		runUnitTest(t, testCaseForUpdate, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CredentialsCreate(&awsCredCreatePayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CredentialsUpdate(returnValues.Id, &updateAwsCredCreatePayload).Times(1).Return(updateReturnValues, nil),
			)
			gomock.InOrder(
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(2).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(updateReturnValues.Id).Times(1).Return(updateReturnValues, nil),
			)
			mock.EXPECT().CloudCredentialsDelete(updateReturnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("throw error when don't enter any valid options", func(t *testing.T) {
		runUnitTest(t, testCaseFormMissingValidInputError, func(mock *client.MockApiClientInterface) {
		})
	})

	t.Run("throw error when don't enter duration valid values", func(t *testing.T) {
		runUnitTest(t, resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      resourceConfigCreate(resourceType, resourceName, invalidDurationAwsCredentialResource),
					ExpectError: regexp.MustCompile("Error: must be one of"),
				},
			},
		}, func(mock *client.MockApiClientInterface) {
		})
	})
}

func TestUnitAzureCostCredentialsResource(t *testing.T) {
	resourceType := "env0_azure_cost_credentials"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)
	azureCredentialResource := map[string]any{
		"name":            "test",
		"client_id":       "11111",
		"client_secret":   "client-secret1",
		"subscription_id": "subscription-id1",
		"tenant_id":       "tenant-id1",
	}

	updateAzureCredentialResource := map[string]any{
		"name":            "testUpdate",
		"client_id":       "22222",
		"client_secret":   "client-secret2",
		"subscription_id": "subscription-id2",
		"tenant_id":       "tenant-id2",
	}

	azureCredCreatePayload := client.AzureCredentialsCreatePayload{
		Name: azureCredentialResource["name"].(string),
		Value: client.AzureCredentialsValuePayload{
			ClientId:       azureCredentialResource["client_id"].(string),
			ClientSecret:   azureCredentialResource["client_secret"].(string),
			SubscriptionId: azureCredentialResource["subscription_id"].(string),
			TenantId:       azureCredentialResource["tenant_id"].(string),
		},
		Type: client.AzureCostCredentialsType,
	}

	updateAzureCredCreatePayload := client.AzureCredentialsCreatePayload{
		Name: updateAzureCredentialResource["name"].(string),
		Value: client.AzureCredentialsValuePayload{
			ClientId:       updateAzureCredentialResource["client_id"].(string),
			ClientSecret:   updateAzureCredentialResource["client_secret"].(string),
			SubscriptionId: updateAzureCredentialResource["subscription_id"].(string),
			TenantId:       updateAzureCredentialResource["tenant_id"].(string),
		},
		Type: client.AzureCostCredentialsType,
	}

	returnValues := client.Credentials{
		Id:             "id",
		Name:           "test",
		OrganizationId: "id",
		Type:           string(client.AzureCostCredentialsType),
	}

	updateReturnValues := client.Credentials{
		Id:             "id",
		Name:           "update",
		OrganizationId: "id",
		Type:           string(client.AzureCostCredentialsType),
	}

	testCaseForCreate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, azureCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", azureCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "client_id", azureCredentialResource["client_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "client_secret", azureCredentialResource["client_secret"].(string)),
					resource.TestCheckResourceAttr(accessor, "subscription_id", azureCredentialResource["subscription_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "tenant_id", azureCredentialResource["tenant_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", returnValues.Id),
				),
			},
		},
	}

	testCaseForUpdate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, azureCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", azureCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "client_id", azureCredentialResource["client_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "client_secret", azureCredentialResource["client_secret"].(string)),
					resource.TestCheckResourceAttr(accessor, "subscription_id", azureCredentialResource["subscription_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "tenant_id", azureCredentialResource["tenant_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", returnValues.Id),
				),
			},
			{
				Config: resourceConfigCreate(resourceType, resourceName, updateAzureCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", updateAzureCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "client_id", updateAzureCredentialResource["client_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "client_secret", updateAzureCredentialResource["client_secret"].(string)),
					resource.TestCheckResourceAttr(accessor, "subscription_id", updateAzureCredentialResource["subscription_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "tenant_id", updateAzureCredentialResource["tenant_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", updateReturnValues.Id),
				),
			},
		},
	}

	t.Run("create", func(t *testing.T) {
		runUnitTest(t, testCaseForCreate, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CredentialsCreate(&azureCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("any update cause a destroy before a new create", func(t *testing.T) {
		runUnitTest(t, testCaseForUpdate, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CredentialsCreate(&azureCredCreatePayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CredentialsUpdate(returnValues.Id, &updateAzureCredCreatePayload).Times(1).Return(updateReturnValues, nil),
			)
			gomock.InOrder(
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(2).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(updateReturnValues.Id).Times(1).Return(updateReturnValues, nil),
			)
			mock.EXPECT().CloudCredentialsDelete(updateReturnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("validate missing arguments", func(t *testing.T) {
		missingArgumentsTestCases := []resource.TestCase{
			missingArgumentTestCaseForCostCred(resourceType, resourceName, map[string]any{}),
			missingArgumentTestCaseForCostCred(resourceType, resourceName, map[string]any{}),
			missingArgumentTestCaseForCostCred(resourceType, resourceName, map[string]any{}),
			missingArgumentTestCaseForCostCred(resourceType, resourceName, map[string]any{}),
			missingArgumentTestCaseForCostCred(resourceType, resourceName, map[string]any{}),
		}
		for _, testCase := range missingArgumentsTestCases {
			tc := testCase

			t.Run("validate specific argument", func(t *testing.T) {
				runUnitTest(t, tc, func(mock *client.MockApiClientInterface) {})
			})
		}
	})
}

func TestUnitGoogleCostCredentialsResource(t *testing.T) {
	resourceType := "env0_gcp_cost_credentials"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	googleCostCredentialResource := map[string]any{
		"name":     "test",
		"table_id": "11111",
		"secret":   "22222",
	}

	updateGoogleCostCredentialResource := map[string]any{
		"name":     "testUpdate",
		"table_id": "333_update",
		"secret":   "444_update",
	}

	googleCostCredCreatePayload := client.GoogleCostCredentialsCreatePayload{
		Name: googleCostCredentialResource["name"].(string),
		Value: client.GoogleCostCredentialsValuePayload{
			TableId: googleCostCredentialResource["table_id"].(string),
			Secret:  googleCostCredentialResource["secret"].(string),
		},
		Type: client.GoogleCostCredentialsType,
	}

	updateGoogleCostCredCreatePayload := client.GoogleCostCredentialsCreatePayload{
		Name: updateGoogleCostCredentialResource["name"].(string),
		Value: client.GoogleCostCredentialsValuePayload{
			TableId: updateGoogleCostCredentialResource["table_id"].(string),
			Secret:  updateGoogleCostCredentialResource["secret"].(string),
		},
		Type: client.GoogleCostCredentialsType,
	}

	returnValues := client.Credentials{
		Id:             "id",
		Name:           "test",
		OrganizationId: "id",
		Type:           string(client.GoogleCostCredentialsType),
	}

	updateReturnValues := client.Credentials{
		Id:             "id",
		Name:           "update",
		OrganizationId: "id",
		Type:           string(client.GoogleCostCredentialsType),
	}

	testCaseForCreate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, googleCostCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", googleCostCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "table_id", googleCostCredentialResource["table_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "secret", googleCostCredentialResource["secret"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", returnValues.Id),
				),
			},
		},
	}

	testCaseForUpdate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, googleCostCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", googleCostCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "table_id", googleCostCredentialResource["table_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "secret", googleCostCredentialResource["secret"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", returnValues.Id),
				),
			},
			{
				Config: resourceConfigCreate(resourceType, resourceName, updateGoogleCostCredentialResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", updateGoogleCostCredentialResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "table_id", updateGoogleCostCredentialResource["table_id"].(string)),
					resource.TestCheckResourceAttr(accessor, "secret", updateGoogleCostCredentialResource["secret"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", updateReturnValues.Id),
				),
			},
		},
	}

	t.Run("create", func(t *testing.T) {
		runUnitTest(t, testCaseForCreate, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CredentialsCreate(&googleCostCredCreatePayload).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil)
			mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("any update cause a destroy before a new create", func(t *testing.T) {
		runUnitTest(t, testCaseForUpdate, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CredentialsCreate(&googleCostCredCreatePayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CredentialsUpdate(returnValues.Id, &updateGoogleCostCredCreatePayload).Times(1).Return(updateReturnValues, nil),
			)
			gomock.InOrder(
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(2).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(updateReturnValues.Id).Times(1).Return(updateReturnValues, nil),
			)
			mock.EXPECT().CloudCredentialsDelete(updateReturnValues.Id).Times(1).Return(nil)
		})
	})

	t.Run("create with project_id", func(t *testing.T) {
		projectId := "project-xyz"
		googleCostCredentialResourceWithProject := map[string]any{
			"name":       "test",
			"table_id":   "11111",
			"secret":     "22222",
			"project_id": projectId,
		}
		googleCostCredCreatePayloadWithProject := client.GoogleCostCredentialsCreatePayload{
			Name: googleCostCredentialResourceWithProject["name"].(string),
			Value: client.GoogleCostCredentialsValuePayload{
				TableId: googleCostCredentialResourceWithProject["table_id"].(string),
				Secret:  googleCostCredentialResourceWithProject["secret"].(string),
			},
			Type:      client.GoogleCostCredentialsType,
			ProjectId: projectId,
		}
		returnValuesWithProject := client.Credentials{
			Id:        "id",
			Name:      "test",
			Type:      string(client.GoogleCostCredentialsType),
			ProjectId: projectId,
		}
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, googleCostCredentialResourceWithProject),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "name", googleCostCredentialResourceWithProject["name"].(string)),
						resource.TestCheckResourceAttr(accessor, "table_id", googleCostCredentialResourceWithProject["table_id"].(string)),
						resource.TestCheckResourceAttr(accessor, "secret", googleCostCredentialResourceWithProject["secret"].(string)),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "id", returnValuesWithProject.Id),
					),
				},
			},
		}
		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CredentialsCreate(&googleCostCredCreatePayloadWithProject).Times(1).Return(returnValuesWithProject, nil)
			mock.EXPECT().CloudCredentials(returnValuesWithProject.Id).Times(1).Return(returnValuesWithProject, nil)
			mock.EXPECT().CloudCredentialsDelete(returnValuesWithProject.Id).Times(1).Return(nil)
		})
	})

	t.Run("update with project_id", func(t *testing.T) {
		projectId := "project-xyz"
		googleCostCredentialResourceWithProject := map[string]any{
			"name":       "test",
			"table_id":   "11111",
			"secret":     "22222",
			"project_id": projectId,
		}
		googleCostCredentialResourceWithProjectUpdated := map[string]any{
			"name":       "test-updated",
			"table_id":   "11111-updated",
			"secret":     "22222-updated",
			"project_id": projectId, // project_id remains the same
		}
		googleCostCredCreatePayloadWithProject := client.GoogleCostCredentialsCreatePayload{
			Name: googleCostCredentialResourceWithProject["name"].(string),
			Value: client.GoogleCostCredentialsValuePayload{
				TableId: googleCostCredentialResourceWithProject["table_id"].(string),
				Secret:  googleCostCredentialResourceWithProject["secret"].(string),
			},
			Type:      client.GoogleCostCredentialsType,
			ProjectId: projectId,
		}
		googleCostCredUpdatePayloadWithProject := client.GoogleCostCredentialsCreatePayload{
			Name: googleCostCredentialResourceWithProjectUpdated["name"].(string),
			Value: client.GoogleCostCredentialsValuePayload{
				TableId: googleCostCredentialResourceWithProjectUpdated["table_id"].(string),
				Secret:  googleCostCredentialResourceWithProjectUpdated["secret"].(string),
			},
			Type:      client.GoogleCostCredentialsType,
			ProjectId: projectId,
		}
		returnValuesWithProject := client.Credentials{
			Id:        "id",
			Name:      "test",
			Type:      string(client.GoogleCostCredentialsType),
			ProjectId: projectId,
		}
		updateReturnValuesWithProject := client.Credentials{
			Id:        "id",
			Name:      "test-updated",
			Type:      string(client.GoogleCostCredentialsType),
			ProjectId: projectId,
		}
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, googleCostCredentialResourceWithProject),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "name", googleCostCredentialResourceWithProject["name"].(string)),
						resource.TestCheckResourceAttr(accessor, "table_id", googleCostCredentialResourceWithProject["table_id"].(string)),
						resource.TestCheckResourceAttr(accessor, "secret", googleCostCredentialResourceWithProject["secret"].(string)),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "id", returnValuesWithProject.Id),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, googleCostCredentialResourceWithProjectUpdated),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "name", googleCostCredentialResourceWithProjectUpdated["name"].(string)),
						resource.TestCheckResourceAttr(accessor, "table_id", googleCostCredentialResourceWithProjectUpdated["table_id"].(string)),
						resource.TestCheckResourceAttr(accessor, "secret", googleCostCredentialResourceWithProjectUpdated["secret"].(string)),
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "id", updateReturnValuesWithProject.Id),
					),
				},
			},
		}
		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CredentialsCreate(&googleCostCredCreatePayloadWithProject).Times(1).Return(returnValuesWithProject, nil),
				mock.EXPECT().CredentialsUpdate(returnValuesWithProject.Id, &googleCostCredUpdatePayloadWithProject).Times(1).Return(updateReturnValuesWithProject, nil),
			)
			gomock.InOrder(
				mock.EXPECT().CloudCredentials(returnValuesWithProject.Id).Times(2).Return(returnValuesWithProject, nil),
				mock.EXPECT().CloudCredentials(updateReturnValuesWithProject.Id).Times(1).Return(updateReturnValuesWithProject, nil),
			)
			mock.EXPECT().CloudCredentialsDelete(updateReturnValuesWithProject.Id).Times(1).Return(nil)
		})
	})
}
