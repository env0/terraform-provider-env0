package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEnvironmentImport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentImportCreate,
		ReadContext:   resourceEnvironmentImportRead,
		UpdateContext: resourceEnvironmentImportUpdate,
		DeleteContext: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "id of the environment import",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "name to give the environment",
				Optional:    true,
			},
			"variables": {
				Type:        schema.TypeList,
				Description: "key value pairs to set as environment variables",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "the variable name",
							Required:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "the variable value",
						},
						"isSensitive": {
							Type:        schema.TypeBool,
							Description: "is the variable sensitive",
						},
						"type": {
							Type:        schema.TypeString,
							Description: "the variable type \"string\" | \"JSON\"",
							Optional:    true,
							Default:     "string",
						},
					},
				},
			},
			"path": {
				Type:        schema.TypeString,
				Description: "path to the tofu configuration",
				Optional:    true,
			},
			"revision": {
				Type:        schema.TypeString,
				Description: "revision of the environment",
				Optional:    true,
			},
			"repository": {
				Type:        schema.TypeString,
				Description: "repository of the environment",
				Optional:    true,
			},
			"provider": {
				Type:        schema.TypeString,
				Description: "vcs provider of the environment ( one of \"github\" | \"gitlab\" | \"bitbucket\" )",
				Optional:    true,
			},
			"workspace": {
				Type:        schema.TypeString,
				Description: "workspace of the environment",
				Optional:    true,
			},
			"iacType": {
				Type:        schema.TypeString,
				Description: "iac type of the environment ( one of \"opentofu\" | \"terraform\" )",
				Optional:    true,
			},
			"iacVersion": {
				Type:        schema.TypeString,
				Description: "iac version of the environment",
				Optional:    true,
			},
		},
	}
}

func resourceEnvironmentImportCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.EnvironmentImportCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	environmentImport, err := apiClient.EnvironmentImportCreate(&payload)
	if err != nil {
		return diag.Errorf("could not create project: %v", err)
	}

	// how are the other fields set?
	d.SetId(environmentImport.Id)

	return nil
}

func resourceEnvironmentImportRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	environmentImport, err := apiClient.EnvironmentImportGet(d.Id())
	if err != nil {
		return ResourceGetFailure(ctx, "environment import", d, err)
	}

	if err := writeResourceData(&environmentImport, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	return nil
}

func resourceEnvironmentImportUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	var payload client.EnvironmentImportUpdatePayload

	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if _, err := apiClient.EnvironmentImportUpdate(id, &payload); err != nil {
		return diag.Errorf("could not update environment import: %v", err)
	}

	return nil
}

// should not actually delete the environment import
func resourceEnvironmentImportDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
