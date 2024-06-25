resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_project" "test_project" {
  name        = "Test-Project-${random_string.random.result}"
  description = "Test Description ${var.second_run ? "after update" : ""}"
  wait        = true
  force_destroy = true
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

resource "env0_environment" "the_trigger" {
  depends_on                 = [env0_template_project_assignment.assignment]
  force_destroy              = true
  name                       = "the_trigger-${random_string.random.result}"
  project_id                 = env0_project.test_project.id
  template_id                = env0_template.template.id
  approve_plan_automatically = true
}

resource "env0_environment" "downstream_environment" {
  depends_on                 = [env0_template_project_assignment.assignment]
  force_destroy              = true
  name                       = "downstream_environment-${random_string.random.result}"
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

resource "env0_environment" "the_trigger_2" {
  depends_on                 = [env0_template_project_assignment.assignment]
  force_destroy              = true
  name                       = "the_trigger-${random_string.random.result}-2"
  project_id                 = env0_project.test_project.id
  template_id                = env0_template.template.id
  approve_plan_automatically = true
}

resource "env0_environment" "downstream_environment_2" {
  depends_on                 = [env0_template_project_assignment.assignment]
  force_destroy              = true
  name                       = "downstream_environment-${random_string.random.result}-2"
  project_id                 = env0_project.test_project.id
  template_id                = env0_template.template.id
  approve_plan_automatically = true
}

resource "env0_workflow_trigger" "trigger_link_2" {
  environment_id            = env0_environment.the_trigger_2.id
  downstream_environment_id = env0_environment.downstream_environment_2.id
}
