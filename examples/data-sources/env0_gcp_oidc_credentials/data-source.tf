resource "env0_gcp_oidc_credentials" "credentials" {
  name = "example"
  credential_configuration_file_content = jsonencode({
    "key" : "value"
  })
}

data "env0_gcp_oidc_credentials" "by_id" {
  id = env0_gcp_oidc_credentials.example.id
}

data "env0_gcp_oidc_credentials" "by_name" {
  name = env0_gcp_oidc_credentials.example.name
}
