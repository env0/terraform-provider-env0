resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_vault_oidc_credentials" "oidc_credentials" {
  name                  = "test vault oidc credentials ${random_string.random.result}"
  address               = var.second_run ? "http://fake2.com:80" : "http://fake1.com:80"
  version               = "version"
  role_name             = "role_name"
  jwt_auth_backend_path = var.second_run ? "path2" : "path1"
  namespace             = "namespace"
}

data "env0_vault_oidc_credentials" "oidc_credentials" {
  name       = "test vault oidc credentials ${random_string.random.result}"
  depends_on = [env0_vault_oidc_credentials.oidc_credentials]
}

