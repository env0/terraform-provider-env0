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
  number_of_environments        = 1
  number_of_environments_total  = 1
  requires_approval_default     = true
  include_cost_estimation       = false
  skip_apply_when_plan_is_empty = false
  disable_destroy_environments  = false
  skip_redundant_deployments    = false
}

resource "env0_project_policy" "test_policy_2" {
  project_id                    = env0_project.test_project.id
  number_of_environments        = 1
  number_of_environments_total  = 1
  requires_approval_default     = true
  include_cost_estimation       = true
  skip_apply_when_plan_is_empty = true
  disable_destroy_environments  = true
  skip_redundant_deployments    = true
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
  max_ttl                       = "3-d"
  default_ttl                   = "12-h"
}
