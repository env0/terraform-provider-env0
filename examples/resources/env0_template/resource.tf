data "env0_project" "default_project" {
  name = "Default Organization Project"
}

data "env0_ssh_key" "my_key" {
  name = "Secret Key"
}

resource "env0_template" "example" {
  name = "example"
  description = "Example template"
  repository = "https://github.com/env0/templates"
  path = "aws/hello-world"
  ssh_keys = [ data.ssh_keys.my_key]
}
