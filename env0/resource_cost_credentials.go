package env0

import (
	"context"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const AWS = "aws"
const AZURE = "azure"
const GOOGLE = "google"

func resourceCostCredentials(providerName string) *schema.Resource {
	getSchema := func() map[string]*schema.Schema {
		switch providerName {
		case AWS:
			return map[string]*schema.Schema{
				"arn": {
					Type:        schema.TypeString,
					Description: "the aws role arn",
					Required:    true,
				},
				"duration": {
					Type:             schema.TypeInt,
					Description:      "the session duration in seconds. If set must be one of the following: 3600 (1h), 7200 (2h), 14400 (4h), 18000 (5h default), 28800 (8h), 43200 (12h)",
					Optional:         true,
					ValidateDiagFunc: NewIntInValidator([]int{3600, 7200, 14400, 18000, 28800, 43200}),
				},
			}
		case AZURE:
			return map[string]*schema.Schema{
				"client_id": {
					Type:        schema.TypeString,
					Description: "the azure client id",
					Required:    true,
				},
				"client_secret": {
					Type:        schema.TypeString,
					Description: "the azure client secret",
					Sensitive:   true,
					Required:    true,
				},
				"tenant_id": {
					Type:        schema.TypeString,
					Description: "the azure tenant id",
					Required:    true,
				},
				"subscription_id": {
					Type:        schema.TypeString,
					Description: "the azure subscription id",
					Required:    true,
				},
			}
		case GOOGLE:
			return map[string]*schema.Schema{
				"table_id": {
					Type:        schema.TypeString,
					Description: "the full BigQuery table id of the exported billing data",
					Required:    true,
				},
				"secret": {
					Type:        schema.TypeString,
					Description: "the GCP service account key",
					Sensitive:   true,
					Required:    true,
				},
			}
		default:
			panic("unhandled provider name: " + providerName)
		}
	}

	getPayload := func(d *schema.ResourceData) (client.CredentialCreatePayload, error) {
		var payload client.CredentialCreatePayload

		var value any

		switch providerName {
		case AWS:
			payload = &client.AwsCredentialsCreatePayload{
				Type: client.AwsCostCredentialsType,
			}
			value = &payload.(*client.AwsCredentialsCreatePayload).Value
		case AZURE:
			payload = &client.AzureCredentialsCreatePayload{
				Type: client.AzureCostCredentialsType,
			}
			value = &payload.(*client.AzureCredentialsCreatePayload).Value
		case GOOGLE:
			payload = &client.GoogleCostCredentialsCreatePayload{
				Type: client.GoogleCostCredentialsType,
			}
			value = &payload.(*client.GoogleCostCredentialsCreatePayload).Value
		default:
			panic("unhandled provider name: " + providerName)
		}

		if err := readResourceData(value, d); err != nil {
			return nil, fmt.Errorf("schema resource data deserialization failed: %w", err)
		}

		if err := readResourceData(payload, d); err != nil {
			return nil, fmt.Errorf("schema resource data deserialization failed: %w", err)
		}

		return payload, nil
	}

	getResourceCreate := func() func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
			payload, err := getPayload(d)
			if err != nil {
				return diag.FromErr(err)
			}

			apiClient := meta.(client.ApiClientInterface)

			res, err := apiClient.CredentialsCreate(payload)
			if err != nil {
				return diag.FromErr(err)
			}

			d.SetId(res.Id)

			return nil
		}
	}

	getResourceUpdate := func() func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
			payload, err := getPayload(d)
			if err != nil {
				return diag.FromErr(err)
			}

			apiClient := meta.(client.ApiClientInterface)

			if _, err := apiClient.CredentialsUpdate(d.Id(), payload); err != nil {
				return diag.FromErr(err)
			}

			return nil
		}
	}

	resourceSchema := getSchema()

	resourceSchema["name"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "the name for the credentials",
		Required:    true,
	}

	return &schema.Resource{
		CreateContext: getResourceCreate(),
		UpdateContext: getResourceUpdate(),
		ReadContext:   resourceCostCredentialsRead,
		DeleteContext: resourceCostCredentialsDelete,
		Schema:        resourceSchema,
	}
}

func resourceCostCredentialsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if _, err := apiClient.CloudCredentials(d.Id()); err != nil {
		return ResourceGetFailure(ctx, "cost credentials", d, err)
	}

	return nil
}

func resourceCostCredentialsDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.CloudCredentialsDelete(d.Id()); err != nil {
		return diag.Errorf("could not delete credentials: %v", err)
	}

	return nil
}
