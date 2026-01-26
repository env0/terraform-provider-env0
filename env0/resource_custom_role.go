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

func resourceCustomRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomRoleCreate,
		ReadContext:   resourceCustomRoleRead,
		UpdateContext: resourceCustomRoleUpdate,
		DeleteContext: resourceCustomRoleDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceCustomRoleImport},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the custom role",
				Required:    true,
			},
			"permissions": {
				Type:        schema.TypeList,
				Description: "The list of permissions assigned to the role. The allowed values are available within the 'permissions' attribute at this URL: https://docs.envzero.com/api-reference/roles/create-a-role#create-a-role",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_default_role": {
				Type:        schema.TypeBool,
				Description: "Defaults to 'false'",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceCustomRoleCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.RoleCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	role, err := apiClient.RoleCreate(payload)
	if err != nil {
		return diag.Errorf("could not create a custom role: %v", err)
	}

	d.SetId(role.Id)

	return nil
}

func resourceCustomRoleRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	role, err := apiClient.Role(d.Id())
	if err != nil {
		return ResourceGetFailure(ctx, "role", d, err)
	}

	if err := writeResourceData(role, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceCustomRoleUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.RoleUpdatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if _, err := apiClient.RoleUpdate(d.Id(), payload); err != nil {
		return diag.Errorf("could not update custom role: %v", err)
	}

	return nil
}

func resourceCustomRoleDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	err := apiClient.RoleDelete(d.Id())
	if err != nil {
		return diag.Errorf("could not delete role: %v", err)
	}

	return nil
}

func getCustomRoleById(id string, meta any) (*client.Role, error) {
	apiClient := meta.(client.ApiClientInterface)

	return apiClient.Role(id)
}

func getCustomRoleByName(name string, meta any) (*client.Role, error) {
	apiClient := meta.(client.ApiClientInterface)

	roles, err := apiClient.Roles()
	if err != nil {
		return nil, err
	}

	var foundRoles []client.Role

	for _, role := range roles {
		if role.Name == name {
			foundRoles = append(foundRoles, role)
		}
	}

	if len(foundRoles) == 0 {
		return nil, fmt.Errorf("role with name %v not found", name)
	}

	if len(foundRoles) > 1 {
		return nil, fmt.Errorf("found multiple custom roles with name: %s. Use id instead or make sure role names are unique %v", name, foundRoles)
	}

	return &foundRoles[0], nil
}

func getCustomRole(ctx context.Context, id string, meta any) (*client.Role, error) {
	_, err := uuid.Parse(id)
	if err == nil {
		tflog.Info(ctx, "Resolving custom role by id", map[string]any{"id": id})

		return getCustomRoleById(id, meta)
	} else {
		tflog.Info(ctx, "Resolving custom role by name", map[string]any{"name": id})

		return getCustomRoleByName(id, meta)
	}
}

func resourceCustomRoleImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	role, err := getCustomRole(ctx, d.Id(), meta)

	if err != nil {
		return nil, err
	}

	if role == nil {
		return nil, fmt.Errorf("custom role with id %v not found", d.Id())
	}

	if err := writeResourceData(role, d); err != nil {
		return nil, fmt.Errorf("schema resource data serialization failed: %w", err)
	}

	return []*schema.ResourceData{d}, nil
}
