package env0

import (
	"errors"
	"regexp"
	"strconv"
	"testing"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
		MaxTtl:                              stringPtr("3-d"),
		DefaultTtl:                          stringPtr("12-h"),
		DoNotReportSkippedStatusChecks:      false,
		DoNotConsiderMergeCommitsForPrPlans: true,
	}

	organizationUpdated := client.Organization{
		Id:                                  organizationId,
		Name:                                "name",
		DefaultTtl:                          stringPtr("1-M"),
		DoNotReportSkippedStatusChecks:      true,
		DoNotConsiderMergeCommitsForPrPlans: false,
	}

	t.Run("Success", func(t *testing.T) {
		testCase := resource.TestCase{
			Steps: []resource.TestStep{
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
					),
				},
				{
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"default_ttl":                         *organizationUpdated.DefaultTtl,
						"do_not_report_skipped_status_checks": organizationUpdated.DoNotReportSkippedStatusChecks,
					}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(accessor, "id", organization.Id),
						resource.TestCheckResourceAttr(accessor, "default_ttl", *organizationUpdated.DefaultTtl),
						resource.TestCheckResourceAttr(accessor, "do_not_report_skipped_status_checks", strconv.FormatBool(organizationUpdated.DoNotReportSkippedStatusChecks)),
						resource.TestCheckResourceAttr(accessor, "do_not_consider_merge_commits_for_pr_plans", strconv.FormatBool(organizationUpdated.DoNotConsiderMergeCommitsForPrPlans)),
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
				}).Times(1).Return(&organization, nil),
				mock.EXPECT().Organization().Times(2).Return(organization, nil),
				mock.EXPECT().OrganizationPolicyUpdate(client.OrganizationPolicyUpdatePayload{
					DefaultTtl:                     organizationUpdated.DefaultTtl,
					DoNotReportSkippedStatusChecks: &organizationUpdated.DoNotReportSkippedStatusChecks,
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
						"max_ttl":     "12-h",
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
					Config: resourceConfigCreate(resourceType, resourceName, map[string]interface{}{
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
			}).Times(1).Return(nil, errors.New("error"))
		})
	})
}
