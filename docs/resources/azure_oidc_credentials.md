---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "env0_azure_oidc_credentials Resource - terraform-provider-env0"
subcategory: ""
description: |-
  
---

# env0_azure_oidc_credentials (Resource)



## Example Usage

```terraform
resource "env0_azure_oidc_credentials" "credentials" {
  name            = "example"
  tenant_id       = "4234-2343-24234234234-42343"
  client_id       = "fff333-345555-4444"
  subscription_id = "f1111-222-2222"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `client_id` (String) the azure client id
- `name` (String) name for the oidc credentials
- `subscription_id` (String) the azure subscription id
- `tenant_id` (String) the azure tenant id

### Optional

- `project_id` (String) the env0 project id to associate the credentials with

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import env0_aws_oidc_credentials.by_id d31a6b30-5f69-4d24-937c-22322754934e
terraform import env0_aws_oidc_credentials.by_name "credentials name"
```
