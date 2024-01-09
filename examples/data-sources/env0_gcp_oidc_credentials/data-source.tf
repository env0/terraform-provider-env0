resource "env0_azure_oidc_credentials" "credentials" {
  name            = "example"
  tenant_id       = "4234-2343-24234234234-42343"
  client_id       = "fff333-345555-4444"
  subscription_id = "f1111-222-2222"
}


data "env0_azure_oidc_credentials" "by_id" {
  id = env0_azure_oidc_credentials.example.id
}

data "env0_azure_oidc_credentials" "by_name" {
  name = env0_azure_oidc_credentials.example.name
}
