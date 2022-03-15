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
		ReadContext:   resourceAzureCredentialsRead,
		DeleteContext: resourceAzureCredentialsDelete,

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
				Description: "the azure client secret",
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
		},
	}
}

func resourceAzureCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	value := client.AzureCredentialsValuePayload{
		ClientId:       d.Get("client_id").(string),
		ClientSecret:   d.Get("client_secret").(string),
		SubscriptionId: d.Get("subscription_id").(string),
		TenantId:       d.Get("tenant_id").(string),
	}
	requestType := client.AzureServicePrincipalCredentialsType

	request := client.AzureCredentialsCreatePayload{
		Name:  d.Get("name").(string),
		Value: value,
		Type:  requestType,
	}
	credentials, err := apiClient.AzureCredentialsCreate(request)
	if err != nil {
		return diag.Errorf("could not create credentials key: %v", err)
	}

	d.SetId(credentials.Id)

	return nil
}

func resourceAzureCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	_, err := apiClient.CloudCredentials(id)
	if err != nil {
		return diag.Errorf("could not get credentials: %v", err)
	}
	return nil
}

func resourceAzureCredentialsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	err := apiClient.CloudCredentialsDelete(id)
	if err != nil {
		return diag.Errorf("could not delete credentials: %v", err)
	}
	return nil
}
