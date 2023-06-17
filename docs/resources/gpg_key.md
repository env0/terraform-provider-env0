---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "env0_gpg_key Resource - terraform-provider-env0"
subcategory: ""
description: |-
  
---

# env0_gpg_key (Resource)



## Example Usage

```terraform
resource "env0_gpg_key" "example" {
  name    = "gpg-key-example"
  key_id  = "ABCDABCDABCDABCD"
  content = "key block"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `content` (String) the gpg public key block
- `key_id` (String) the gpg key id
- `name` (String) the gpg key name

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import env0_gpg_key.by_id ddda7b30-6789-4d24-937c-22322754934e
terraform import env0_gpg_key.by_name gpg-key-name
```