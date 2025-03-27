data "env0_template" "github_template" {
  name = "github_template"
}

resource "env0_custom_flow" "custom_flow" {
  name                   = "Custom Flow"
  repository             = data.env0_template.github_template.repository
  github_installation_id = data.env0_template.github_template.github_installation_id // The installation ID is taken from an existing authorized template
  path                   = "custom-flows/my-custom-flow.yaml"
}


// Self Hosted VCS
resource "env0_custom_flow" "ghe_custom_flow" {
  name                 = "GHE Custom Flow"
  revision             = "my-revision"
  repository           = "https://mycompany.github.com/myorg/myrepo"
  path                 = "custom-flows/my-custom-flow.yaml"
  is_github_enterprise = true
}
