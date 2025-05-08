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

func TestUnitOciCredentialsResource(t *testing.T) {
	resourceType := "env0_oci_credentials"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	ociCredentialsResource := map[string]any{
		"name":         "test",
		"tenancy_ocid": "tenancy1",
		"user_ocid":    "user1",
		"fingerprint":  "fingerprint1",
		"private_key":  "privatekey1",
		"region":       "region1",
	}

	updatedOciCredentialsResource := map[string]any{
		"name":         "test",
		"tenancy_ocid": "tenancy2",
		"user_ocid":    "user2",
		"fingerprint":  "fingerprint2",
		"private_key":  "privatekey1", // unchanged to avoid drift
		"region":       "region2",
	}

	createPayload := client.OciCredentialsCreatePayload{
		Name: ociCredentialsResource["name"].(string),
		Value: client.OciCredentialsValuePayload{
			TenancyOcid: ociCredentialsResource["tenancy_ocid"].(string),
			UserOcid:    ociCredentialsResource["user_ocid"].(string),
			Fingerprint: ociCredentialsResource["fingerprint"].(string),
			PrivateKey:  ociCredentialsResource["private_key"].(string),
			Region:      ociCredentialsResource["region"].(string),
		},
		Type: client.OciApiKeyCredentialsType,
	}

	updatePayload := client.OciCredentialsCreatePayload{
		Value: client.OciCredentialsValuePayload{
			TenancyOcid: updatedOciCredentialsResource["tenancy_ocid"].(string),
			UserOcid:    updatedOciCredentialsResource["user_ocid"].(string),
			Fingerprint: updatedOciCredentialsResource["fingerprint"].(string),
			PrivateKey:  updatedOciCredentialsResource["private_key"].(string),
			Region:      updatedOciCredentialsResource["region"].(string),
		},
		Type: client.OciApiKeyCredentialsType,
	}

	returnValues := client.Credentials{
		Id:             "f595c4b6-0a24-4c22-89f7-7030045de31f",
		Name:           "test",
		OrganizationId: "id",
		Type:           string(client.OciApiKeyCredentialsType),
	}

	otherTypeReturnValues := client.Credentials{
		Id:             "f595c4b6-0a24-4c22-89f7-7030045de31a",
		Name:           "test",
		OrganizationId: "id",
		Type:           "AWS_....",
	}

	updateReturnValues := client.Credentials{
		Id:             returnValues.Id,
		Name:           returnValues.Name,
		OrganizationId: "id",
		Type:           string(client.OciApiKeyCredentialsType),
	}

	testCaseForCreateAndUpdate := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, ociCredentialsResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "id", returnValues.Id),
					resource.TestCheckResourceAttr(accessor, "name", ociCredentialsResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "tenancy_ocid", ociCredentialsResource["tenancy_ocid"].(string)),
					resource.TestCheckResourceAttr(accessor, "user_ocid", ociCredentialsResource["user_ocid"].(string)),
					resource.TestCheckResourceAttr(accessor, "fingerprint", ociCredentialsResource["fingerprint"].(string)),
					resource.TestCheckResourceAttr(accessor, "region", ociCredentialsResource["region"].(string)),
					// private_key is sensitive and should not be checked directly
				),
			},
			{
				Config: resourceConfigCreate(resourceType, resourceName, updatedOciCredentialsResource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "id", updateReturnValues.Id),
					resource.TestCheckResourceAttr(accessor, "name", updatedOciCredentialsResource["name"].(string)),
					resource.TestCheckResourceAttr(accessor, "tenancy_ocid", updatedOciCredentialsResource["tenancy_ocid"].(string)),
					resource.TestCheckResourceAttr(accessor, "user_ocid", updatedOciCredentialsResource["user_ocid"].(string)),
					resource.TestCheckResourceAttr(accessor, "fingerprint", updatedOciCredentialsResource["fingerprint"].(string)),
					resource.TestCheckResourceAttr(accessor, "region", updatedOciCredentialsResource["region"].(string)),
					// private_key is sensitive and should not be checked directly
				),
			},
		},
	}

	t.Run("create and update", func(t *testing.T) {
		runUnitTest(t, testCaseForCreateAndUpdate, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CredentialsCreate(&createPayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(2).Return(returnValues, nil),
				mock.EXPECT().CredentialsUpdate(returnValues.Id, &updatePayload).Times(1).Return(updateReturnValues, nil),
				mock.EXPECT().CloudCredentials(updateReturnValues.Id).Times(1).Return(updateReturnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("drift", func(t *testing.T) {
		stepConfig := resourceConfigCreate(resourceType, resourceName, ociCredentialsResource)

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
				mock.EXPECT().CredentialsCreate(&createPayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, http.NewMockFailedResponseError(404)),
				mock.EXPECT().CredentialsCreate(&createPayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("import by name", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, ociCredentialsResource),
				},
				{
					ResourceName:            resourceNameImport,
					ImportState:             true,
					ImportStateId:           ociCredentialsResource["name"].(string),
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"private_key", "tenancy_ocid", "user_ocid", "fingerprint", "region"},
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CredentialsCreate(&createPayload).Times(1).Return(returnValues, nil),
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
					Config: resourceConfigCreate(resourceType, resourceName, ociCredentialsResource),
				},
				{
					ResourceName:            resourceNameImport,
					ImportState:             true,
					ImportStateId:           returnValues.Id,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"private_key", "tenancy_ocid", "user_ocid", "fingerprint", "region"},
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CredentialsCreate(&createPayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(3).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("import by id not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, ociCredentialsResource),
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
				mock.EXPECT().CredentialsCreate(&createPayload).Times(1).Return(returnValues, nil),
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
					Config: resourceConfigCreate(resourceType, resourceName, ociCredentialsResource),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     ociCredentialsResource["name"].(string),
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(fmt.Sprintf("credentials with name %v not found", ociCredentialsResource["name"].(string))),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CredentialsCreate(&createPayload).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentials(returnValues.Id).Times(1).Return(returnValues, nil),
				mock.EXPECT().CloudCredentialsList().Times(1).Return([]client.Credentials{otherTypeReturnValues}, nil),
				mock.EXPECT().CloudCredentialsDelete(returnValues.Id).Times(1).Return(nil),
			)
		})
	})
}
