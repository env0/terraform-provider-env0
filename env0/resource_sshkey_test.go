package env0

import (
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitSshKeyResource(t *testing.T) {
	resourceType := "env0_ssh_key"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)
	sshKey := &client.SshKey{
		Id:    "id0",
		Name:  "name0",
		Value: "Key🔑",
	}
	sshKeyCreatePayload := client.SshKeyCreatePayload{
		Name:  sshKey.Name,
		Value: sshKey.Value,
	}
	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":  sshKey.Name,
						"value": sshKey.Value,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", sshKey.Id),
						resource.TestCheckResourceAttr(accessor, "name", sshKey.Name),
						resource.TestCheckResourceAttr(accessor, "value", sshKey.Value),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().SshKeyCreate(sshKeyCreatePayload).Times(1).Return(sshKey, nil)
			mock.EXPECT().SshKeys().Times(1).Return([]client.SshKey{*sshKey}, nil)
			mock.EXPECT().SshKeyDelete(sshKey.Id).Times(1).Return(nil)
		})
	})

	t.Run("SSH Key removed in UI", func(t *testing.T) {
		stepConfig := resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
			"name":  sshKey.Name,
			"value": sshKey.Value,
		})

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
				mock.EXPECT().SshKeyCreate(sshKeyCreatePayload).Times(1).Return(sshKey, nil),
				mock.EXPECT().SshKeys().Times(1).Return([]client.SshKey{*sshKey}, nil),
				mock.EXPECT().SshKeys().Times(1).Return(nil, nil),
				mock.EXPECT().SshKeyCreate(sshKeyCreatePayload).Times(1).Return(sshKey, nil),
				mock.EXPECT().SshKeys().Times(1).Return([]client.SshKey{*sshKey}, nil),
				mock.EXPECT().SshKeyDelete(sshKey.Id).Times(1).Return(nil),
			)
		})
	})
}
