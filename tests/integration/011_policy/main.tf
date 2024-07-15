resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_project" "test_project" {
  name = "Test-Project-For_policy-${random_string.random.result}"
}

resource "env0_project_policy" "test_policy" {
  project_id                    = env0_project.test_project.id
  number_of_environments        = var.second_run ? null : 1
  number_of_environments_total  = var.second_run ? null : 1
  requires_approval_default     = true
  include_cost_estimation       = false
  skip_apply_when_plan_is_empty = false
  disable_destroy_environments  = false
  skip_redundant_deployments    = false
  drift_detection_cron          = var.second_run ? "0 4 * * *" : "0 3 * * *"

}

resource "env0_project_policy" "test_policy_2" {
  project_id                      = env0_project.test_project.id
  number_of_environments          = 1
  number_of_environments_total    = 1
  requires_approval_default       = true
  include_cost_estimation         = true
  skip_apply_when_plan_is_empty   = true
  disable_destroy_environments    = true
  skip_redundant_deployments      = true
  vcs_pr_comments_enabled_default = true
}

resource "env0_project_policy" "test_policy_ttl" {
  project_id                    = env0_project.test_project.id
  number_of_environments        = 1
  number_of_environments_total  = 1
  requires_approval_default     = true
  include_cost_estimation       = true
  skip_apply_when_plan_is_empty = true
  disable_destroy_environments  = true
  skip_redundant_deployments    = true
  max_ttl                       = "4-d"
  default_ttl                   = "14-h"
}

resource "env0_project_policy" "test_policy_infinite" {
  project_id                    = env0_project.test_project.id
  number_of_environments        = 1
  number_of_environments_total  = 1
  requires_approval_default     = true
  include_cost_estimation       = true
  skip_apply_when_plan_is_empty = true
  disable_destroy_environments  = true
  skip_redundant_deployments    = true
  max_ttl                       = "Infinite"
  default_ttl                   = var.second_run ? "4-d" : "Infinite"
}
