data "env0_template" "template" {
  name = "Template Name"
}

data "env0_project" "project" {
  name = "Default Organization Project"
}

resource "env0_environment" "environment" {
  name        = "environment"
  project_id  = data.env0_project.project.id
  template_id = data.env0_template.template.id
}

resource "env0_api_key" "api_key" {
  name              = "user-api-key-sample"
  organization_role = "User"
}

resource "env0_custom_role" "custom_role" {
  name = "custom-role-sample"
  permissions = [
    "VIEW_PROJECT",
    "EDIT_PROJECT_SETTINGS"
  ]
}

resource "env0_user_environment_assignment" "assignment" {
  user_id        = env0_api_key.api_key.id
  environment_id = env0_environment.environment.id
  role_id        = env0_custom_role.custom_role.id
}
