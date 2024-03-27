provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_team" "team_resource" {
  name        = "Test-Team-009-${random_string.random.result}"
  description = var.second_run ? "second description" : "first description"
}

data "env0_team" "team_data" {
  name       = env0_team.team_resource.name
  depends_on = [env0_team.team_resource]
}

output "team_resource_description" {
  value = env0_team.team_resource.description
}

output "team_data_description" {
  value = data.env0_team.team_data.description
}

resource "env0_custom_role" "custom_role1" {
  name = "custom-role-${random_string.random.result}"
  permissions = [
    "CREATE_PROJECT"
  ]
}

resource "env0_custom_role" "custom_role2" {
  name = "custom-role-${random_string.random.result}2"
  permissions = [
    "RUN_PLAN"
  ]
}

resource "env0_team_organization_assignment" "organization_role_assignment" {
  team_id = env0_team.team_resource.id
  role_id = var.second_run ? env0_custom_role.custom_role1.id : env0_custom_role.custom_role2.id
}

