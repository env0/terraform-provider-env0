resource "env0_project" "test_project" {
  name        = "Test-Project"
  description = "Test Description ${var.second_run ? "after update" : ""}"
}

resource "env0_template" "template" {
  name              = "Template for environment resource"
  type              = "terraform"
  repository        = "https://github.com/env0/templates"
  path              = "misc/null-resource"
  terraform_version = "0.15.1"
}

resource "env0_environment" "the_trigger" {
  force_destroy              = true
  name                       = "the_trigger"
  project_id                 = env0_project.test_project.id
  template_id                = env0_template.template.id
  approve_plan_automatically = true
}

resource "env0_environment" "downstream_environment" {
  force_destroy              = true
  name                       = "downstream_environment"
  project_id                 = env0_project.test_project.id
  template_id                = env0_template.template.id
  approve_plan_automatically = true
}

resource "env0_workflow_triggers" "trigger_link" {
  environment_id = env0_environment.the_trigger.id
  downstream_environment_ids = [
    env0_environment.downstream_environment.id
  ]
}
