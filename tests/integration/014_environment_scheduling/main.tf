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

resource "env0_environment" "environment" {
  force_destroy              = true
  name                       = "the_trigger"
  project_id                 = data.env0_project.default.id
  template_id                = env0_template.template.id
  approve_plan_automatically = true
}

resource "env0_environment_scheduling" "scheduling" {
  environment_id = env0_environment.environment.id
  deploy_cron = "5 * * * *"
  destroy_cron = "10 * * * *"
}
