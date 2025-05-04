package env0

import (
	"context"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOciCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOciCredentialsCreate,
		UpdateContext: resourceOciCredentialsUpdate,
		ReadContext:   resourceCredentialsRead(OCI_TYPE),
		DeleteContext: resourceCredentialsDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceCredentialsImport(OCI_TYPE)},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name for the oci credentials",
				Required:    true,
				ForceNew:    true,
			},
			"tenancy_ocid": {
				Type:        schema.TypeString,
				Description: "OCI tenancy OCID",
				Required:    true,
			},
			"user_ocid": {
				Type:        schema.TypeString,
				Description: "OCI user OCID",
				Required:    true,
			},
			"fingerprint": {
				Type:        schema.TypeString,
				Description: "OCI API key fingerprint",
				Required:    true,
			},
			"private_key": {
				Type:        schema.TypeString,
				Description: "OCI API private key",
				Required:    true,
				Sensitive:   true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "OCI region",
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

func ociCredentialsGetValue(d *schema.ResourceData) (client.OciCredentialsValuePayload, error) {
	value := client.OciCredentialsValuePayload{}
	if err := readResourceData(&value, d); err != nil {
		return value, fmt.Errorf("schema resource data deserialization failed: %w", err)
	}
	return value, nil
}

func resourceOciCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	value, err := ociCredentialsGetValue(d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := client.OciCredentialsCreatePayload{
		Name:      d.Get("name").(string),
		Value:     value,
		Type:      client.OciApiKeyCredentialsType,
		ProjectId: d.Get("project_id").(string),
	}

	credentials, err := apiClient.CredentialsCreate(&request)
	if err != nil {
		return diag.Errorf("could not create oci credentials: %v", err)
	}

	d.SetId(credentials.Id)

	return nil
}

func resourceOciCredentialsUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	value, err := ociCredentialsGetValue(d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := client.OciCredentialsCreatePayload{
		Value: value,
		Type:  client.OciApiKeyCredentialsType,
	}

	if _, err := apiClient.CredentialsUpdate(d.Id(), &request); err != nil {
		return diag.Errorf("could not update oci credentials: %s %v", d.Id(), err)
	}

	return nil
}
