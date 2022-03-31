package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataModule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataModuleRead,

		Schema: map[string]*schema.Schema{
			"module_name": {
				Type:         schema.TypeString,
				Description:  "the name of the module",
				Optional:     true,
				ExactlyOneOf: []string{"module_name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "id of the module",
				Optional:     true,
				ExactlyOneOf: []string{"module_name", "id"},
			},
			"module_provider": {
				Type:        schema.TypeString,
				Description: "the provider name in the module source",
				Computed:    true,
			},
			"repository": {
				Type:        schema.TypeString,
				Description: "template source code repository url",
				Computed:    true,
			},
			"token_id": {
				Type:        schema.TypeString,
				Description: "the token id used for integration with GitLab",
				Optional:    true,
			},
			"token_name": {
				Type:        schema.TypeString,
				Description: "the token name used for integration with GitLab",
				Optional:    true,
			},
			"github_installation_id": {
				Type:        schema.TypeInt,
				Description: "the env0 application installation id on the relevant github repository",
				Optional:    true,
			},
			"bitbucket_client_key": {
				Type:        schema.TypeString,
				Description: "the client key used for integration with Bitbucket",
				Optional:    true,
			},
		},
	}
}

func dataModuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var module *client.Module
	var err error

	id, ok := d.GetOk("id")
	if ok {
		module, err = meta.(client.ApiClientInterface).Module(id.(string))
	} else {
		name := d.Get("module_name").(string)
		module, err = getModuleByName(name, meta)
	}

	if err != nil {
		return DataGetFailure("module", id, err)
	}

	if err := writeResourceData(module, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}
