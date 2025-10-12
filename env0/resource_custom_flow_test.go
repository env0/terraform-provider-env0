package env0

import (
	"errors"
	"regexp"
	"strconv"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitCustomFlowResource(t *testing.T) {
	resourceType := "env0_custom_flow"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	customFlow := client.CustomFlow{
		Id:                   uuid.NewString(),
		Name:                 "name",
		Repository:           "repository",
		Path:                 "path",
		Revision:             "revision",
		TokenId:              "token_id",
		GithubInstallationId: 1,
		IsGithubEnterprise:   true,
	}

	updatedCustomFlow := client.CustomFlow{
		Id:            customFlow.Id,
		Name:          "name2",
		Repository:    "repository2",
		Path:          "path2",
		Revision:      "revision2",
		TokenId:       "token_id2",
		IsAzureDevOps: true,
	}

	createPayload := client.CustomFlowCreatePayload{
		Name:                 customFlow.Name,
		Repository:           customFlow.Repository,
		Path:                 customFlow.Path,
		Revision:             customFlow.Revision,
		TokenId:              customFlow.TokenId,
		GithubInstallationId: customFlow.GithubInstallationId,
		IsGithubEnterprise:   customFlow.IsGithubEnterprise,
	}

	updatePayload := client.CustomFlowCreatePayload{
		Name:          updatedCustomFlow.Name,
		Repository:    updatedCustomFlow.Repository,
		Path:          updatedCustomFlow.Path,
		Revision:      updatedCustomFlow.Revision,
		TokenId:       updatedCustomFlow.TokenId,
		IsAzureDevOps: updatedCustomFlow.IsAzureDevOps,
	}

	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":                   customFlow.Name,
						"repository":             customFlow.Repository,
						"path":                   customFlow.Path,
						"revision":               customFlow.Revision,
						"token_id":               customFlow.TokenId,
						"github_installation_id": customFlow.GithubInstallationId,
						"is_github_enterprise":   "true",
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", customFlow.Id),
						resource.TestCheckResourceAttr(accessor, "repository", customFlow.Repository),
						resource.TestCheckResourceAttr(accessor, "path", customFlow.Path),
						resource.TestCheckResourceAttr(accessor, "revision", customFlow.Revision),
						resource.TestCheckResourceAttr(accessor, "token_id", customFlow.TokenId),
						resource.TestCheckResourceAttr(accessor, "github_installation_id", strconv.Itoa(customFlow.GithubInstallationId)),
						resource.TestCheckResourceAttr(accessor, "is_github_enterprise", strconv.FormatBool(customFlow.IsGithubEnterprise)),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":            updatedCustomFlow.Name,
						"repository":      updatedCustomFlow.Repository,
						"path":            updatedCustomFlow.Path,
						"revision":        updatedCustomFlow.Revision,
						"token_id":        updatedCustomFlow.TokenId,
						"is_azure_devops": "true",
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedCustomFlow.Id),
						resource.TestCheckResourceAttr(accessor, "repository", updatedCustomFlow.Repository),
						resource.TestCheckResourceAttr(accessor, "path", updatedCustomFlow.Path),
						resource.TestCheckResourceAttr(accessor, "revision", updatedCustomFlow.Revision),
						resource.TestCheckResourceAttr(accessor, "token_id", updatedCustomFlow.TokenId),
						resource.TestCheckResourceAttr(accessor, "github_installation_id", "0"),
						resource.TestCheckResourceAttr(accessor, "is_github_enterprise", "false"),
						resource.TestCheckResourceAttr(accessor, "is_azure_devops", strconv.FormatBool(updatedCustomFlow.IsAzureDevOps)),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CustomFlowCreate(createPayload).Times(1).Return(&customFlow, nil),
				mock.EXPECT().CustomFlow(customFlow.Id).Times(2).Return(&customFlow, nil),
				mock.EXPECT().CustomFlowUpdate(customFlow.Id, updatePayload).Times(1).Return(&updatedCustomFlow, nil),
				mock.EXPECT().CustomFlow(updatedCustomFlow.Id).Times(1).Return(&updatedCustomFlow, nil),
				mock.EXPECT().CustomFlowDelete(customFlow.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Create Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":                   customFlow.Name,
						"repository":             customFlow.Repository,
						"path":                   customFlow.Path,
						"revision":               customFlow.Revision,
						"token_id":               customFlow.TokenId,
						"github_installation_id": customFlow.GithubInstallationId,
						"is_github_enterprise":   "true",
					}),
					ExpectError: regexp.MustCompile("could not create custom flow: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CustomFlowCreate(createPayload).Times(1).Return(nil, errors.New("error"))
		})
	})

	t.Run("Update Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":                   customFlow.Name,
						"repository":             customFlow.Repository,
						"path":                   customFlow.Path,
						"revision":               customFlow.Revision,
						"token_id":               customFlow.TokenId,
						"github_installation_id": customFlow.GithubInstallationId,
						"is_github_enterprise":   "true",
					}),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":            updatedCustomFlow.Name,
						"repository":      updatedCustomFlow.Repository,
						"path":            updatedCustomFlow.Path,
						"revision":        updatedCustomFlow.Revision,
						"token_id":        updatedCustomFlow.TokenId,
						"is_azure_devops": "true",
					}),
					ExpectError: regexp.MustCompile("could not update custom flow: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CustomFlowCreate(createPayload).Times(1).Return(&customFlow, nil),
				mock.EXPECT().CustomFlow(customFlow.Id).Times(2).Return(&customFlow, nil),
				mock.EXPECT().CustomFlowUpdate(customFlow.Id, updatePayload).Times(1).Return(nil, errors.New("error")),
				mock.EXPECT().CustomFlowDelete(customFlow.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Import By Name", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":                   customFlow.Name,
						"repository":             customFlow.Repository,
						"path":                   customFlow.Path,
						"revision":               customFlow.Revision,
						"token_id":               customFlow.TokenId,
						"github_installation_id": customFlow.GithubInstallationId,
						"is_github_enterprise":   "true",
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     customFlow.Name,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CustomFlowCreate(createPayload).Times(1).Return(&customFlow, nil),
				mock.EXPECT().CustomFlow(customFlow.Id).Times(1).Return(&customFlow, nil),
				mock.EXPECT().CustomFlows(customFlow.Name).Times(1).Return([]client.CustomFlow{customFlow}, nil),
				mock.EXPECT().CustomFlow(customFlow.Id).Times(1).Return(&customFlow, nil),
				mock.EXPECT().CustomFlowDelete(customFlow.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Import By Id", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":                   customFlow.Name,
						"repository":             customFlow.Repository,
						"path":                   customFlow.Path,
						"revision":               customFlow.Revision,
						"token_id":               customFlow.TokenId,
						"github_installation_id": customFlow.GithubInstallationId,
						"is_github_enterprise":   "true",
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     customFlow.Id,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().CustomFlowCreate(createPayload).Times(1).Return(&customFlow, nil),
				mock.EXPECT().CustomFlow(customFlow.Id).Times(3).Return(&customFlow, nil),
				mock.EXPECT().CustomFlowDelete(customFlow.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Success with vcs_connection_id", func(t *testing.T) {
		vcsCustomFlow := client.CustomFlow{
			Id:              uuid.NewString(),
			Name:            "vcs-custom-flow",
			Repository:      "repository",
			Path:            "path",
			Revision:        "revision",
			VcsConnectionId: "vcs-conn-123",
		}

		vcsCreatePayload := client.CustomFlowCreatePayload{
			Name:            vcsCustomFlow.Name,
			Repository:      vcsCustomFlow.Repository,
			Path:            vcsCustomFlow.Path,
			Revision:        vcsCustomFlow.Revision,
			VcsConnectionId: vcsCustomFlow.VcsConnectionId,
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":              vcsCustomFlow.Name,
						"repository":        vcsCustomFlow.Repository,
						"path":              vcsCustomFlow.Path,
						"revision":          vcsCustomFlow.Revision,
						"vcs_connection_id": vcsCustomFlow.VcsConnectionId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", vcsCustomFlow.Id),
						resource.TestCheckResourceAttr(accessor, "name", vcsCustomFlow.Name),
						resource.TestCheckResourceAttr(accessor, "vcs_connection_id", vcsCustomFlow.VcsConnectionId),
						resource.TestCheckNoResourceAttr(accessor, "github_installation_id"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CustomFlowCreate(vcsCreatePayload).Times(1).Return(&vcsCustomFlow, nil)
			mock.EXPECT().CustomFlow(vcsCustomFlow.Id).Times(1).Return(&vcsCustomFlow, nil)
			mock.EXPECT().CustomFlowDelete(vcsCustomFlow.Id).Times(1).Return(nil)
		})
	})

	t.Run("vcs_connection_id ignores github_installation_id from backend to avoid drift", func(t *testing.T) {
		vcsCustomFlow := client.CustomFlow{
			Id:              uuid.NewString(),
			Name:            "vcs-custom-flow",
			Repository:      "repository",
			VcsConnectionId: "vcs-conn-123",
		}

		vcsCustomFlowFromBackend := client.CustomFlow{
			Id:                   vcsCustomFlow.Id,
			Name:                 vcsCustomFlow.Name,
			Repository:           vcsCustomFlow.Repository,
			VcsConnectionId:      vcsCustomFlow.VcsConnectionId,
			GithubInstallationId: 456,
		}

		vcsCreatePayload := client.CustomFlowCreatePayload{
			Name:            vcsCustomFlow.Name,
			Repository:      vcsCustomFlow.Repository,
			VcsConnectionId: vcsCustomFlow.VcsConnectionId,
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":              vcsCustomFlow.Name,
						"repository":        vcsCustomFlow.Repository,
						"vcs_connection_id": vcsCustomFlow.VcsConnectionId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", vcsCustomFlow.Id),
						resource.TestCheckResourceAttr(accessor, "name", vcsCustomFlow.Name),
						resource.TestCheckResourceAttr(accessor, "vcs_connection_id", vcsCustomFlow.VcsConnectionId),
						resource.TestCheckNoResourceAttr(accessor, "github_installation_id"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().CustomFlowCreate(vcsCreatePayload).Times(1).Return(&vcsCustomFlow, nil)
			mock.EXPECT().CustomFlow(vcsCustomFlow.Id).Times(1).Return(&vcsCustomFlowFromBackend, nil)
			mock.EXPECT().CustomFlowDelete(vcsCustomFlow.Id).Times(1).Return(nil)
		})
	})

	t.Run("vcs_connection_id and github_installation_id are mutually exclusive", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":                   "custom-flow",
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
