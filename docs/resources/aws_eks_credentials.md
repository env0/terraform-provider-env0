---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "env0_aws_eks_credentials Resource - terraform-provider-env0"
subcategory: ""
description: |-
  
---

# env0_aws_eks_credentials (Resource)



## Example Usage

```terraform
resource "env0_aws_eks_credentials" "credentials" {
  name           = "example"
  cluster_name   = "my-cluster"
  cluster_region = "us-east-2"
}

data "env0_project" "project" {
  name = "my-project"
}

resource "env0_cloud_credentials_project_assignment" "assignment" {
  credential_id = env0_aws_eks_credentials.credentials.id
  project_id    = data.env0_project.project.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cluster_name` (String) eks cluster name
- `cluster_region` (String) the AWS region of the eks cluster
- `name` (String) name for the credentials

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import env0_aws_eks_credentials.by_id d31a6b30-5f69-4d24-937c-22322754934e
terraform import env0_aws_eks_credentials.by_name "credentials name"
```
