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
		CreateContext: resourceTeamOrganizationAssignmentCreateOrUpdate,
		UpdateContext: resourceTeamOrganizationAssignmentCreateOrUpdate,
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
				ValidateDiagFunc: NewRoleValidator([]string{"User", "Admin"}),
			},
		},
	}
}

func resourceTeamOrganizationAssignmentCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	organizationId, err := apiClient.OrganizationId()
	if err != nil {
		return diag.Errorf("could not get organization id: %v", err)

	}

	var payload client.TeamRoleAssignmentCreateOrUpdatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}
	payload.OrganizationId = organizationId

	assignment, err := apiClient.TeamRoleAssignmentCreateOrUpdate(&payload)
	if err != nil {
		return diag.Errorf("could not create assignment: %v", err)
	}

	d.SetId(assignment.Id)

	return nil
}

func resourceTeamOrganizationAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	organizationId, err := apiClient.OrganizationId()
	if err != nil {
		return diag.Errorf("could not get organization id: %v", err)

	}

	var payload client.TeamRoleAssignmentListPayload
	payload.OrganizationId = organizationId

	id := d.Id()

	assignments, err := apiClient.TeamRoleAssignments(&payload)
	if err != nil {
		return diag.Errorf("could not get assignments: %v", err)
	}

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

func resourceTeamOrganizationAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.TeamRoleAssignmentDeletePayload

	organizationId, err := apiClient.OrganizationId()
	if err != nil {
		return diag.Errorf("could not get organization id: %v", err)

	}

	payload.OrganizationId = organizationId
	payload.TeamId = d.Get("team_id").(string)

	if err := apiClient.TeamRoleAssignmentDelete(&payload); err != nil {
		return diag.Errorf("could not delete assignment: %v", err)
	}

	return nil
}
