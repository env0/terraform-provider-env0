package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAzureCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureCredentialsCreate,
		ReadContext:   resourceCredentialsRead(AZURE_TYPE),
		DeleteContext: resourceCredentialsDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceCredentialsImport(AZURE_TYPE)},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name for the credentials",
				Required:    true,
				ForceNew:    true,
			},
			"client_id": {
				Type:        schema.TypeString,
				Description: "the azure client id",
				Required:    true,
				ForceNew:    true,
			},
			"client_secret": {
				Type:        schema.TypeString,
				Description: "the azure client secret. In case your organization is self-hosted, please use a secret reference in the shape of ${azure:<secret-name>@<vault-name>}",
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Description: "the azure subscription id",
				Required:    true,
				ForceNew:    true,
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Description: "the azure tenant id",
				Required:    true,
				ForceNew:    true,
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

func resourceAzureCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	value := client.AzureCredentialsValuePayload{}
	if err := readResourceData(&value, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	apiClient := meta.(client.ApiClientInterface)

	requestType := client.AzureServicePrincipalCredentialsType

	request := client.AzureCredentialsCreatePayload{
		Name:      d.Get("name").(string),
		Value:     value,
		Type:      requestType,
		ProjectId: d.Get("project_id").(string),
	}

	credentials, err := apiClient.CredentialsCreate(&request)
	if err != nil {
		return diag.Errorf("could not create credentials key: %v", err)
	}

	d.SetId(credentials.Id)

	return nil
}
