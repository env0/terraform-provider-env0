---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "env0_gcp_cloud_configuration Resource - terraform-provider-env0"
subcategory: ""
description: |-
  configure a GCP cloud account (Cloud Compass)
---

# env0_gcp_cloud_configuration (Resource)

configure a GCP cloud account (Cloud Compass)

## Example Usage

```terraform
resource "env0_gcp_cloud_configuration" "example" {
  name                                 = "example-gcp-config"
  gcp_project_id                       = "your-gcp-project-id"
  credential_configuration_file_content = file("path/to/your-gcp-service-account.json")
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `credential_configuration_file_content` (String, Sensitive) the GCP credential configuration file content (JSON)
- `gcp_project_id` (String) the GCP project ID
- `name` (String) name for the cloud configuration for insights

### Read-Only

- `health` (Boolean) an indicator if the configuration is valid
- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash
# Example import script for env0_gcp_cloud_configuration
terraform import env0_gcp_cloud_configuration.example <cloud_configuration_id_or_name>
```
