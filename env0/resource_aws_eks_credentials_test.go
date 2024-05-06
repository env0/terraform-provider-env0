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

func TestUnitAwsEksCredentialsResource(t *testing.T) {
	resourceType := "env0_aws_eks_credentials"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	credentialsResource := map[string]interface{}{
		"name":           "test",
		"cluster_name":   "my-cluster",
		"cluster_region": "us-east-2",
	}

	updatedCredentialsResource := map[string]interface{}{
		"name":           "test",
		"cluster_name":   "my-cluster2",
		"cluster_region": "us-east-2",
	}

	createPayload := client.KubernetesCredentialsCreatePayload{
		Name: credentialsResource["name"].(string),
		Value: client.AwsEksValue{
			ClusterName:   credentialsResource["cluster_name"].(string),
			ClusterRegion: credentialsResource["cluster_region"].(string),
		},
		Type: client.AwsEksCredentialsType,
	}

	updatePayload := client.KubernetesCredentialsCreatePayload{
		Name: updatedCredentialsResource["name"].(string),
		Value: client.AwsEksValue{
			ClusterName:   updatedCredentialsResource["cluster_name"].(string),
			ClusterRegion: updatedCredentialsResource["cluster_region"].(string),
		},
		Type: client.AwsEksCredentialsType,
	}

	returnValues := client.Credentials{
		Id:             "f595c4b6-0a24-4c22-89f7-7030045de30f",
		Name:           "test",
		OrganizationId: "id",
		Type:           string(client.AwsEksCredentialsType),
	}

	otherTypeReturnValues := client.Credentials{
		Id:             "f595c4b6-0a24-4c22-89f7-7030045de30a",
		Name:           "test",
		OrganizationId: "id",
		Type:           "AWS_....",
	}

	updateReturnValues := client.Credentials{
		Id:             "dsdsdsd-0a24-4c22-89f7-7030045de30f",
		Name:           returnValues.Name,
		OrganizationId: "id",
		Type:           string(client.AwsEksCredentialsType),
	}

	testCaseForCreateAndUpdate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, credentialsResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", credentialsResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "cluster_name", credentialsResource["cluster_name"].(string)),
					resource.TestCheckResourceAttr(accessor, "cluster_region", credentialsResource["cluster_region"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", returnValues.Id),
				),
			},
			{
				Config: resourceConfigCreate(resourceType, resourceName, updatedCredentialsResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "name", updatedCredentialsResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "cluster_name", updatedCredentialsResource["cluster_name"].(string)),
					resource.TestCheckResourceAttr(accessor, "cluster_region", updatedCredentialsResource["cluster_region"].(string)),
					resource.TestCheckResourceAttr(accessor, "id", updateReturnValues.Id),
				),
			},
		},
	}

	t.Run("create and update", func(t *testing.T) {
		runUnitTest(t, testCaseForCreateAndUpdate, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().KubernetesCredentialsCreate(&createPayload).Times(1).Return(&returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(2).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
				mock.EXPECT().KubernetesCredentialsCreate(&updatePayload).Times(1).Return(&updateReturnValues, nil),
				mock.EXPECT().CloudCredentials(updateReturnValues.Id).Times(1).Return(updateReturnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(updateReturnValues.Id).Times(1).Return(nil),
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
					ImportStateVerifyIgnore: []string{"cluster_name", "cluster_region"},
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
					ImportStateVerifyIgnore: []string{"cluster_name", "cluster_region"},
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
