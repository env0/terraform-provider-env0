provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

data "env0_template" "github_template" {
  name = "Github Integrated Template"
}

resource "env0_template" "github_template1" {
  name                                    = "Github Test ${random_string.random.result}-1"
  description                             = "Template description - GitHub"
  type                                    = "terraform"
  repository                              = data.env0_template.github_template.repository
  github_installation_id                  = data.env0_template.github_template.github_installation_id
  path                                    = "misc/null-resource"
  retries_on_deploy                       = 3
  retry_on_deploy_only_when_matches_regex = "abc"
  retries_on_destroy                      = 1
  terraform_version                       = "0.15.1"
}

resource "env0_template" "github_template2" {
  name                                    = "Github Test ${random_string.random.result}-2"
  description                             = "Template description - GitHub"
  type                                    = "terraform"
  repository                              = data.env0_template.github_template.repository
  github_installation_id                  = data.env0_template.github_template.github_installation_id
  path                                    = "misc/null-resource"
  retries_on_deploy                       = 3
  retry_on_deploy_only_when_matches_regex = "abc"
  retries_on_destroy                      = 1
  terraform_version                       = "0.15.1"
}

data "env0_templates" "all_templates" {}

# This is removed temporarily until https://github.com/env0/terraform-provider-env0/issues/350 is fixed
#data "env0_template" "templates" {
#  for_each = toset(data.env0_templates.all_templates.names)
#  name     = each.value
#}
