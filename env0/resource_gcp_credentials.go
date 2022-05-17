package env0

import (
	"context"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGcpCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGcpCredentialsCreate,
		ReadContext:   resourceGcpCredentialsRead,
		DeleteContext: resourceGcpCredentialsDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceGcpCredentialsImport},

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

	credentials, err := apiClient.GcpCredentialsCreate(request)
	if err != nil {
		return diag.Errorf("could not create credentials key: %v", err)
	}

	d.SetId(credentials.Id)

	return nil
}

func resourceGcpCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	credentials, err := apiClient.CloudCredentials(id)
	if err != nil {
		return ResourceGetFailure("gcp credentials", d, err)
	}

	if err := writeResourceData(&credentials, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceGcpCredentialsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	err := apiClient.CloudCredentialsDelete(id)
	if err != nil {
		return diag.Errorf("could not delete credentials: %v", err)
	}
	return nil
}

func resourceGcpCredentialsImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	credentials, err := getCredentials(d.Id(), "GCP_", meta)
	if err != nil {
		if _, ok := err.(*client.NotFoundError); ok {
			return nil, fmt.Errorf("gcp credentials resource with id %v not found", d.Id())
		}
		return nil, err
	}

	if err := writeResourceData(&credentials, d); err != nil {
		return nil, fmt.Errorf("schema resource data serialization failed: %v", err)
	}

	return []*schema.ResourceData{d}, nil
}
