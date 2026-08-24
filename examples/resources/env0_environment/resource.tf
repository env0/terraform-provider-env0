data "env0_template" "example" {
  name = "Template Name"
}

data "env0_project" "default_project" {
  name = "Default Organization Project"
}

resource "env0_environment" "example" {
  name        = "environment"
  project_id  = data.env0_project.default_project.id
  template_id = data.env0_template.example.id
}

resource "env0_environment" "example_with_tags" {
  name        = "environment with tags"
  project_id  = data.env0_project.default_project.id
  template_id = data.env0_template.example.id

  tags = {
    owner         = "alice"
    "cost-center" = "1234"
    # a key may hold several values - join them with a comma and no space
    team = "eng,payments"
  }
}

resource "env0_environment" "example_with_hcl_configuration" {
  name        = "environment with hcl"
  project_id  = data.env0_project.default_project.id
  template_id = data.env0_template.example.id

  configuration {
    name          = "TEST1234"
    type          = "terraform"
    value         = <<EOF
      {
        a = "world11111"
        b = {
          c = "d"
        }
      }
    EOF
    schema_format = "HCL"
  }
}
