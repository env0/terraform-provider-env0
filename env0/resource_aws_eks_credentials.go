package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAwsEksCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAwsEksCredentialsCreate,
		ReadContext:   resourceCredentialsRead(AWS_EKS_TYPE),
		DeleteContext: resourceCredentialsDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceCredentialsImport(AWS_EKS_TYPE)},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name for the credentials",
				Required:    true,
				ForceNew:    true,
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Description: "eks cluster name",
				Required:    true,
				ForceNew:    true,
			},
			"cluster_region": {
				Type:        schema.TypeString,
				Description: "the AWS region of the eks cluster",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceAwsEksCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	value := client.AwsEksValue{}
	if err := readResourceData(&value, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	apiClient := meta.(client.ApiClientInterface)

	request := client.KubernetesCredentialsCreatePayload{
		Name:  d.Get("name").(string),
		Value: value,
		Type:  client.AwsEksCredentialsType,
	}

	credentials, err := apiClient.KubernetesCredentialsCreate(&request)
	if err != nil {
		return diag.Errorf("could not create credentials: %v", err)
	}

	d.SetId(credentials.Id)

	return nil
}
