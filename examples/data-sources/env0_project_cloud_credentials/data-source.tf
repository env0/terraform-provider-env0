data "env0_project_cloud_credentials" "project_cloud_credentials" {
  project_id = "pid1"
}

output "pid1_credential_0_id" {
  value = data.env0_project_cloud_credentials.project_cloud_credentials.ids.0
}

output "pid1_credential_1_id" {
  value = data.env0_project_cloud_credentials.project_cloud_credentials.ids.1
}
