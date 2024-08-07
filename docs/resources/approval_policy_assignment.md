---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "env0_approval_policy_assignment Resource - terraform-provider-env0"
subcategory: ""
description: |-
  
---

# env0_approval_policy_assignment (Resource)



## Example Usage

```terraform
data "env0_project" "project" {
  name = "project"
}

resource "env0_approval_policy" "approval_policy" {
  name                   = "approval policy"
  repository             = "reopository"
  github_installation_id = 4234234234
  path                   = "misc/null-resource"

}

resource "env0_approval_policy_assignment" "approval_policy_assignment" {
  scope        = "PROJECT"
  scope_id     = data.env0_project.project.id
  blueprint_id = env0_approval_policy.approval_policy.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `blueprint_id` (String) the id of the approval policy
- `scope_id` (String) the id of the scope (E.g. project id or template id)

### Optional

- `scope` (String) the type of the scope. Valid values: PROJECT or BLUEPRINT. Default value: PROJECT

### Read-Only

- `id` (String) The ID of this resource.
