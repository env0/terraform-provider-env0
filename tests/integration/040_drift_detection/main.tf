resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_project" "test_project" {
  name          = "Test-Project-for-drift-detection-${random_string.random.result}"
  wait          = true
  force_destroy = true
}

resource "env0_template" "template" {
  name              = "Template for drift detection ${random_string.random.result}"
  type              = "terraform"
  repository        = "https://github.com/env0/templates"
  path              = "misc/null-resource"
  terraform_version = "0.15.1"
}

resource "env0_template_project_assignment" "assignment" {
  template_id = env0_template.template.id
  project_id  = env0_project.test_project.id
}

# Environment 1: DISABLED drift remediation
resource "env0_environment" "environment_disabled" {
  depends_on                 = [env0_template_project_assignment.assignment]
  force_destroy              = true
  name                       = "drift-disabled-${random_string.random.result}"
  project_id                 = env0_project.test_project.id
  template_id                = env0_template.template.id
  approve_plan_automatically = true
}

resource "env0_environment_drift_detection" "drift_disabled" {
  environment_id         = env0_environment.environment_disabled.id
  cron                   = "0 2 * * *"
  auto_drift_remediation = "DISABLED"
}

# Environment 2: CODE_TO_CLOUD drift remediation
resource "env0_environment" "environment_code_to_cloud" {
  depends_on                 = [env0_template_project_assignment.assignment]
  force_destroy              = true
  name                       = "drift-code-to-cloud-${random_string.random.result}"
  project_id                 = env0_project.test_project.id
  template_id                = env0_template.template.id
  approve_plan_automatically = true
}

resource "env0_environment_drift_detection" "drift_code_to_cloud" {
  environment_id         = env0_environment.environment_code_to_cloud.id
  cron                   = "0 3 * * *"
  auto_drift_remediation = "CODE_TO_CLOUD"
}

# Environment 3: CLOUD_TO_CODE drift remediation
resource "env0_environment" "environment_cloud_to_code" {
  depends_on                 = [env0_template_project_assignment.assignment]
  force_destroy              = true
  name                       = "drift-cloud-to-code-${random_string.random.result}"
  project_id                 = env0_project.test_project.id
  template_id                = env0_template.template.id
  approve_plan_automatically = true
}

resource "env0_environment_drift_detection" "drift_cloud_to_code" {
  environment_id         = env0_environment.environment_cloud_to_code.id
  cron                   = "0 4 * * *"
  auto_drift_remediation = "CLOUD_TO_CODE"
}

# Environment 4: SMART_REMEDIATION drift remediation
resource "env0_environment" "environment_smart_remediation" {
  depends_on                 = [env0_template_project_assignment.assignment]
  force_destroy              = true
  name                       = "drift-smart-remediation-${random_string.random.result}"
  project_id                 = env0_project.test_project.id
  template_id                = env0_template.template.id
  approve_plan_automatically = true
}

resource "env0_environment_drift_detection" "drift_smart_remediation" {
  environment_id         = env0_environment.environment_smart_remediation.id
  cron                   = "0 5 * * *"
  auto_drift_remediation = "SMART_REMEDIATION"
}

# Outputs to verify drift remediation types
output "drift_disabled_type" {
  value = env0_environment_drift_detection.drift_disabled.auto_drift_remediation
}

output "drift_code_to_cloud_type" {
  value = env0_environment_drift_detection.drift_code_to_cloud.auto_drift_remediation
}

output "drift_cloud_to_code_type" {
  value = env0_environment_drift_detection.drift_cloud_to_code.auto_drift_remediation
}

output "drift_smart_remediation_type" {
  value = env0_environment_drift_detection.drift_smart_remediation.auto_drift_remediation
}