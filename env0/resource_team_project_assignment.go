package env0

import (
	"context"
	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTeamProjectAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeamProjectAssignmentCreateOrUpdate,
		ReadContext:   resourceSshKeyRead,
		UpdateContext: resourceTeamProjectAssignmentCreateOrUpdate,
		DeleteContext: resourceTeamProjectAssignmentDelete,

		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:        schema.TypeString,
				Description: "id of the team",
				Required:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "id of the project",
				Required:    true,
			},
			"role": {
				Type:        schema.TypeString,
				Description: "the assigned role",
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					_ = client.Role(val.(string))
					return
				},
			},
		},
	}
}

func resourceTeamProjectAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	//TODO: complete this
	return nil
}

func resourceTeamProjectAssignmentCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	request := client.TeamProjectAssignmentPayload{
		TeamId:      d.Get("team_id").(string),
		ProjectId:   d.Get("project_id").(string),
		ProjectRole: d.Get("role").(client.Role),
	}
	_, err := apiClient.TeamProjectAssignmentCreateOrUpdate(request)
	if err != nil {
		return diag.Errorf("could not Create or Update TeamProjectAssignment: %v", err)
	}
	return nil
}

func resourceTeamProjectAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	err := apiClient.TeamProjectAssignmentDelete(d.Id())
	if err != nil {
		return diag.Errorf("could not delete TeamProjectAssignment: %v", err)
	}

	return nil
}
