package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataAwsCostCredentials() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataAwsCostCredentialsRead,

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

func dataAwsCostCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var err diag.Diagnostics
	var credentials client.Credentials

	id, ok := d.GetOk("id")
	if ok {
		credentials, err = getAwsCostCredentialsById(id.(string), meta)
		if err != nil {
			return err
		}
	} else {
		name, ok := d.GetOk("name")
		if !ok {
			return diag.Errorf("Either 'name' or 'id' must be specified")
		}
		credentials, err = getAwsCostCredentialsByName(name.(string), meta)
		if err != nil {
			return err
		}
	}

	d.SetId(credentials.Id)
	d.Set("name", credentials.Name)

	return nil
}

func getAwsCostCredentialsByName(name interface{}, meta interface{}) (client.Credentials, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)
	credentialsList, err := apiClient.CloudCredentialsList()
	if err != nil {
		return client.Credentials{}, diag.Errorf("Could not query AWS Cost Credentials by name: %v", err)
	}

	credentialsByNameAndType := make([]client.Credentials, 0)
	for _, candidate := range credentialsList {
		if candidate.Name == name.(string) && candidate.Type == string(client.AwsCostCredentialsType) {
			credentialsByNameAndType = append(credentialsByNameAndType, candidate)
		}
	}

	if len(credentialsByNameAndType) > 1 {
		return client.Credentials{}, diag.Errorf("Found multiple AWS Cost Credentials for name: %s", name)
	}
	if len(credentialsByNameAndType) == 0 {
		return client.Credentials{}, diag.Errorf("Could not find AWS Cost Credentials with name: %s", name)
	}
	return credentialsByNameAndType[0], nil
}

func getAwsCostCredentialsById(id string, meta interface{}) (client.Credentials, diag.Diagnostics) {
	apiClient := meta.(client.ApiClientInterface)
	credentials, err := apiClient.CloudCredentials(id)
	if credentials.Type != string(client.AwsCostCredentialsType) {
		return client.Credentials{}, diag.Errorf("Found credentials which are not AWS Cost Credentials: %v", credentials)
Write
	}
	if err != nil {
		return client.Credentials{}, diag.Errorf("Could not query AWS Cost Credentials: %v", err)
	}
	return credentials, nil
}
