---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "env0_environment_state_access Resource - terraform-provider-env0"
subcategory: ""
description: |-
  
---

# env0_environment_state_access (Resource)



## Example Usage

```terraform
data "env0_environment" "environment" {
  name = "Environment Name"
}

data "env0_project" "project" {
  name = "Project Name"
}

resource "env0_environment_state_access" "example_allowed_projects" {
  environment_id      = data.env0_environment.environment.id
  allowed_project_ids = [data.env0_project.project.id]
}

resource "env0_environment_state_access" "example_entire_organization" {
  environment_id                      = data.env0_environment.environment.id
  accessible_from_entire_organization = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) id of the environment

### Optional

- `accessible_from_entire_organization` (Boolean) when this parameter is 'false', allowed_project_ids should be provided. Defaults to 'false'
- `allowed_project_ids` (List of String) list of allowed project_ids. Used when 'accessible_from_entire_organization' is 'false'

### Read-Only

- `id` (String) The ID of this resource.

