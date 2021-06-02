package env0

import (
	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestUnitSshKeyDataSourceById(t *testing.T) {
	resourceType := "env0_ssh_key"
	resourceName := "test"
	resourceFullName := dataSourceAccessor(resourceType, resourceName)
	sshKey := client.SshKey{
		Id:    "id0",
		Name:  "name0",
		Value: "Key🔑",
	}

	testCase := resource.TestCase{
		ProviderFactories: testUnitProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigCreate(resourceType, resourceName, map[string]string{
					"id": sshKey.Id,
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "id", sshKey.Id),
					resource.TestCheckResourceAttr(resourceFullName, "name", sshKey.Name),
				),
			},
		},
	}

	runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().SshKeys().AnyTimes().Return([]client.SshKey{sshKey}, nil)
	})
}

func TestUnitSshKeyDataSourceByName(t *testing.T) {
	resourceType := "env0_ssh_key"
	resourceName := "test"
	resourceFullName := dataSourceAccessor(resourceType, resourceName)
	sshKey := client.SshKey{
		Id:    "id0",
		Name:  "name0",
		Value: "Key🔑",
	}

	testCase := resource.TestCase{
		ProviderFactories: testUnitProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigCreate(resourceType, resourceName, map[string]string{
					"name": sshKey.Name,
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceFullName, "id", sshKey.Id),
					resource.TestCheckResourceAttr(resourceFullName, "name", sshKey.Name),
				),
			},
		},
	}

	runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().SshKeys().AnyTimes().Return([]client.SshKey{sshKey}, nil)
	})
}
