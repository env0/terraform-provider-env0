package env0

import (
	"context"
	"log"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserProjectAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserProjectAssignmentCreate,
		UpdateContext: resourceUserProjectAssignmentUpdate,
		ReadContext:   resourceUserProjectAssignmentRead,
		DeleteContext: resourceUserProjectAssignmentDelete,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Description: `id of the user. Note: can also be an id of a "User" API key`,
				Required:    true,
				ForceNew:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "id of the project",
				Required:    true,
				ForceNew:    true,
			},
			"role": {
				Type:             schema.TypeString,
				Description:      "the assigned role (Admin, Planner, Viewer, Deployer)",
				Optional:         true,
				ValidateDiagFunc: ValidateRole,
				ExactlyOneOf:     []string{"custom_role_id", "role"},
			},
			"custom_role_id": {
				Type:         schema.TypeString,
				Description:  "id of the assigned custom role",
				Optional:     true,
				ExactlyOneOf: []string{"custom_role_id", "role"},
			},
		},
	}
}

func resourceUserProjectAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var newAssignment client.AssignUserToProjectPayload
	if err := readResourceData(&newAssignment, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	role, ok := d.GetOk("role")
	if !ok {
		role = d.Get("custom_role_id")
	}
	newAssignment.Role = role.(string)

	projectId := d.Get("project_id").(string)

	apiClient := meta.(client.ApiClientInterface)

	assignment, err := apiClient.AssignUserToProject(projectId, &newAssignment)
	if err != nil {
		return diag.Errorf("could not create assignment: %v", err)
	}

	d.SetId(assignment.Id)

	return nil
}

func resourceUserProjectAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	id := d.Id()
	projectId := d.Get("project_id").(string)

	assignments, err := apiClient.UserProjectAssignments(projectId)
	if err != nil {
		return diag.Errorf("could not get UserProjectAssignments: %v", err)
	}

	for _, assignment := range assignments {
		if assignment.Id == id {
			if err := writeResourceData(&assignment, d); err != nil {
				return diag.Errorf("schema resource data serialization failed: %v", err)
			}

			if client.IsBuiltinProjectRole(assignment.Role) {
				d.Set("role", assignment.Role)
			} else {
				d.Set("custom_role_id", assignment.Role)
			}

			return nil
		}
	}

	log.Printf("[WARN] Drift Detected: Terraform will remove %s from state", d.Id())
	d.SetId("")

	return nil
}

func resourceUserProjectAssignmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectId := d.Get("project_id").(string)
	userId := d.Get("user_id").(string)

	var payload client.UpdateUserProjectAssignmentPayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	role, ok := d.GetOk("role")
	if !ok {
		role = d.Get("custom_role_id")
	}
	payload.Role = role.(string)

	apiClient := meta.(client.ApiClientInterface)
	if _, err := apiClient.UpdateUserProjectAssignment(projectId, userId, &payload); err != nil {
		return diag.Errorf("could not update role for UserProjectAssignment: %v", err)
	}

	return nil
}

func resourceUserProjectAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectId := d.Get("project_id").(string)
	userId := d.Get("user_id").(string)

	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.RemoveUserFromProject(projectId, userId); err != nil {
		return diag.Errorf("could not delete UserProjectAssignment: %v", err)
	}

	return nil
}
