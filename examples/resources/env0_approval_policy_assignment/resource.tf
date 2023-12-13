data "env0_project" "project" {
  name = "project"
}

resource "env0_approval_policy" "approval_policy" {
  name                   = "approval policy"
  repository             = "reopository"
  github_installation_id = 4234234234
  path                   = "misc/null-resource"

}

resource "env0_approval_policy_assignment" "approval_policy_assignment" {
  scope        = "PROJECT"
  scope_id     = data.env0_project.project.id
  blueprint_id = env0_approval_policy.approval_policy.id
}
