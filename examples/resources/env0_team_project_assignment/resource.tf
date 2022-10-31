data "env0_user" "user_example" {
  email = "example@email.com"
}

resource "env0_project" "project_example" {
  name = "project-example"
}

resource "env0_user_project_assignment" "role_assignment_example" {
  user_id    = data.env0_user.user_example.id
  project_id = env0_project.project_example.id
  role       = "Viewer"
}

resource "env0_custom_role" "custom_role_example" {
  name = "my custom role 1"
  permissions = [
    "VIEW_PROJECT",
    "EDIT_PROJECT_SETTINGS"
  ]
}

resource "env0_user_project_assignment" "role_assignment_custom_role_example" {
  user_id        = data.env0_user.user_example.id
  project_id     = env0_project.project_example.id
  custom_role_id = env0_custom_role.custom_role_example.id
}
