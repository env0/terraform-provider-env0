package env0

import (
	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestUnitSshKeyResource(t *testing.T) {
	resourceType := "env0_ssh_key"
	resourceName := "test"
	resourceFullName := resourceAccessor(resourceType, resourceName)
	sshKey := client.SshKey{
		Id:    "id0",
		Name:  "name0",
		Value: "Key🔑",
	}

	testCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: resourceConfigCreate(resourceType, resourceName, map[string]string{
					"name":  sshKey.Name,
					"value": sshKey.Value,
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "id", sshKey.Id),
					resource.TestCheckResourceAttr(resourceFullName, "name", sshKey.Name),
					resource.TestCheckResourceAttr(resourceFullName, "value", sshKey.Value),
				),
			},
		},
	}

	runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().SshKeyCreate(client.SshKeyCreatePayload{Name: sshKey.Name, Value: sshKey.Value}).Times(1).Return(sshKey, nil).Return(sshKey, nil)
		mock.EXPECT().SshKeys().Times(1).Return([]client.SshKey{sshKey}, nil)
		mock.EXPECT().SshKeyDelete(sshKey.Id).Times(1).Return(nil)
	})
}
