data "env0_template" "github_template" {
  name = "github_template"
}

resource "env0_custom_flow" "custom_flow" {
  name                   = "Custom Flow"
  repository             = data.env0_template.github_template.repository
  github_installation_id = data.env0_template.github_template.github_installation_id
  path                   = "custom-flows/opa.yaml"
}
