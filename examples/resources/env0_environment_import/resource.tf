resource "env0_environment_import" "new_environment_import" {
  name         = "name"
  git_provider = "github"
  path         = "path/to/tf/config"
  repository   = "reponame"
  revision     = "revision"
  workspace    = "workspace"
  iac_version  = "1.7.1"
  iac_type     = "opentofu"
}

