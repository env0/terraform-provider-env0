resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_project" "test_project" {
  name = "Test-Project-for-environment-scheduling-${random_string.random.result}"
  wait = true
}

resource "env0_template" "template" {
  name              = "Template for environment resource ${random_string.random.result}"
  type              = "terraform"
  repository        = "https://github.com/env0/templates"
  path              = "misc/null-resource"
  terraform_version = "0.15.1"
}

resource "env0_template_project_assignment" "assignment" {
  template_id = env0_template.template.id
  project_id  = env0_project.test_project.id
}

resource "env0_environment" "environment" {
  depends_on                 = [env0_template_project_assignment.assignment]
  force_destroy              = true
  name                       = "the_trigger-${random_string.random.result}"
  project_id                 = env0_project.test_project.id
  template_id                = env0_template.template.id
  approve_plan_automatically = true
}

resource "env0_environment_scheduling" "scheduling" {
  environment_id = env0_environment.environment.id
  deploy_cron    = "15 * * * *"
  destroy_cron   = var.second_run ? null : "10 * * * *"
}
