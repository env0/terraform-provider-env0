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
    "EDIT_PROJECT_SETTINGS"
  ]
}

data "env0_custom_roles" "all_roles" {}

data "env0_custom_role" "roles" {
  for_each = toset(data.env0_custom_roles.all_roles.names)
  name     = each.value
}

resource "env0_api_key" "test_api_key" {
  name = "api-key-${random_string.random.result}"
}

resource "env0_user_organization_assignment" "user_org" {
  user_id        = env0_api_key.test_api_key.id
  custom_role_id = var.second_run ? null : env0_custom_role.custom_role1.id
  role           = var.second_run ? "Viewer" : null
}
