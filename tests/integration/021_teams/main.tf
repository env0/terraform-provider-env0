provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_team" "team_resource1" {
  name        = "team1-021-${random_string.random.result}"
  description = "team 1 description"
}

resource "env0_team" "team_resource2" {
  name        = "team2-021-${random_string.random.result}"
  description = "team 2 description"
}

data "env0_teams" "all_teams" {
  depends_on = [env0_team.team_resource1, env0_team.team_resource2]
}

data "env0_team" "team_resource1" {
  name = data.env0_teams.all_teams.names[index(data.env0_teams.all_teams.names, env0_team.team_resource1.name)]
}

data "env0_team" "team_resource2" {
  name = data.env0_teams.all_teams.names[index(data.env0_teams.all_teams.names, env0_team.team_resource2.name)]
}

output "team1_description" {
  value = var.second_run ? data.env0_team.team_resource1.description : ""
}

output "team2_description" {
  value = var.second_run ? data.env0_team.team_resource2.description : ""
}
