resource "env0_team" "team" {
  name = "team"
}

resource "env0_custom_role" "custom_role" {
  name = "custom-role-sample"
  permissions = [
    "VIEW_ORGANIZATION",
    "VIEW_DASHBOARD"
  ]
}

resource "env0_team_organization_assignment" "assignment" {
  team_id = env0_team.team.id
  role_id = env0_custom_role.custom_role.id
}
