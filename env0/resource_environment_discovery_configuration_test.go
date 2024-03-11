package env0

import (
	"strconv"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitEnvironmentDiscoveryConfigurationResource(t *testing.T) {
	resourceType := "env0_environment_discovery_configuration"
	resourceName := "test"

	accessor := resourceAccessor(resourceType, resourceName)

	projectId := "pid"
	id := "id"

	t.Run("default", func(t *testing.T) {
		putPayload := client.EnvironmentDiscoveryPutPayload{
			GlobPattern:          "**",
			Repository:           "https://re.po",
			Type:                 "opentofu",
			EnvironmentPlacement: "topProject",
			WorkspaceNaming:      "default",
			OpentofuVersion:      "1.6.2",
			GithubInstallationId: 12345,
		}

		getPayload := client.EnvironmentDiscoveryPayload{
			Id:                   id,
			GlobPattern:          putPayload.GlobPattern,
			Repository:           putPayload.Repository,
			Type:                 putPayload.Type,
			EnvironmentPlacement: putPayload.EnvironmentPlacement,
			WorkspaceNaming:      putPayload.WorkspaceNaming,
			OpentofuVersion:      putPayload.OpentofuVersion,
			GithubInstallationId: putPayload.GithubInstallationId,
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"project_id":       projectId,
						"glob_pattern":     putPayload.GlobPattern,
						"repository":       putPayload.Repository,
						"opentofu_version": putPayload.OpentofuVersion,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "project_id", projectId),
						resource.TestCheckResourceAttr(accessor, "glob_pattern", putPayload.GlobPattern),
						resource.TestCheckResourceAttr(accessor, "repository", putPayload.Repository),
						resource.TestCheckResourceAttr(accessor, "type", putPayload.Type),
						resource.TestCheckResourceAttr(accessor, "environment_placement", putPayload.EnvironmentPlacement),
						resource.TestCheckResourceAttr(accessor, "workspace_naming", putPayload.WorkspaceNaming),
						resource.TestCheckResourceAttr(accessor, "opentofu_version", putPayload.OpentofuVersion),
						resource.TestCheckResourceAttr(accessor, "github_installation_id", strconv.Itoa(putPayload.GithubInstallationId)),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().EnableUpdateEnvironmentDiscovery(projectId, &putPayload).Times(1).Return(&getPayload, nil),
				mock.EXPECT().GetEnvironmentDiscovery(projectId).Times(1).Return(&getPayload, nil),
				mock.EXPECT().DeleteEnvironmentDiscovery(projectId).Times(1).Return(nil),
			)
		})
	})
}
