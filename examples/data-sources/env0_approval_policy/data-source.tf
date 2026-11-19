data "env0_approval_policy" "approval_policy_by_name" {
  name = "approval policy 12345"
}

data "env0_approval_policy" "approval_policy_by_id" {
  id = "12345676safsd"
}

resource "env0_approval_policy_assignment" "assignment" {
  scope_id     = "project-id"
  blueprint_id = data.env0_approval_policy.approval_policy_by_name.id
}
