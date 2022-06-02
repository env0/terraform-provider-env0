package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var credentialsTypeToPrefixList map[string][]string = map[string][]string{
	"gcp":        {string(client.GcpServiceAccountCredentialsType)},
	"azure":      {string(client.AzureServicePrincipalCredentialsType)},
	"aws":        {string(client.AwsAssumedRoleCredentialsType), string(client.AwsAccessKeysCredentialsType)},
	"gcp_cost":   {string(client.GoogleCostCredentialsType)},
	"azure_cost": {string(client.AzureCostCredentialsType)},
	"aws_cost":   {string(client.AwsCostCredentialsType)},
}

func dataCredentials(cloudType string) *schema.Resource {
	return &schema.Resource{
		ReadContext: dataCredentialsRead(cloudType),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "the name of the credentials",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "the id of the credentials",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
		},
	}
}

func dataCredentialsRead(cloudType string) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		var err error
		var credentials client.Credentials

		prefixList := credentialsTypeToPrefixList[cloudType]

		id, ok := d.GetOk("id")
		if ok {
			credentials, err = getCredentialsById(id.(string), prefixList, meta)
			if err != nil {
				return diag.Errorf("could not query %s credentials by id: %v", cloudType, err)
			}
		} else {
			name, ok := d.GetOk("name")
			if !ok {
				return diag.Errorf("either 'name' or 'id' must be specified")
			}
			credentials, err = getCredentialsByName(name.(string), prefixList, meta)
			if err != nil {
				return diag.Errorf("could not query %s credentials by name: %v", cloudType, err)
			}
		}

		if err := writeResourceData(&credentials, d); err != nil {
			return diag.Errorf("schema resource data serialization failed: %v", err)
		}

		return nil
	}
}
