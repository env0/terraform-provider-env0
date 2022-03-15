package env0

import (
	"context"
	"errors"
	"log"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAzureCostCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureCostCredentialsCreate,
		ReadContext:   resourceAzureCostCredentialsRead,
		DeleteContext: resourceAzureCostCredentialsDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceAwsCredentialsImport},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name for the credentials",
				Required:    true,
				ForceNew:    true,
			},
			"client_Id": {
				Type:        schema.TypeString,
				Description: "the azure client id",
				Required:    true,
			},
			"client_Secret": {
				Type:        schema.TypeString,
				Description: "azure client secret",
				Required:    true,
				Sensitive:   true,
			},
			"tenant_Id": {
				Type:        schema.TypeString,
				Description: "azure tenant id",
				Required:    true,
			},
			"subscription_Id": {
				Type:        schema.TypeString,
				Description: "azure subscription id",
				Required:    true,
			},
		},
	}
}

func resourceAzureCostCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	value := client.AzureCredentialsValuePayload{
		ClientId:       d.Get("client_id").(string),
		ClientSecret:   d.Get("client_Secret").(string),
		TenantId:       d.Get("tenant_id").(string),
		SubscriptionId: d.Get("subscription_id").(string),
	}
	request := client.AzureCredentialsCreatePayload{
		Name:  d.Get("name").(string),
		Type:  client.AzureCostCredentialsType,
		Value: value,
	}

	apiKey, err := apiClient.AzureCredentialsCreate(request)

	if err != nil {
		return diag.Errorf("could not create azure credentials: %v", err)
	}

	d.SetId(apiKey.Id)
	return nil
}

func resourceAzureCostCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	_, err := apiClient.CloudCredentials(id)
	if err != nil {
		return diag.Errorf("could not get credentials: %v", err)
	}
	return nil

}

func resourceAzureCostCredentialsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	err := apiClient.CloudCredentialsDelete(id)
	if err != nil {
		return diag.Errorf("could not delete credentials: %v", err)
	}
	return nil
}

func resourceAzureCostCredentialsImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()
	var getErr diag.Diagnostics
	_, uuidErr := uuid.Parse(id)
	if uuidErr == nil {
		log.Println("[INFO] Resolving AWS Credentials by id: ", id)
		_, getErr = getAwsCostCredentialsById(id, meta)
	} else {
		log.Println("[DEBUG] ID is not a valid env0 id ", id)
		log.Println("[INFO] Resolving AWS Credentials by name: ", id)
		var awsCredential client.ApiKey
		awsCredential, getErr = getAwsCostCredentialsByName(id, meta)
		d.SetId(awsCredential.Id)
	}
	if getErr != nil {
		return nil, errors.New(getErr[0].Summary)
	} else {
		return []*schema.ResourceData{d}, nil
	}
}
