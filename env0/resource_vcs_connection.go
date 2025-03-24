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

func resourceVcsConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVcsConnectionCreate,
		ReadContext:   resourceVcsConnectionRead,
		UpdateContext: resourceVcsConnectionUpdate,
		DeleteContext: resourceVcsConnectionDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceVcsConnectionImport},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the VCS type (BitBucketServer, GitLabEnterprise, or GitHubEnterprise)",
				ValidateDiagFunc: NewStringInValidator([]string{
					"BitBucketServer",
					"GitLabEnterprise",
					"GitHubEnterprise",
				}),
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of the VCS connection",
			},
			"url": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "URL of the VCS server. This can either be a 'VCS URL' (e.g.: https://github.com) or 'Repository URL' (E.g.: https://github.com/env0/myrepo)",
				ValidateDiagFunc: ValidateUrl,
			},
			"vcs_agent_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VCS agent key. Use a custom agent key or 'ENV0_DEFAULT' to use the default env0 agent",
			},
		},
	}
}

func resourceVcsConnectionCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.VcsConnectionCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	vcsConnection, err := apiClient.VcsConnectionCreate(payload)
	if err != nil {
		return diag.Errorf("could not create VCS connection: %v", err)
	}

	d.SetId(vcsConnection.Id)

	return nil
}

func resourceVcsConnectionRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	vcsConnection, err := apiClient.VcsConnection(d.Id())
	if err != nil {
		return diag.Errorf("could not get VCS connection: %v", err)
	}

	if err := writeResourceData(vcsConnection, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceVcsConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.VcsConnectionUpdatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if _, err := apiClient.VcsConnectionUpdate(d.Id(), payload); err != nil {
		return diag.Errorf("could not update VCS connection: %v", err)
	}

	return resourceVcsConnectionRead(ctx, d, meta)
}

func resourceVcsConnectionDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.VcsConnectionDelete(d.Id()); err != nil {
		return diag.Errorf("could not delete VCS connection: %v", err)
	}

	return nil
}

func getVcsConnectionByName(name string, meta any) (*client.VcsConnection, error) {
	apiClient := meta.(client.ApiClientInterface)

	vcsConnections, err := apiClient.VcsConnections()
	if err != nil {
		return nil, err
	}

	var foundConnections []client.VcsConnection

	for _, connection := range vcsConnections {
		if connection.Name == name {
			foundConnections = append(foundConnections, connection)
		}
	}

	if len(foundConnections) == 0 {
		return nil, fmt.Errorf("VCS connection with name %v not found", name)
	}

	if len(foundConnections) > 1 {
		return nil, fmt.Errorf("found multiple VCS connections with name: %s. Use id instead or make sure VCS connection names are unique %v", name, foundConnections)
	}

	return &foundConnections[0], nil
}

func getVcsConnection(ctx context.Context, id string, meta any) (*client.VcsConnection, error) {
	if _, err := uuid.Parse(id); err == nil {
		tflog.Info(ctx, "Resolving VCS connection by id", map[string]any{"id": id})

		return meta.(client.ApiClientInterface).VcsConnection(id)
	}

	tflog.Info(ctx, "Resolving VCS connection by name", map[string]any{"name": id})

	return getVcsConnectionByName(id, meta)
}

func resourceVcsConnectionImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	vcsConnection, err := getVcsConnection(ctx, d.Id(), meta)
	if err != nil {
		return nil, err
	}

	if err := writeResourceData(vcsConnection, d); err != nil {
		return nil, fmt.Errorf("schema resource data serialization failed: %w", err)
	}

	return []*schema.ResourceData{d}, nil
}
