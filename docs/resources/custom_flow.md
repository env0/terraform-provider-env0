---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "env0_custom_flow Resource - terraform-provider-env0"
subcategory: ""
description: |-
  
---

# env0_custom_flow (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) name for the custom flow
- `repository` (String) repository url for the custom flow source code

### Optional

- `bitbucket_client_key` (String) the bitbucket client key used for integration
- `github_installation_id` (Number) the env0 application installation id on the relevant github repository
- `gitlab_project_id` (Number) the project id of the relevant repository
- `is_azure_devops` (Boolean) true if this custom flow integrates with azure dev ops repository
- `is_bitbucket_server` (Boolean) true if this custom flow uses bitbucket server repository
- `is_github_enterprise` (Boolean) true if this custom flow uses github enterprise repository
- `is_gitlab` (Boolean) true if this custom flow integrates with gitlab repository
- `is_gitlab_enterprise` (Boolean) true if this custom flow uses gitlab enterprise repository
- `path` (String) terraform / terragrunt file folder inside source code. Should be the full path including the .yaml/.yml file
- `revision` (String) source code revision (branch / tag) to use
- `ssh_keys` (List of Map of String) an array of references to 'data_ssh_key' to use when accessing git over ssh
- `token_id` (String) the git token id to be used

### Read-Only

- `id` (String) id of the custom flow

