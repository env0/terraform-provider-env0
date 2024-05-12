package env0

import (
	"encoding/json"
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitSshKeyDataSourceById(t *testing.T) {
	testUnitSshKeyDataSource(t, "id")
}

func TestUnitSshKeyDataSourceByName(t *testing.T) {
	testUnitSshKeyDataSource(t, "name")
}

func TestUnitSshKeyDataSourceByIdNotFound(t *testing.T) {
	resourceType := "env0_ssh_key"
	resourceName := "test"

	testCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{
					"id": "id123",
				}),
				ExpectError: regexp.MustCompile("could not read ssh key: id id123 not found"),
			},
		},
	}

	runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
		mock.EXPECT().SshKeys().Times(1).Return(nil, nil)
	})
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

func TestUnitSshKeyDataSourceRetryByName(t *testing.T) {
	resourceType := "env0_ssh_key"
	resourceName := "test"
	accessor := dataSourceAccessor(resourceType, resourceName)

	sshKey := client.SshKey{
		Id:    "id0",
		Name:  "name0",
		Value: "Key🔑",
	}

	otherSshKey := client.SshKey{
		Id:    "id1",
		Name:  "name1",
		Value: "value1",
	}

	testCase := resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfigCreate(resourceType, resourceName, map[string]interface{}{
					"name": sshKey.Name,
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(accessor, "id", sshKey.Id),
					resource.TestCheckResourceAttr(accessor, "name", sshKey.Name),
				),
			},
		},
	}

	runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
		gomock.InOrder(
			mock.EXPECT().SshKeys().Times(1).Return(nil, errors.New("error")),
			mock.EXPECT().SshKeys().Times(1).Return([]client.SshKey{otherSshKey}, nil),
			mock.EXPECT().SshKeys().AnyTimes().Return([]client.SshKey{sshKey, otherSshKey}, nil),
		)
	})
}
