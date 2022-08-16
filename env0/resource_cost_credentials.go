package env0

import (
	"context"
	"errors"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const AWS = "aws"
const AZURE = "azure"
const GOOGLE = "google"

func resourceCostCredentials(providerName string) *schema.Resource {

	awsSchema := map[string]*schema.Schema{
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
			Description: "the azure client secret",
			Sensitive:   true,
			ForceNew:    true,
			Required:    true,
		},
		"tenant_id": {
			Type:        schema.TypeString,
			Description: "the azure tenant id",
			ForceNew:    true,
			Required:    true,
		},
		"subscription_id": {
			Type:        schema.TypeString,
			Description: "the azure subscription id",
			ForceNew:    true,
			Required:    true,
		},
	}

	googleSchema := map[string]*schema.Schema{
		"table_id": {
			Type:        schema.TypeString,
			Description: "the full BigQuery table id of the exported billing data",
			ForceNew:    true,
			Required:    true,
		},
		"secret": {
			Type:        schema.TypeString,
			Description: "the GCP service account key",
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
		Schema:        extendSchema(schemaMap[providerName]),
	}
}

func extendSchema(schemaToReadFrom map[string]*schema.Schema) map[string]*schema.Schema {

	resultsSchema := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "the name for the credentials",
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

	apiKey, err := sendApiCallToCreateCred(d, meta)
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

	if err := apiClient.CloudCredentialsDelete(d.Id()); err != nil {
		return diag.Errorf("could not delete credentials: %v", err)
	}

	return nil
}

func sendApiCallToCreateCred(d *schema.ResourceData, meta interface{}) (client.Credentials, error) {
	apiClient := meta.(client.ApiClientInterface)
	_, awsOk := d.GetOk("arn")
	_, azureOk := d.GetOk("client_id")
	_, googleOk := d.GetOk("table_id")
	switch {
	case awsOk:
		var value client.AwsCredentialsValuePayload
		if err := readResourceData(&value, d); err != nil {
			return client.Credentials{}, fmt.Errorf("schema resource data deserialization failed: %v", err)
		}

		return apiClient.CredentialsCreate(&client.AwsCredentialsCreatePayload{
			Name:  d.Get("name").(string),
			Type:  client.AwsCostCredentialsType,
			Value: value,
		})
	case azureOk:
		var value client.AzureCredentialsValuePayload
		if err := readResourceData(&value, d); err != nil {
			return client.Credentials{}, fmt.Errorf("schema resource data deserialization failed: %v", err)
		}

		return apiClient.CredentialsCreate(&client.AzureCredentialsCreatePayload{
			Name:  d.Get("name").(string),
			Type:  client.AzureCostCredentialsType,
			Value: value,
		})
	case googleOk:
		var value client.GoogleCostCredentialsValuePayload
		if err := readResourceData(&value, d); err != nil {
			return client.Credentials{}, fmt.Errorf("schema resource data deserialization failed: %v", err)
		}

		return apiClient.CredentialsCreate(&client.GoogleCostCredentialsCreatePayload{
			Name:  d.Get("name").(string),
			Type:  client.GoogleCostCredentialsType,
			Value: value,
		})
	default:
		return client.Credentials{}, errors.New("error in schema, no required value defined")
	}
}
