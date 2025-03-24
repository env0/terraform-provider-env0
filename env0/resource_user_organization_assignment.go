package env0

import (
	"context"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func isBuiltinOrganizationRole(role string) bool {
	return role == "Admin" || role == "User"
}

func resourceUserOrganizationAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserOrganizationAssignmentCreateOrUpdate,
		ReadContext:   resourceUserOrganizationAssignmentRead,
		UpdateContext: resourceUserOrganizationAssignmentCreateOrUpdate,
		DeleteContext: resourceUserOrganizationAssignmentDelete,

		Description: "note: if removed the user role will be reset to 'Organization User'",

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Description: "id of the user",
				Required:    true,
				ForceNew:    true,
			},
			"custom_role_id": {
				Type:             schema.TypeString,
				Description:      "id of the custom role",
				Optional:         true,
				ValidateDiagFunc: ValidateNotEmptyString,
				ExactlyOneOf:     []string{"custom_role_id", "role"},
			},
			"role": {
				Type:             schema.TypeString,
				Description:      "the assigned built-in roles (User or Admin)",
				Optional:         true,
				ValidateDiagFunc: NewStringInValidator([]string{"User", "Admin"}),
				ExactlyOneOf:     []string{"custom_role_id", "role"},
			},
		},
	}
}

func resourceUserOrganizationAssignmentRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	userId := d.Get("user_id").(string)

	users, err := apiClient.Users()
	if err != nil {
		return diag.Errorf("could not get list of users: %v", err)
	}

	var user *client.OrganizationUser

	for i := range users {
		if users[i].User.UserId == userId {
			user = &users[i]

			break
		}
	}

	if user == nil {
		tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]any{"id": d.Id()})
		d.SetId("")

		return nil
	}

	if isBuiltinOrganizationRole(user.Role) {
		d.Set("role", user.Role)
	} else {
		d.Set("custom_role_id", user.Role)
	}

	return nil
}

func resourceUserOrganizationAssignmentCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	userId := d.Get("user_id").(string)

	role, ok := d.GetOk("role")
	if !ok {
		role = d.Get("custom_role_id")
	}

	if err := apiClient.OrganizationUserUpdateRole(userId, role.(string)); err != nil {
		return diag.Errorf("failed to update user role organization: %v", err)
	}

	d.SetId(userId)

	return nil
}

func resourceUserOrganizationAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	userId := d.Get("user_id").(string)

	if err := apiClient.OrganizationUserUpdateRole(userId, "User"); err != nil {
		return diag.Errorf("failed to update user role organization: %v", err)
	}

	return nil
}
