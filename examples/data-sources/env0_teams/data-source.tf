data "env0_teams" "all_teams" {}

data "env0_team" "teams" {
  for_each = toset(data.env0_teams.all_teams.names)
  name     = each.value
}

output "team1_name" {
  value = data.env0_team.teams["team1"].name
}

output "team2_name" {
  value = data.env0_team.teams["team2"].name
}