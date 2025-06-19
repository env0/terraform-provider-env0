resource "env0_gcp_cloud_configuration" "example" {
  name                                 = "example-gcp-config"
  gcp_project_id                       = "your-gcp-project-id"
  credential_configuration_file_content = file("path/to/your-gcp-service-account.json")
}
