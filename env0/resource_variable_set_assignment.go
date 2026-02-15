package env0

import (
	"context"
	"slices"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type variableSetAssignmentSchema struct {
	Scope   string
	ScopeId string
	SetIds  []string
}

func resourceVariableSetAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVariableSetAssignmentCreate,
		UpdateContext: resourceVariableSetAssignmentUpdate,
		ReadContext:   resourceVariableSetAssignmentRead,
		DeleteContext: resourceVariableSetAssignmentDelete,

		Description: "note: avoid assigning to environments using 'variable_sets' (see environment schema).",

		Schema: map[string]*schema.Schema{
			"scope": {
				Type:             schema.TypeString,
				Description:      "the resource(scope) type to assign to. Valid values: 'template', 'environment', 'module', 'organization', 'project', 'deployment'",
				Required:         true,
				ValidateDiagFunc: NewStringInValidator([]string{"template", "environment", "module", "organization", "project", "deployment"}),
				ForceNew:         true,
			},
			"scope_id": {
				Type:        schema.TypeString,
				Description: "the resource(scope)id (e.g. template id)",
				Required:    true,
				ForceNew:    true,
			},
			"set_ids": {
				Type:        schema.TypeList,
				Description: "list of variable sets",
				Required:    true,
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "the variable set id",
				},
			},
		},
	}
}

func resourceVariableSetAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var assignmentSchema variableSetAssignmentSchema

	if err := readResourceData(&assignmentSchema, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if len(assignmentSchema.SetIds) > 0 {
		if err := apiClient.AssignConfigurationSets(assignmentSchema.Scope, assignmentSchema.ScopeId, assignmentSchema.SetIds); err != nil {
			return diag.Errorf("failed to assign variable sets to the scope: %v", err)
		}
	}

	d.SetId(assignmentSchema.ScopeId)

	return nil
}

func resourceVariableSetAssignmentUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var assignmentSchema variableSetAssignmentSchema

	if err := readResourceData(&assignmentSchema, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	apiConfigurationSets, err := apiClient.ConfigurationSetsAssignments(assignmentSchema.Scope, assignmentSchema.ScopeId)
	if err != nil {
		return diag.Errorf("failed to get variable sets assignments: %v", err)
	}

	// Compare between apiSetIds and schemaSetIds to find what to set ids to delete and what set ids to add.
	var toDelete, toAdd []string

	// In API but not in Schema - delete.
	for _, apiConfigurationSet := range apiConfigurationSets {
		if apiConfigurationSet.AssignmentScopeId != assignmentSchema.ScopeId {
			continue
		}

		found := false

		apiSetId := apiConfigurationSet.Id
		if slices.Contains(assignmentSchema.SetIds, apiSetId) {
			found = true
		}

		if !found {
			toDelete = append(toDelete, apiSetId)
		}
	}

	// In Schema but not in API - add.
	for _, schemaSetId := range assignmentSchema.SetIds {
		found := false

		for _, apiConfigurationSet := range apiConfigurationSets {
			apiSetId := apiConfigurationSet.Id
			if schemaSetId == apiSetId {
				found = true

				break
			}
		}

		if !found {
			toAdd = append(toAdd, schemaSetId)
		}
	}

	if len(toDelete) > 0 {
		if err := apiClient.UnassignConfigurationSets(assignmentSchema.Scope, assignmentSchema.ScopeId, toDelete); err != nil {
			return diag.Errorf("failed to unassign variable sets: %v", err)
		}
	}

	if len(toAdd) > 0 {
		if err := apiClient.AssignConfigurationSets(assignmentSchema.Scope, assignmentSchema.ScopeId, toAdd); err != nil {
			return diag.Errorf("failed to assign variable sets: %v", err)
		}
	}

	return nil
}

func resourceVariableSetAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var assignmentSchema variableSetAssignmentSchema

	if err := readResourceData(&assignmentSchema, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if len(assignmentSchema.SetIds) > 0 {
		if err := apiClient.UnassignConfigurationSets(assignmentSchema.Scope, assignmentSchema.ScopeId, assignmentSchema.SetIds); err != nil {
			return diag.Errorf("failed to unassign variable sets: %v", err)
		}
	}

	return nil
}

func resourceVariableSetAssignmentRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var assignmentSchema variableSetAssignmentSchema

	if err := readResourceData(&assignmentSchema, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	apiConfigurationSets, err := apiClient.ConfigurationSetsAssignments(assignmentSchema.Scope, assignmentSchema.ScopeId)
	if err != nil {
		return diag.Errorf("failed to get variable sets assignments: %v", err)
	}

	newSchemaSetIds := []string{}

	// To avoid drifts keep the schema order as much as possible.
	for _, schemaSetId := range assignmentSchema.SetIds {
		for _, apiConfigurationSet := range apiConfigurationSets {
			apiSetId := apiConfigurationSet.Id

			if schemaSetId == apiSetId {
				newSchemaSetIds = append(newSchemaSetIds, schemaSetId)

				break
			}
		}
	}

	for _, apiConfigurationSet := range apiConfigurationSets {
		// Filter out inherited assignments (e.g parent project).
		if apiConfigurationSet.AssignmentScopeId != assignmentSchema.ScopeId {
			continue
		}

		apiSetId := apiConfigurationSet.Id
		found := slices.Contains(assignmentSchema.SetIds, apiSetId)

		if !found {
			newSchemaSetIds = append(newSchemaSetIds, apiSetId)
		}
	}

	d.Set("set_ids", newSchemaSetIds)

	return nil
}
