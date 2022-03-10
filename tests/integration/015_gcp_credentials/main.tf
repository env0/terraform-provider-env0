resource "env0_gcp_credentials" "gcp_cred" {
  name                = "test gcp credentials 1"
  service_account_key = "service account key example"
}

data "env0_gcp_credentials" "gcp_cred" {
  name = env0_gcp_credentials.gcp_cred.name
}

resource "env0_gcp_credentials" "gcp_cred_with_project_id" {
  name                = "Test GCP credentials with project_id 2"
  service_account_key = "example service_account_key"
  project_id          = "example project id"
}

data "env0_gcp_credentials" "gcp_cred_with_project_id" {
  name = env0_gcp_credentials.gcp_cred_with_project_id.name
}


output "gcp_cred_name" {
  value = data.env0_gcp_credentials.gcp_cred.name
}

output "gcp_cred_name_with_project_id" {
  value = data.env0_gcp_credentials.gcp_cred_with_project_id.name
}
