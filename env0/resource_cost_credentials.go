package env0

import (
	"context"
	"errors"

	"github.com/env0/terraform-provider-env0/client"
	. "github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const AWS = "aws"
const AZURE = "azure"
const GOOGLE = "google"

func resourceCostCredentials(providerName string) *schema.Resource {

	awsSchema := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "name for the credentials",
			Required:    true,
			ForceNew:    true,
		},
		"arn": {
			Type:        schema.TypeString,
			Description: "the aws role arn",
			ForceNew:    true,
			Required:    true,
		},
		"external_id": {
			Type:        schema.TypeString,
			Description: "the aws role external id",
			Sensitive:   true,
			ForceNew:    true,
			Required:    true,
		},
	}

	azureSchema := map[string]*schema.Schema{
		"client_id": {
			Type:        schema.TypeString,
			Description: "the azure client id",
			ForceNew:    true,
			Required:    true,
		},
		"client_secret": {
			Type:        schema.TypeString,
			Description: "azure client secret",
			Sensitive:   true,
			ForceNew:    true,
			Required:    true,
		},
		"tenant_id": {
			Type:        schema.TypeString,
			Description: "azure tenant id",
			ForceNew:    true,
			Required:    true,
		},
		"subscription_id": {
			Type:        schema.TypeString,
			Description: "azure subscription id",
			ForceNew:    true,
			Required:    true,
		},
	}

	googleSchema := map[string]*schema.Schema{
		"table_id": {
			Type:        schema.TypeString,
			Description: "the table id of this credentials ",
			ForceNew:    true,
			Required:    true,
		},
		"secret": {
			Type:        schema.TypeString,
			Description: "the secret of this credentials",
			Sensitive:   true,
			ForceNew:    true,
			Required:    true,
		},
	}

	schemaMap := map[string]map[string]*schema.Schema{AWS: awsSchema, AZURE: azureSchema, GOOGLE: googleSchema}

	return &schema.Resource{
		CreateContext: resourceCostCredentialsCreate,
		ReadContext:   resourceCostCredentialsRead,
		DeleteContext: resourceCostCredentialsDelete,
		Schema:        expendSchema(schemaMap[providerName]),
	}
}

func expendSchema(schemaToReadFrom map[string]*schema.Schema) map[string]*schema.Schema {

	resultsSchema := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "name for the credentials",
			Required:    true,
			ForceNew:    true,
		},
	}

	for index, element := range schemaToReadFrom {
		resultsSchema[index] = element
	}

	return resultsSchema
}

func resourceCostCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	apiClient := meta.(ApiClientInterface)
	var apiKey client.Credentials
	var err error
	payLoad, credType, err := getCredentailsPayload(d)
	if err != nil {
		return diag.Errorf("ERROR: %v", err)
	}

	switch credType {
	case string(client.AwsCostCredentialsType):
		payLoadToSend := payLoad.(client.AwsCredentialsCreatePayload)
		apiKey, err = apiClient.AwsCredentialsCreate(payLoadToSend)
	case string(client.AzureCostCredentialsType):
		payLoadToSend := payLoad.(client.AzureCredentialsCreatePayload)
		apiKey, err = apiClient.AzureCredentialsCreate(payLoadToSend)
	case string(client.GoogleCostCredentialsType):
		payLoadToSend := payLoad.(client.GoogleCostCredentialsCreatePayload)
		apiKey, err = apiClient.GoogleCostCredentialsCreate(payLoadToSend)
	}

	if err != nil {
		return diag.Errorf("Cost credential failed: %v", err)
	}
	d.SetId(apiKey.Id)
	return nil
}

func resourceCostCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	_, err := apiClient.CloudCredentials(id)
	if err != nil {
		return ResourceGetFailure("cost credentials", d, err)
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
		return string(client.GoogleCostCredentialsType), nil
	default:
		return "error", errors.New("error in schema, no required value defined")
	}

}

func getCredentailsPayload(d *schema.ResourceData) (interface{}, string, error) {
	credType, err := getCredType(d)
	if err != nil {
		return nil, "", err
	}
	switch credType {
	case string(client.AwsCostCredentialsType):

		return client.AwsCredentialsCreatePayload{
			Name: d.Get("name").(string),
			Type: client.AwsCostCredentialsType,
			Value: client.AwsCredentialsValuePayload{
				RoleArn:    d.Get("arn").(string),
				ExternalId: d.Get("external_id").(string),
			},
		}, credType, nil
	case string(client.AzureCostCredentialsType):
		return client.AzureCredentialsCreatePayload{
			Name: d.Get("name").(string),
			Type: client.AzureCostCredentialsType,
			Value: client.AzureCredentialsValuePayload{
				ClientId:       d.Get("client_id").(string),
				ClientSecret:   d.Get("client_secret").(string),
				TenantId:       d.Get("tenant_id").(string),
				SubscriptionId: d.Get("subscription_id").(string),
			},
		}, credType, nil
	case string(client.GoogleCostCredentialsType):
		return client.GoogleCostCredentialsCreatePayload{
			Name: d.Get("name").(string),
			Type: client.GoogleCostCredentialsType,
			Value: client.GoogleCostCredentialsValeuPayload{
				TableId: d.Get("table_id").(string),
				Secret:  d.Get("secret").(string),
			},
		}, credType, nil

	default:
		return "", "", errors.New("Failed to build credentials payload, please contact env0")
	}

}

