data "env0_project" "default" {
  name = "Default Organization Project"
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
  project_id                 = data.env0_project.default.id
  template_id                = env0_template.template.id
  approve_plan_automatically = true
}

resource "env0_environment" "downstream_environment" {
  force_destroy              = true
  name                       = "downstream_environment"
  project_id                 = data.env0_project.default.id
  template_id                = env0_template.template.id
  approve_plan_automatically = true
}

resource "env0_workflow_trigger" "trigger_link" {
  environment_id            = env0_environment.the_trigger.id
  downstream_environment_id = env0_environment.downstream_environment.id
}
