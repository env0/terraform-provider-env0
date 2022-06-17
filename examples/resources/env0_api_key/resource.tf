resource "env0_api_key" "api_key_example" {
  name = "api-key-example"
}

resource "env0_project" "project_resource" {
  name = "project-resource"
}

resource "env0_user_project_assignment" "api_key_project_assignment_example" {
  user_id    = env0_api_key.api_key_example.id
  project_id = env0_project.project_resource.id
  role       = "Viewer"
}

resource "env0_team" "team_resource" {
  name = "team-resource"
}

resource "env0_user_team_assignment" "api_key_team_assignment_example" {
  user_id = env0_api_key.api_key_example.id
  team_id = env0_team.team_resource.id
}
