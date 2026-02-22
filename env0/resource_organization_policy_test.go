package env0

import (
	"errors"
	"regexp"
	"strconv"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"go.uber.org/mock/gomock"
)

func TestUnitOrganizationPolicyResource(t *testing.T) {
	resourceType := "env0_organization_policy"
	resourceName := "test"
	organizationId := "org"
	accessor := resourceAccessor(resourceType, resourceName)

	defaultOrganization := client.Organization{
		Id:   organizationId,
		Name: "name",
	}

	organization := client.Organization{
		Id:                                  organizationId,
		Name:                                "name",
		MaxTtl:                              new("4-d"),
		DefaultTtl:                          new("13-h"),
		DoNotReportSkippedStatusChecks:      false,
		DoNotConsiderMergeCommitsForPrPlans: true,
		EnableOidc:                          false,
		EnforcePrCommenterPermissions:       false,
		AllowMergeableBypassForPrApply:      false,
	}

	organizationUpdated := client.Organization{
		Id:                                  organizationId,
		Name:                                "name",
		DefaultTtl:                          new("2-M"),
		DoNotReportSkippedStatusChecks:      true,
		DoNotConsiderMergeCommitsForPrPlans: false,
		EnableOidc:                          true,
		EnforcePrCommenterPermissions:       true,
		AllowMergeableBypassForPrApply:      true,
	}

	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"max_ttl":     *organization.MaxTtl,
						"default_ttl": *organization.DefaultTtl,
						"do_not_consider_merge_commits_for_pr_plans": organization.DoNotConsiderMergeCommitsForPrPlans,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", organization.Id),
						resource.TestCheckResourceAttr(accessor, "max_ttl", *organization.MaxTtl),
						resource.TestCheckResourceAttr(accessor, "default_ttl", *organization.DefaultTtl),
						resource.TestCheckResourceAttr(accessor, "do_not_report_skipped_status_checks", strconv.FormatBool(organization.DoNotReportSkippedStatusChecks)),
						resource.TestCheckResourceAttr(accessor, "do_not_consider_merge_commits_for_pr_plans", strconv.FormatBool(organization.DoNotConsiderMergeCommitsForPrPlans)),
						resource.TestCheckResourceAttr(accessor, "enable_oidc", strconv.FormatBool(organization.EnableOidc)),
						resource.TestCheckResourceAttr(accessor, "enforce_pr_commenter_permissions", strconv.FormatBool(organization.EnforcePrCommenterPermissions)),
						resource.TestCheckResourceAttr(accessor, "allow_mergeable_bypass_for_pr_apply", strconv.FormatBool(organization.AllowMergeableBypassForPrApply)),
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"default_ttl":                         *organizationUpdated.DefaultTtl,
						"do_not_report_skipped_status_checks": organizationUpdated.DoNotReportSkippedStatusChecks,
						"enable_oidc":                         organizationUpdated.EnableOidc,
						"enforce_pr_commenter_permissions":    organizationUpdated.EnforcePrCommenterPermissions,
						"allow_mergeable_bypass_for_pr_apply": organizationUpdated.AllowMergeableBypassForPrApply,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", organization.Id),
						resource.TestCheckResourceAttr(accessor, "default_ttl", *organizationUpdated.DefaultTtl),
						resource.TestCheckResourceAttr(accessor, "do_not_report_skipped_status_checks", strconv.FormatBool(organizationUpdated.DoNotReportSkippedStatusChecks)),
						resource.TestCheckResourceAttr(accessor, "do_not_consider_merge_commits_for_pr_plans", strconv.FormatBool(organizationUpdated.DoNotConsiderMergeCommitsForPrPlans)),
						resource.TestCheckResourceAttr(accessor, "enable_oidc", strconv.FormatBool(organizationUpdated.EnableOidc)),
						resource.TestCheckResourceAttr(accessor, "enforce_pr_commenter_permissions", strconv.FormatBool(organizationUpdated.EnforcePrCommenterPermissions)),
						resource.TestCheckResourceAttr(accessor, "allow_mergeable_bypass_for_pr_apply", strconv.FormatBool(organizationUpdated.AllowMergeableBypassForPrApply)),
					),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			gomock.InOrder(
				mock.EXPECT().OrganizationPolicyUpdate(client.OrganizationPolicyUpdatePayload{
					MaxTtl:                              organization.MaxTtl,
					DefaultTtl:                          organization.DefaultTtl,
					DoNotConsiderMergeCommitsForPrPlans: &organization.DoNotConsiderMergeCommitsForPrPlans,
					DoNotReportSkippedStatusChecks:      new(false),
					EnableOidc:                          new(false),
					EnforcePrCommenterPermissions:       new(false),
					AllowMergeableBypassForPrApply:      new(false),
				}).Times(1).Return(&organization, nil),
				mock.EXPECT().Organization().Times(2).Return(organization, nil),
				mock.EXPECT().OrganizationPolicyUpdate(client.OrganizationPolicyUpdatePayload{
					DefaultTtl:                          organizationUpdated.DefaultTtl,
					DoNotReportSkippedStatusChecks:      &organizationUpdated.DoNotReportSkippedStatusChecks,
					EnableOidc:                          &organizationUpdated.EnableOidc,
					EnforcePrCommenterPermissions:       &organizationUpdated.EnforcePrCommenterPermissions,
					AllowMergeableBypassForPrApply:      &organizationUpdated.AllowMergeableBypassForPrApply,
					DoNotConsiderMergeCommitsForPrPlans: new(false),
					MaxTtl:                              new(""),
				}).Times(1).Return(&organizationUpdated, nil),
				mock.EXPECT().Organization().Times(1).Return(organizationUpdated, nil),
				mock.EXPECT().OrganizationPolicyUpdate(client.OrganizationPolicyUpdatePayload{}).Times(1).Return(&defaultOrganization, nil),
			)
		})
	})

	t.Run("Create Failure - max smaller than default", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"max_ttl":     "23-h",
						"default_ttl": "1-d",
					}),
					ExpectError: regexp.MustCompile("default ttl must not be larger than max ttl"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {})
	})

	t.Run("Create/Update Failure", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]any{
						"max_ttl":     *organization.MaxTtl,
						"default_ttl": *organization.DefaultTtl,
						"do_not_consider_merge_commits_for_pr_plans": organization.DoNotConsiderMergeCommitsForPrPlans,
					}),
					ExpectError: regexp.MustCompile("could not update organization policy: error"),
				},
			},
		}

		runUnitTest(t, testCase, func(mock *client.MockApiClientInterface) {
			mock.EXPECT().OrganizationPolicyUpdate(client.OrganizationPolicyUpdatePayload{
				MaxTtl:                              organization.MaxTtl,
				DefaultTtl:                          organization.DefaultTtl,
				DoNotConsiderMergeCommitsForPrPlans: &organization.DoNotConsiderMergeCommitsForPrPlans,
				DoNotReportSkippedStatusChecks:      new(false),
				EnableOidc:                          new(false),
				EnforcePrCommenterPermissions:       new(false),
				AllowMergeableBypassForPrApply:      new(false),
			}).Times(1).Return(nil, errors.New("error"))
		})
	})
}
