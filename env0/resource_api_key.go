package env0

import (
	"context"
	"errors"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	organizationRoleAdmin = "Admin"
	organizationRoleUser  = "User"
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
				Type:        schema.TypeString,
				Description: "the api key type. 'Admin', 'User' or a custom role id. Defaults to 'Admin'. For more details check https://docs.env0.com/docs/api-keys",
				Default:     organizationRoleAdmin,
				Optional:    true,
				ForceNew:    true,
				ValidateDiagFunc: func(i any, p cty.Path) diag.Diagnostics {
					val := i.(string)
					if val != organizationRoleAdmin && val != organizationRoleUser {
						_, err := uuid.Parse(val)
						if err != nil {
							return diag.Errorf("Organization role must be either 'Admin', 'User' or a custom role id, got: %s", val)
						}
					}

					return nil
				},
			},
			"project_permissions": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "Project-specific permissions. Only valid when organization_role is 'User' or a custom role id",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The project ID to assign permissions to",
						},
						"project_role": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							Description:      "The role for this project. Must be one of: Planner, Viewer, Deployer, Admin",
							ValidateDiagFunc: NewStringInValidator([]string{"Planner", "Viewer", "Deployer", "Admin"}),
						},
					},
				},
			},
			"omit_api_key_secret": {
				Type:        schema.TypeBool,
				Description: "if set to 'true' will omit the api_key_secret from the state. This would mean that the api_key_secret cannot be used",
				Optional:    true,
				ForceNew:    true,
			},
			"api_key_secret": {
				Type:        schema.TypeString,
				Description: "the api key secret. This attribute is not computed for imported resources. Note that this will be written to the state file. To omit the secret: set 'omit_api_key_secret' to 'true'",
				Computed:    true,
				Sensitive:   true,
			},
			"api_key_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the api key id",
			},
		},
	}
}

func resourceApiKeyCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	organizationRole := d.Get("organization_role").(string)
	if organizationRole == organizationRoleAdmin && len(d.Get("project_permissions").(*schema.Set).List()) > 0 {
		return diag.Errorf("project_permissions cannot be set when organization_role is Admin")
	}

	var payload client.ApiKeyCreatePayload

	payload.Name = d.Get("name").(string)

	if err := readResourceData(&payload.Permissions, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

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

func resourceApiKeyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiKey, err := getApiKeyById(d.Id(), meta)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]any{"id": d.Id()})
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

func resourceApiKeyDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.ApiKeyDelete(d.Id()); err != nil {
		return diag.Errorf("could not delete api key: %v", err)
	}

	return nil
}

func getApiKeyById(id string, meta any) (*client.ApiKey, error) {
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

func getApiKeyByName(name string, meta any) (*client.ApiKey, error) {
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

func getApiKey(ctx context.Context, id string, meta any) (*client.ApiKey, error) {
	_, err := uuid.Parse(id)
	if err == nil {
		tflog.Info(ctx, "Resolving api key by id", map[string]any{"id": id})

		return getApiKeyById(id, meta)
	} else {
		tflog.Info(ctx, "Resolving api key by name", map[string]any{"name": id})

		return getApiKeyByName(id, meta)
	}
}

func resourceApiKeyImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
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
