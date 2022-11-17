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
