package env0

import (
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
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

	updatedSshKey := *sshKey
	updatedSshKey.Value = "new-valuw"

	sshKeyCreatePayload := client.SshKeyCreatePayload{
		Name:  sshKey.Name,
		Value: sshKey.Value,
	}

	sshKeyUpdatePayload := client.SshKeyUpdatePayload{
		Value: updatedSshKey.Value,
	}

	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":  sshKey.Name,
						"value": sshKey.Value,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", sshKey.Id),
						resource.TestCheckResourceAttr(accessor, "name", sshKey.Name),
						resource.TestCheckResourceAttr(accessor, "value", sshKey.Value),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":  updatedSshKey.Name,
						"value": updatedSshKey.Value,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", sshKey.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedSshKey.Name),
						resource.TestCheckResourceAttr(accessor, "value", updatedSshKey.Value),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().SshKeyCreate(sshKeyCreatePayload).Times(1).Return(sshKey, nil),
				mock.EXPECT().SshKeys().Times(2).Return([]client.SshKey{*sshKey}, nil),
				mock.EXPECT().SshKeyUpdate(sshKey.Id, &sshKeyUpdatePayload).Times(1).Return(&updatedSshKey, nil),
				mock.EXPECT().SshKeys().Times(1).Return([]client.SshKey{updatedSshKey}, nil),
				mock.EXPECT().SshKeyDelete(sshKey.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("SSH Key removed in UI", func(t *testing.T) {
		stepConfig := resourceConfigCreate(resourceType, resourceName, map[string]any{
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
