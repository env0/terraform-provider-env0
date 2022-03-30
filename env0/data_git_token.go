package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataGitToken() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataGitTokenRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "the name of the git token",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "id of the git token",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
			},
		},
	}
}

func dataGitTokenRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var gitToken *client.GitToken
	var err error

	id, ok := d.GetOk("id")
	if ok {
		gitToken, err = meta.(client.ApiClientInterface).GitToken(id.(string))
	} else {
		gitToken, err = getGitTokenByName(d.Get("name").(string), meta)
	}

	if err != nil {
		return DataGetFailure("git token", id, err)
	}

	if err := writeResourceData(gitToken, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}
