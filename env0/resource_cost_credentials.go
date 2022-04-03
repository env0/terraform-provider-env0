package env0

import (
	"context"
	"errors"

	"github.com/env0/terraform-provider-env0/client"
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCostCredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCostCredentialsCreate,
		ReadContext:   resourceCostCredentialsRead,
		DeleteContext: resourceCostCredentialsDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name for the credentials",
				Required:    true,
				ForceNew:    true,
			},
			"arn": {
				Type:          schema.TypeString,
				Description:   "the aws role arn",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"client_id", "client_secret", "tenant_id", "subscription_id", "table_id", "secret"},
				ExactlyOneOf:  []string{"client_id", "table_id"},
			},
			"external_id": {
				Type:          schema.TypeString,
				Description:   "the aws role external id",
				Optional:      true,
				Sensitive:     true,
				ForceNew:      true,
				RequiredWith:  []string{"arn"},
				ConflictsWith: []string{"client_id", "client_secret", "tenant_id", "subscription_id", "table_id", "secret"},
			},
			"client_id": {
				Type:          schema.TypeString,
				Description:   "the azure client id",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"external_id", "arn", "table_id", "secret"},
				ExactlyOneOf:  []string{"arn", "table_id"},
			},
			"client_secret": {
				Type:          schema.TypeString,
				Description:   "azure client secret",
				Optional:      true,
				Sensitive:     true,
				ForceNew:      true,
				RequiredWith:  []string{"client_id"},
				ConflictsWith: []string{"external_id", "arn", "table_id", "secret"},
			},
			"tenant_id": {
				Type:          schema.TypeString,
				Description:   "azure tenant id",
				Optional:      true,
				ForceNew:      true,
				RequiredWith:  []string{"client_id"},
				ConflictsWith: []string{"external_id", "arn", "table_id", "secret"},
			},
			"subscription_id": {
				Type:          schema.TypeString,
				Description:   "azure subscription id",
				Optional:      true,
				ForceNew:      true,
				RequiredWith:  []string{"client_id"},
				ConflictsWith: []string{"external_id", "arn", "table_id", "secret"},
			},
			"table_id": {
				Type:          schema.TypeString,
				Description:   "the table id of this credentials ",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"client_id", "client_secret", "tenant_id", "subscription_id", "arn", "external_id"},
				ExactlyOneOf:  []string{"arn", "client_id"},
			},
			"secret": {
				Type:          schema.TypeString,
				Description:   "the secret of this credentials",
				Optional:      true,
				Sensitive:     true,
				ForceNew:      true,
				RequiredWith:  []string{"table_id"},
				ConflictsWith: []string{"client_id", "client_secret", "tenant_id", "subscription_id", "arn", "external_id"},
			},
		},
	}
}

func resourceCostCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	apiClient := meta.(ApiClientInterface)
	var apikey client.Credentials
	var err error
	payLoad, credType, err := setPayload(d)
	if err != nil {
		return diag.Errorf("ERROR: %v", err)
	}

	switch credType {
	case string(client.AwsCostCredentialsType):
		payLoadToSent := payLoad.(client.AwsCredentialsCreatePayload)
		payLoadToSent.Name = d.Get("name").(string)
		apikey, err = apiClient.AwsCredentialsCreate(payLoadToSent)
	case string(client.AzureCostCredentialsType):
		payLoadToSent := payLoad.(client.AzureCredentialsCreatePayload)
		payLoadToSent.Name = d.Get("name").(string)
		apikey, err = apiClient.AzureCredentialsCreate(payLoadToSent)
	case string(client.GoogleCostCredentiassType):
		payLoadToSent := payLoad.(client.GoogleCostCredentialsCreatePayload)
		payLoadToSent.Name = d.Get("name").(string)
		apikey, err = apiClient.GoogleCostCredentialsCreate(payLoadToSent)
	}

	if err != nil {
		return diag.Errorf("cost credential fialed: %v", err)
	}
	d.SetId(apikey.Id)
	return nil
}

func resourceCostCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	_, err := apiClient.CloudCredentials(id)
	if err != nil {
		return diag.Errorf("could not get credentials: %v", err)
	}
	return nil

}

func resourceCostCredentialsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	err := apiClient.CloudCredentialsDelete(id)
	if err != nil {
		return diag.Errorf("could not delete credentials: %v", err)
	}
	return nil

}

func getCredType(d *schema.ResourceData) (string, error) {
	_, awsOk := d.GetOk("arn")
	_, azureOk := d.GetOk("client_id")
	_, googleOk := d.GetOk("table_id")
	switch {
	case awsOk:
		return string(client.AwsCostCredentialsType), nil
	case azureOk:
		return string(client.AzureCostCredentialsType), nil
	case googleOk:
		return string(client.GoogleCostCredentiassType), nil
	default:
		return "error", errors.New("error in schema, no required value  defined")
	}

}

func setPayload(d *schema.ResourceData) (interface{}, string, error) {
	credType, err := getCredType(d)
	if err != nil {
		return nil, "", err
	}
	switch credType {
	case string(client.AwsCostCredentialsType):

		return client.AwsCredentialsCreatePayload{
			Type: client.AwsCostCredentialsType,
			Value: client.AwsCredentialsValuePayload{
				RoleArn:    d.Get("arn").(string),
				ExternalId: d.Get("external_id").(string),
			},
		}, credType, nil
	case string(client.AzureCostCredentialsType):
		return client.AzureCredentialsCreatePayload{
			Type: client.AzureCostCredentialsType,
			Value: client.AzureCredentialsValuePayload{
				ClientId:       d.Get("client_id").(string),
				ClientSecret:   d.Get("client_secret").(string),
				TenantId:       d.Get("tenant_id").(string),
				SubscriptionId: d.Get("subscription_id").(string),
			},
		}, credType, nil
	case string(client.GoogleCostCredentiassType):
		return client.GoogleCostCredentialsCreatePayload{
			Type: client.GoogleCostCredentiassType,
			Value: client.GoogleCostCredentialsValeuPayload{
				TableId: d.Get("table_id").(string),
				Secret:  d.Get("secret").(string),
			},
		}, credType, nil

	default:
		return "", "", errors.New("cant create payload")
	}

}
