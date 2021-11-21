data "env0_project" "default" {
  name = "Default Organization Project"
}

resource "env0_project_policy" "test_policy" {
  project_id                    = data.env0_project.default.id
  number_of_environments        = 1
  number_of_environments_total  = 1
  requires_approval_default     = true
  include_cost_estimation       = false
  skip_apply_when_plan_is_empty = false
  disable_destroy_environments  = false
  skip_redundant_deployments    = false
}

resource "env0_project_policy" "test_policy_2" {
  project_id                    = data.env0_project.default.id
  number_of_environments        = 1
  number_of_environments_total  = 1
  requires_approval_default     = true
  include_cost_estimation       = true
  skip_apply_when_plan_is_empty = true
  disable_destroy_environments  = true
  skip_redundant_deployments    = true
}
