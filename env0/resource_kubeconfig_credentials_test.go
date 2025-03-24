package env0

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitKubeconfigCredentialsResource(t *testing.T) {
	resourceType := "env0_kubeconfig_credentials"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	credentialsResource := map[string]any{
		"name":        "test",
		"kube_config": "file content...",
	}

	updatedCredentialsResource := map[string]any{
		"name":        "test",
		"kube_config": "new file content...",
	}

	createPayload := client.KubernetesCredentialsCreatePayload{
		Name: credentialsResource["name"].(string),
		Value: client.KubeconfigFileValue{
			KubeConfig: credentialsResource["kube_config"].(string),
		},
		Type: client.KubeconfigCredentialsType,
	}

	updatePayload := client.KubernetesCredentialsUpdatePayload{
		Value: client.KubeconfigFileValue{
			KubeConfig: updatedCredentialsResource["kube_config"].(string),
		},
		Type: client.KubeconfigCredentialsType,
	}

	returnValues := client.Credentials{
		Id:             "f595c4b6-0a24-4c22-89f7-7030045de30f",
		Name:           "test",
		OrganizationId: "id",
		Type:           string(client.KubeconfigCredentialsType),
	}

	otherTypeReturnValues := client.Credentials{
		Id:             "f595c4b6-0a24-4c22-89f7-7030045de30a",
		Name:           "test",
		OrganizationId: "id",
		Type:           "AWS_....",
	}

	testCaseForCreateAndUpdate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, credentialsResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", credentialsResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "kube_config", credentialsResource["kube_config"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", returnValues.Id),
				),
			},
			{
				Config: resourceConfigCreate(resourceType, resourceName, updatedCredentialsResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", updatedCredentialsResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "kube_config", updatedCredentialsResource["kube_config"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", returnValues.Id),
				),
			},
		},
	}

	t.Run("create and update", func(t *testing.T) {
		runUnitTest(t, testCaseForCreateAndUpdate, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().KubernetesCredentialsCreate(&createPayload).Times(1).Return(&returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(2).Return(returnValues, nil),
				mock.EXPECT().KubernetesCredentialsUpdate(returnValues.Id, &updatePayload).Times(1).Return(&returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("drift", func(t *testing.T) {
		stepConfig := resourceConfigCreate(resourceType, resourceName, credentialsResource)

		createTestCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: stepConfig,
				},
				{
					Config: stepConfig,
				},
			},
		}

		runUnitTest(t, createTestCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().KubernetesCredentialsCreate(&createPayload).Times(1).Return(&returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, http.NewMockFailedResponseError(404)),
				mock.EXPECT().KubernetesCredentialsCreate(&createPayload).Times(1).Return(&returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("import by name", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, credentialsResource),
				},
				{
					ResourceName:            resourceNameImport,
					ImportState:             true,
					ImportStateId:           credentialsResource["name"].(string),
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"kube_config"},
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().KubernetesCredentialsCreate(&createPayload).Times(1).Return(&returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsList().Times(1).Return([]client.Credentials{otherTypeReturnValues, returnValues}, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("import by id", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, credentialsResource),
				},
				{
					ResourceName:            resourceNameImport,
					ImportState:             true,
					ImportStateId:           returnValues.Id,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"kube_config"},
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().KubernetesCredentialsCreate(&createPayload).Times(1).Return(&returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(3).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("import by id not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, credentialsResource),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     otherTypeReturnValues.Id,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile("credentials not found"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().KubernetesCredentialsCreate(&createPayload).Times(1).Return(&returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(otherTypeReturnValues.Id).Times(1).Return(client.Credentials{}, &client.NotFoundError{}),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("import by name not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, credentialsResource),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     credentialsResource["name"].(string),
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(fmt.Sprintf("credentials with name %v not found", credentialsResource["name"].(string))),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().KubernetesCredentialsCreate(&createPayload).Times(1).Return(&returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsList().Times(1).Return([]client.Credentials{otherTypeReturnValues}, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})
}
