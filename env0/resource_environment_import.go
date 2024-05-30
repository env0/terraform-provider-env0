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
		DeleteContext: resourceEnvironmentImportDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "name to give the environment",
				Optional:    true,
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
			"git_provider": {
				Type:             schema.TypeString,
				Description:      "vcs provider of the environment ( one of \"github\" | \"gitlab\" | \"bitbucket\" | \"azure\" | \"other\" )",
				Optional:         true,
				ValidateDiagFunc: NewStringInValidator([]string{"github", "gitlab", "bitbucket", "azure", "other"}),
			},
			"workspace": {
				Type:        schema.TypeString,
				Description: "workspace of the environment",
				Optional:    true,
			},
			"iac_type": {
				Type:             schema.TypeString,
				Description:      "iac type of the environment ( one of \"opentofu\" | \"terraform\" )",
				Optional:         true,
				ValidateDiagFunc: NewStringInValidator([]string{"opentofu", "terraform"}),
			},
			"iac_version": {
				Type:        schema.TypeString,
				Description: "iac version of the environment",
				Optional:    true,
			},
			"soft_delete": {
				Type:        schema.TypeBool,
				Description: "soft delete the configuration variable, once removed from the configuration it won't be deleted from env0",
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

	if d.Get("git_provider").(string) != "" {
		if err := readResourceData(&payload.GitConfig, d); err != nil {
			return diag.Errorf("schema resource data deserialization failed: %v", err)
		}
	}

	environmentImport, err := apiClient.EnvironmentImportCreate(&payload)
	if err != nil {
		return diag.Errorf("could not create environment import: %v", err)
	}

	d.SetId(environmentImport.Id)

	return nil
}

func resourceEnvironmentImportRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	environmentImport, err := apiClient.EnvironmentImportGet(d.Id())
	if err != nil {
		return ResourceGetFailure(ctx, "environment_import", d, err)
	}

	if err := writeResourceData(environmentImport, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if err := writeResourceData(&environmentImport.GitConfig, d); err != nil {
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

	if err := readResourceData(&payload.GitConfig, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if _, err := apiClient.EnvironmentImportUpdate(id, &payload); err != nil {
		return diag.Errorf("could not update environment import: %v", err)
	}

	return nil
}

func resourceEnvironmentImportDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// don't delete if soft delete is set
	if softDelete, ok := d.GetOk("soft_delete"); ok && softDelete.(bool) {
		return nil
	}

	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.EnvironmentImportDelete(d.Id()); err != nil {
		return diag.Errorf("could not delete environment import: %v", err)
	}

	return nil
}
