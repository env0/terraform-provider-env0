---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "env0_aws_oidc_credentials Data Source - terraform-provider-env0"
subcategory: ""
description: |-
  
---

# env0_aws_oidc_credentials (Data Source)



## Example Usage

```terraform
resource "env0_aws_oidc_credentials" "example" {
  name     = "name"
  role_arn = "role_arn"
}

data "env0_aws_oidc_credentials" "by_id" {
  id = env0_aws_oidc_credentials.example.id
}

data "env0_aws_oidc_credentials" "by_name" {
  name = env0_aws_oidc_credentials.example.name
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (String) the id of the aws_oidc oidc credentials
- `name` (String) the name of the aws_oidc oidc credentials
