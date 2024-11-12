provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_project" "project" {
  name          = "project-environment-without_template-${random_string.random.result}"
  force_destroy = true
}

resource "env0_environment" "environment" {
  name                      = "environment-without_template-${random_string.random.result}"
  project_id                = env0_project.project.id
  is_remote_backend         = false
  deploy_on_push            = true
  run_plan_on_pull_requests = true
  vcs_pr_comments_enabled   = true
  force_destroy             = true

  without_template_settings {
    description           = "Core Azure"
    repository            = "https://github.com/env0/templates"
    path                  = "terragrunt/misc/null-resource"
    opentofu_version      = "1.8.1"
    terragrunt_version    = "0.64.5"
    type                  = "terragrunt"
    revision              = "master"
    is_terragrunt_run_all = true
    terragrunt_tf_binary  = "opentofu"
  }
}