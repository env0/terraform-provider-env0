package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataAwsCredentials() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataAwsCredentialsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "the name of the credential",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "the id of the credential",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
		},
	}
}

func dataAwsCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var err diag.Diagnostics
	var credentials client.ApiKey

	id, ok := d.GetOk("id")
	if ok {
		credentials, err = getAwsCredentialsById(id.(string), meta)
		if err != nil {
			return err
		}
	} else {
		name, ok := d.GetOk("name")
		if !ok {
			return diag.Errorf("Either 'name' or 'id' must be specified")
		}
		credentials, err = getAwsCredentialsByName(name.(string), meta)
		if err != nil {
			return err
		}
	}

	d.SetId(credentials.Id)
	d.Set("name", credentials.Name)

	return nil
}

func getAwsCredentialsByName(name interface{}, meta interface{}) (client.ApiKey, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)
	credentialsList, err := apiClient.AwsCredentialsList()
	if err != nil {
		return client.ApiKey{}, diag.Errorf("Could not query AWS Credentials by name: %v", err)
	}

	credentialsByNameAndType := make([]client.ApiKey, 0)
	for _, candidate := range credentialsList {
		if candidate.Name == name.(string) && candidate.Type == "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT" {
			credentialsByNameAndType = append(credentialsByNameAndType, candidate)
		}
	}

	if len(credentialsByNameAndType) > 1 {
		return client.ApiKey{}, diag.Errorf("Found multiple AWS Credentials for name: %s", name)
	}
	if len(credentialsByNameAndType) == 0 {
		return client.ApiKey{}, diag.Errorf("Could not find AWS Credentials with name: %s", name)
	}
	return credentialsByNameAndType[0], nil
}

func getAwsCredentialsById(id string, meta interface{}) (client.ApiKey, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)
	credentials, err := apiClient.AwsCredentials(id)
	if credentials.Type != "aws_assumed_role" {
		return client.ApiKey{}, diag.Errorf("Found credentials which are not AWS Credentials: %v", credentials)
	}
	if err != nil {
		return client.ApiKey{}, diag.Errorf("Could not query AWS Credentials: %v", err)
	}
	return credentials, nil
}
