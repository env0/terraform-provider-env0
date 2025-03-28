package env0

import (
	"context"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceModule() *schema.Resource {
	vcsTypeConflicts := func(fieldName string, allFields []string) []string {
		conflicts := make([]string, 0, len(allFields)-1)

		for _, field := range allFields {
			if field != fieldName {
				conflicts = append(conflicts, field)
			}
		}

		return conflicts
	}

	vcsTypes := []string{
		"is_azure_devops",
		"is_gitlab",
		"is_bitbucket_server",
		"is_github_enterprise",
		"is_gitlab_enterprise",
	}

	return &schema.Resource{
		CreateContext: resourceModuleCreate,
		ReadContext:   resourceModuleRead,
		UpdateContext: resourceModuleUpdate,
		DeleteContext: resourceModuleDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceModuleImport},

		Schema: map[string]*schema.Schema{
			"module_name": {
				Type:             schema.TypeString,
				Description:      "name of the module (Match pattern: ^[0-9A-Za-z](?:[0-9A-Za-z-_]{0,62}[0-9A-Za-z])?$)",
				Required:         true,
				ValidateDiagFunc: NewRegexValidator(`^[0-9A-Za-z](?:[0-9A-Za-z-_]{0,62}[0-9A-Za-z])?$`),
			},
			"module_provider": {
				Type:             schema.TypeString,
				Description:      "the provider name in the module source (Match pattern: ^[0-9a-z]{0,64}$)",
				Required:         true,
				ValidateDiagFunc: NewRegexValidator(`^[0-9a-z]{0,64}$`),
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
				Type:          schema.TypeString,
				Description:   "the git token id to be used",
				Optional:      true,
				ConflictsWith: []string{"github_installation_id", "bitbucket_client_key"},
			},
			"token_name": {
				Type:         schema.TypeString,
				Description:  "the token name used for integration with GitLab",
				Optional:     true,
				RequiredWith: []string{"token_id"},
			},
			"github_installation_id": {
				Type:          schema.TypeInt,
				Description:   "the env0 application installation id on the relevant Github repository",
				Optional:      true,
				ConflictsWith: []string{"token_id", "bitbucket_client_key"},
			},
			"bitbucket_client_key": {
				Type:          schema.TypeString,
				Description:   "the client key used for integration with Bitbucket",
				Optional:      true,
				ConflictsWith: []string{"token_id", "github_installation_id"},
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
			"path": {
				Type:        schema.TypeString,
				Description: "the folder in the repository to create the module from",
				Optional:    true,
			},
			"tag_prefix": {
				Type:        schema.TypeString,
				Description: "a tag prefix for the module",
				Optional:    true,
			},
			"module_test_enabled": {
				Type:        schema.TypeBool,
				Description: "set to 'true' to enable module test (defaults to 'false')",
				Optional:    true,
				Default:     false,
			},
			"run_tests_on_pull_request": {
				Type:         schema.TypeBool,
				Description:  "set to 'true' to run tests on pull request (defaults to 'false'). Can only be enabled if 'module_test_enabled' is enabled",
				Optional:     true,
				Default:      false,
				RequiredWith: []string{"module_test_enabled"},
			},
			"opentofu_version": {
				Type:             schema.TypeString,
				Description:      "the opentofu version to use, Can only be set if 'module_test_enabled' is enabled",
				Optional:         true,
				Default:          "",
				RequiredWith:     []string{"module_test_enabled"},
				ValidateDiagFunc: NewOpenTofuVersionValidator(),
			},
			"is_azure_devops": {
				Type:          schema.TypeBool,
				Description:   "true if this module integrates with azure dev ops",
				Optional:      true,
				Default:       false,
				ConflictsWith: vcsTypeConflicts("is_azure_devops", vcsTypes),
			},
			"is_gitlab": {
				Type:          schema.TypeBool,
				Description:   "true if this module integrates with GitLab",
				Optional:      true,
				Default:       false,
				ConflictsWith: vcsTypeConflicts("is_gitlab", vcsTypes),
			},
			"is_bitbucket_server": {
				Type:          schema.TypeBool,
				Description:   "true if this module integrates with Bitbucket Server",
				Optional:      true,
				Default:       false,
				ConflictsWith: vcsTypeConflicts("is_bitbucket_server", vcsTypes),
			},
			"is_github_enterprise": {
				Type:          schema.TypeBool,
				Description:   "true if this module integrates with GitHub Enterprise",
				Optional:      true,
				Default:       false,
				ConflictsWith: vcsTypeConflicts("is_github_enterprise", vcsTypes),
			},
			"is_gitlab_enterprise": {
				Type:          schema.TypeBool,
				Description:   "true if this module integrates with GitLab Enterprise",
				Optional:      true,
				Default:       false,
				ConflictsWith: vcsTypeConflicts("is_gitlab_enterprise", vcsTypes),
			},
		},
	}
}

func resourceModuleCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.ModuleCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if !payload.ModuleTestEnabled && (payload.RunTestsOnPullRequest || payload.OpentofuVersion != "") {
		return diag.Errorf("'run_tests_on_pull_request' and/or 'opentofu_version' may only be set if 'module_test_enabled' is enabled (set to 'true')")
	}

	module, err := apiClient.ModuleCreate(payload)
	if err != nil {
		return diag.Errorf("could not create module: %v", err)
	}

	d.SetId(module.Id)

	return nil
}

func resourceModuleRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	module, err := apiClient.Module(d.Id())
	if err != nil {
		return ResourceGetFailure(ctx, "module", d, err)
	}

	if module.IsDeleted {
		tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]any{"id": d.Id()})
		d.SetId("")

		return nil
	}

	if err := writeResourceData(module, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceModuleUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.ModuleUpdatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if !payload.ModuleTestEnabled && (payload.RunTestsOnPullRequest || payload.OpentofuVersion != "") {
		return diag.Errorf("'run_tests_on_pull_request' and/or 'opentofu_version' may only be set if 'module_test_enabled' is enabled (set to 'true')")
	}

	if _, err := apiClient.ModuleUpdate(d.Id(), payload); err != nil {
		return diag.Errorf("could not update module: %v", err)
	}

	return nil
}

func resourceModuleDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.ModuleDelete(d.Id()); err != nil {
		return diag.Errorf("could not delete module: %v", err)
	}

	return nil
}

func getModuleByName(name string, meta any) (*client.Module, error) {
	apiClient := meta.(client.ApiClientInterface)

	modules, err := apiClient.Modules()
	if err != nil {
		return nil, err
	}

	var foundModules []client.Module

	for _, module := range modules {
		if !module.IsDeleted && module.ModuleName == name {
			foundModules = append(foundModules, module)
		}
	}

	if len(foundModules) == 0 {
		return nil, fmt.Errorf("module with name %v not found", name)
	}

	if len(foundModules) > 1 {
		return nil, fmt.Errorf("found multiple modules with name: %s. Use id instead or make sure module names are unique %v", name, foundModules)
	}

	return &foundModules[0], nil
}

func getModule(ctx context.Context, id string, meta any) (*client.Module, error) {
	_, err := uuid.Parse(id)
	if err == nil {
		tflog.Info(ctx, "Resolving module by id", map[string]any{"id": id})

		return meta.(client.ApiClientInterface).Module(id)
	} else {
		tflog.Info(ctx, "Resolving module by name", map[string]any{"name": id})

		return getModuleByName(id, meta)
	}
}

func resourceModuleImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	module, err := getModule(ctx, d.Id(), meta)
	if err != nil {
		return nil, err
	}

	if err := writeResourceData(module, d); err != nil {
		return nil, fmt.Errorf("schema resource data serialization failed: %w", err)
	}

	return []*schema.ResourceData{d}, nil
}
