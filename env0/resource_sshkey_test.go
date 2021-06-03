package env0

import (
	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestUnitSshKeyResource(t *testing.T) {
	resourceType := "env0_ssh_key"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)
	sshKey := client.SshKey{
		Id:    "id0",
		Name:  "name0",
		Value: "Key🔑",
	}

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
		mock.EXPECT().SshKeyCreate(client.SshKeyCreatePayload{Name: sshKey.Name, Value: sshKey.Value}).Times(1).Return(sshKey, nil)
		mock.EXPECT().SshKeys().Times(1).Return([]client.SshKey{sshKey}, nil)
		mock.EXPECT().SshKeyDelete(sshKey.Id).Times(1).Return(nil)
	})
}
