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
		readGitFields(&payload.GitConfig, d)
	}

	environmentImport, err := apiClient.EnvironmentImportCreate(&payload)
	if err != nil {
		return diag.Errorf("could not create environment import: %v", err)
	}

	// how are the other fields set?
	d.SetId(environmentImport.Id)

	return nil
}

func readGitFields(gitConfig *client.GitConfig, d *schema.ResourceData) {
	gitConfig.Path = d.Get("path").(string)
	gitConfig.Revision = d.Get("revision").(string)
	gitConfig.Repository = d.Get("repository").(string)
	gitConfig.Provider = d.Get("git_provider").(string)
}

func resourceEnvironmentImportRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	environmentImport, err := apiClient.EnvironmentImportGet(d.Id())
	if err != nil {
		return ResourceGetFailure(ctx, "environment import", d, err)
	}

	if err := writeResourceData(environmentImport, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	d.Set("path", environmentImport.GitConfig.Path)
	d.Set("revision", environmentImport.GitConfig.Revision)
	d.Set("repository", environmentImport.GitConfig.Repository)
	d.Set("git_provider", environmentImport.GitConfig.Provider)

	return nil
}

func resourceEnvironmentImportUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	var payload client.EnvironmentImportUpdatePayload

	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	readGitFields(&payload.GitConfig, d)

	if _, err := apiClient.EnvironmentImportUpdate(id, &payload); err != nil {
		return diag.Errorf("could not update environment import: %v", err)
	}

	return nil
}

// should not actually delete the environment import
func resourceEnvironmentImportDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
