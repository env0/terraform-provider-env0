package env0

import (
	"context"
	"errors"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CloudType string

const (
	AWS_TYPE        CloudType = "aws"
	AWS_COST_TYPE   CloudType = "aws_cost"
	AWS_OIDC_TYPE   CloudType = "aws_oidc"
	AZURE_TYPE      CloudType = "azure"
	AZURE_COST_TYPE CloudType = "azure_cost"
	AZURE_OIDC_TYPE CloudType = "azure_oidc"
	GCP_TYPE        CloudType = "gcp"
	GCP_COST_TYPE   CloudType = "google_cost"
)

var credentialsTypeToPrefixList map[CloudType][]string = map[CloudType][]string{
	AWS_TYPE:        {string(client.AwsAssumedRoleCredentialsType), string(client.AwsAccessKeysCredentialsType)},
	AWS_COST_TYPE:   {string(client.AwsCostCredentialsType)},
	AWS_OIDC_TYPE:   {string(client.AwsOidcCredentialsType)},
	AZURE_TYPE:      {string(client.AzureServicePrincipalCredentialsType)},
	AZURE_COST_TYPE: {string(client.AzureCostCredentialsType)},
	AZURE_OIDC_TYPE: {string(client.AzureOidcCredentialsType)},
	GCP_TYPE:        {string(client.GcpServiceAccountCredentialsType)},
	GCP_COST_TYPE:   {string(client.GoogleCostCredentialsType)},
}

func getCredentialsByName(name string, prefixList []string, meta interface{}) (client.Credentials, error) {
	apiClient := meta.(client.ApiClientInterface)

	credentialsList, err := apiClient.CloudCredentialsList()
	if err != nil {
		return client.Credentials{}, err
	}

	var foundCredentials []client.Credentials
	for _, credentials := range credentialsList {
		if credentials.Name == name && credentials.HasPrefix(prefixList) {
			foundCredentials = append(foundCredentials, credentials)
		}
	}

	if len(foundCredentials) == 0 {
		return client.Credentials{}, fmt.Errorf("credentials with name %v not found", name)
	}

	if len(foundCredentials) > 1 {
		return client.Credentials{}, fmt.Errorf("found multiple credentials with name: %s. Use id instead or make sure credential names are unique %v", name, foundCredentials)
	}

	return foundCredentials[0], nil
}

func getCredentialsById(id string, prefixList []string, meta interface{}) (client.Credentials, error) {
	apiClient := meta.(client.ApiClientInterface)
	credentials, err := apiClient.CloudCredentials(id)
	if err != nil {
		if _, ok := err.(*client.NotFoundError); ok {
			return client.Credentials{}, errors.New("credentials not found")
		}
		return client.Credentials{}, err
	}

	if !credentials.HasPrefix(prefixList) {
		return client.Credentials{}, fmt.Errorf("credentials type mistmatch %s", credentials.Type)
	}

	return credentials, nil
}

func getCredentials(ctx context.Context, id string, prefixList []string, meta interface{}) (client.Credentials, error) {
	_, err := uuid.Parse(id)
	if err == nil {
		tflog.Info(ctx, "Resolving credentials by id", map[string]interface{}{"id": id})
		return getCredentialsById(id, prefixList, meta)
	} else {
		tflog.Info(ctx, "Resolving credentials by name", map[string]interface{}{"name": id})
		return getCredentialsByName(id, prefixList, meta)
	}
}

func resourceCredentialsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	err := apiClient.CloudCredentialsDelete(id)
	if err != nil {
		return diag.Errorf("could not delete credentials: %v", err)
	}
	return nil
}

func resourceCredentialsRead(cloudType CloudType) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(client.ApiClientInterface)

		credentials, err := apiClient.CloudCredentials(d.Id())
		if err != nil {
			return ResourceGetFailure(ctx, string(cloudType)+" credentials", d, err)
		}

		if err := writeResourceData(&credentials, d); err != nil {
			return diag.Errorf("schema resource data serialization failed: %v", err)
		}

		return nil
	}
}

func resourceCredentialsImport(cloudType CloudType) schema.StateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
		credentials, err := getCredentials(ctx, d.Id(), credentialsTypeToPrefixList[cloudType], meta)
		if err != nil {
			if _, ok := err.(*client.NotFoundError); ok {
				return nil, fmt.Errorf(string(cloudType)+" credentials resource with id %v not found", d.Id())
			}
			return nil, err
		}

		if err := writeResourceData(&credentials, d); err != nil {
			return nil, fmt.Errorf("schema resource data serialization failed: %v", err)
		}

		return []*schema.ResourceData{d}, nil
	}
}
