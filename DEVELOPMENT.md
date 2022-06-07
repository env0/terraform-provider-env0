# Development Guidelines

This document intends to give guidelines for developing new resources for the env0 Terraform provider.

## Project Structure

* client - Contains code for using the [env0 API](https://developer.env0.com/docs/api) + unit tests.
* env0 - Contains code that implements the env0 Terraform "resources" and "data" + acceptance tests.
* docs - Contains auto-generated md files derived from the "resources" and "data" go files. Deployed as part of the terraform provider documentation.
* tests - Contains integration tests (executed by harness.go).
* examples - Contains examples files for resources and data. Deployed as part of the documentation. Deployed as part of the terraform provider documentation.

## Client

The first step to adding a new resource is implementing the client API calls under the folder client.

Review the API documentation. And pay attention to the following details:
* URL Path (path parameters) and HTTP method (GET, POST, PUT, PATCH, etc...).
* Request Body (pay special attention to required fields).
* Response Body.

The env0 website uses the API as well. Therefore, if examples are required. The easiest way to understand the API is to make the calls from the GUI itself. Use Chrome developer tools and review the relevant API requests and responses.

Create a new file under the client directory that describes the API call(s) being implemented. Use existing implementations as templates for implementing the new API call(s). Check [notification.go](./client/notification.go) for reference.

Each go file defines the models and implements all the relevant client API calls.

Add the new functions to the [APIClient interface](./client/api_client.go)

Finally, create unit tests for added functionality. Check [notification_test.go](./client/notification_test.go)

## Resource

After the client API calls are implemented, the next step is implementing the resource.

Create a new file under the env0 directory. Use existing implementations as templates for implementing the new resource. Check [resource_module.go](./env0/resource_module.go) for reference.

Start by defining the schema. Use the API documentation to identify:
* what fields are required vs. optional.
* what fields require custom validators (see existing [validators](./env0/validators.go))
* what fields force a new resource if modified.

For all resources: CreateContext, ReadContext and DeleteContext are required.
Most resources also implement UpdateContext and Importer (for imports).
Without UpdateContext, the resource is destroyed and created for every change.
The Importer is used for the Terraform import command.

Finally, create acceptance tests for the added functionality. Check [resource_module_test.go](./env0/resource_module_test.go)

### readResourceData and writeResourceData

The file [utils.go](./env0/utils.go) contains some very useful functions.
Especially useful are readResourceData and writeResourceData.

The readResourceData function receives a golang struct and a Terraform configuration. It reads the Terraform configuration values and copies them to the golang struct.

The writeResourceData function receives a golang struct and a Terraform configuration. It reads the golang struct values and copies them to the Terraform configuration.

Check [resource_module.go](./env0/resource_module.go) that uses the utilities vs [resource_environment.go](./env0/resource_environment.go) that does not.

Pay attention to the following caveats:
* The golang fields are in CamalCase, while the terraform fields are in snake_case. They must match. E.g., ProjectName (golang) == project_name (Terraform). To override the default CamalCase to snake_case conversion you may use the tag `tfschema`. To ignore a field set the `tfschema` tag value to `-`.

#### writeResourceDataSlice

The writeResourceDataSlice function receives a golang slice, a field name (of type list) and a terraform configuration.
It will try to automatically write the slice structs to the terraform configuration under the field name.

#### Important Notes

When using any of these functions be sure to test them well.
These are "best-effort" helpers that leverage golang's refelection. They will work well for most basic cases, but may fall short in complex scenarios.

### Handling drifts

If ReadContext is called and the resource isn't found by the current ID, it's required to reset the ID.
This will re-create the resource.

```
apiClient := meta.(client.ApiClientInterface)

module, err := apiClient.Module(d.Id())
if err != nil {
    return ResourceGetFailure("module", d, err)
}

.
.
.

func ResourceGetFailure(resourceName string, d *schema.ResourceData, err error) diag.Diagnostics {
    if frerr, ok := err.(*http.FailedResponseError); ok && frerr.NotFound() {
        log.Printf("[WARN] Drift Detected: Terraform will remove %s from state", d.Id())
        d.SetId("")
        return nil
    }

    return diag.Errorf("could not get %s: %v", resourceName, err)
}
```

## Data

In most cases, Terraform data is required as well.

Create the data file under the env0 directory.
Check [data_gcp_credentials.go](./env0/data_gcp_credentials.go) for an example.

Finally create acceptance tests for data. Check [data_gcp_credentials_test.go](./env0/data_gcp_credentials_test.go) for an example.

**Note** In most cases, data is calculated by name or id. Names may not be unique. If searching by name and more than one resource is returned, it's considered an error.
```
func getGcpCredentialsByName(name interface{}, meta interface{}) (client.Credentials, diag.Diagnostics) {
    apiClient := meta.(client.ApiClientInterface)
    credentialsList, err := apiClient.CloudCredentialsList()
    if err != nil {
        return client.Credentials{}, diag.Errorf("Could not query GCP Credentials by name: %v", err)
    }

    credentialsByNameAndType := make([]client.Credentials, 0)
    for _, candidate := range credentialsList {
        if candidate.Name == name.(string) && isValidGcpCredentialsType(candidate.Type) {
            credentialsByNameAndType = append(credentialsByNameAndType, candidate)
        }
    }

    if len(credentialsByNameAndType) > 1 {
        return client.Credentials{}, diag.Errorf("Found multiple GCP Credentials for name: %s", name)
    }
    if len(credentialsByNameAndType) == 0 {
        return client.Credentials{}, diag.Errorf("Could not find GCP Credentials with name: %s", name)
    }
    return credentialsByNameAndType[0], nil
}
```

## Integration Tests and Examples

If applicable, create an integration test for the new resource and data.
The folder tests/integration contains a list of folders. One folder for each resource.

Each numbered resource folder contains the following files:
* conf.tf - The provider configuration. In most cases, this doesn't change.
* main.tf - The terraform instructions to run.
* expected_outputs.json - The Terraform outputs (using the `output` Terraform functionality).

Note: if there is no expected output, set expected_outputs.json contents to `{}`.

Under the examples directory, add examples for the resource (resource.tf) and data (data-source.tf).
For `import` add the file `import.sh` in the resource directory (see example [here](./examples/resources/env0_api_key/import.sh)).
