package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKubeconfigCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKubeconfigCredentialsCreate,
		UpdateContext: resourceKubeconfigCredentialsUpdate,
		ReadContext:   resourceCredentialsRead(KUBECONFIG_TYPE),
		DeleteContext: resourceCredentialsDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceCredentialsImport(KUBECONFIG_TYPE)},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name for the credentials",
				Required:    true,
				ForceNew:    true,
			},
			"kube_config": {
				Type:        schema.TypeString,
				Description: "A valid kubeconfig file content",
				Required:    true,
			},
		},
	}
}

func resourceKubeconfigCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	value := client.KubeconfigFileValue{}
	if err := readResourceData(&value, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	apiClient := meta.(client.ApiClientInterface)

	request := client.KubernetesCredentialsCreatePayload{
		Name:  d.Get("name").(string),
		Value: value,
		Type:  client.KubeconfigCredentialsType,
	}

	credentials, err := apiClient.KubernetesCredentialsCreate(&request)
	if err != nil {
		return diag.Errorf("could not create credentials: %v", err)
	}

	d.SetId(credentials.Id)

	return nil
}

func resourceKubeconfigCredentialsUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	value := client.KubeconfigFileValue{}
	if err := readResourceData(&value, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	apiClient := meta.(client.ApiClientInterface)

	request := client.KubernetesCredentialsUpdatePayload{
		Value: value,
		Type:  client.KubeconfigCredentialsType,
	}

	if _, err := apiClient.KubernetesCredentialsUpdate(d.Id(), &request); err != nil {
		return diag.Errorf("could not create credentials: %v", err)
	}

	return nil
}
