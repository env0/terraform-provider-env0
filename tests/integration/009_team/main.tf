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