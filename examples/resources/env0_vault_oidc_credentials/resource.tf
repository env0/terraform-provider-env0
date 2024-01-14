resource "env0_vault_oidc_credentials" "example" {
  name                  = "example"
  address               = "http://fake1.com:80"
  version               = "version"
  role_name             = "role_name"
  jwt_auth_backend_path = "path"
  namespace             = "namespace"
}
