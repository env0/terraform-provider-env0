provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_project" "test_project" {
  name        = "Test-Project-${random_string.random.result}"
  description = "Test Description ${var.second_run ? "after update" : ""}"
}

resource "env0_team" "team_resource" {
  name        = "Test-Team-010-${random_string.random.result}"
  description = var.second_run ? "second description" : "first description"
}

resource "env0_team" "team_resource2" {
  name        = "Test-Team-010-${random_string.random.result}2"
  description = var.second_run ? "second description" : "first description"
}

resource "env0_team_project_assignment" "assignment" {
  depends_on = [env0_team.team_resource, env0_project.test_project]
  project_id = env0_project.test_project.id
  team_id    = env0_team.team_resource.id
  role       = "Admin"
}

resource "env0_custom_role" "custom_role" {
  name = "custom-role-${random_string.random.result}"
  permissions = [
    "VIEW_PROJECT",
    "EDIT_PROJECT_SETTINGS"
  ]
}

data "env0_custom_role" "custom_role" {
  name       = "custom-role-${random_string.random.result}"
  depends_on = [env0_custom_role.custom_role]
}

resource "env0_team_project_assignment" "custom_assignment" {
  depends_on     = [env0_team.team_resource, env0_project.test_project]
  project_id     = env0_project.test_project.id
  team_id        = env0_team.team_resource2.id
  custom_role_id = var.second_run ? data.env0_custom_role.custom_role.id : null
  role           = var.second_run ? null : "Viewer"
}
