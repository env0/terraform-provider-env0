package env0

import (
	"errors"
	"regexp"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitModuleResource(t *testing.T) {
	resourceType := "env0_module"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	author := client.User{
		Name: "user0",
	}

	module := client.Module{
		ModuleName:            "name1",
		ModuleProvider:        "provider1",
		Repository:            "repository1",
		Description:           "description1",
		TokenId:               "tokenid",
		TokenName:             "tokenname",
		IsGitlab:              false,
		OrganizationId:        "org1",
		Author:                author,
		AuthorId:              "author1",
		Id:                    uuid.NewString(),
		Path:                  "path1",
		TagPrefix:             "prefix1",
		ModuleTestEnabled:     true,
		OpentofuVersion:       "1.7.0",
		RunTestsOnPullRequest: true,
	}

	updatedModule := client.Module{
		ModuleName:            "name2",
		ModuleProvider:        "provider1",
		Repository:            "repository1",
		Description:           "description1",
		BitbucketClientKey:    new("1234"),
		OrganizationId:        "org1",
		Author:                author,
		AuthorId:              "author1",
		Id:                    module.Id,
		Path:                  "path2",
		TagPrefix:             "prefix2",
		ModuleTestEnabled:     true,
		OpentofuVersion:       "1.8.0",
		RunTestsOnPullRequest: false,
		IsAzureDevOps:         true,
	}

	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"module_name":               module.ModuleName,
						"module_provider":           module.ModuleProvider,
						"repository":                module.Repository,
						"description":               module.Description,
						"token_id":                  module.TokenId,
						"token_name":                module.TokenName,
						"path":                      module.Path,
						"tag_prefix":                module.TagPrefix,
						"module_test_enabled":       module.ModuleTestEnabled,
						"run_tests_on_pull_request": module.RunTestsOnPullRequest,
						"opentofu_version":          module.OpentofuVersion,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", module.Id),
						resource.TestCheckResourceAttr(accessor, "module_name", module.ModuleName),
						resource.TestCheckResourceAttr(accessor, "module_provider", module.ModuleProvider),
						resource.TestCheckResourceAttr(accessor, "repository", module.Repository),
						resource.TestCheckResourceAttr(accessor, "description", module.Description),
						resource.TestCheckResourceAttr(accessor, "token_id", module.TokenId),
						resource.TestCheckResourceAttr(accessor, "token_name", module.TokenName),
						resource.TestCheckResourceAttr(accessor, "path", module.Path),
						resource.TestCheckResourceAttr(accessor, "tag_prefix", module.TagPrefix),
						resource.TestCheckResourceAttr(accessor, "module_test_enabled", "true"),
						resource.TestCheckResourceAttr(accessor, "run_tests_on_pull_request", "true"),
						resource.TestCheckResourceAttr(accessor, "opentofu_version", module.OpentofuVersion),
						resource.TestCheckResourceAttr(accessor, "is_azure_devops", "false"),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"module_name":          updatedModule.ModuleName,
						"module_provider":      updatedModule.ModuleProvider,
						"repository":           updatedModule.Repository,
						"description":          updatedModule.Description,
						"bitbucket_client_key": *updatedModule.BitbucketClientKey,
						"path":                 updatedModule.Path,
						"tag_prefix":           updatedModule.TagPrefix,
						"module_test_enabled":  updatedModule.ModuleTestEnabled,
						"opentofu_version":     updatedModule.OpentofuVersion,
						"is_azure_devops":      updatedModule.IsAzureDevOps,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedModule.Id),
						resource.TestCheckResourceAttr(accessor, "module_name", updatedModule.ModuleName),
						resource.TestCheckResourceAttr(accessor, "module_provider", updatedModule.ModuleProvider),
						resource.TestCheckResourceAttr(accessor, "repository", updatedModule.Repository),
						resource.TestCheckResourceAttr(accessor, "description", updatedModule.Description),
						resource.TestCheckResourceAttr(accessor, "token_id", ""),
						resource.TestCheckResourceAttr(accessor, "token_name", ""),
						resource.TestCheckResourceAttr(accessor, "bitbucket_client_key", *updatedModule.BitbucketClientKey),
						resource.TestCheckResourceAttr(accessor, "path", updatedModule.Path),
						resource.TestCheckResourceAttr(accessor, "tag_prefix", updatedModule.TagPrefix),
						resource.TestCheckResourceAttr(accessor, "module_test_enabled", "true"),
						resource.TestCheckResourceAttr(accessor, "run_tests_on_pull_request", "false"),
						resource.TestCheckResourceAttr(accessor, "opentofu_version", updatedModule.OpentofuVersion),
						resource.TestCheckResourceAttr(accessor, "is_azure_devops", "true"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ModuleCreate(client.ModuleCreatePayload{
				ModuleName:            module.ModuleName,
				ModuleProvider:        module.ModuleProvider,
				Repository:            module.Repository,
				Description:           module.Description,
				TokenId:               module.TokenId,
				TokenName:             module.TokenName,
				Path:                  module.Path,
				TagPrefix:             module.TagPrefix,
				ModuleTestEnabled:     module.ModuleTestEnabled,
				RunTestsOnPullRequest: module.RunTestsOnPullRequest,
				OpentofuVersion:       module.OpentofuVersion,
			}).Times(1).Return(&module, nil)

			mock.EXPECT().ModuleUpdate(updatedModule.Id, client.ModuleUpdatePayload{
				ModuleName:           updatedModule.ModuleName,
				ModuleProvider:       updatedModule.ModuleProvider,
				Repository:           updatedModule.Repository,
				Description:          updatedModule.Description,
				TokenId:              "",
				TokenName:            "",
				IsGitlab:             false,
				GithubInstallationId: nil,
				BitbucketClientKey:   *updatedModule.BitbucketClientKey,
				Path:                 updatedModule.Path,
				TagPrefix:            updatedModule.TagPrefix,
				ModuleTestEnabled:    updatedModule.ModuleTestEnabled,
				OpentofuVersion:      updatedModule.OpentofuVersion,
				IsAzureDevOps:        updatedModule.IsAzureDevOps,
			}).Times(1).Return(&updatedModule, nil)

			gomock.InOrder(
				mock.EXPECT().Module(module.Id).Times(2).Return(&module, nil),
				mock.EXPECT().Module(module.Id).Times(1).Return(&updatedModule, nil),
			)

			mock.EXPECT().ModuleDelete(module.Id).Times(1)
		})
	})

	t.Run("Create Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"module_name":     module.ModuleName,
						"module_provider": module.ModuleProvider,
						"repository":      module.Repository,
						"description":     module.Description,
						"token_id":        module.TokenId,
						"token_name":      module.TokenName,
					}),
					ExpectError: regexp.MustCompile("could not create module: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ModuleCreate(client.ModuleCreatePayload{
				ModuleName:     module.ModuleName,
				ModuleProvider: module.ModuleProvider,
				Repository:     module.Repository,
				Description:    module.Description,
				TokenId:        module.TokenId,
				TokenName:      module.TokenName,
			}).Times(1).Return(nil, errors.New("error"))
		})
	})

	t.Run("Create Failure - Invalid ModuleName", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"module_name":     "bad!!name^",
						"module_provider": module.ModuleProvider,
						"repository":      module.Repository,
						"description":     module.Description,
						"token_id":        module.TokenId,
						"token_name":      module.TokenName,
					}),
					ExpectError: regexp.MustCompile("must match pattern"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("Create Failure - Invalid ModuleProvider", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"module_name":     module.ModuleName,
						"module_provider": "ab_c",
						"repository":      module.Repository,
						"description":     module.Description,
						"token_id":        module.TokenId,
						"token_name":      module.TokenName,
					}),
					ExpectError: regexp.MustCompile("must match pattern"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("Create Failure - module test not enabled - opentofu_version set", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"module_name":         module.ModuleName,
						"module_provider":     module.ModuleProvider,
						"repository":          module.Repository,
						"description":         module.Description,
						"token_id":            module.TokenId,
						"token_name":          module.TokenName,
						"path":                module.Path,
						"tag_prefix":          module.TagPrefix,
						"module_test_enabled": false,
						"opentofu_version":    module.OpentofuVersion,
					}),
					ExpectError: regexp.MustCompile("'run_tests_on_pull_request' and/or 'opentofu_version' may only be set"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("Create Failure - module test not enabled - run_tests_on_pull_request set", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"module_name":               module.ModuleName,
						"module_provider":           module.ModuleProvider,
						"repository":                module.Repository,
						"description":               module.Description,
						"token_id":                  module.TokenId,
						"token_name":                module.TokenName,
						"path":                      module.Path,
						"tag_prefix":                module.TagPrefix,
						"module_test_enabled":       false,
						"run_tests_on_pull_request": module.RunTestsOnPullRequest,
					}),
					ExpectError: regexp.MustCompile("'run_tests_on_pull_request' and/or 'opentofu_version' may only be set"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("Update Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"module_name":               module.ModuleName,
						"module_provider":           module.ModuleProvider,
						"repository":                module.Repository,
						"description":               module.Description,
						"token_id":                  module.TokenId,
						"token_name":                module.TokenName,
						"path":                      module.Path,
						"tag_prefix":                module.TagPrefix,
						"module_test_enabled":       module.ModuleTestEnabled,
						"run_tests_on_pull_request": module.RunTestsOnPullRequest,
						"opentofu_version":          module.OpentofuVersion,
					}),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"module_name":          updatedModule.ModuleName,
						"module_provider":      updatedModule.ModuleProvider,
						"repository":           updatedModule.Repository,
						"description":          updatedModule.Description,
						"bitbucket_client_key": *updatedModule.BitbucketClientKey,
						"path":                 updatedModule.Path,
					}),
					ExpectError: regexp.MustCompile("could not update module: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ModuleCreate(client.ModuleCreatePayload{
				ModuleName:            module.ModuleName,
				ModuleProvider:        module.ModuleProvider,
				Repository:            module.Repository,
				Description:           module.Description,
				TokenId:               module.TokenId,
				TokenName:             module.TokenName,
				Path:                  module.Path,
				TagPrefix:             module.TagPrefix,
				ModuleTestEnabled:     module.ModuleTestEnabled,
				RunTestsOnPullRequest: module.RunTestsOnPullRequest,
				OpentofuVersion:       module.OpentofuVersion,
			}).Times(1).Return(&module, nil)

			mock.EXPECT().ModuleUpdate(updatedModule.Id, client.ModuleUpdatePayload{
				ModuleName:           updatedModule.ModuleName,
				ModuleProvider:       updatedModule.ModuleProvider,
				Repository:           updatedModule.Repository,
				Description:          updatedModule.Description,
				TokenId:              "",
				TokenName:            "",
				IsGitlab:             false,
				GithubInstallationId: nil,
				BitbucketClientKey:   *updatedModule.BitbucketClientKey,
				Path:                 updatedModule.Path,
			}).Times(1).Return(nil, errors.New("error"))

			mock.EXPECT().Module(module.Id).Times(2).Return(&module, nil)
			mock.EXPECT().ModuleDelete(module.Id).Times(1)
		})
	})

	t.Run("Update Failure - module test not enabled", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"module_name":               module.ModuleName,
						"module_provider":           module.ModuleProvider,
						"repository":                module.Repository,
						"description":               module.Description,
						"token_id":                  module.TokenId,
						"token_name":                module.TokenName,
						"path":                      module.Path,
						"tag_prefix":                module.TagPrefix,
						"module_test_enabled":       module.ModuleTestEnabled,
						"run_tests_on_pull_request": module.RunTestsOnPullRequest,
						"opentofu_version":          module.OpentofuVersion,
					}),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"module_name":               updatedModule.ModuleName,
						"module_provider":           updatedModule.ModuleProvider,
						"repository":                updatedModule.Repository,
						"description":               updatedModule.Description,
						"bitbucket_client_key":      *updatedModule.BitbucketClientKey,
						"path":                      updatedModule.Path,
						"module_test_enabled":       false,
						"run_tests_on_pull_request": true,
					}),
					ExpectError: regexp.MustCompile("'run_tests_on_pull_request' and/or 'opentofu_version' may only be set"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ModuleCreate(client.ModuleCreatePayload{
				ModuleName:            module.ModuleName,
				ModuleProvider:        module.ModuleProvider,
				Repository:            module.Repository,
				Description:           module.Description,
				TokenId:               module.TokenId,
				TokenName:             module.TokenName,
				Path:                  module.Path,
				TagPrefix:             module.TagPrefix,
				ModuleTestEnabled:     module.ModuleTestEnabled,
				RunTestsOnPullRequest: module.RunTestsOnPullRequest,
				OpentofuVersion:       module.OpentofuVersion,
			}).Times(1).Return(&module, nil)

			mock.EXPECT().Module(module.Id).Times(2).Return(&module, nil)
			mock.EXPECT().ModuleDelete(module.Id).Times(1)
		})
	})

	t.Run("Import By Name", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"module_name":               module.ModuleName,
						"module_provider":           module.ModuleProvider,
						"repository":                module.Repository,
						"description":               module.Description,
						"token_id":                  module.TokenId,
						"token_name":                module.TokenName,
						"path":                      module.Path,
						"tag_prefix":                module.TagPrefix,
						"module_test_enabled":       module.ModuleTestEnabled,
						"run_tests_on_pull_request": module.RunTestsOnPullRequest,
						"opentofu_version":          module.OpentofuVersion,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     module.ModuleName,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ModuleCreate(client.ModuleCreatePayload{
				ModuleName:            module.ModuleName,
				ModuleProvider:        module.ModuleProvider,
				Repository:            module.Repository,
				Description:           module.Description,
				TokenId:               module.TokenId,
				TokenName:             module.TokenName,
				Path:                  module.Path,
				TagPrefix:             module.TagPrefix,
				ModuleTestEnabled:     module.ModuleTestEnabled,
				RunTestsOnPullRequest: module.RunTestsOnPullRequest,
				OpentofuVersion:       module.OpentofuVersion,
			}).Times(1).Return(&module, nil)
			mock.EXPECT().Modules().Times(1).Return([]client.Module{module}, nil)
			mock.EXPECT().Module(module.Id).Times(2).Return(&module, nil)
			mock.EXPECT().ModuleDelete(module.Id).Times(1)
		})
	})

	t.Run("Import By Id", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"module_name":               module.ModuleName,
						"module_provider":           module.ModuleProvider,
						"repository":                module.Repository,
						"description":               module.Description,
						"token_id":                  module.TokenId,
						"token_name":                module.TokenName,
						"path":                      module.Path,
						"tag_prefix":                module.TagPrefix,
						"module_test_enabled":       module.ModuleTestEnabled,
						"run_tests_on_pull_request": module.RunTestsOnPullRequest,
						"opentofu_version":          module.OpentofuVersion,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     module.Id,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ModuleCreate(client.ModuleCreatePayload{
				ModuleName:            module.ModuleName,
				ModuleProvider:        module.ModuleProvider,
				Repository:            module.Repository,
				Description:           module.Description,
				TokenId:               module.TokenId,
				TokenName:             module.TokenName,
				Path:                  module.Path,
				TagPrefix:             module.TagPrefix,
				ModuleTestEnabled:     module.ModuleTestEnabled,
				RunTestsOnPullRequest: module.RunTestsOnPullRequest,
				OpentofuVersion:       module.OpentofuVersion,
			}).Times(1).Return(&module, nil)
			mock.EXPECT().Module(module.Id).Times(3).Return(&module, nil)
			mock.EXPECT().ModuleDelete(module.Id).Times(1)
		})
	})

	t.Run("Import By Name - Module Is Deleted", func(t *testing.T) {
		deletedModule := module
		deletedModule.IsDeleted = true

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"module_name":               module.ModuleName,
						"module_provider":           module.ModuleProvider,
						"repository":                module.Repository,
						"description":               module.Description,
						"token_id":                  module.TokenId,
						"token_name":                module.TokenName,
						"path":                      module.Path,
						"tag_prefix":                module.TagPrefix,
						"module_test_enabled":       module.ModuleTestEnabled,
						"run_tests_on_pull_request": module.RunTestsOnPullRequest,
						"opentofu_version":          module.OpentofuVersion,
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     deletedModule.ModuleName,
					ImportStateVerify: false,
					ExpectError:       regexp.MustCompile("module with name .* not found"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ModuleCreate(gomock.Any()).Times(1).Return(&module, nil)
			mock.EXPECT().Module(module.Id).Times(1).Return(&module, nil)
			mock.EXPECT().Modules().Times(1).Return([]client.Module{deletedModule}, nil)
			mock.EXPECT().ModuleDelete(module.Id).Times(1)
		})
	})

	t.Run("Detect Drift When Module Is Deleted", func(t *testing.T) {
		deletedModule := module
		deletedModule.IsDeleted = true

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"module_name":               module.ModuleName,
						"module_provider":           module.ModuleProvider,
						"repository":                module.Repository,
						"description":               module.Description,
						"token_id":                  module.TokenId,
						"token_name":                module.TokenName,
						"path":                      module.Path,
						"tag_prefix":                module.TagPrefix,
						"module_test_enabled":       module.ModuleTestEnabled,
						"run_tests_on_pull_request": module.RunTestsOnPullRequest,
						"opentofu_version":          module.OpentofuVersion,
					}),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"module_name":               module.ModuleName,
						"module_provider":           module.ModuleProvider,
						"repository":                module.Repository,
						"description":               module.Description,
						"token_id":                  module.TokenId,
						"token_name":                module.TokenName,
						"path":                      module.Path,
						"tag_prefix":                module.TagPrefix,
						"module_test_enabled":       module.ModuleTestEnabled,
						"run_tests_on_pull_request": module.RunTestsOnPullRequest,
						"opentofu_version":          module.OpentofuVersion,
					}),
					ExpectNonEmptyPlan: true,
					PlanOnly:           true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ModuleCreate(gomock.Any()).Times(1).Return(&module, nil)
			mock.EXPECT().Module(module.Id).Times(1).Return(&module, nil)
			mock.EXPECT().Module(module.Id).Times(2).Return(&deletedModule, nil)
			mock.EXPECT().ModuleDelete(module.Id).Times(1)
		})
	})

	t.Run("Success with vcs_connection_id", func(t *testing.T) {
		vcsModule := client.Module{
			ModuleName:      "vcs-module",
			ModuleProvider:  "provider1",
			Repository:      "repository1",
			Description:     "description1",
			VcsConnectionId: new("vcs-conn-123"),
			OrganizationId:  "org1",
			Author:          author,
			AuthorId:        "author1",
			Id:              uuid.NewString(),
			Path:            "path1",
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"module_name":       vcsModule.ModuleName,
						"module_provider":   vcsModule.ModuleProvider,
						"repository":        vcsModule.Repository,
						"description":       vcsModule.Description,
						"vcs_connection_id": *vcsModule.VcsConnectionId,
						"path":              vcsModule.Path,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", vcsModule.Id),
						resource.TestCheckResourceAttr(accessor, "module_name", vcsModule.ModuleName),
						resource.TestCheckResourceAttr(accessor, "vcs_connection_id", *vcsModule.VcsConnectionId),
						resource.TestCheckNoResourceAttr(accessor, "github_installation_id"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ModuleCreate(client.ModuleCreatePayload{
				ModuleName:      vcsModule.ModuleName,
				ModuleProvider:  vcsModule.ModuleProvider,
				Repository:      vcsModule.Repository,
				Description:     vcsModule.Description,
				VcsConnectionId: *vcsModule.VcsConnectionId,
				Path:            vcsModule.Path,
			}).Times(1).Return(&vcsModule, nil)
			mock.EXPECT().Module(vcsModule.Id).Times(1).Return(&vcsModule, nil)
			mock.EXPECT().ModuleDelete(vcsModule.Id).Times(1)
		})
	})

	t.Run("vcs_connection_id ignores github_installation_id from backend to avoid drift", func(t *testing.T) {
		vcsModule := client.Module{
			ModuleName:      "vcs-module",
			ModuleProvider:  "provider1",
			Repository:      "repository1",
			VcsConnectionId: new("vcs-conn-123"),
			OrganizationId:  "org1",
			Author:          author,
			AuthorId:        "author1",
			Id:              uuid.NewString(),
		}

		vcsModuleFromBackend := client.Module{
			ModuleName:           vcsModule.ModuleName,
			ModuleProvider:       vcsModule.ModuleProvider,
			Repository:           vcsModule.Repository,
			VcsConnectionId:      vcsModule.VcsConnectionId,
			GithubInstallationId: new(456),
			OrganizationId:       vcsModule.OrganizationId,
			Author:               vcsModule.Author,
			AuthorId:             vcsModule.AuthorId,
			Id:                   vcsModule.Id,
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"module_name":       vcsModule.ModuleName,
						"module_provider":   vcsModule.ModuleProvider,
						"repository":        vcsModule.Repository,
						"vcs_connection_id": *vcsModule.VcsConnectionId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", vcsModule.Id),
						resource.TestCheckResourceAttr(accessor, "module_name", vcsModule.ModuleName),
						resource.TestCheckResourceAttr(accessor, "vcs_connection_id", *vcsModule.VcsConnectionId),
						resource.TestCheckNoResourceAttr(accessor, "github_installation_id"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ModuleCreate(client.ModuleCreatePayload{
				ModuleName:      vcsModule.ModuleName,
				ModuleProvider:  vcsModule.ModuleProvider,
				Repository:      vcsModule.Repository,
				VcsConnectionId: *vcsModule.VcsConnectionId,
			}).Times(1).Return(&vcsModule, nil)
			mock.EXPECT().Module(vcsModule.Id).Times(1).Return(&vcsModuleFromBackend, nil)
			mock.EXPECT().ModuleDelete(vcsModule.Id).Times(1)
		})
	})

	t.Run("vcs_connection_id and github_installation_id are mutually exclusive", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"module_name":            "module",
						"module_provider":        "provider1",
						"repository":             "repo",
						"github_installation_id": 123,
						"vcs_connection_id":      "vcs-conn-456",
					}),
					ExpectError: regexp.MustCompile("conflicts with"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})
}
