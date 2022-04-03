package env0

import (
	"context"
	"fmt"
	"log"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApiKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApiKeyCreate,
		ReadContext:   resourceApiKeyRead,
		DeleteContext: resourceApiKeyDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceApiKeyImport},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "the api key name",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceApiKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.ApiKeyCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	apiKey, err := apiClient.ApiKeyCreate(payload)
	if err != nil {
		return diag.Errorf("could not create api key: %v", err)
	}

	d.SetId(apiKey.Id)

	return nil
}

func resourceApiKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiKey, err := getApiKeyById(d.Id(), meta)
	if err != nil {
		return diag.Errorf("could not get api key: %v", err)
	}
	if apiKey == nil {
		log.Printf("[WARN] Drift Detected: Terraform will remove %s from state", d.Id())
		d.SetId("")
		return nil
	}

	if err := writeResourceData(apiKey, d); err != nil {
		diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceApiKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.ApiKeyDelete(d.Id()); err != nil {
		return diag.Errorf("could not delete api key: %v", err)
	}

	return nil
}

func getApiKeyById(id string, meta interface{}) (*client.ApiKey, error) {
	apiClient := meta.(client.ApiClientInterface)

	apiKeys, err := apiClient.ApiKeys()
	if err != nil {
		return nil, err
	}
	for _, apiKey := range apiKeys {
		if apiKey.Id == id {
			return &apiKey, nil
		}
	}
	return nil, nil
}

func getApiKeyByName(name string, meta interface{}) (*client.ApiKey, error) {
	apiClient := meta.(client.ApiClientInterface)

	apiKeys, err := apiClient.ApiKeys()
	if err != nil {
		return nil, err
	}

	var foundApiKeys []client.ApiKey
	for _, apiKey := range apiKeys {
		if apiKey.Name == name {
			foundApiKeys = append(foundApiKeys, apiKey)
		}
	}

	if len(foundApiKeys) == 0 {
		return nil, fmt.Errorf("api key with name %v not found", name)
	}

	if len(foundApiKeys) > 1 {
		return nil, fmt.Errorf("found multiple api keys with name: %s. Use id instead or make sure api key names are unique %v", name, foundApiKeys)
	}

	return &foundApiKeys[0], nil
}

func getApiKey(id string, meta interface{}) (*client.ApiKey, error) {
	_, err := uuid.Parse(id)
	if err == nil {
		log.Println("[INFO] Resolving api key by id: ", id)
		return getApiKeyById(id, meta)
	} else {
		log.Println("[INFO] Resolving api key by name: ", id)
		return getApiKeyByName(id, meta)
	}
}

func resourceApiKeyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	apiKey, err := getApiKey(d.Id(), meta)
	if err != nil {
		return nil, err
	}
	if apiKey == nil {
		return nil, fmt.Errorf("api key with id %v not found", d.Id())
	}

	if err := writeResourceData(apiKey, d); err != nil {
		diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return []*schema.ResourceData{d}, nil
}
