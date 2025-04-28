package env0

import (
	"context"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAzureOidcCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureOidcCredentialsCreate,
		UpdateContext: resourceAzureOidcCredentialsUpdate,
		ReadContext:   resourceCredentialsRead(AZURE_OIDC_TYPE),
		DeleteContext: resourceCredentialsDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceCredentialsImport(AZURE_OIDC_TYPE)},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name for the oidc credentials",
				Required:    true,
				ForceNew:    true,
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Description: "the azure subscription id",
				Required:    true,
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Description: "the azure tenant id",
				Required:    true,
			},
			"client_id": {
				Type:        schema.TypeString,
				Description: "the azure client id",
				Required:    true,
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

func azureOidcCredentialsGetValue(d *schema.ResourceData) (client.AzureCredentialsValuePayload, error) {
	value := client.AzureCredentialsValuePayload{}

	if err := readResourceData(&value, d); err != nil {
		return value, fmt.Errorf("schema resource data deserialization failed: %w", err)
	}

	return value, nil
}

func resourceAzureOidcCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	value, err := azureOidcCredentialsGetValue(d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := client.AzureCredentialsCreatePayload{
		Name:      d.Get("name").(string),
		Value:     value,
		Type:      client.AzureOidcCredentialsType,
		ProjectId: d.Get("project_id").(string),
	}

	credentials, err := apiClient.CredentialsCreate(&request)
	if err != nil {
		return diag.Errorf("could not create azure oidc credentials: %v", err)
	}

	d.SetId(credentials.Id)

	return nil
}

func resourceAzureOidcCredentialsUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	value, err := azureOidcCredentialsGetValue(d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := client.AzureCredentialsCreatePayload{
		Value: value,
		Type:  client.AzureOidcCredentialsType,
	}

	if _, err := apiClient.CredentialsUpdate(d.Id(), &request); err != nil {
		return diag.Errorf("could not update azure oidc credentials: %s %v", d.Id(), err)
	}

	return nil
}
