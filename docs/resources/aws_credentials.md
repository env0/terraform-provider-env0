---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "env0_aws_credentials Resource - terraform-provider-env0"
subcategory: ""
description: |-
  
---

# env0_aws_credentials (Resource)



## Example Usage

```terraform
resource "env0_aws_credentials" "credentials" {
  name = "example"
  arn  = "Example role ARN"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) name for the credentials

### Optional

- `access_key_id` (String, Sensitive) the aws access key id
- `arn` (String) the aws role arn
- `duration` (Number) the session duration in seconds for AWS_ASSUMED_ROLE_FOR_DEPLOYMENT. If set must be one of the following: 3600 (1h), 7200 (2h), 14400 (4h), 18000 (5h default), 28800 (8h), 43200 (12h)
- `project_id` (String) the env0 project id to associate the credentials with
- `secret_access_key` (String, Sensitive) the aws access key secret. In case your organization is self-hosted, please use a secret reference in the shape of ${ssm:<secret-id>}

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import env0_aws_credentials.by_id d31a6b30-5f69-4d24-937c-22322754934e
terraform import env0_aws_credentials.by_name "credentials name"
```
