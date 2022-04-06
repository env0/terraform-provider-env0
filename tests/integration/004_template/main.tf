provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

# Github Integration must be done manually - so we expect an existing Github Template with this name -
# It must be for https://github.com/env0/templates - We validate that in the outputs
data "env0_template" "github_template" {
  name = "Github Integrated Template"
}

# Gitlab Integration must be done manually - so we expect an existing Gitlab Template with this name
# It must be for https://gitlab.com/env0/gitlab-vcs-integration-tests - the gitlab_project_id is still static
data "env0_template" "gitlab_template" {
  name = "Gitlab Integrated Template"
}

resource "env0_template" "tested1" {
  name                                    = "tested1-${random_string.random.result}"
  description                             = "Tested 1 description"
  type                                    = "terraform"
  repository                              = data.env0_template.github_template.repository
  github_installation_id                  = data.env0_template.github_template.github_installation_id
  path                                    = var.second_run ? "second" : "misc/null-resource"
  retries_on_deploy                       = 3
  retry_on_deploy_only_when_matches_regex = "abc"
  retries_on_destroy                      = 1
  terraform_version                       = "0.15.1"
}

resource "env0_template" "tested2" {
  name                                    = "GitLab Test-${random_string.random.result}"
  description                             = "Tested 2 description - Gitlab"
  type                                    = "terraform"
  repository                              = data.env0_template.gitlab_template.repository
  token_id                                = data.env0_template.gitlab_template.token_id
  gitlab_project_id                       = 32315446
  path                                    = var.second_run ? "second" : "misc/null-resource"
  retries_on_deploy                       = 3
  retry_on_deploy_only_when_matches_regex = "abc"
  retries_on_destroy                      = 1
  terraform_version                       = "0.15.1"
}

resource "env0_template" "template_tg" {
  name               = "Template for environment resource - tg-${random_string.random.result}"
  type               = "terragrunt"
  repository         = "https://github.com/env0/templates"
  path               = "terragrunt/misc/null-resource"
  terraform_version  = "0.15.1"
  terragrunt_version = "0.35.0"
}

resource "env0_configuration_variable" "in_a_template" {
  name        = "fake_key"
  value       = "fake value"
  template_id = env0_template.tested1.id
}

resource "env0_configuration_variable" "in_a_template2" {
  name        = "fake_key_2"
  value       = "fake value 2"
  template_id = env0_template.tested1.id
  type        = "terraform"
}

output "tested2_template_id" {
  value = env0_template.tested2.id
}
output "tested2_template_type" {
  value = env0_template.tested2.type
}
output "tested2_template_name" {
  value = replace(env0_template.tested2.name, random_string.random.result, "")
}
output "tested2_template_repository" {
  value = env0_template.tested2.repository
}
output "tested1_template_repository" {
  value = env0_template.tested1.repository
}
output "tested2_template_path" {
  value = env0_template.tested2.path
}
output "tg_tg_version" {
  value = env0_template.template_tg.terragrunt_version
}

output "github_template_type" {
  value = data.env0_template.github_template.type
}
