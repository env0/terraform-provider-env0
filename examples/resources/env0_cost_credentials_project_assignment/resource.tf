resource "env0_azure_cost_credentials" "azure_cost_credentials" {
  name            = "cost credentials"
  client_id       = "client id"
  client_secret   = "client secret"
  subscription_id = "43242342dsdfsdfsdf"
  tenant_id       = "fsdf-fsdfdsfs-fsdfsdfsd-fsdfsd"
}

resource "env0_project" "project" {
  name = "project"
}

resource "env0_cost_credentials_project_assignment" "cost_project_assignment" {
  credential_id = env0_azure_cost_credentials.azure_cost_credentials.id
  project_id    = env0_project.project.id
}
