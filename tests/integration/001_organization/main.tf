data "env0_organization" "my_organization" {}

output "organization_name" {
  value = data.env0_organization.my_organization.name
}

resource "env0_organization_policy" "my_organization_policy" {
  max_ttl                                    = "2-M"
  default_ttl                                = var.second_run ? "7-h" : "13-h"
  do_not_consider_merge_commits_for_pr_plans = var.second_run ? false : true
  enable_oidc                                = var.second_run ? false : true
}
