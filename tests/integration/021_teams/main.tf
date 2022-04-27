resource "env0_team" "team_resource1" {
  name = "team1"
}

resource "env0_team" "team_resource2" {
  name = "team2"
}

data "env0_teams" "all_teams" {}

data "env0_team" "teams" {
  for_each = toset(data.env0_teams.all_teams.names)
  name     = each.value
}

output "team1_name" {
  value = var.second_run ? data.env0_team.teams["team1"].name : ""
}

output "team2_name" {
  value = var.second_run ? data.env0_team.teams["team2"].name : ""
}