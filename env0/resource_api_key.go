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
			"organization_role": {
				Type:             schema.TypeString,
				Description:      "the api key type. 'Admin' or 'User'. Defaults to 'Admin'. For more details check https://docs.env0.com/docs/api-keys",
				Default:          "Admin",
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: NewStringInValidator([]string{"Admin", "User"}),
			},
			"api_key_secret": {
				Type:        schema.TypeString,
				Description: "the api key secret. This attribute is not computed for imported resources. Note that this will be written to the state file. To omit the secret: set 'omit_api_key_secret' to 'true'",
				Computed:    true,
				Sensitive:   true,
			},
			"omit_api_key_secret": {
				Type:        schema.TypeBool,
				Description: "if set to 'true' will omit the api_key_secret from the state. This would mean that the api_key_secret cannot be used",
				Optional:    true,
				ForceNew:    true,
			},
			"api_key_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the api key id",
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

	organizationRole := d.Get("organization_role").(string)
	payload.Permissions.OrganizationRole = organizationRole

	apiKey, err := apiClient.ApiKeyCreate(payload)
	if err != nil {
		return diag.Errorf("could not create api key: %v", err)
	}

	if err := writeResourceData(apiKey, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	if omit, ok := d.GetOk("omit_api_key_secret"); ok && omit.(bool) {
		d.Set("api_key_secret", "omitted")
	}

	return nil
}

func resourceApiKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiKey, err := getApiKeyById(d.Id(), meta)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]interface{}{"id": d.Id()})
			d.SetId("")

			return nil
		}

		return diag.Errorf("could not get api key: %v", err)
	}

	apiKey.ApiKeySecret = "" // Don't override the api key secret currently in the state.

	if err := writeResourceData(apiKey, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
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

	return nil, ErrNotFound
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

func getApiKey(ctx context.Context, id string, meta interface{}) (*client.ApiKey, error) {
	_, err := uuid.Parse(id)
	if err == nil {
		tflog.Info(ctx, "Resolving api key by id", map[string]interface{}{"id": id})

		return getApiKeyById(id, meta)
	} else {
		tflog.Info(ctx, "Resolving api key by name", map[string]interface{}{"name": id})

		return getApiKeyByName(id, meta)
	}
}

func resourceApiKeyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	apiKey, err := getApiKey(ctx, d.Id(), meta)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, fmt.Errorf("api key with id %v not found", d.Id())
		}

		return nil, err
	}

	if err := writeResourceData(apiKey, d); err != nil {
		return nil, fmt.Errorf("schema resource data serialization failed: %w", err)
	}

	return []*schema.ResourceData{d}, nil
}
