resource "env0_gcp_oidc_credentials" "credentials" {
  name = "example"
  credential_configuration_file_content = jsonencode({
    "key" : "value"
  })
}
