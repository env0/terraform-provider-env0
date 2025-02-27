provider "random" {}

resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_azure_cloud_configuration" "azure_cloud_configuration" {
  name                       = "azure-${random_string.random.result}"
  tenant_id                  = var.second_run ? "22222222-2222-2222-2222-222222222222" : "11111111-1111-1111-1111-111111111111"
  client_id                  = var.second_run ? "44444444-4444-4444-4444-444444444444" : "33333333-3333-3333-3333-333333333333"
  log_analytics_workspace_id = var.second_run ? "66666666-6666-6666-6666-666666666666" : "55555555-5555-5555-5555-555555555555"
}
