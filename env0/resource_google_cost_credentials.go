package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGoogleCostCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGoogleCostCredentialsCreate,
		ReadContext:   resourceGoogleCostCredentialsRead,
		DeleteContext: resourceGoogleCostCredentialsDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name for the credentials",
				Required:    true,
				ForceNew:    true,
			},
			"table_id": {
				Type:        schema.TypeString,
				Description: "the table id of this credentials ",
				Required:    true,
				ForceNew:    true,
			},
			"secret": {
				Type:        schema.TypeString,
				Description: "the secret of this credentials",
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
			},
		},
	}
}

func resourceGoogleCostCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	value := client.GoogleCostCredentialsValeuPayload{
		TableId: d.Get("table_Id").(string),
		Secret:  d.Get("secret").(string),
	}

	request := client.GoogleCostCredentialsCreatePayload{
		Name:  d.Get("name").(string),
		Value: value,
		Type:  client.GoogleCostCredentiassType,
	}
	credentials, err := apiClient.GoogleCostCredentialsCreate(request)
	if err != nil {
		return diag.Errorf("could not create credentials: %v", err)
	}

	d.SetId(credentials.Id)

	return nil
}

func resourceGoogleCostCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	_, err := apiClient.CloudCredentials(id)
	if err != nil {
		return diag.Errorf("could not get credentials: %v", err)
	}
	return nil
}

func resourceGoogleCostCredentialsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	err := apiClient.CloudCredentialsDelete(id)
	if err != nil {
		return diag.Errorf("could not delete credentials: %v", err)
	}
	return nil
}
