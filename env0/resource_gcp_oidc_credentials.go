package env0

import (
	"context"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGcpOidcCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGcpOidcCredentialsCreate,
		UpdateContext: resourceGcpOidcCredentialsUpdate,
		ReadContext:   resourceCredentialsRead(GCP_OIDC_TYPE),
		DeleteContext: resourceCredentialsDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceCredentialsImport(GCP_OIDC_TYPE)},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name for the oidc credentials",
				Required:    true,
				ForceNew:    true,
			},
			"credential_configuration_file_content": {
				Type:        schema.TypeString,
				Description: "the JSON content of the JWT configuration file",
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

func gcpOidcCredentialsGetValue(d *schema.ResourceData) (client.GcpCredentialsValuePayload, error) {
	value := client.GcpCredentialsValuePayload{}

	if err := readResourceData(&value, d); err != nil {
		return value, fmt.Errorf("schema resource data deserialization failed: %w", err)
	}

	return value, nil
}

func resourceGcpOidcCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	value, err := gcpOidcCredentialsGetValue(d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := client.GcpCredentialsCreatePayload{
		Name:      d.Get("name").(string),
		Value:     value,
		Type:      client.GcpOidcCredentialsType,
		ProjectId: d.Get("project_id").(string),
	}

	credentials, err := apiClient.CredentialsCreate(&request)
	if err != nil {
		return diag.Errorf("could not create gcp oidc credentials: %v", err)
	}

	d.SetId(credentials.Id)

	return nil
}

func resourceGcpOidcCredentialsUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	value, err := gcpOidcCredentialsGetValue(d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := client.GcpCredentialsCreatePayload{
		Value: value,
		Type:  client.GcpOidcCredentialsType,
	}

	if _, err := apiClient.CredentialsUpdate(d.Id(), &request); err != nil {
		return diag.Errorf("could not update gcp oidc credentials: %s %v", d.Id(), err)
	}

	return nil
}
