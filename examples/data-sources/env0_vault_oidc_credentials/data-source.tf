resource "env0_vault_oidc_credentials" "example" {
  name                  = "example"
  address               = "http://fake1.com:80"
  version               = "version"
  role_name             = "role_name"
  jwt_auth_backend_path = "path"
  namespace             = "namespace"
}

data "env0_vault_oidc_credentials" "by_id" {
  id = env0_vault_oidc_credentials.example.id
}

data "env0_vault_oidc_credentials" "by_name" {
  name = env0_vault_oidc_credentials.example.name
}
