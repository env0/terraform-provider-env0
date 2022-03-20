package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceModule() *schema.Resource {
	vcsExcatlyOneOf := []string{
		"token_id",
		"github_installation_id",
		"bitbucket_client_key",
	}

	return &schema.Resource{
		CreateContext: resourceModuleCreate,
		ReadContext:   resourceModuleRead,
		UpdateContext: resourceModuleUpdate,
		DeleteContext: resourceModuleDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceModuleImport},

		Schema: map[string]*schema.Schema{
			"module_name": {
				Type:        schema.TypeString,
				Description: "name of the module",
				Required:    true,
			},
			"module_provider": {
				Type:        schema.TypeString,
				Description: "the provider name in the module source",
				Optional:    true,
			},
			"repository": {
				Type:        schema.TypeString,
				Description: "the repository containing the module files",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "description of the module",
				Optional:    true,
			},
			"token_id": {
				Type:         schema.TypeString,
				Description:  "the token id used for integration with GitLab",
				Optional:     true,
				ExactlyOneOf: vcsExcatlyOneOf,
			},
			"token_name": {
				Type:         schema.TypeString,
				Description:  "the token name used for integration with GitLab",
				Optional:     true,
				RequiredWith: []string{"token_id"},
			},
			"github_installation_id": {
				Type:         schema.TypeInt,
				Description:  "The env0 application installation id on the relevant Github repository",
				Optional:     true,
				ExactlyOneOf: vcsExcatlyOneOf,
			},
			"bitbucket_client_key": {
				Type:         schema.TypeString,
				Description:  "The client key used for integration with Bitbucket",
				Optional:     true,
				ExactlyOneOf: vcsExcatlyOneOf,
			},
			"ssh_keys": {
				Type:        schema.TypeList,
				Description: "an array of references to 'data_ssh_key' to use when accessing git over ssh",
				Optional:    true,
				Elem: &schema.Schema{
					Type:        schema.TypeMap,
					Description: "a map of env0_ssh_key.id and env0_ssh_key.name for each module",
				},
			},
		},
	}
}

func resourceModuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.ModuleCreatePayload
	if err := deserializeResourceData(&payload, d); err != nil {
		diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if len(payload.TokenId) > 0 {
		payload.IsGitlab = boolPtr(true)
	}

	module, err := apiClient.ModuleCreate(payload)
	if err != nil {
		return diag.Errorf("could not create module: %v", err)
	}

	d.SetId(module.Id)

	return nil
}

func resourceModuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	module, err := apiClient.Module(d.Id())
	if err != nil {
		return diag.Errorf("could not get module: %v", err)
	}

	if err := serializeResourceData(module, d); err != nil {
		diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceModuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.ModuleUpdatePayload
	if err := deserializeResourceData(&payload, d); err != nil {
		diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if _, err := apiClient.ModuleUpdate(d.Id(), payload); err != nil {
		return diag.Errorf("could not update module: %v", err)
	}

	return nil
}

func resourceModuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.ModuleDelete(d.Id()); err != nil {
		return diag.Errorf("could not delete module: %v", err)
	}

	return nil
}

func resourceModuleImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return nil, nil
}
