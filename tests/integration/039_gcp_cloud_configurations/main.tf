provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_gcp_cloud_configuration" "test" {
  name                                  = "test-gcp-config-${random_string.random.result}"
  gcp_project_id                        = "test-gcp-project"
  credential_configuration_file_content = <<EOF
{
  "audience": "//iam.googleapis.com/projects/578187717855/locations/global/workloadIdentityPools/cloudcompass-wif-pool/providers/cloudcompass-oidc-provider",
  "credential_source": {
    "file": "example.json",
    "format": {
      "type": "json",
      "subject_token_field_name": "access_token"
    }
  }
}
EOF
}
