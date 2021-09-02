resource "env0_team" "team_resource" {
  name        = "Test-Team"
  description = var.second_run ? "second description" : "first description"
}

data "env0_team" "team_data" {
  name       = env0_team.team_resource.name
  depends_on = [env0_team.team_resource]
}

output "team_resource_name" {
  value = env0_team.team_resource.name
}

output "team_data_name" {
  value = data.env0_team.team_data.name
}