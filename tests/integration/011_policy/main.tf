data "env0_policy" "default" {
  project_id = "Project0"
}

resource "env0_policy" "test_policy" {
  id         = data.env0_policy.default.id
  project_id = data.env0_policy.default.project_id
}

output "policy_project_id" {
  value = data.env0_policy.default.project_id
}

