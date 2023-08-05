package env0

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jinzhu/copier"
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
	copier.Copy(&template, &approvalPolicy)
	template.Type = APPROVAL_POLICY

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
	copier.Copy(&updatedTemplate, &updatedApprovalPolicy)
	template.Type = APPROVAL_POLICY

	createPayload := client.TemplateCreatePayload{
		Name:                 approvalPolicy.Name,
		Repository:           approvalPolicy.Repository,
		Path:                 approvalPolicy.Path,
		Revision:             approvalPolicy.Revision,
		TokenId:              approvalPolicy.TokenId,
		GithubInstallationId: approvalPolicy.GithubInstallationId,
		IsGithubEnterprise:   approvalPolicy.IsGithubEnterprise,
		Type:                 APPROVAL_POLICY,
	}

	updatePayload := client.TemplateCreatePayload{
		Name:          updatedApprovalPolicy.Name,
		Repository:    updatedApprovalPolicy.Repository,
		Path:          updatedApprovalPolicy.Path,
		Revision:      updatedApprovalPolicy.Revision,
		TokenId:       updatedApprovalPolicy.TokenId,
		IsAzureDevOps: updatedApprovalPolicy.IsAzureDevOps,
		Type:          APPROVAL_POLICY,
	}

	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
				mock.EXPECT().TemplateCreate(createPayload).Times(1).Return(template, nil),
				mock.EXPECT().Template(approvalPolicy.Id).Times(2).Return(template, nil),
				mock.EXPECT().TemplateUpdate(approvalPolicy.Id, updatePayload).Times(1).Return(updatedTemplate, nil),
				mock.EXPECT().Template(approvalPolicy.Id).Times(1).Return(updatedTemplate, nil),
				mock.EXPECT().TemplateDelete(approvalPolicy.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Drift detected - deleted", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
				mock.EXPECT().TemplateCreate(createPayload).Times(1).Return(template, nil),
				mock.EXPECT().Template(template.Id).Times(2).Return(deletedTemplate, nil),
			)
		})
	})

	t.Run("Drift detected - not found", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
				mock.EXPECT().TemplateCreate(createPayload).Times(1).Return(template, nil),
				mock.EXPECT().Template(template.Id).Times(2).Return(client.Template{}, http.NewMockFailedResponseError(404)),
			)
		})
	})

	t.Run("Create Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":                   approvalPolicy.Name,
						"repository":             approvalPolicy.Repository,
						"path":                   approvalPolicy.Path,
						"revision":               approvalPolicy.Revision,
						"token_id":               approvalPolicy.TokenId,
						"github_installation_id": approvalPolicy.GithubInstallationId,
						"is_github_enterprise":   "true",
					}),
					ExpectError: regexp.MustCompile("could not create approval policy template: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().TemplateCreate(createPayload).Times(1).Return(client.Template{}, errors.New("error"))
		})
	})

	t.Run("Update Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"name":            updatedApprovalPolicy.Name,
						"repository":      updatedApprovalPolicy.Repository,
						"path":            updatedApprovalPolicy.Path,
						"revision":        updatedApprovalPolicy.Revision,
						"token_id":        updatedApprovalPolicy.TokenId,
						"is_azure_devops": "true",
					}),
					ExpectError: regexp.MustCompile("could not update approval policy template: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().TemplateCreate(createPayload).Times(1).Return(template, nil),
				mock.EXPECT().Template(approvalPolicy.Id).Times(2).Return(template, nil),
				mock.EXPECT().TemplateUpdate(approvalPolicy.Id, updatePayload).Times(1).Return(client.Template{}, errors.New("error")),
				mock.EXPECT().TemplateDelete(approvalPolicy.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Import By Id", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
				mock.EXPECT().TemplateCreate(createPayload).Times(1).Return(template, nil),
				mock.EXPECT().Template(approvalPolicy.Id).Times(3).Return(template, nil),
				mock.EXPECT().TemplateDelete(approvalPolicy.Id).Times(1).Return(nil),
			)
		})
	})

	t.Run("Import By Id - not approval policy type", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
				mock.EXPECT().TemplateCreate(createPayload).Times(1).Return(template, nil),
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
				mock.EXPECT().TemplateCreate(createPayload).Times(1).Return(template, nil),
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
				mock.EXPECT().TemplateCreate(createPayload).Times(1).Return(template, nil),
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
				mock.EXPECT().TemplateCreate(createPayload).Times(1).Return(template, nil),
				mock.EXPECT().Template(approvalPolicy.Id).Times(1).Return(template, nil),
				mock.EXPECT().ApprovalPolicies(approvalPolicy.Name).Times(1).Return([]client.ApprovalPolicy{approvalPolicy, approvalPolicy}, nil),
				mock.EXPECT().TemplateDelete(approvalPolicy.Id).Times(1).Return(nil),
			)
		})
	})
}
