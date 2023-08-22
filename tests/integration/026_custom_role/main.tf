provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_custom_role" "custom_role1" {
  name = "custom-role-${random_string.random.result}1"
  permissions = [
    "VIEW_PROJECT",
  ]
}

resource "env0_custom_role" "custom_role2" {
  name = "custom-role-${random_string.random.result}2"
  permissions = [
    "VIEW_ORGANIZATION",
    "EDIT_ORGANIZATION_SETTINGS",
    "CREATE_AND_EDIT_TEMPLATES",
    "CREATE_AND_EDIT_MODULES",
    "CREATE_PROJECT",
    "VIEW_PROJECT",
    "EDIT_PROJECT_SETTINGS",
    "MANAGE_PROJECT_TEMPLATES",
    "EDIT_ENVIRONMENT_SETTINGS",
    "ARCHIVE_ENVIRONMENT",
    "OVERRIDE_MAX_TTL",
    "CREATE_CROSS_PROJECT_ENVIRONMENTS",
    "OVERRIDE_MAX_ENVIRONMENT_PROJECT_LIMITS",
    "RUN_PLAN",
    "RUN_APPLY",
    "ABORT_DEPLOYMENT",
    "RUN_TASK",
    "CREATE_CUSTOM_ROLES",
    "VIEW_DASHBOARD",
    "VIEW_MODULES",
    "READ_STATE",
    "WRITE_STATE",
    "FORCE_UNLOCK_WORKSPACE",
    "MANAGE_BILLING",
    "VIEW_AUDIT_LOGS",
    "MANAGE_ENVIRONMENT_LOCK",
  ]
}

data "env0_custom_roles" "all_roles" {
  depends_on = [env0_custom_role.custom_role1, env0_custom_role.custom_role2]
}

data "env0_team" "team_resource1" {
  name = data.env0_custom_roles.all_roles.names[index(data.env0_custom_roles.all_roles.names, env0_custom_role.custom_role1.name)]
}

data "env0_team" "team_resource2" {
  name = data.env0_custom_roles.all_roles.names[index(data.env0_custom_roles.all_roles.names, env0_custom_role.custom_role2.name)]
}


resource "env0_api_key" "test_api_key" {
  name = "api-key-${random_string.random.result}"
}

resource "env0_user_organization_assignment" "user_org" {
  user_id        = env0_api_key.test_api_key.id
  custom_role_id = var.second_run ? null : env0_custom_role.custom_role1.id
  role           = var.second_run ? "Admin" : null
}
