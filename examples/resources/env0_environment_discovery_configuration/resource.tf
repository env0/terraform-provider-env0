data "env0_project" "project" {
  name = "existing-project"
}

resource "env0_project" "new_project" {
  name = "new-project"
}

resource "env0_environment_discovery_configuration" "example" {
  project_id             = data.env0_project.project.id
  glob_pattern           = "**"
  repository             = "https://github.com/env0/templates"
  opentofu_version       = "1.6.7"
  github_installation_id = 12345
}

resource "env0_environment_discovery_configuration" "terragrunt_example" {
  project_id             = env0_project.new_project.id
  glob_pattern           = "**"
  repository             = "https://github.com/env0/blueprints"
  type                   = "terragrunt"
  terraform_version      = "1.7.1"
  terragrunt_version     = "0.67.4"
  terragrunt_tf_binary   = "terraform"
  github_installation_id = 12345
}

resource "env0_environment_discovery_configuration" "discovery_file_example" {
  project_id       = data.env0_project.project.id
  repository_regex = "env0-example/.*|acme-corp/web.*|company/web-frontend"
}
