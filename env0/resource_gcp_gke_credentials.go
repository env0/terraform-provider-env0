package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGcpGkeCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGcpGkeCredentialsCreate,
		UpdateContext: resourceGcpGkeCredentialsUpdate,
		ReadContext:   resourceCredentialsRead(GCP_GKE_TYPE),
		DeleteContext: resourceCredentialsDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceCredentialsImport(GCP_GKE_TYPE)},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name for the credentials",
				Required:    true,
				ForceNew:    true,
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Description: "gke cluster name",
				Required:    true,
			},
			"compute_region": {
				Type:        schema.TypeString,
				Description: "the GCP gke compute region",
				Required:    true,
			},
		},
	}
}

func resourceGcpGkeCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	value := client.GcpGkeValue{}
	if err := readResourceData(&value, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	apiClient := meta.(client.ApiClientInterface)

	request := client.KubernetesCredentialsCreatePayload{
		Name:  d.Get("name").(string),
		Value: value,
		Type:  client.GcpGkeCredentialsType,
	}

	credentials, err := apiClient.KubernetesCredentialsCreate(&request)
	if err != nil {
		return diag.Errorf("could not create credentials: %v", err)
	}

	d.SetId(credentials.Id)

	return nil
}

func resourceGcpGkeCredentialsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	value := client.GcpGkeValue{}
	if err := readResourceData(&value, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	apiClient := meta.(client.ApiClientInterface)

	request := client.KubernetesCredentialsUpdatePayload{
		Value: value,
		Type:  client.GcpGkeCredentialsType,
	}

	if _, err := apiClient.KubernetesCredentialsUpdate(d.Id(), &request); err != nil {
		return diag.Errorf("could not create credentials: %v", err)
	}

	return nil
}
