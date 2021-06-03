package env0

import (
	"encoding/json"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestUnitSshKeyDataSourceById(t *testing.T) {
	testUnitSshKeyDataSource(t, "id")
}

func TestUnitSshKeyDataSourceByName(t *testing.T) {
	testUnitSshKeyDataSource(t, "name")
}

func testUnitSshKeyDataSource(t *testing.T, byKey string) {
	resourceType := "env0_ssh_key"
	resourceName := "test"
	accessor := dataSourceAccessor(resourceType, resourceName)
	sshKey := client.SshKey{
		Id:    "id0",
		Name:  "name0",
		Value: "Key🔑",
	}

	sshKeyAsJson, _ := json.Marshal(sshKey)
	var jsonData map[string]string
	_ = json.Unmarshal(sshKeyAsJson, &jsonData)

	testCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{
					byKey: jsonData[byKey],
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "id", sshKey.Id),
					resource.TestCheckResourceAttr(accessor, "name", sshKey.Name),
				),
			},
		},
	}

	runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
		// TODO: AnyTimes because we find that READ runs for 5 times. need investigation.
		mock.EXPECT().SshKeys().AnyTimes().Return([]client.SshKey{sshKey}, nil)
	})
}
