package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataGithubInstallationId() *schema.Resource {
	return &schema.Resource{
		Description: "returns the github installation id of a git hub repositroy",

		ReadContext: dataGithubInstallationIdRead,

		Schema: map[string]*schema.Schema{
			"repository": {
				Type:        schema.TypeString,
				Description: "the name of the repository",
				Required:    true,
			},
			"github_installation_id": {
				Type:        schema.TypeInt,
				Description: "the github installation id",
				Computed:    true,
			},
		},
	}
}

func dataGithubInstallationIdRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	repositroy := d.Get("repository").(string)

	token, err := apiClient.VcsToken("github", repositroy)
	if err != nil {
		return diag.Errorf("failed to get github installation id: %v", err)
	}

	d.Set("github_installation_id", token.Token)
	d.SetId(repositroy)

	return nil
}
