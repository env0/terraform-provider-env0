resource "env0_project" "test_project" {
  name = "test-project"
}

resource "env0_team" "team_resource" {
  name = "test-team"
}

resource "env0_team_project_assignment" "role_assignment_example" {
  project_id = env0_project.test_project.id
  team_id    = env0_team.team_resource.id
  role       = "Admin"
}

resource "env0_custom_role" "custom_role_example" {
  name = "my custom role 1"
  permissions = [
    "VIEW_PROJECT",
    "EDIT_PROJECT_SETTINGS"
  ]
}

resource "env0_team_project_assignment" "role_assignment_custom_role_example" {
  team_id        = env0_team.team_resource.id
  project_id     = env0_project.project_example.id
  custom_role_id = env0_custom_role.custom_role_example.id
}
