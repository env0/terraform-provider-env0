package env0

import (
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestEnvironmentImportResource(t *testing.T) {
	resourceType := "env0_environment_import"
	resourceName := "test"
	accessor := resourceAccessor(resourceType, resourceName)

	environmentImport := client.EnvironmentImport{
		Id:        "id0",
		Name:      "name0",
		IacType:   "opentofu",
		Workspace: "workspace0",
		GitConfig: client.GitConfig{
			Repository: "repo0",
			Path:       "path0",
			Revision:   "rev0",
			Provider:   "github",
		},
		IacVersion: "1.0",
	}

	updatedEnvironmentImport := client.EnvironmentImport{
		Id:        environmentImport.Id,
		Name:      "new name",
		Workspace: "new workspace",
		IacType:   "terraform",
		GitConfig: client.GitConfig{
			Repository: "new repo",
			Path:       "new path",
			Revision:   "new rev",
			Provider:   "gitlab",
		},
		IacVersion: "2.0",
	}

	t.Run("Test environment import", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":         environmentImport.Name,
						"iac_type":     environmentImport.IacType,
						"workspace":    environmentImport.Workspace,
						"repository":   environmentImport.GitConfig.Repository,
						"path":         environmentImport.GitConfig.Path,
						"revision":     environmentImport.GitConfig.Revision,
						"git_provider": environmentImport.GitConfig.Provider,
						"iac_version":  environmentImport.IacVersion,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", environmentImport.Id),
						resource.TestCheckResourceAttr(accessor, "name", environmentImport.Name),
						resource.TestCheckResourceAttr(accessor, "iac_type", environmentImport.IacType),
						resource.TestCheckResourceAttr(accessor, "workspace", environmentImport.Workspace),
						resource.TestCheckResourceAttr(accessor, "repository", environmentImport.GitConfig.Repository),
						resource.TestCheckResourceAttr(accessor, "path", environmentImport.GitConfig.Path),
						resource.TestCheckResourceAttr(accessor, "revision", environmentImport.GitConfig.Revision),
						resource.TestCheckResourceAttr(accessor, "git_provider", environmentImport.GitConfig.Provider),
						resource.TestCheckResourceAttr(accessor, "iac_version", environmentImport.IacVersion),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":         updatedEnvironmentImport.Name,
						"iac_type":     updatedEnvironmentImport.IacType,
						"workspace":    updatedEnvironmentImport.Workspace,
						"repository":   updatedEnvironmentImport.GitConfig.Repository,
						"path":         updatedEnvironmentImport.GitConfig.Path,
						"revision":     updatedEnvironmentImport.GitConfig.Revision,
						"git_provider": updatedEnvironmentImport.GitConfig.Provider,
						"iac_version":  updatedEnvironmentImport.IacVersion,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedEnvironmentImport.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedEnvironmentImport.Name),
						resource.TestCheckResourceAttr(accessor, "iac_type", updatedEnvironmentImport.IacType),
						resource.TestCheckResourceAttr(accessor, "workspace", updatedEnvironmentImport.Workspace),
						resource.TestCheckResourceAttr(accessor, "repository", updatedEnvironmentImport.GitConfig.Repository),
						resource.TestCheckResourceAttr(accessor, "path", updatedEnvironmentImport.GitConfig.Path),
						resource.TestCheckResourceAttr(accessor, "revision", updatedEnvironmentImport.GitConfig.Revision),
						resource.TestCheckResourceAttr(accessor, "git_provider", updatedEnvironmentImport.GitConfig.Provider),
						resource.TestCheckResourceAttr(accessor, "iac_version", updatedEnvironmentImport.IacVersion),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().EnvironmentImportCreate(&client.EnvironmentImportCreatePayload{
				Name:       environmentImport.Name,
				IacType:    environmentImport.IacType,
				Workspace:  environmentImport.Workspace,
				GitConfig:  environmentImport.GitConfig,
				IacVersion: environmentImport.IacVersion,
			}).Times(1).Return(&environmentImport, nil)
			mock.EXPECT().EnvironmentImportUpdate(updatedEnvironmentImport.Id, &client.EnvironmentImportUpdatePayload{
				Name: updatedEnvironmentImport.Name,
				GitConfig: client.GitConfig{
					Repository: updatedEnvironmentImport.GitConfig.Repository,
					Path:       updatedEnvironmentImport.GitConfig.Path,
					Revision:   updatedEnvironmentImport.GitConfig.Revision,
					Provider:   updatedEnvironmentImport.GitConfig.Provider,
				},
				IacType:    updatedEnvironmentImport.IacType,
				IacVersion: updatedEnvironmentImport.IacVersion,
			}).Times(1).Return(&updatedEnvironmentImport, nil)

			gomock.InOrder(
				mock.EXPECT().EnvironmentImportGet(gomock.Any()).Times(2).Return(&environmentImport, nil),        // 1 after create, 1 before update
				mock.EXPECT().EnvironmentImportGet(gomock.Any()).Times(1).Return(&updatedEnvironmentImport, nil), // 1 after update
				mock.EXPECT().EnvironmentImportDelete(environmentImport.Id).Times(1),                             // 1 after update
			)
		})
	})

	t.Run("Environment Import soft delete", func(t *testing.T) {
		environmentImport := client.EnvironmentImport{
			Id:   "id0",
			Name: "name0",
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":        environmentImport.Name,
						"soft_delete": true,
					})},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().EnvironmentImportCreate(&client.EnvironmentImportCreatePayload{
				Name: environmentImport.Name,
			}).Times(1).Return(&environmentImport, nil)

			gomock.InOrder(
				mock.EXPECT().EnvironmentImportGet(gomock.Any()).Times(2).Return(&environmentImport, nil),
				mock.EXPECT().EnvironmentImportDelete(environmentImport.Id).Times(0),
			)
		})
	})

}
