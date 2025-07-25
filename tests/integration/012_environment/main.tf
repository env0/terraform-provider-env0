provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_project" "test_project" {
  name          = "Test-Project-for-environment-${random_string.random.result}"
  force_destroy = true
}

resource "env0_project" "test_project2" {
  name          = "Test-Project2-for-environment-${random_string.random.result}"
  force_destroy = true
}

data "env0_template" "github_template_for_environment" {
  name = "Github Integrated Template"
}

resource "env0_template" "template" {
  repository             = data.env0_template.github_template_for_environment.repository
  github_installation_id = data.env0_template.github_template_for_environment.github_installation_id
  name                   = "Template for environment resource-${random_string.random.result}"
  type                   = "terraform"
  path                   = "misc/null-resource"
  terraform_version      = "0.15.1"
}

resource "env0_template" "template2" {
  repository             = data.env0_template.github_template_for_environment.repository
  github_installation_id = data.env0_template.github_template_for_environment.github_installation_id
  name                   = "Template for environment resource 2-${random_string.random.result}"
  type                   = "terraform"
  path                   = "misc/null-resource"
  terraform_version      = "0.16.1"
}

resource "env0_template_project_assignment" "assignment" {
  template_id = env0_template.template.id
  project_id  = env0_project.test_project.id
}

resource "env0_template_project_assignment" "assignment2" {
  template_id = env0_template.template.id
  project_id  = env0_project.test_project2.id
}

resource "env0_environment" "auto_glob_environment" {
  depends_on                       = [env0_template_project_assignment.assignment]
  name                             = "environment-auto-glob-${random_string.random.result}"
  project_id                       = env0_project.test_project.id
  template_id                      = env0_template.template.id
  auto_deploy_by_custom_glob       = var.second_run ? null : "//*"
  auto_deploy_on_path_changes_only = true
  approve_plan_automatically       = true
  deploy_on_push                   = true
  force_destroy                    = true
  wait_for_destroy                 = true
}

resource "env0_environment" "example" {
  depends_on    = [env0_template_project_assignment.assignment]
  force_destroy = true
  name          = "environment-${random_string.random.result}"
  project_id    = env0_project.test_project.id
  template_id   = env0_template.template.id
  configuration {
    name  = "environment configuration variable"
    value = "value"
  }
  approve_plan_automatically = true
  revision                   = "master"
  vcs_commands_alias         = "alias"
  drift_detection_cron       = var.second_run ? "*/5 * * * *" : "*/10 * * * *"
}

resource "env0_environment" "move_environment" {
  depends_on          = [env0_template_project_assignment.assignment]
  force_destroy       = true
  name                = "environment-move-${random_string.random.result}"
  project_id          = var.second_run ? env0_project.test_project2.id : env0_project.test_project.id
  template_id         = env0_template.template.id
  prevent_auto_deploy = true
}

resource "env0_environment" "modify_template" {
  depends_on                 = [env0_template_project_assignment.assignment]
  force_destroy              = true
  name                       = "environment-modify-template-${random_string.random.result}"
  project_id                 = env0_project.test_project.id
  template_id                = var.second_run ? env0_template.template2.id : env0_template.template.id
  prevent_auto_deploy        = true
  approve_plan_automatically = true
}


resource "env0_custom_role" "custom_role1" {
  name = "custom-role-${random_string.random.result}"
  permissions = [
    "VIEW_PROJECT",
    "EDIT_ENVIRONMENT_SETTINGS"
  ]
}

resource "env0_custom_role" "custom_role2" {
  name = "custom-role-${random_string.random.result}2"
  permissions = [
    "EDIT_ENVIRONMENT_SETTINGS"
  ]
}

resource "env0_api_key" "test_user_api_key" {
  name              = "my-little-user-api-key-${random_string.random.result}"
  organization_role = "User"
}

resource "env0_user_environment_assignment" "user_role_environment_assignment" {
  user_id        = env0_api_key.test_user_api_key.id
  environment_id = env0_environment.example.id
  role_id        = var.second_run ? env0_custom_role.custom_role1.id : env0_custom_role.custom_role2.id
}

resource "env0_team" "team" {
  name = "environment-team-${random_string.random.result}"
}

resource "env0_team_environment_assignment" "team_role_environment_assignment" {
  team_id        = env0_team.team.id
  environment_id = env0_environment.example.id
  role_id        = var.second_run ? env0_custom_role.custom_role1.id : env0_custom_role.custom_role2.id
}

/* TODO
resource "env0_environment_state_access" "disallow" {
  environment_id                      = env0_environment.example.id
  accessible_from_entire_organization = false
  allowed_project_ids                 = []
}
*/

resource "env0_template" "terragrunt_template" {
  name                 = "Terragrunt template for environment resource-${random_string.random.result}"
  type                 = "terragrunt"
  repository           = "https://github.com/env0/templates"
  path                 = "misc/null-resource"
  terragrunt_tf_binary = "terraform"
  terraform_version    = "0.15.1"
  terragrunt_version   = "0.35.0"
}

resource "env0_template_project_assignment" "terragrunt_assignment" {
  template_id = env0_template.terragrunt_template.id
  project_id  = env0_project.test_project.id
}

resource "env0_environment" "terragrunt_environment" {
  depends_on                       = [env0_template_project_assignment.terragrunt_assignment]
  force_destroy                    = true
  name                             = "environment-${random_string.random.result}"
  project_id                       = env0_project.test_project.id
  template_id                      = env0_template.terragrunt_template.id
  approve_plan_automatically       = true
  revision                         = "master"
  terragrunt_working_directory     = var.second_run ? "/second-dir" : "/first-dir"
  auto_deploy_on_path_changes_only = false
}

data "env0_environment" "test" {
  id = env0_environment.example.id
}

output "revision" {
  value = data.env0_environment.test.revision
}

output "terragrunt_working_directory" {
  value = env0_environment.terragrunt_environment.terragrunt_working_directory
}

data "env0_template" "github_template" {
  name = "Github Integrated Template"
}

resource "env0_environment" "environment-without-template" {
  force_destroy                    = true
  name                             = "environment-without-template-${random_string.random.result}"
  project_id                       = env0_project.test_project.id
  approve_plan_automatically       = true
  auto_deploy_on_path_changes_only = false

  without_template_settings {
    description                             = "Template description - GitHub"
    type                                    = "terraform"
    revision                                = "master"
    repository                              = data.env0_template.github_template.repository
    github_installation_id                  = data.env0_template.github_template.github_installation_id
    path                                    = var.second_run ? "second" : "misc/null-resource"
    retries_on_deploy                       = 3
    retry_on_deploy_only_when_matches_regex = "abc"
    retries_on_destroy                      = 1
    terraform_version                       = "0.15.1"
  }
}

data "env0_template" "gitlab_template" {
  name = "Gitlab Integrated Template"
}

resource "env0_environment" "environment-without-template-start-with-no-vcs" {
  name          = "start-with-non-vsc-${random_string.random.result}"
  project_id    = env0_project.test_project.id
  force_destroy = true

  auto_deploy_on_path_changes_only = var.second_run ? true : false
  deploy_on_push                   = var.second_run ? true : false
  run_plan_on_pull_requests        = var.second_run ? true : false

  without_template_settings {
    type             = "opentofu"
    is_gitlab        = var.second_run ? true : false
    repository       = data.env0_template.gitlab_template.repository
    token_id         = data.env0_template.gitlab_template.token_id
    token_name       = data.env0_template.gitlab_template.token_name
    opentofu_version = "latest"
    path             = "null-resource"
  }
}

resource "env0_environment" "inactive" {
  depends_on                 = [env0_template_project_assignment.assignment]
  force_destroy              = true
  name                       = "environment-${random_string.random.result}-inactive"
  project_id                 = env0_project.test_project.id
  template_id                = env0_template.template.id
  approve_plan_automatically = true
  revision                   = "master"
  vcs_commands_alias         = "alias"
  is_inactive                = var.second_run ? "true" : "false"
}

# workflow environment

resource "env0_template" "workflow_template" {
  name              = "Template for workflow environment-${random_string.random.result}"
  type              = "workflow"
  repository        = "https://github.com/env0/templates"
  path              = "misc/single-environment-workflow"
  terraform_version = "1.1.5"
}

data "env0_template" "sub_environment_null_template" {
  name = "null resource"
}

resource "env0_template_project_assignment" "assignment_sub_environment_null_template" {
  template_id = data.env0_template.sub_environment_null_template.id
  project_id  = env0_project.test_project.id
}

resource "env0_template_project_assignment" "assignment_workflow" {
  template_id = env0_template.workflow_template.id
  project_id  = env0_project.test_project.id
}

resource "env0_variable_set" "variable_set1" {
  name        = "variable-set-project-${random_string.random.result}"
  description = "description123"
  scope       = "project"
  scope_id    = env0_project.test_project.id

  variable {
    name   = "n1"
    value  = "v1"
    type   = "terraform"
    format = "text"
  }
}

resource "env0_variable_set" "variable_set2" {
  name        = "variable-set-project2-${random_string.random.result}"
  description = "description123"
  scope       = "project"
  scope_id    = env0_project.test_project.id

  variable {
    name   = "n2"
    value  = "v2"
    type   = "terraform"
    format = "text"
  }
}

resource "env0_environment" "workflow-environment" {
  depends_on = [
    env0_template_project_assignment.assignment_workflow,
    env0_template_project_assignment.assignment_sub_environment_null_template
  ]
  force_destroy = true
  name          = "environment-workflow-${random_string.random.result}"
  project_id    = env0_project.test_project.id
  template_id   = env0_template.workflow_template.id

  configuration {
    name  = "n1"
    value = "v1"
  }

  variable_sets = var.second_run ? [env0_variable_set.variable_set2.id] : [env0_variable_set.variable_set1.id]

  sub_environment_configuration {
    alias                      = "rootService1"
    revision                   = "master"
    approve_plan_automatically = var.second_run ? false : true
    configuration {
      name  = "sub_env1_var1"
      value = "hello"
    }
    configuration {
      name  = "sub_env1_var2"
      value = "world"
    }
  }
}

resource "env0_configuration_variable" "sub_environment_example" {
  depends_on            = [env0_environment.workflow-environment]
  name                  = "SUB_ENVIRONMENT_VARIABLE-${random_string.random.result}"
  value                 = "sub env value"
  template_id           = env0_template.workflow_template.id
  sub_environment_alias = "rootService1"
  description           = "Variable for a sub environment scope"
}

data "env0_environment" "test-workflow" {
  id = env0_environment.workflow-environment.id
}

resource "env0_environment" "mark_as_archived" {
  depends_on       = [env0_template_project_assignment.assignment]
  name             = "environment-mark-as-archived-${random_string.random.result}"
  project_id       = env0_project.test_project.id
  template_id      = env0_template.template.id
  removal_strategy = "mark_as_archived"
}

resource "env0_environment" "environment_output_format" {
  depends_on    = [env0_template_project_assignment.assignment]
  name          = "environment-output-format-${random_string.random.result}"
  project_id    = env0_project.test_project.id
  template_id   = env0_template.template.id
  force_destroy = true

  configuration {
    name = "vpc_id"
    value = jsonencode({
      environmentId = "vpc"
      outputName    = "vpc_id"
    })
    schema_format = "ENVIRONMENT_OUTPUT"
    type          = "terraform"
  }
}

