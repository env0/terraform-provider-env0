package env0

import (
	"context"
	"fmt"
	"strings"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var allowedCredentialTypes = []string{
	"AWS_ASSUMED_ROLE",
	"AWS_ASSUMED_ROLE_FOR_DEPLOYMENT",
	"AWS_ACCESS_KEYS_FOR_DEPLOYMENT",
	"GCP_CREDENTIALS",
	"GCP_SERVICE_ACCOUNT_FOR_DEPLOYMENT",
	"AZURE_CREDENTIALS",
	"AZURE_SERVICE_PRINCIPAL_FOR_DEPLOYMENT",
}

func dataCloudCredentials() *schema.Resource {
	allowedCredentialTypesStr := fmt.Sprintf("(allowed values: %s)", strings.Join(allowedCredentialTypes, ", "))

	return &schema.Resource{
		ReadContext: dataCloudCredentialsRead,

		Schema: map[string]*schema.Schema{
			"names": {
				Type:        schema.TypeList,
				Description: "list of all cloud credentials (by name), optionally filtered by credential_type",
				Computed:    true,
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "the credential name",
				},
			},
			"credential_type": {
				Type:             schema.TypeString,
				Description:      "the type of cloud credential to filter by " + allowedCredentialTypesStr,
				Optional:         true,
				ValidateDiagFunc: NewStringInValidator(allowedCredentialTypes),
			},
		},
	}
}

func dataCloudCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	credentialsList, err := apiClient.CloudCredentialsList()
	if err != nil {
		return diag.Errorf("Could not get cloud credentials list: %v", err)
	}

	credential_type, filter := d.GetOk("credential_type")

	data := []string{}

	for _, credentials := range credentialsList {
		if filter && credential_type != credentials.Type {
			continue
		}
		data = append(data, credentials.Name)
	}

	d.Set("names", data)

	// Not really needed. But required by Terraform SDK - https://github.com/hashicorp/terraform-plugin-sdk/issues/541
	d.SetId("all_cloud_credential_names")

	return nil
}
