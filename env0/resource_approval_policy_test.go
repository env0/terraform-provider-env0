package env0

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUnitApprovalPolicyResource(t *testing.T) {
	resourceType := "env0_approval_policy"
	resourceName := "test"
	resourceNameImport := resourceType + "." + resourceName
	accessor := resourceAccessor(resourceType, resourceName)

	approvalPolicy := client.ApprovalPolicy{
		Id:                   uuid.NewString(),
		Name:                 "name",
		Repository:           "repository",
		Path:                 "path",
		Revision:             "revision",
		TokenId:              "token_id",
		GithubInstallationId: 1,
		IsGithubEnterprise:   true,
	}

	var template client.Template

	require.NoError(t, copier.Copy(&template, &approvalPolicy))
	template.Type = string(ApprovalPolicy)

	deletedTemplate := template
	deletedTemplate.IsDeleted = true

	notApprovalPolicyTemplate := template
	notApprovalPolicyTemplate.Type = "terraform"

	updatedApprovalPolicy := client.ApprovalPolicy{
		Id:            approvalPolicy.Id,
		Name:          "name",
		Repository:    "repository2",
		Path:          "path2",
		Revision:      "revision2",
		TokenId:       "token_id2",
		IsAzureDevOps: true,
	}

	var updatedTemplate client.Template

	require.NoError(t, copier.Copy(&updatedTemplate, &updatedApprovalPolicy))
	updatedTemplate.Type = string(ApprovalPolicy)

	createPayload := client.ApprovalPolicyCreatePayload{
		Name:                 approvalPolicy.Name,
		Repository:           approvalPolicy.Repository,
		Path:                 approvalPolicy.Path,
		Revision:             approvalPolicy.Revision,
		TokenId:              approvalPolicy.TokenId,
		GithubInstallationId: approvalPolicy.GithubInstallationId,
		IsGithubEnterprise:   approvalPolicy.IsGithubEnterprise,
	}

	updatePayload := client.ApprovalPolicyUpdatePayload{
		Name:          updatedApprovalPolicy.Name,
		Repository:    updatedApprovalPolicy.Repository,
		Path:          updatedApprovalPolicy.Path,
		Revision:      updatedApprovalPolicy.Revision,
		TokenId:       updatedApprovalPolicy.TokenId,
		IsAzureDevOps: updatedApprovalPolicy.IsAzureDevOps,
		Id:            approvalPolicy.Id,
	}

	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":                   approvalPolicy.Name,
						"repository":             approvalPolicy.Repository,
						"path":                   approvalPolicy.Path,
						"revision":               approvalPolicy.Revision,
						"token_id":               approvalPolicy.TokenId,
						"github_installation_id": approvalPolicy.GithubInstallationId,
						"is_github_enterprise":   "true",
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", approvalPolicy.Id),
						resource.TestCheckResourceAttr(accessor, "name", approvalPolicy.Name),
						resource.TestCheckResourceAttr(accessor, "repository", approvalPolicy.Repository),
						resource.TestCheckResourceAttr(accessor, "path", approvalPolicy.Path),
						resource.TestCheckResourceAttr(accessor, "revision", approvalPolicy.Revision),
						resource.TestCheckResourceAttr(accessor, "token_id", approvalPolicy.TokenId),
						resource.TestCheckResourceAttr(accessor, "github_installation_id", strconv.Itoa(approvalPolicy.GithubInstallationId)),
						resource.TestCheckResourceAttr(accessor, "is_github_enterprise", strconv.FormatBool(approvalPolicy.IsGithubEnterprise)),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":            updatedApprovalPolicy.Name,
						"repository":      updatedApprovalPolicy.Repository,
						"path":            updatedApprovalPolicy.Path,
						"revision":        updatedApprovalPolicy.Revision,
						"token_id":        updatedApprovalPolicy.TokenId,
						"is_azure_devops": "true",
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedApprovalPolicy.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedApprovalPolicy.Name),
						resource.TestCheckResourceAttr(accessor, "repository", updatedApprovalPolicy.Repository),
						resource.TestCheckResourceAttr(accessor, "path", updatedApprovalPolicy.Path),
						resource.TestCheckResourceAttr(accessor, "revision", updatedApprovalPolicy.Revision),
						resource.TestCheckResourceAttr(accessor, "token_id", updatedApprovalPolicy.TokenId),
						resource.TestCheckResourceAttr(accessor, "github_installation_id", "0"),
						resource.TestCheckResourceAttr(accessor, "is_github_enterprise", "false"),
						resource.TestCheckResourceAttr(accessor, "is_azure_devops", strconv.FormatBool(updatedApprovalPolicy.IsAzureDevOps)),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ApprovalPolicyCreate(&createPayload).Times(1).Return(&approvalPolicy, nil),
				mock.EXPECT().Template(approvalPolicy.Id).Times(2).Return(template, nil),
				mock.EXPECT().ApprovalPolicyUpdate(&updatePayload).Times(1).Return(&updatedApprovalPolicy, nil),
				mock.EXPECT().Template(approvalPolicy.Id).Times(1).Return(updatedTemplate, nil),
				mock.EXPECT().TemplateDelete(approvalPolicy.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Drift detected - deleted", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":                   approvalPolicy.Name,
						"repository":             approvalPolicy.Repository,
						"path":                   approvalPolicy.Path,
						"revision":               approvalPolicy.Revision,
						"token_id":               approvalPolicy.TokenId,
						"github_installation_id": approvalPolicy.GithubInstallationId,
						"is_github_enterprise":   "true",
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", approvalPolicy.Id),
						resource.TestCheckResourceAttr(accessor, "name", approvalPolicy.Name),
						resource.TestCheckResourceAttr(accessor, "repository", approvalPolicy.Repository),
						resource.TestCheckResourceAttr(accessor, "path", approvalPolicy.Path),
						resource.TestCheckResourceAttr(accessor, "revision", approvalPolicy.Revision),
						resource.TestCheckResourceAttr(accessor, "token_id", approvalPolicy.TokenId),
						resource.TestCheckResourceAttr(accessor, "github_installation_id", strconv.Itoa(approvalPolicy.GithubInstallationId)),
						resource.TestCheckResourceAttr(accessor, "is_github_enterprise", strconv.FormatBool(approvalPolicy.IsGithubEnterprise)),
					),
					ExpectNonEmptyPlan: true,
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":            updatedApprovalPolicy.Name,
						"repository":      updatedApprovalPolicy.Repository,
						"path":            updatedApprovalPolicy.Path,
						"revision":        updatedApprovalPolicy.Revision,
						"token_id":        updatedApprovalPolicy.TokenId,
						"is_azure_devops": "true",
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedApprovalPolicy.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedApprovalPolicy.Name),
						resource.TestCheckResourceAttr(accessor, "repository", updatedApprovalPolicy.Repository),
						resource.TestCheckResourceAttr(accessor, "path", updatedApprovalPolicy.Path),
						resource.TestCheckResourceAttr(accessor, "revision", updatedApprovalPolicy.Revision),
						resource.TestCheckResourceAttr(accessor, "token_id", updatedApprovalPolicy.TokenId),
						resource.TestCheckResourceAttr(accessor, "github_installation_id", "0"),
						resource.TestCheckResourceAttr(accessor, "is_github_enterprise", "false"),
						resource.TestCheckResourceAttr(accessor, "is_azure_devops", strconv.FormatBool(updatedApprovalPolicy.IsAzureDevOps)),
					),
					PlanOnly:           true,
					ExpectNonEmptyPlan: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ApprovalPolicyCreate(&createPayload).Times(1).Return(&approvalPolicy, nil),
				mock.EXPECT().Template(template.Id).Times(2).Return(deletedTemplate, nil),
			)
		})
	})

	t.Run("Drift detected - not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":                   approvalPolicy.Name,
						"repository":             approvalPolicy.Repository,
						"path":                   approvalPolicy.Path,
						"revision":               approvalPolicy.Revision,
						"token_id":               approvalPolicy.TokenId,
						"github_installation_id": approvalPolicy.GithubInstallationId,
						"is_github_enterprise":   "true",
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", approvalPolicy.Id),
						resource.TestCheckResourceAttr(accessor, "name", approvalPolicy.Name),
						resource.TestCheckResourceAttr(accessor, "repository", approvalPolicy.Repository),
						resource.TestCheckResourceAttr(accessor, "path", approvalPolicy.Path),
						resource.TestCheckResourceAttr(accessor, "revision", approvalPolicy.Revision),
						resource.TestCheckResourceAttr(accessor, "token_id", approvalPolicy.TokenId),
						resource.TestCheckResourceAttr(accessor, "github_installation_id", strconv.Itoa(approvalPolicy.GithubInstallationId)),
						resource.TestCheckResourceAttr(accessor, "is_github_enterprise", strconv.FormatBool(approvalPolicy.IsGithubEnterprise)),
					),
					ExpectNonEmptyPlan: true,
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":            updatedApprovalPolicy.Name,
						"repository":      updatedApprovalPolicy.Repository,
						"path":            updatedApprovalPolicy.Path,
						"revision":        updatedApprovalPolicy.Revision,
						"token_id":        updatedApprovalPolicy.TokenId,
						"is_azure_devops": "true",
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", updatedApprovalPolicy.Id),
						resource.TestCheckResourceAttr(accessor, "name", updatedApprovalPolicy.Name),
						resource.TestCheckResourceAttr(accessor, "repository", updatedApprovalPolicy.Repository),
						resource.TestCheckResourceAttr(accessor, "path", updatedApprovalPolicy.Path),
						resource.TestCheckResourceAttr(accessor, "revision", updatedApprovalPolicy.Revision),
						resource.TestCheckResourceAttr(accessor, "token_id", updatedApprovalPolicy.TokenId),
						resource.TestCheckResourceAttr(accessor, "github_installation_id", "0"),
						resource.TestCheckResourceAttr(accessor, "is_github_enterprise", "false"),
						resource.TestCheckResourceAttr(accessor, "is_azure_devops", strconv.FormatBool(updatedApprovalPolicy.IsAzureDevOps)),
					),
					PlanOnly:           true,
					ExpectNonEmptyPlan: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ApprovalPolicyCreate(&createPayload).Times(1).Return(&approvalPolicy, nil),
				mock.EXPECT().Template(template.Id).Times(2).Return(client.Template{}, http.NewMockFailedResponseError(404)),
			)
		})
	})

	t.Run("Create Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":                   approvalPolicy.Name,
						"repository":             approvalPolicy.Repository,
						"path":                   approvalPolicy.Path,
						"revision":               approvalPolicy.Revision,
						"token_id":               approvalPolicy.TokenId,
						"github_installation_id": approvalPolicy.GithubInstallationId,
						"is_github_enterprise":   "true",
					}),
					ExpectError: regexp.MustCompile("failed to create approval policy: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().ApprovalPolicyCreate(&createPayload).Times(1).Return(nil, errors.New("error"))
		})
	})

	t.Run("Update Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":                   approvalPolicy.Name,
						"repository":             approvalPolicy.Repository,
						"path":                   approvalPolicy.Path,
						"revision":               approvalPolicy.Revision,
						"token_id":               approvalPolicy.TokenId,
						"github_installation_id": approvalPolicy.GithubInstallationId,
						"is_github_enterprise":   "true",
					}),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":            updatedApprovalPolicy.Name,
						"repository":      updatedApprovalPolicy.Repository,
						"path":            updatedApprovalPolicy.Path,
						"revision":        updatedApprovalPolicy.Revision,
						"token_id":        updatedApprovalPolicy.TokenId,
						"is_azure_devops": "true",
					}),
					ExpectError: regexp.MustCompile("failed to update approval policy: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ApprovalPolicyCreate(&createPayload).Times(1).Return(&approvalPolicy, nil),
				mock.EXPECT().Template(approvalPolicy.Id).Times(2).Return(template, nil),
				mock.EXPECT().ApprovalPolicyUpdate(&updatePayload).Times(1).Return(nil, errors.New("error")),
				mock.EXPECT().TemplateDelete(approvalPolicy.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Import By Id", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":                   approvalPolicy.Name,
						"repository":             approvalPolicy.Repository,
						"path":                   approvalPolicy.Path,
						"revision":               approvalPolicy.Revision,
						"token_id":               approvalPolicy.TokenId,
						"github_installation_id": approvalPolicy.GithubInstallationId,
						"is_github_enterprise":   "true",
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     approvalPolicy.Id,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ApprovalPolicyCreate(&createPayload).Times(1).Return(&approvalPolicy, nil),
				mock.EXPECT().Template(approvalPolicy.Id).Times(3).Return(template, nil),
				mock.EXPECT().TemplateDelete(approvalPolicy.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Import By Id - not approval policy type", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":                   approvalPolicy.Name,
						"repository":             approvalPolicy.Repository,
						"path":                   approvalPolicy.Path,
						"revision":               approvalPolicy.Revision,
						"token_id":               approvalPolicy.TokenId,
						"github_installation_id": approvalPolicy.GithubInstallationId,
						"is_github_enterprise":   "true",
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     approvalPolicy.Id,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile("template type requires type approval-policy but received type terraform"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ApprovalPolicyCreate(&createPayload).Times(1).Return(&approvalPolicy, nil),
				mock.EXPECT().Template(approvalPolicy.Id).Times(1).Return(template, nil),
				mock.EXPECT().Template(approvalPolicy.Id).Times(1).Return(notApprovalPolicyTemplate, nil),
				mock.EXPECT().TemplateDelete(approvalPolicy.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Import By Name", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":                   approvalPolicy.Name,
						"repository":             approvalPolicy.Repository,
						"path":                   approvalPolicy.Path,
						"revision":               approvalPolicy.Revision,
						"token_id":               approvalPolicy.TokenId,
						"github_installation_id": approvalPolicy.GithubInstallationId,
						"is_github_enterprise":   "true",
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     approvalPolicy.Name,
					ImportStateVerify: true,
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ApprovalPolicyCreate(&createPayload).Times(1).Return(&approvalPolicy, nil),
				mock.EXPECT().Template(approvalPolicy.Id).Times(1).Return(template, nil),
				mock.EXPECT().ApprovalPolicies(approvalPolicy.Name).Times(1).Return([]client.ApprovalPolicy{approvalPolicy}, nil),
				mock.EXPECT().Template(approvalPolicy.Id).Times(1).Return(template, nil),
				mock.EXPECT().TemplateDelete(approvalPolicy.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Import By Name - not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":                   approvalPolicy.Name,
						"repository":             approvalPolicy.Repository,
						"path":                   approvalPolicy.Path,
						"revision":               approvalPolicy.Revision,
						"token_id":               approvalPolicy.TokenId,
						"github_installation_id": approvalPolicy.GithubInstallationId,
						"is_github_enterprise":   "true",
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     approvalPolicy.Name,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(fmt.Sprintf("approval policy with name %v not found", approvalPolicy.Name)),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ApprovalPolicyCreate(&createPayload).Times(1).Return(&approvalPolicy, nil),
				mock.EXPECT().Template(approvalPolicy.Id).Times(1).Return(template, nil),
				mock.EXPECT().ApprovalPolicies(approvalPolicy.Name).Times(1).Return([]client.ApprovalPolicy{}, nil),
				mock.EXPECT().TemplateDelete(approvalPolicy.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Import By Name - too many results", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":                   approvalPolicy.Name,
						"repository":             approvalPolicy.Repository,
						"path":                   approvalPolicy.Path,
						"revision":               approvalPolicy.Revision,
						"token_id":               approvalPolicy.TokenId,
						"github_installation_id": approvalPolicy.GithubInstallationId,
						"is_github_enterprise":   "true",
					}),
				},
				{
					ResourceName:      resourceNameImport,
					ImportState:       true,
					ImportStateId:     approvalPolicy.Name,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile("found multiple approval policies with"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ApprovalPolicyCreate(&createPayload).Times(1).Return(&approvalPolicy, nil),
				mock.EXPECT().Template(approvalPolicy.Id).Times(1).Return(template, nil),
				mock.EXPECT().ApprovalPolicies(approvalPolicy.Name).Times(1).Return([]client.ApprovalPolicy{approvalPolicy, approvalPolicy}, nil),
				mock.EXPECT().TemplateDelete(approvalPolicy.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Success with vcs_connection_id", func(t *testing.T) {
		vcsApprovalPolicy := client.ApprovalPolicy{
			Id:              uuid.NewString(),
			Name:            "name",
			Repository:      "repository",
			Path:            "path",
			Revision:        "revision",
			VcsConnectionId: "vcs-conn-123",
		}

		var vcsTemplate client.Template

		require.NoError(t, copier.Copy(&vcsTemplate, &vcsApprovalPolicy))
		vcsTemplate.Type = string(ApprovalPolicy)

		vcsCreatePayload := client.ApprovalPolicyCreatePayload{
			Name:            vcsApprovalPolicy.Name,
			Repository:      vcsApprovalPolicy.Repository,
			Path:            vcsApprovalPolicy.Path,
			Revision:        vcsApprovalPolicy.Revision,
			VcsConnectionId: vcsApprovalPolicy.VcsConnectionId,
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":              vcsApprovalPolicy.Name,
						"repository":        vcsApprovalPolicy.Repository,
						"path":              vcsApprovalPolicy.Path,
						"revision":          vcsApprovalPolicy.Revision,
						"vcs_connection_id": vcsApprovalPolicy.VcsConnectionId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", vcsApprovalPolicy.Id),
						resource.TestCheckResourceAttr(accessor, "name", vcsApprovalPolicy.Name),
						resource.TestCheckResourceAttr(accessor, "vcs_connection_id", vcsApprovalPolicy.VcsConnectionId),
						resource.TestCheckNoResourceAttr(accessor, "github_installation_id"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ApprovalPolicyCreate(&vcsCreatePayload).Times(1).Return(&vcsApprovalPolicy, nil),
				mock.EXPECT().Template(vcsApprovalPolicy.Id).Times(1).Return(vcsTemplate, nil),
				mock.EXPECT().TemplateDelete(vcsApprovalPolicy.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("vcs_connection_id ignores github_installation_id from backend to avoid drift", func(t *testing.T) {
		vcsApprovalPolicy := client.ApprovalPolicy{
			Id:              uuid.NewString(),
			Name:            "name",
			Repository:      "repository",
			VcsConnectionId: "vcs-conn-123",
		}

		vcsApprovalPolicyFromBackend := client.ApprovalPolicy{
			Id:                   vcsApprovalPolicy.Id,
			Name:                 vcsApprovalPolicy.Name,
			Repository:           vcsApprovalPolicy.Repository,
			VcsConnectionId:      vcsApprovalPolicy.VcsConnectionId,
			GithubInstallationId: 456,
		}

		var vcsTemplate client.Template

		require.NoError(t, copier.Copy(&vcsTemplate, &vcsApprovalPolicyFromBackend))
		vcsTemplate.Type = string(ApprovalPolicy)

		vcsCreatePayload := client.ApprovalPolicyCreatePayload{
			Name:            vcsApprovalPolicy.Name,
			Repository:      vcsApprovalPolicy.Repository,
			VcsConnectionId: vcsApprovalPolicy.VcsConnectionId,
		}

		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":              vcsApprovalPolicy.Name,
						"repository":        vcsApprovalPolicy.Repository,
						"vcs_connection_id": vcsApprovalPolicy.VcsConnectionId,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", vcsApprovalPolicy.Id),
						resource.TestCheckResourceAttr(accessor, "name", vcsApprovalPolicy.Name),
						resource.TestCheckResourceAttr(accessor, "vcs_connection_id", vcsApprovalPolicy.VcsConnectionId),
						resource.TestCheckNoResourceAttr(accessor, "github_installation_id"),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().ApprovalPolicyCreate(&vcsCreatePayload).Times(1).Return(&vcsApprovalPolicy, nil),
				mock.EXPECT().Template(vcsApprovalPolicy.Id).Times(1).Return(vcsTemplate, nil),
				mock.EXPECT().TemplateDelete(vcsApprovalPolicy.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("vcs_connection_id and github_installation_id are mutually exclusive", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"name":                   "test",
						"repository":             "env0/repo",
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
