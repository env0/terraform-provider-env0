package env0

import (
	"context"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataOidcCredentials(credentialsType CloudType) *schema.Resource {
	return &schema.Resource{
		ReadContext: dataOidcCredentialRead(credentialsType),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  fmt.Sprintf("the name of the %s oidc credentials", credentialsType),
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Description:  fmt.Sprintf("the id of the %s oidc credentials", credentialsType),
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
		},
	}
}

func dataOidcCredentialRead(credentialsType CloudType) func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		var credentials client.Credentials

		var err error

		id, ok := d.GetOk("id")
		if ok {
			credentials, err = getCredentialsById(id.(string), credentialsTypeToPrefixList[credentialsType], meta)
		} else {
			credentials, err = getCredentialsByName(d.Get("name").(string), credentialsTypeToPrefixList[credentialsType], meta)
		}

		if err != nil {
			return DataGetFailure(fmt.Sprintf("%s oidc credentials", credentialsType), id, err)
		}

		if err := writeResourceData(&credentials, d); err != nil {
			return diag.Errorf("schema resource data serialization failed: %v", err)
		}

		return nil
	}
}
