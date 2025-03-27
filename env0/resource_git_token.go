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

func resourceGitToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGitTokenCreate,
		ReadContext:   resourceGitTokenRead,
		DeleteContext: resourceGitTokenDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceGitTokenImport},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "the git token name",
				Required:    true,
				ForceNew:    true,
			},
			"value": {
				Type:        schema.TypeString,
				Description: "the git token value",
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceGitTokenCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.GitTokenCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	gitToken, err := apiClient.GitTokenCreate(payload)
	if err != nil {
		return diag.Errorf("could not create git token: %v", err)
	}

	d.SetId(gitToken.Id)

	return nil
}

func resourceGitTokenRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	gitToken, err := apiClient.GitToken(d.Id())
	if err != nil {
		return ResourceGetFailure(ctx, "git token", d, err)
	}

	if err := writeResourceData(gitToken, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceGitTokenDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.GitTokenDelete(d.Id()); err != nil {
		return diag.Errorf("could not delete git token: %v", err)
	}

	return nil
}

func getGitTokenByName(name string, meta any) (*client.GitToken, error) {
	apiClient := meta.(client.ApiClientInterface)

	gitTokens, err := apiClient.GitTokens()
	if err != nil {
		return nil, err
	}

	var foundGitTokens []client.GitToken

	for _, gitToken := range gitTokens {
		if gitToken.Name == name {
			foundGitTokens = append(foundGitTokens, gitToken)
		}
	}

	if len(foundGitTokens) == 0 {
		return nil, fmt.Errorf("git token with name %v not found", name)
	}

	if len(foundGitTokens) > 1 {
		return nil, fmt.Errorf("found multiple git tokens with name: %s. Use id instead or make sure git token names are unique %v", name, foundGitTokens)
	}

	return &foundGitTokens[0], nil
}

func getGitToken(ctx context.Context, id string, meta any) (*client.GitToken, error) {
	_, err := uuid.Parse(id)
	if err == nil {
		tflog.Info(ctx, "Resolving git token by id", map[string]any{"id": id})

		return meta.(client.ApiClientInterface).GitToken(id)
	} else {
		tflog.Info(ctx, "Resolving git token by name", map[string]any{"name": id})

		return getGitTokenByName(id, meta)
	}
}

func resourceGitTokenImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	gitToken, err := getGitToken(ctx, d.Id(), meta)
	if err != nil {
		return nil, err
	}

	if err := writeResourceData(gitToken, d); err != nil {
		return nil, fmt.Errorf("schema resource data serialization failed: %w", err)
	}

	return []*schema.ResourceData{d}, nil
}
