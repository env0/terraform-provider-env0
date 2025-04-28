package env0

import (
	"context"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVaultOidcCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVaultOidcCredentialsCreate,
		UpdateContext: resourceVaultOidcCredentialsUpdate,
		ReadContext:   resourceCredentialsRead(VAULT_OIDC_TYPE),
		DeleteContext: resourceCredentialsDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceCredentialsImport(VAULT_OIDC_TYPE)},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name for the oidc credentials",
				Required:    true,
				ForceNew:    true,
			},
			"address": {
				Type:        schema.TypeString,
				Description: "the vault address, including port",
				Required:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "the vault version to use",
				Required:    true,
			},
			"role_name": {
				Type:        schema.TypeString,
				Description: "the vault role name",
				Required:    true,
			},
			"jwt_auth_backend_path": {
				Type:        schema.TypeString,
				Description: "path to the new authentication method",
				Required:    true,
			},
			"namespace": {
				Type:        schema.TypeString,
				Description: "an optional vault namespace",
				Optional:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "the env0 project id to associate the credentials with",
				Optional:    true,
				ForceNew:    true,
			},
		},
	}
}

func vaultOidcCredentialsGetValue(d *schema.ResourceData) (client.VaultCredentialsValuePayload, error) {
	var value client.VaultCredentialsValuePayload

	if err := readResourceData(&value, d); err != nil {
		return value, fmt.Errorf("schema resource data deserialization failed: %w", err)
	}

	return value, nil
}

func resourceVaultOidcCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	value, err := vaultOidcCredentialsGetValue(d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := client.VaultCredentialsCreatePayload{
		Name:      d.Get("name").(string),
		Value:     value,
		Type:      client.VaultOidcCredentialsType,
		ProjectId: d.Get("project_id").(string),
	}

	credentials, err := apiClient.CredentialsCreate(&request)
	if err != nil {
		return diag.Errorf("could not create vault oidc credentials: %v", err)
	}

	d.SetId(credentials.Id)

	return nil
}

func resourceVaultOidcCredentialsUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	value, err := vaultOidcCredentialsGetValue(d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := client.VaultCredentialsCreatePayload{
		Value: value,
		Type:  client.VaultOidcCredentialsType,
	}

	if _, err := apiClient.CredentialsUpdate(d.Id(), &request); err != nil {
		return diag.Errorf("could not update vault oidc credentials: %s %v", d.Id(), err)
	}

	return nil
}
