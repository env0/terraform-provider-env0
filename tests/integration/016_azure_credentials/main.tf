resource "env0_azure_credentials" "azure_cred" {
  name            = "test azure credentials 1"
  client_id       = "client_id"
  client_secret   = "client_secret"
  subscription_id = "subscription_id"
  tenant_id       = "tenant_id"
}

data "env0_azure_credentials" "azure_cred" {
  name = env0_azure_credentials.azure_cred.name
}

output "azure_cred_name" {
  value = data.env0_azure_credentials.azure_cred.name
}


