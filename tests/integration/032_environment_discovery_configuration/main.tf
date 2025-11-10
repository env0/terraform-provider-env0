provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_project" "project" {
  name = "project-${random_string.random.result}"
  wait = true
}

data "env0_template" "github_template" {
  name = "Github Integrated Template"
}

resource "env0_template_project_assignment" "assignment" {
  template_id = data.env0_template.github_template.id
  project_id  = env0_project.project.id
}

resource "env0_environment_discovery_configuration" "example" {
  project_id             = env0_project.project.id
  glob_pattern           = var.second_run ? "**" : "**/**"
  opentofu_version       = "1.6.2"
  repository             = data.env0_template.github_template.repository
  github_installation_id = data.env0_template.github_template.github_installation_id
  create_new_environments_from_pull_requests = true

  depends_on = [env0_template_project_assignment.assignment]
}
