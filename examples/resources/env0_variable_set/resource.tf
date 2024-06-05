data "env0_project" "project" {
  name = "project"
}

data "env0_environment" "environment" {
  name = "environment"
}

resource "env0_variable_set" "organization_scope_example" {
  name        = "variable-set-example1"
  description = "description123"

  variable {
    name   = "n1"
    value  = "v1"
    format = "text"
  }

  variable {
    name         = "n1"
    value        = "v2"
    type         = "environment"
    format       = "text"
    is_sensitive = true
  }

  variable {
    name   = "n3"
    value  = "v3"
    type   = "terraform"
    format = "hcl"
  }

  variable {
    name   = "n4"
    value  = "{}"
    type   = "terraform"
    format = "json"
  }

  variable {
    name            = "n5"
    dropdown_values = ["o3", "o2"]
    type            = "terraform"
    format          = "dropdown"
  }
}

resource "env0_variable_set" "project_scope_example" {
  name        = "variable-set-example2"
  description = "description123"
  scope       = "project"
  scope_id    = data.env0_project.project.id

  variable {
    name   = "n1"
    value  = "v1"
    format = "text"
  }
}

resource "env0_variable_set_assignment" "assignment" {
  scope    = "environment"
  scope_id = data.env0_environment.id
  set_ids = [
    env0_variable_set.project_scope_example.id,
    env0_variable_set.organization_scope_example.id,
  ]
}
