package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTeamOrganizationAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeamOrganizationAssignmentCreate,
		UpdateContext: resourceTeamOrganizationAssignmentUpdate,
		ReadContext:   resourceTeamOrganizationAssignmentRead,
		DeleteContext: resourceTeamOrganizationAssignmentDelete,

		Description: "assigns an organization role to a team",

		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:        schema.TypeString,
				Description: "id of the team",
				Required:    true,
				ForceNew:    true,
			},
			"role_id": {
				Type:             schema.TypeString,
				Description:      "id of the assigned custom role. The following built-in roles can be passed as well: `User`, `Admin`",
				Required:         true,
				ValidateDiagFunc: ValidateNotEmptyString,
			},
		},
	}
}

func resourceTeamOrganizationAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api_client := meta.(client.ApiClientInterface)

	var newAssignment client.AssignOrganizationRoleToTeamPayload
	if err := readResourceData(&newAssignment, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	assignment, err := api_client.AssignOrganizationRoleToTeam(&newAssignment)
	if err != nil {
		return diag.Errorf("could not create assignment: %v", err)
	}

	d.SetId(assignment.Id)

	return nil
}

func resourceTeamOrganizationAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api_client := meta.(client.ApiClientInterface)

	assignments, err := api_client.OrganizationRoleTeamAssignments()
	if err != nil {
		return diag.Errorf("could not get assignments: %v", err)
	}

	id := d.Id()

	for _, assignment := range assignments {
		if assignment.Id == id {
			if err := writeResourceData(&assignment, d); err != nil {
				return diag.Errorf("schema resource data serialization failed: %v", err)
			}

			d.Set("role_id", assignment.Role)

			return nil
		}
	}

	tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]interface{}{"id": d.Id()})
	d.SetId("")

	return nil
}

func resourceTeamOrganizationAssignmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api_client := meta.(client.ApiClientInterface)

	var payload client.AssignOrganizationRoleToTeamPayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	assignment, err := api_client.AssignOrganizationRoleToTeam(&payload)
	if err != nil {
		return diag.Errorf("could not update assignment: %v", err)
	}

	d.SetId(assignment.Id)

	return nil
}

func resourceTeamOrganizationAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api_client := meta.(client.ApiClientInterface)

	teamId := d.Get("team_id").(string)

	if err := api_client.RemoveOrganizationRoleFromTeam(teamId); err != nil {
		return diag.Errorf("could not delete assignment: %v", err)
	}

	return nil
}
