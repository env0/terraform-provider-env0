resource "env0_gcp_cloud_configuration" "test" {
  name                                 = "test-gcp-config"
  gcp_project_id                       = "test-gcp-project"
  credential_configuration_file_content = "{\"type\":\"service_account\",...}"
}
