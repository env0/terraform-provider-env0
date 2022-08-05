provider "random" {}

resource "random_string" "random" {
  length    = 5
  special   = false
  min_lower = 5
}

resource "env0_api_key" "test_api_key" {
  name = "my-little-api-key-${random_string.random.result}"
}

resource "env0_api_key" "test_user_api_key" {
  name              = "my-little-user-api-key-${random_string.random.result}"
  organization_role = "User"
}

resource "env0_team" "team_resource" {
  name        = "team-with-api-key-020-${random_string.random.result}"
  description = "description"
}

resource "env0_user_team_assignment" "api_key_team_assignment" {
  user_id = env0_api_key.test_user_api_key.id
  team_id = env0_team.team_resource.id
}

resource "env0_project" "project_resource" {
  name        = "Test-Project-API-${random_string.random.result}"
  description = "Test Description"
}

resource "env0_user_project_assignment" "api_key_project_assignment" {
  user_id    = env0_api_key.test_user_api_key.id
  project_id = env0_project.project_resource.id
  role       = var.second_run ? "Viewer" : "Planner"
}

data "env0_api_key" "test_api_key1" {
  name       = env0_api_key.test_api_key.name
  depends_on = [env0_api_key.test_api_key]
}

data "env0_api_key" "test_api_key2" {
  id = env0_api_key.test_api_key.id
}

resource "env0_api_key" "test_api_key_omitted" {
  name = "omitted-api-key-secret-${random_string.random.result}"
}

output "omitted_secret" {
  value = env0_api_key.test_api_key_omitted.api_key_secret
}
