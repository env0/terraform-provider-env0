---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "env0_gcp_credentials Data Source - terraform-provider-env0"
subcategory: ""
description: |-
  
---

# env0_gcp_credentials (Data Source)



## Example Usage

```terraform
data "env0_gcp_credentials" "gcp_credentials_by_name" {
  name = "gcp credentials"
}

data "env0_gcp_credentials" "gcp_credentials_by_id" {
  id = "12345676safsd"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (String) the id of the credentials
- `name` (String) the name of the credentials
