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

resource "env0_project" "project_resource2" {
  name        = "Test-Project-API-${random_string.random.result}2"
  description = "Test Description"
}

resource "env0_user_project_assignment" "api_key_project_assignment" {
  user_id    = env0_api_key.test_user_api_key.id
  project_id = env0_project.project_resource.id
  role       = var.second_run ? "Viewer" : "Planner"
}

resource "time_sleep" "wait_15_seconds" {
  depends_on = [env0_api_key.test_api_key]

  create_duration = "15s"
}

data "env0_api_key" "test_api_key1" {
  name       = env0_api_key.test_api_key.name
  depends_on = [time_sleep.wait_15_seconds]
}

data "env0_api_key" "test_api_key2" {
  id         = env0_api_key.test_api_key.id
  depends_on = [time_sleep.wait_15_seconds]
}

resource "env0_api_key" "test_api_key_omitted" {
  name                = "omitted-api-key-secret-${random_string.random.result}"
  omit_api_key_secret = true
}

resource "env0_custom_role" "custom_role" {
  name = "custom-role-${random_string.random.result}"
  permissions = [
    "VIEW_PROJECT",
    "EDIT_PROJECT_SETTINGS"
  ]
}

resource "env0_user_project_assignment" "api_key_project_assignment_custom_role" {
  user_id        = env0_api_key.test_user_api_key.id
  project_id     = env0_project.project_resource2.id
  custom_role_id = var.second_run ? null : env0_custom_role.custom_role.id
  role           = var.second_run ? "Viewer" : null
}

resource "env0_api_key" "test_api_key_with_permissions" {
  name              = "api-key-with-permissions-${random_string.random.result}"
  organization_role = "User"

  project_permissions {
    project_id   = env0_project.project_resource.id
    project_role = "Deployer"
  }

  project_permissions {
    project_id   = env0_project.project_resource2.id
    project_role = "Viewer"
  }
}

resource "env0_api_key" "test_api_key_custom_role" {
  name              = "api-key-custom-role-${random_string.random.result}"
  organization_role = env0_custom_role.custom_role.id

  project_permissions {
    project_id   = env0_project.project_resource.id
    project_role = "Planner"
  }
}
