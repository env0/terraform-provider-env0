resource "env0_organization_policy" "policy_example" {
  max_ttl                                    = "1-M"
  default_ttl                                = "12-h"
  do_not_consider_merge_commits_for_pr_plans = true
}
