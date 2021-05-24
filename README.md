# terraform-provider-env0

Terraform provider to interact with env0

Available in the [Terraform Registry](https://registry.terraform.io/providers/env0/env0/latest)

The full list of supported resources is available [here](#resources).

## Example usage

```terraform
terraform {
  required_providers {
    env0 = {
      source = "env0/env0"
      version = "0.0.2"
    }
  }
}

provider "env0" {}

data "env0_project" "default_project" {
  name = "Default Organization Project"
}

resource "env0_template" "example" {
  name        = "example"
  description = "Example template"
  repository  = "https://github.com/env0/templates"
  path        = "aws/hello-world"
  project_ids = [data.env0_project.default_project.id]
}

resource "env0_configuration_variable" "in_a_template" {
  name        = "VARIABLE_NAME"
  value       = "some value"
  template_id = env0_template.tested1.id
}
```

## Authentication

First, generate an `api_key` and `api_secret` from the organization settings page.
See [here](https://docs.env0.com/reference#authentication).

These can be provided by one of two methods. First method consists of setting `ENV0_API_KEY` and `ENV0_API_SECRET` environment variables, and just declaring the provider with no parameters:

```terraform
provider "env0" {}
```

The second method would be to specify these fields as parameters to the provider:

```terraform
variable "env0_api_key" {}
variable "env0_api_secret" {}

provider "env0" {
    api_key = var.env0_api_key
    api_secret = var.env0_api_secret
}
```

## Resources

The env0 Terraform provider provides the following building blocks:

- `env0_organization` - [data source](#env0_organization-data-source)
- `env0_project` - [data source](#env0_project-data-source) and [resource](#env0_project-resource)
- `env0_configuration_variable` - [data source](#env0_configuration_variable-data-source) and [resource](#env0_configuration_variable-resource)
- `env0_template` - [data source](#env0_template-data-source) and [resource](#env0_template-resource)
- `ssh_key` - [data source](#env0_ssh_key-data-source) and [resource](#env0_ssh_key-resource)

### `env0_ssh_key` resource

Define a new ssh key.

#### Example usage

```terraform
resource "tls_private_key" "throwaway" {
  algorithm = "RSA"
}
output "public_key_you_need_to_add_to_github_ssh_keys" {
  value = tls_private_key.throwaway.public_key_openssh
}

resource "env0_ssh_key" "tested" {
  name  = "test key"
  value = tls_private_key.throwaway.private_key_pem
}

data "env0_ssh_key" "tested" {
  name       = "test key"
  depends_on = [env0_ssh_key.tested]
}
```

#### Argument reference

The following arguments are supported:

- `name` - Name for the ssh key;
- `value` - Value of the key;


[^ Back to all resources](#resources)

### `env0_ssh_key` data source

Fetch metadata associated with an existing ssh key.

#### Example usage

```terraform
data "env0_ssh_key" "my_key" {
  name = "Secret Key"
}

resource "env0_template" "example" {
  # ...
  ssh_keys = [data.env0_ssh_key.my_keys]
}
```

#### Argument reference

The following arguments are supported:

- `id` - (Required if name is not set) - Fetch ssh key by id;
- `name` - (Required if id not set, mutally exclusive) - Look for the first ssh key that matches said name;


[^ Back to all resources](#resources)

### `env0_template` resource

Define a new template in the organization

#### Example usage

```terraform
data "env0_project" "default_project" {
  name = "Default Organization Project"
}
resource "env0_template" "example" {
  name        = "example"
  description = "Example template"
  repository  = "https://github.com/env0/templates"
  path        = "aws/hello-world"
  project_ids = [data.env0_project.default_project.id]
  ssh_keys = [data.ssh_keys.my_keys]
}
```

#### Argument reference

The following arguments are supported:

- `name` - (Required) - name to give the template;
- `description` - (Optional) - description for the template;
- `repository` - (Required) - git repository for the template source code;
- `path` - (Optional, default "/") - terraform / terragrunt file folder inside source code;
- `revision` - (Optional) - source code revision (branch / tag) to use;
- `type` - (Optional, default "terraform") - `terraform` or `terragrunt`;
- `project_ids` - (Optional) - a list of which projects may access this template (id of project);
- `ssh_keys` - (Optional) - an array of references to [`env0_ssh_key`](#env0_ssh_key-data-source) terraform data source to be assigned to this template;
- `retries_on_deploy` - (Optional) - number of times to retry when deploying an environment based on this template (between 1 and 3)
- `retry_on_deploy_only_when_matches_regex` - (Optional) - if specified, will only retry (on deploy) if error matches specified regex;
- `retries_on_destroy` - (Optional) - number of times to retry when destroying an environment based on this template (between 1 and 3)
- `retry_on_destroy_only_when_matches_regex` - (Optional) - if specified, will only retry (on destroy) if error matches specified regex;

#### Attributes reference

There are no additional attributes other than the arguments above.

[^ Back to all resources](#resources)

### `env0_configuration_variable` resource

A configuration variable is either an environment variable or a terraform variable. Configuration variables can configuration at the organization scope, project scope, template scope or environment scope. If two variables exists with the same name in two different scope, the more specific of the scopes is the value that will be used.

#### Example usage

```terraform
resource "env0_configuration_variable" "example" {
  name  = "ENVIRONMENT_VARIABLE_NAME"
  value = "example value"
}
```

#### Argument reference

The following arguments are supported:

- `name` - (Required) - Name of the variable;
- `value` - (Required) - Value for the variable;
- `is_sensitive` - (Optional, default false) - set variable to be sensitive;
- `type` - (Optional, default 'environment') - either `environment` or `terraform`;
- `enum` - (Optional) - list of strings, for possible values allowed for this variable, when overriding the value through the UI;
- `project_id` - (Optional, mutually exclusive) - define the variable under the project scope (by default, variable are created under the organization scope);
- `template_id` - (Optional, mutually exclusive) - define the variable under the template scope;
- `environment_id` - (Optional, mutually exclusive) - define the variable under the environment scope;

#### Attributes reference

In addition to all arguments above, the following attributes are exported:

- `name` - The name of the organization;

[^ Back to all resources](#resources)


### `env0_organization` data source

Each api key is associated with a single organization, so this resource can be used to fetch
that organization metadata.

#### Example usage

```terraform
data "env0_organization" "my_organization" {}

output "organization_name" {
  value = data.env0_organization.my_organization.name
}
```

#### Argument reference

No argument are supported

#### Attributes reference

In addition to all arguments above, the following attributes are exported:

- `name` - The name of the organization;
- `role` - The role of the api key in the organization;
- `is_self_hosted` - Is the organizaton self hosted;

[^ Back to all resources](#resources)

### `env0_project` data source

Fetch metadata associated with a existing project.

#### Example usage

```terraform
data "env0_project" "default_project" {
  name = "Default Organization Project"
}

output "project_id" {
  value = data.env0_project.default_project.id
}
```

#### Argument reference

The following arguments are supported:

- `id` - (Required if name is not set) - Fetch project by the project id;
- `name` - (Required if id not set, mutally exclusive) - Look for the first project that matches said name;

#### Attributes reference

In addition to all arguments above, the following attributes are exported:

- `role` - The role of the api_key in this project;

[^ Back to all resources](#resources)

### `env0_template` data source

Fetch metadata of already defined template.

#### Example usage

```terraform
data "env0_template" "example" {
  name = "Template Name"
}

output "template_id" {
  value = data.env0_template.example.id
}
```

#### Argument reference

The following arguments are supported:

- `id` - (Required if name is not set) - Fetch template by the template id;
- `name` - (Required if id not set, mutally exclusive) - Look for the first template that matches said name;

#### Attributes reference

In addition to all arguments above, the following attributes are exported:

- `repository` - template source code repository url;
- `path` - terraform / terrgrunt folder inside source code repository;
- `revision` - source code revision (branch / tag) to use;
- `type` - `terraform` or `terragrunt`;
- `project_ids` - which projects may access this template (id of project);
- `retries_on_deploy` - number of times to retry when deploying an environment based on this template;
- `retry_on_deploy_only_when_matches_regex` - will only retry (on deploy) if error matches specified regex;
- `retries_on_destroy` - number of times to retry when destroying an environment based on this template;
- `retry_on_destroy_only_when_matches_regex` - will only retry (on destroy) if error matches specified regex;

[^ Back to all resources](#resources)

### `env0_configuration_variable` data source

A configuration variable is either an environment variable or a terraform variable. Configuration variables can configuration at the organization scope, project scope, template scope or environment scope. If two variables exists with the same name in two different scope, the more specific of the scopes is the value that will be used.

This data source allows fetching existing configuration variables, and their values. Note that
fetching sensitive configuration variables will result in "******" as the variable value.

#### Example usage

```terraform
data "env0_configuration_variable" "aws_default_region" {
  name = "AWS_DEFAULT_REGION"
}

output "aws_default_region" {
  value = data.env0_configuration_variable.aws_default_region.value
}
```

#### Argument reference

The following arguments are supported:

- `id` - (Required, mutually exclusive) - the id of the variable;
- `name` - (Required, mutually exclusive) - the variable name;

#### Attributes reference

In addition to all arguments above, the following attributes are exported:

- `value` - value of the variable. will be '*********' if configuration variable is sensitive;
- `is_sensitive` - `true` if configuration variable is sensitive

[^ Back to all resources](#resources)

## Dev setup

To build locally, you can use the `./build.sh` script.
The rest of the steps below are relevant if you would like to use the provider outside of the test harness (the test harness performs the steps described here).

To use local binary version, you'll need to create a local terraform provider repository.
The simplest way to do so would be to create a folder on your disk.
Under that folder, copy the built provider binary to `terraform-registry.env0.com/env0/env0/6.6.6/linux_amd64/terraform-provider-env0` (note to replace linux_amd64 if using a different platform).
Then, create a `terraform.rc` file (location doesn't matter).
The content of said file should look like so:

```
provider_installation {
  filesystem_mirror {
    path    = "<absolute path to local repository folder>"
    include = ["terraform-registry.env0.com/*/*"]
  }
  direct {}
}
```

Finally, set an environment variable `TF_CLI_CONFIG_FILE` to point to the `terraform.rc` file created.
After that, `terraform init` should be able to locate the provider on disk.
To define this variable only when running terraform once, you can, in bash shell:

```bash
TF_CLI_CONFIG_FILE=<terraform.rc path> terraform init
```

## Testing

If you have `ENV0_API_KEY` and `ENV0_API_SECRET` environment variables defined, after building the provider locally, just run `go run tests/harness.go` to run all the tests. Make sure to run from the project root folder.

Use `go run tests/harness.go 003_configuration_variable` to run a specific test.
The last argument can also be specified as a full path, e.g., `tests/003_configuration_variable/`.

Each tests performs the following steps:

- `terraform init`
- `terraform apply -auto-approve -var second_run=0`
- `terraform apply -auto-approve -var second_run=1`
- `terraform outputs -json` - and verifies expected outputs from `expected_outputs.json`
- `terraform destroy`

The harness has two convineint modes to help while developing: If an environment variable `DESTROY_MODE` exists and it's value is `NO_DESTROY`, the harness will avoid calling `terraform destroy`, allowing the developer to inspect the resources created, through the dashboard, for example.
Afterwards, when cleanup is required, just set `DESTROY_MODE` to `DESTROY_ONLY` and *only* `terraform destroy` will run.
