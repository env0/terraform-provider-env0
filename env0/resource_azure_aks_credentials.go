package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAzureAksCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureAksCredentialsCreate,
		UpdateContext: resourceAzureAksCredentialsUpdate,
		ReadContext:   resourceCredentialsRead(AZURE_AKS_TYPE),
		DeleteContext: resourceCredentialsDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceCredentialsImport(AZURE_AKS_TYPE)},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name for the credentials",
				Required:    true,
				ForceNew:    true,
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Description: "aks cluster name",
				Required:    true,
			},
			"resource_group": {
				Type:        schema.TypeString,
				Description: "the resource group of the aks",
				Required:    true,
			},
		},
	}
}

func resourceAzureAksCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	value := client.AzureAksValue{}
	if err := readResourceData(&value, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	apiClient := meta.(client.ApiClientInterface)

	request := client.KubernetesCredentialsCreatePayload{
		Name:  d.Get("name").(string),
		Value: value,
		Type:  client.AzureAksCredentialsType,
	}

	credentials, err := apiClient.KubernetesCredentialsCreate(&request)
	if err != nil {
		return diag.Errorf("could not create credentials: %v", err)
	}

	d.SetId(credentials.Id)

	return nil
}

func resourceAzureAksCredentialsUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	value := client.AzureAksValue{}
	if err := readResourceData(&value, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	apiClient := meta.(client.ApiClientInterface)

	request := client.KubernetesCredentialsUpdatePayload{
		Value: value,
		Type:  client.AzureAksCredentialsType,
	}

	if _, err := apiClient.KubernetesCredentialsUpdate(d.Id(), &request); err != nil {
		return diag.Errorf("could not create credentials: %v", err)
	}

	return nil
}
