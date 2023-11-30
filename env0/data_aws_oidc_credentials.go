package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataAwsOidcCredentials() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataAwsOidcCredentialRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "the name of the aws oidc credentials",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "the id of the aws oidc credentials",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"oidc_sub": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the jwt oidc sub",
			},
		},
	}
}

func dataAwsOidcCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var credentials client.Credentials
	var err error

	id, ok := d.GetOk("id")
	if ok {
		credentials, err = getCredentialsById(id.(string), credentialsTypeToPrefixList[AWS_OIDC_TYPE], meta)
	} else {
		credentials, err = getCredentialsByName(d.Get("name").(string), credentialsTypeToPrefixList[AWS_OIDC_TYPE], meta)
	}

	if err != nil {
		return DataGetFailure("aws oidc credentials", id, err)
	}

	if err := writeResourceData(&credentials, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	apiClient := meta.(client.ApiClientInterface)

	oidcSub, err := apiClient.OidcSub()
	if err != nil {
		return diag.Errorf("failed to get oidc sub: %v", err)
	}

	d.Set("oidc_sub", oidcSub)

	return nil
}
