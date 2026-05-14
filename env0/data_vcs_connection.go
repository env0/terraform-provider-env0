package env0

import (
	"context"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataVcsConnection() *schema.Resource {
	return &schema.Resource{
		Description: "Returns a VCS connection by access scope and connection type. The access_scope value can be found in the env0 UI under Organization Settings > VCS.",
		ReadContext: dataVcsConnectionRead,

		Schema: map[string]*schema.Schema{
			"access_scope": {
				Type:        schema.TypeString,
				Description: "the access scope of the VCS connection. The format depends on the VCS type: 'Organization:<org-name>' or 'User:<username>' for GitHub, 'Workspace:<workspace-slug>' for Bitbucket, '<username>' or '<token-name>' for GitLab, '<display-name>' for Azure DevOps, or 'url:<base-url>' for self-hosted VCS",
				Required:    true,
			},
			"connection_type": {
				Type:             schema.TypeString,
				Description:      "the connection type. Valid values: 'DeploymentPipeline' or 'CodeWrite'",
				Required:         true,
				ValidateDiagFunc: NewStringInValidator([]string{"DeploymentPipeline", "CodeWrite"}),
			},
			"id": {
				Type:        schema.TypeString,
				Description: "the id of the VCS connection",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "the name of the VCS connection",
				Computed:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "the VCS type (e.g. GitHub, GitLab, BitBucket, AzureDevOps, GitHubEnterprise, GitLabEnterprise, BitBucketServer)",
				Computed:    true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "the URL of the VCS server (for self-hosted VCS connections)",
				Computed:    true,
			},
			"vcs_agent_key": {
				Type:        schema.TypeString,
				Description: "the VCS agent key (only applicable to self-hosted VCS connections)",
				Computed:    true,
			},
		},
	}
}

func dataVcsConnectionRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	accessScope := d.Get("access_scope").(string)
	connectionType := d.Get("connection_type").(string)

	vcsConnection, err := getVcsConnectionByAccessScope(accessScope, connectionType, apiClient)
	if err != nil {
		return diag.Errorf("could not read VCS connection: %v", err)
	}

	if err := writeResourceData(vcsConnection, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func getVcsConnectionByAccessScope(accessScope string, connectionType string, apiClient client.ApiClientInterface) (*client.VcsConnection, error) {
	vcsConnections, err := apiClient.VcsConnections()
	if err != nil {
		return nil, err
	}

	var foundConnections []client.VcsConnection

	for _, connection := range vcsConnections {
		if connection.AccessScope != accessScope {
			continue
		}

		if connection.ConnectionType != connectionType {
			continue
		}

		foundConnections = append(foundConnections, connection)
	}

	if len(foundConnections) == 0 {
		return nil, fmt.Errorf("VCS connection with access scope '%s' and connection type '%s' not found", accessScope, connectionType)
	}

	if len(foundConnections) > 1 {
		return nil, fmt.Errorf("found multiple VCS connections with access scope '%s' and connection type '%s'", accessScope, connectionType)
	}

	return &foundConnections[0], nil
}
