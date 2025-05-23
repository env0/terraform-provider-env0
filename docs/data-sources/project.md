---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "env0_project Data Source - terraform-provider-env0"
subcategory: ""
description: |-
  
---

# env0_project (Data Source)



## Example Usage

```terraform
data "env0_project" "default_project" {
  name = "Default Organization Project"
}

data "env0_project" "with_parent_name_filter" {
  name                = "Default Organization Project"
  parent_project_name = "parent project name"
}

data "env0_project" "with_parent_id_filter" {
  name              = "Default Organization Project"
  parent_project_id = "parent-project-id"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (String) id of the project
- `name` (String) the name of the project
- `parent_project_id` (String) the id of the parent project. Can be used as a filter when there are multiple subprojects with the same name under different parent projects
- `parent_project_name` (String) the name of the parent project. Can be used as a filter when there are multiple subprojects with the same name under different parent projects
- `parent_project_path` (String) a path of ancestors projects divided by the prefix '|'. Can be used as a filter when there are multiple subprojects with the same name under different parent projects. For example: 'App|Dev|us-east-1' will search for a project with the hierarchy 'App -> Dev -> us-east-1' ('us-east-1' being the parent)

### Read-Only

- `created_by` (String) textual description of the entity who created the project
- `description` (String) textual description of the project
- `hierarchy` (String) the hierarchy of the project
- `role` (String) role of the authenticated user (through api key) in the project
