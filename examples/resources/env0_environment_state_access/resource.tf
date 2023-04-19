data "env0_environment" "environment" {
  name = "Environment Name"
}

data "env0_project" "project" {
  name = "Project Name"
}

resource "env0_environment_state_access" "example_allowed_projects" {
  environment_id      = data.env0_environment.environment.id
  allowed_project_ids = [data.env0_project.project.id]
}

resource "env0_environment_state_access" "example_entire_organization" {
  environment_id                      = data.env0_environment.environment.id
  accessible_from_entire_organization = true
}
