package env0

import (
	"context"
	"errors"
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
		CustomizeDiff: resourceVcsConnectionCustomizeDiff,

		Importer: &schema.ResourceImporter{StateContext: resourceVcsConnectionImport},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the VCS type (BitBucketServer, GitLabEnterprise, GitHubEnterprise, or AzureDevOps)",
				ValidateDiagFunc: NewStringInValidator([]string{
					"BitBucketServer",
					"GitLabEnterprise",
					"GitHubEnterprise",
					"AzureDevOps",
				}),
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "name of the VCS connection. Connections created in the env0 UI have no name, so omit it when importing one. An AzureDevOps connection cannot be renamed",
			},
			"url": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "URL of the VCS server. This can either be a 'VCS URL' (e.g.: https://github.com) or 'Repository URL' (E.g.: https://github.com/env0/myrepo). Required for BitBucketServer, GitLabEnterprise and GitHubEnterprise",
				ValidateDiagFunc: ValidateUrl,
			},
			"vcs_agent_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VCS agent key. Use a custom agent key or 'ENV0_DEFAULT' to use the default env0 agent",
			},
			"token_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "the id of an Azure DevOps OAuth token, created by authorizing Azure DevOps in the env0 UI. Required for AzureDevOps. Destroying the connection also deletes this token",
				// Older connections don't return their token id, so an imported one may have none in state.
				DiffSuppressFunc: func(_, oldValue, _ string, d *schema.ResourceData) bool {
					return oldValue == "" && d.Id() != ""
				},
			},
			"oauth_provider": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the OAuth method of an AzureDevOps connection's token ('legacy' or 'entra'). Empty means legacy",
			},
		},
	}
}

func resourceVcsConnectionCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ any) error {
	if d.Get("type").(string) == "AzureDevOps" {
		if d.NewValueKnown("token_id") && d.Get("token_id").(string) == "" {
			return errors.New("token_id is required for AzureDevOps")
		}

		// The backend ignores AzureDevOps updates, and replacing the connection would delete its token.
		if d.HasChange("name") && d.Id() != "" {
			return errors.New("an AzureDevOps VCS connection cannot be renamed")
		}

		return nil
	}

	if d.NewValueKnown("url") && d.Get("url").(string) == "" {
		return fmt.Errorf("url is required for %s", d.Get("type"))
	}

	return nil
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

	if err := d.Set("oauth_provider", vcsConnection.OauthProvider); err != nil {
		return diag.Errorf("failed to set oauth_provider: %v", err)
	}

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
