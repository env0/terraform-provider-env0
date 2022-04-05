<a href="https://env0.com">
    <img src=".github/env0_logo.svg" alt="env0 logo" title="env0" align="right" height="40" />
</a>

# Terraform Provider for env0

[![Go Report Card](https://goreportcard.com/badge/github.com/env0/terraform-provider-env0)](https://goreportcard.com/report/github.com/env0/terraform-provider-env0)

- [Documentation](https://registry.terraform.io/providers/env0/env0/latest/docs)
- [Usage Examples](https://github.com/env0/terraform-provider-env0/tree/main/examples)

## Quick Start

```terraform
terraform {
  required_providers {
    env0 = {
      source = "env0/env0"
    }
  }
}

provider "env0" {}

data "env0_project" "default_project" {
  name = "My First Project"
}

resource "env0_template" "example" {
  name        = "example"
  description = "Example template"
  repository  = "https://github.com/env0/templates"
  path        = "aws/hello-world"
}

resource "env0_configuration_variable" "in_a_template" {
  name        = "VARIABLE_NAME"
  value       = "some value"
  template_id = env0_template.tested1.id
}
```

## Authentication

1. Generate an `api_key` and `api_secret` from the Organization Settings page.
See [here](https://developer.env0.com/docs/api/YXBpOjY4Njc2-env0-api#creating-an-api-key).

2. These can be provided by one of two methods:

   1. Set `ENV0_API_KEY` and `ENV0_API_SECRET` environment variables, and just declaring the provider with no parameters:

    ```terraform
       provider "env0" {}
    ```

   2. Specify these fields as parameters to the provider:

    ```terraform
    variable "env0_api_key" {}
    variable "env0_api_secret" {}

    provider "env0" {
      api_key = var.env0_api_key
      api_secret = var.env0_api_secret
    }
    ```

## Development Setup

> **Supported Go Version: 1.16**

### Build
- Use the `./build.sh` script.
- The output binary is called `terraform-provider-env0`

### Run local version of the provider
- Build - `./build.sh`
- Create the plugins folder - `mkdir -p ~/.terraform.d/plugins/terraform.env0.com/local/env0/6.6.6/darwin_amd64`
- Copy the built binary - `cp ./terraform-provider-env0 ~/.terraform.d/plugins/terraform.env0.com/local/env0/6.6.6/darwin_amd64` (Replace `darwin` with `linux` on Linux)
- Require the local provider in your `main.tf` -
```
terraform {
  required_providers {
    env0 = {
      version = "6.6.6"
      source  = "terraform.env0.com/local/env0"
    }
  }
}
```

## Testing

### Integration tests
- The integration tests run against the real env0 API
- Have `ENV0_API_KEY` and `ENV0_API_SECRET` environment variables defined.
- Also set `ENV0_API_ENDPOINT` if you want to run against a non-prod environment.
- Run `go run tests/harness.go` (from the project root folder) to run all the tests.
- Use `go run tests/harness.go 003_configuration_variable` to run a specific test.

Each test perform the following steps:
- `terraform init`
- `terraform apply -auto-approve -var second_run=0`
- `terraform apply -auto-approve -var second_run=1`
- `terraform outputs -json` - and verifies expected outputs from `expected_outputs.json`
- `terraform destroy`

The harness has two modes to help while developing: If an environment variable `DESTROY_MODE` exists and it's value is `NO_DESTROY`, the harness will avoid calling `terraform destroy`, allowing the developer to inspect the resources created, through the dashboard, for example.
Afterwards, when cleanup is required, just set `DESTROY_MODE` to `DESTROY_ONLY` and _only_ `terraform destroy` will run.

#### Integration Test Prerequisites
- An env0 organization
- An API Key
- A Github.com integrated template name `Github Integrated Template` for https://github.com/env0/templates
- A Gitlab.com integrated template name `Gitlab Integrated Template` for https://gitlab.com/env0/gitlab-vcs-integration-tests

### Unit Testing
#### How to run tests
Run from root directory:
```shell
go test ./...
```

#### Running a single provider test
```shell
export TEST_PATTERN="TestUnitConfigurationVariableResource/Create" && go test ./env0
```

#### How to use mocks

1. Make sure `GOPATH` is in your `PATH`
```shell
go env GOPATH
echo $PATH
export PATH=$PATH:$(go env GOPATH)  # if not
```
2. Install mockgen
```
go install github.com/golang/mock/mockgen@v1.6.0
```
3. Make sure to add this line in files that include the interface you'd wish to mock:
```
//go:generate mockgen -destination=<file>_mock.go -package=<package> . <interface>
```
4. Run from root directory:
```shell
go generate ./...
```

## Documentation
- Docs are generated using github.com/hashicorp/terraform-plugin-docs
- Run `./generate-docs.sh` to generate docs
- Must be run manually before releasing a version
- Please add an example to `examples/<resources or data-sources>/env0_<name>` dir and make sure it is added to the docs.

## Release
To release a version to the [Terraform Public Registry](https://registry.terraform.io/providers/env0/env0/latest?pollNotifications=true):
1. Validate that all status checks are ✅ on `main` branch (specifically that docs generation is complete)
2. Pull from remote first - `git pull origin main`
3. Create and push a tag **locally**, in semver format - `git tag v0.0.9 && git push origin --tags`
4. New release with binaries **will be automatically generated** by the GitHub action defined in `.github/workflows/release.yml`.
5. The Registry will automatically pick up on the new version.
