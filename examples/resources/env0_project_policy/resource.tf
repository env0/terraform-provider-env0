data "env0_project" "default_project" {
  name = "Default Organization Project"
}

resource "env0_project_policy" "example" {
  project_id                    = data.env0_project.default_project.id
  number_of_environments        = 1
  number_of_environments_total  = 1
  requires_approval_default     = true
  include_cost_estimation       = true
  skip_apply_when_plan_is_empty = true
  disable_destroy_environments  = true
  skip_redundant_deployments    = true
  drift_detection_cron          = "0 4 * * *"  # Run drift detection daily at 4 AM
  auto_drift_remediation        = "CODE_TO_CLOUD"  # Optional, defaults to "DISABLED"
}