provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_gcp_credentials" "gcp_cred" {
  name                = "test gcp credentials 1-${random_string.random.result}"
  service_account_key = "service account key example"
}

data "env0_gcp_credentials" "gcp_cred" {
  name = env0_gcp_credentials.gcp_cred.name
}

resource "env0_gcp_credentials" "gcp_cred_with_project_id" {
  name                = "Test GCP credentials with project_id 2-${random_string.random.result}"
  service_account_key = "example service_account_key"
  project_id          = "example project id"
}

data "env0_gcp_credentials" "gcp_cred_with_project_id" {
  name = env0_gcp_credentials.gcp_cred_with_project_id.name
}

resource "env0_gcp_oidc_credentials" "gcp_credentials" {
  name = "test azure oidc credentials ${random_string.random.result}"
  credential_configuration_file_content = jsonencode({
    "key" : "value"
  })
}

data "env0_gcp_oidc_credentials" "gcp_credentials" {
  name       = "test azure oidc credentials ${random_string.random.result}"
  depends_on = [env0_gcp_oidc_credentials.oidc_credentials]
}



output "gcp_cred_name" {
  value = data.env0_gcp_credentials.gcp_cred.name
}

output "gcp_cred_name_with_project_id" {
  value = data.env0_gcp_credentials.gcp_cred_with_project_id.name
}
