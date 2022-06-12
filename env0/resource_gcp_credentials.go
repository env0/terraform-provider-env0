package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGcpCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGcpCredentialsCreate,
		ReadContext:   resourceCredentialsRead(GCP_TYPE),
		DeleteContext: resourceCredentialsDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceCredentialsImport(GCP_TYPE)},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name for the credentials",
				Required:    true,
				ForceNew:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "the gcp project id",
				Optional:    true,
				Sensitive:   true,
				ForceNew:    true,
			},
			"service_account_key": {
				Type:        schema.TypeString,
				Description: "the gcp service account key",
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
			},
		},
	}
}

func resourceGcpCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	value := client.GcpCredentialsValuePayload{}
	if err := readResourceData(&value, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	apiClient := meta.(client.ApiClientInterface)

	requestType := client.GcpServiceAccountCredentialsType

	request := client.GcpCredentialsCreatePayload{
		Name:  d.Get("name").(string),
		Value: value,
		Type:  requestType,
	}

	credentials, err := apiClient.CredentialsCreate(&request)
	if err != nil {
		return diag.Errorf("could not create credentials key: %v", err)
	}

	d.SetId(credentials.Id)

	return nil
}
