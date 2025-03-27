package env0

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getCloudConfigurationFromSchema(d *schema.ResourceData, provider string) (any, error) {
	var configuration any

	// Find the correct type, create an instance of this type, and deserialize from the schema to the instance.

	switch provider {
	case "AWS":
		configuration = &client.AWSCloudAccountConfiguration{}
	case "AzureLAW":
		configuration = &client.AzureCloudAccountConfiguration{}
	default:
		return nil, fmt.Errorf("unhandled provider: %s", provider)
	}

	if err := readResourceData(configuration, d); err != nil {
		return nil, fmt.Errorf("schema resource data deserialization failed: %w", err)
	}

	return configuration, nil
}

func getCloudConfigurationByNameFromApi(apiClient client.ApiClientInterface, name string) (*client.CloudAccount, error) {
	cloudAccounts, err := apiClient.CloudAccounts()
	if err != nil {
		return nil, err
	}

	for i, cloudAccount := range cloudAccounts {
		if cloudAccount.Name == name {
			return &cloudAccounts[i], nil
		}
	}

	return nil, fmt.Errorf("cloud configuration called '%s' was not found", name)
}

func getCloudConfigurationFromApi(apiClient client.ApiClientInterface, id string) (*client.CloudAccount, error) {
	var err error

	var cloudAccount *client.CloudAccount

	if _, parseErr := uuid.Parse(id); parseErr != nil {
		// Get by name (used by import).
		cloudAccount, err = getCloudConfigurationByNameFromApi(apiClient, id)
	} else {
		// Get by id.
		cloudAccount, err = apiClient.CloudAccount(id)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get clound configuration: %w", err)
	}

	var configuration any

	// Find the correct type, marshal the interface to bytes, and unmarshal the bytes back to an instance of the correct type.

	switch cloudAccount.Provider {
	case "AWS":
		configuration = &client.AWSCloudAccountConfiguration{}
	case "AzureLAW":
		configuration = &client.AzureCloudAccountConfiguration{}
	default:
		return nil, fmt.Errorf("unhandled provider: %s", cloudAccount.Provider)
	}

	b, err := json.Marshal(cloudAccount.Configuration)
	if err != nil {
		return nil, fmt.Errorf("failed to json marshal %s configuration: %w", cloudAccount.Provider, err)
	}

	if err := json.Unmarshal(b, configuration); err != nil {
		return nil, fmt.Errorf("failed to json unmarshal %s configuration: %w", cloudAccount.Provider, err)
	}

	cloudAccount.Configuration = configuration

	return cloudAccount, nil
}

func createCloudConfiguration(d *schema.ResourceData, meta any, provider string) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var createPayload client.CloudAccountCreatePayload

	if err := readResourceData(&createPayload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	configuration, err := getCloudConfigurationFromSchema(d, provider)
	if err != nil {
		return diag.FromErr(err)
	}

	createPayload.Configuration = configuration
	createPayload.Provider = provider

	cloudAccount, err := apiClient.CloudAccountCreate(&createPayload)
	if err != nil {
		return diag.Errorf("failed to create a cloud configuration: %v", err)
	}

	d.Set("health", cloudAccount.Health)

	d.SetId(cloudAccount.Id)

	return nil
}

func readCloudConfiguration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	cloudAccount, err := getCloudConfigurationFromApi(apiClient, d.Id())
	if err != nil {
		return ResourceGetFailure(ctx, "cloud_configuration", d, err)
	}

	if err := writeResourceData(cloudAccount, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	if err := writeResourceData(cloudAccount.Configuration, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func updateCloudConfiguration(d *schema.ResourceData, meta any, provider string) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var updatePayload client.CloudAccountUpdatePayload

	if err := readResourceData(&updatePayload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	configuration, err := getCloudConfigurationFromSchema(d, provider)
	if err != nil {
		return diag.FromErr(err)
	}

	updatePayload.Configuration = configuration

	cloudAccount, err := apiClient.CloudAccountUpdate(d.Id(), &updatePayload)
	if err != nil {
		return diag.Errorf("failed to update cloud configuration: %v", err)
	}

	d.Set("health", cloudAccount.Health)

	return nil
}

func deleteCloudConfiguration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.CloudAccountDelete(d.Id()); err != nil {
		return diag.Errorf("failed to delete cloud configuration: %v", err)
	}

	return nil
}

func importCloudConfiguration(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	apiClient := meta.(client.ApiClientInterface)

	cloudAccount, err := getCloudConfigurationFromApi(apiClient, d.Id())
	if err != nil {
		return nil, err
	}

	if err := writeResourceData(cloudAccount, d); err != nil {
		return nil, fmt.Errorf("schema resource data serialization failed: %w", err)
	}

	if err := writeResourceData(cloudAccount.Configuration, d); err != nil {
		return nil, fmt.Errorf("schema resource data serialization failed: %w", err)
	}

	return []*schema.ResourceData{d}, nil
}
