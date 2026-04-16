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

# Discovery by file - uses repository_regex to match repos by '<owner>/<repo>' pattern.
# Use '|' to separate multiple patterns.
# Examples:
#   "my-org/.*"                     - all repos in my-org
#   "my-org/.*|other-org/.*"        - all repos in both orgs
#   "my-org/web-.*|my-org/api-.*"   - repos starting with 'web-' or 'api-' in my-org
#   "my-org/my-repo"                - a single specific repo
resource "env0_environment_discovery_configuration" "discovery_file_example" {
  project_id       = data.env0_project.project.id
  repository_regex = "my-org/.*|other-org/web-.*|third-org/specific-repo"
}
