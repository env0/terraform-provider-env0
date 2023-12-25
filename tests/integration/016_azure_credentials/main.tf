resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_azure_credentials" "azure_cred" {
  name            = "test azure credentials 1 ${random_string.random.result}"
  client_id       = "client_id"
  client_secret   = "client_secret"
  subscription_id = "subscription_id"
  tenant_id       = "tenant_id"
}

resource "env0_azure_oidc_credentials" "oidc_credentials" {
  name            = "test azure oidc credentials ${random_string.random.result}"
  client_id       = "client_id"
  subscription_id = "subscription_id"
  tenant_id       = "tenant_id"
}

data "env0_azure_oidc_credentials" "oidc_credentials" {
  name       = "test azure oidc credentials ${random_string.random.result}"
  depends_on = [env0_azure_oidc_credentials.oidc_credentials]
}

data "env0_azure_credentials" "azure_cred" {
  name = env0_azure_credentials.azure_cred.name
}

output "azure_cred_name" {
  value = replace(data.env0_azure_credentials.azure_cred.name, random_string.random.result, "")
}


