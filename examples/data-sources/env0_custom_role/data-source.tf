data "env0_custom_role" "my_role" {
  name = "my_role"
}

resource "env0_team_project_assignment" "role_assignment_custom_role_example" {
  team_id        = env0_team.team_resource.id
  project_id     = env0_project.project_example.id
  custom_role_id = data.env0_custom_role.my_role.id
}
