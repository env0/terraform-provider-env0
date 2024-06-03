provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_project" "project" {
  name = "project-for-variable-set-${random_string.random.result}"
}

resource "env0_variable_set" "org_scope" {
  name        = "variable-set-org-${random_string.random.result}"
  description = "description123"

  variable {
    name   = "n1"
    value  = var.second_run ? "v2" : "v1"
    type   = "terraform"
    format = "text"
  }

  variable {
    name         = "n1"
    value        = var.second_run ? "v22" : "v2"
    format       = "text"
    is_sensitive = true
  }

  variable {
    name   = "n3"
    value  = var.second_run ? "v32" : "v3"
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
    dropdown_values = var.second_run ? ["o3", "o2"] : ["o1", "o2"]
    type            = "terraform"
    format          = "dropdown"
  }

  variable {
    name   = "n55555"
    value  = "abcdef"
    type   = var.second_run ? "terraform" : "environment"
    format = "text"
  }
}

resource "env0_variable_set" "project_scope" {
  name        = "variable-set-project-${random_string.random.result}"
  description = "description123"
  scope       = "project"
  scope_id    = env0_project.project.id

  variable {
    name   = "n1"
    value  = "v1"
    type   = "terraform"
    format = "text"
  }
}

resource "env0_variable_set_assignment" "assignment" {
  scope    = "project"
  scope_id = env0_project.project.id
  set_ids  = [env0_variable_set.org_scope.id]
}
