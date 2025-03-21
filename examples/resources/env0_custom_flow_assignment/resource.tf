data "env0_template" "my_template" {
  name = "my-template"
}

resource "env0_project" "my_project" {
  name = "my-project"
}

resource "env0_custom_flow" "my_custom_flow" {
  name                   = "custom-flow"
  repository             = data.env0_template.my_template.repository
  github_installation_id = data.env0_template.my_template.github_installation_id
  path                   = "custom-flows/my-custom-flow.yaml"
}

resource "env0_custom_flow_assignment" "my_assignment" {
  scope_id    = env0_project.my_project.id
  template_id = env0_custom_flow.my_custom_flow.id
}
