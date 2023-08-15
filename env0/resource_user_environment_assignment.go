package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserEnvironmentAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserEnvironmentAssignmentCreate,
		UpdateContext: resourceUserEnvironmentAssignmentUpdate,
		ReadContext:   resourceUserEnvironmentAssignmentRead,
		DeleteContext: resourceUserEnvironmentAssignmentDelete,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Description: "id of the user",
				Required:    true,
				ForceNew:    true,
			},
			"environment_id": {
				Type:        schema.TypeString,
				Description: "id of the environment",
				Required:    true,
				ForceNew:    true,
			},
			"role_id": {
				Type:             schema.TypeString,
				Description:      "id of the assigned role",
				Required:         true,
				ValidateDiagFunc: ValidateNotEmptyString,
			},
		},
	}
}

func resourceUserEnvironmentAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var newAssignment client.AssignUserRoleToEnvironmentPayload
	if err := readResourceData(&newAssignment, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	client := meta.(client.ApiClientInterface)

	assignment, err := client.AssignUserRoleToEnvironment(&newAssignment)
	if err != nil {
		return diag.Errorf("could not create assignment: %v", err)
	}

	d.SetId(assignment.Id)

	return nil
}

func resourceUserEnvironmentAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(client.ApiClientInterface)

	id := d.Id()
	environmentId := d.Get("environment_id").(string)

	assignments, err := client.UserRoleEnvironmentAssignments(environmentId)
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

func resourceUserEnvironmentAssignmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var payload client.AssignUserRoleToEnvironmentPayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	client := meta.(client.ApiClientInterface)

	assignment, err := client.AssignUserRoleToEnvironment(&payload)
	if err != nil {
		return diag.Errorf("could not update assignment: %v", err)
	}

	d.SetId(assignment.Id)

	return nil
}

func resourceUserEnvironmentAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	environmentId := d.Get("environment_id").(string)
	userId := d.Get("user_id").(string)

	client := meta.(client.ApiClientInterface)

	if err := client.RemoveUserRoleFromEnvironment(environmentId, userId); err != nil {
		return diag.Errorf("could not delete assignment: %v", err)
	}

	return nil
}
