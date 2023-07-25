package env0

import (
	"context"
	"fmt"
	"log"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCustomFlow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomFlowCreate,
		ReadContext:   resourceCustomFlowRead,
		UpdateContext: resourceCustomFlowUpdate,
		DeleteContext: resourceCustomFlowDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceCustomFlowImport},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "id of the custom flow",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "name for the custom flow. note: for the UI to render the custom-flow please use `project-<project.id>`",
				Required:    true,
			},
			"repository": {
				Type:        schema.TypeString,
				Description: "repository url for the custom flow source code",
				Required:    true,
			},
			"path": {
				Type:        schema.TypeString,
				Description: "terraform / terragrunt file folder inside source code. Should be the full path including the .yaml/.yml file",
				Optional:    true,
			},
			"revision": {
				Type:        schema.TypeString,
				Description: "source code revision (branch / tag) to use",
				Optional:    true,
			},
			"token_id": {
				Type:        schema.TypeString,
				Description: "the git token id to be used",
				Optional:    true,
			},
			"ssh_keys": {
				Type:        schema.TypeList,
				Description: "an array of references to 'data_ssh_key' to use when accessing git over ssh",
				Optional:    true,
				Elem: &schema.Schema{
					Type:        schema.TypeMap,
					Description: "a map of env0_ssh_key.id and env0_ssh_key.name for each project",
				},
			},
			"gitlab_project_id": {
				Type:         schema.TypeInt,
				Description:  "the project id of the relevant repository",
				Optional:     true,
				RequiredWith: []string{"token_id"},
			},
			"github_installation_id": {
				Type:        schema.TypeInt,
				Description: "the env0 application installation id on the relevant github repository",
				Optional:    true,
			},
			"bitbucket_client_key": {
				Type:        schema.TypeString,
				Description: "the bitbucket client key used for integration",
				Optional:    true,
			},
			"is_bitbucket_server": {
				Type:        schema.TypeBool,
				Description: "true if this custom flow uses bitbucket server repository",
				Optional:    true,
				Default:     false,
			},
			"is_gitlab_enterprise": {
				Type:        schema.TypeBool,
				Description: "true if this custom flow uses gitlab enterprise repository",
				Optional:    true,
				Default:     false,
			},
			"is_github_enterprise": {
				Type:        schema.TypeBool,
				Description: "true if this custom flow uses github enterprise repository",
				Optional:    true,
				Default:     false,
			},
			"is_gitlab": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "true if this custom flow integrates with gitlab repository",
				Default:     false,
			},
			"is_azure_devops": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "true if this custom flow integrates with azure dev ops repository",
				Default:     false,
			},
		},
	}
}

func resourceCustomFlowCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.CustomFlowCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	customFlow, err := apiClient.CustomFlowCreate(payload)
	if err != nil {
		return diag.Errorf("could not create custom flow: %v", err)
	}

	d.SetId(customFlow.Id)

	return nil
}

func resourceCustomFlowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	customFlow, err := apiClient.CustomFlow(d.Id())
	if err != nil {
		return ResourceGetFailure("custom flow", d, err)
	}

	if err := writeResourceData(customFlow, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceCustomFlowUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.CustomFlowCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if _, err := apiClient.CustomFlowUpdate(d.Id(), payload); err != nil {
		return diag.Errorf("could not update custom flow: %v", err)
	}

	return nil
}

func resourceCustomFlowDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.CustomFlowDelete(d.Id()); err != nil {
		return diag.Errorf("could not delete custom flow: %v", err)
	}

	return nil
}

func getCustomFlowByName(name string, meta interface{}) (*client.CustomFlow, error) {
	apiClient := meta.(client.ApiClientInterface)

	customFlows, err := apiClient.CustomFlows(name)
	if err != nil {
		return nil, err
	}

	if len(customFlows) == 0 {
		return nil, fmt.Errorf("custom flow with name %v not found", name)
	}

	if len(customFlows) > 1 {
		return nil, fmt.Errorf("found multiple custom flows with name: %s. Use id instead or make sure custom flow names are unique %v", name, customFlows)
	}

	return &customFlows[0], nil
}

func getCustomFlow(id string, meta interface{}) (*client.CustomFlow, error) {
	if _, err := uuid.Parse(id); err == nil {
		log.Println("[INFO] Resolving custom flow by id: ", id)
		return meta.(client.ApiClientInterface).CustomFlow(id)
	} else {
		log.Println("[INFO] Resolving custom flow by name: ", id)
		return getCustomFlowByName(id, meta)
	}
}

func resourceCustomFlowImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	customFlow, err := getCustomFlow(d.Id(), meta)
	if err != nil {
		return nil, err
	}

	if err := writeResourceData(customFlow, d); err != nil {
		return nil, fmt.Errorf("schema resource data serialization failed: %v", err)
	}

	return []*schema.ResourceData{d}, nil
}
