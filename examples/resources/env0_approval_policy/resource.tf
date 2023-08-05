resource "env0_project" "project" {
  name        = "project-name"
  description = "project-description"
}

resource "env0_approval_policy" "approval-policy" {
  name                   = "approval-policy-PROJECT-${env0_project.project.id}"
  repository             = "repo"
  github_installation_id = 1234
}

resource "env0_approval_policy_assignment" "assignment" {
  scope_id     = env0_project.project.id
  scope        = "PROJECT"
  blueprint_id = env0_approval_policy.approval-policy.id
}
