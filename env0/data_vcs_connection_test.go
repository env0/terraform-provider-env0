package env0

import (
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestVcsConnectionDataSource(t *testing.T) {
	deployConnection := client.VcsConnection{
		Id:             "id0",
		Name:           "connection0",
		Type:           "GitHub",
		AccessScope:    "Organization:env0",
		ConnectionType: "DeploymentPipeline",
	}

	codeWriteConnection := client.VcsConnection{
		Id:             "id2",
		Name:           "connection0-cw",
		Type:           "GitHub",
		AccessScope:    "Organization:env0",
		ConnectionType: "CodeWrite",
	}

	otherConnection := client.VcsConnection{
		Id:             "id1",
		Name:           "connection1",
		Type:           "GitLabEnterprise",
		Url:            "http://glee.dev.env0.com",
		AccessScope:    "url:http://glee.dev.env0.com",
		ConnectionType: "DeploymentPipeline",
	}

	allConnections := []client.VcsConnection{deployConnection, codeWriteConnection, otherConnection}

	resourceType := "env0_vcs_connection"
	resourceName := "test"
	accessor := dataSourceAccessor(resourceType, resourceName)

	getValidTestCase := func(input map[string]any, expected client.VcsConnection) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: dataSourceConfigCreate(resourceType, resourceName, input),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", expected.Id),
						resource.TestCheckResourceAttr(accessor, "name", expected.Name),
						resource.TestCheckResourceAttr(accessor, "type", expected.Type),
						resource.TestCheckResourceAttr(accessor, "access_scope", expected.AccessScope),
						resource.TestCheckResourceAttr(accessor, "connection_type", expected.ConnectionType),
					),
				},
			},
		}
	}

	getErrorTestCase := func(input map[string]any, expectedError string) resource.TestCase {
		return resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config:      dataSourceConfigCreate(resourceType, resourceName, input),
					ExpectError: regexp.MustCompile(expectedError),
				},
			},
		}
	}

	mockList := func(mock *client.MockApiClientInterface) {
		mock.EXPECT().VcsConnections().AnyTimes().Return(allConnections, nil)
	}

	t.Run("DeploymentPipeline", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(map[string]any{
				"access_scope":    deployConnection.AccessScope,
				"connection_type": "DeploymentPipeline",
			}, deployConnection),
			mockList,
		)
	})

	t.Run("CodeWrite", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(map[string]any{
				"access_scope":    codeWriteConnection.AccessScope,
				"connection_type": "CodeWrite",
			}, codeWriteConnection),
			mockList,
		)
	})

	t.Run("Self-hosted URL", func(t *testing.T) {
		runUnitTest(t,
			getValidTestCase(map[string]any{
				"access_scope":    otherConnection.AccessScope,
				"connection_type": "DeploymentPipeline",
			}, otherConnection),
			mockList,
		)
	})

	t.Run("Not found", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(map[string]any{
				"access_scope":    "Organization:nonexistent",
				"connection_type": "DeploymentPipeline",
			}, "not found"),
			mockList,
		)
	})

	t.Run("Invalid connection type", func(t *testing.T) {
		runUnitTest(t,
			getErrorTestCase(map[string]any{
				"access_scope":    deployConnection.AccessScope,
				"connection_type": "Invalid",
			}, "must be one of"),
			func(mock *client.MockApiClientInterface) {},
		)
	})
}
