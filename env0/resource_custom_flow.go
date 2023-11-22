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

func resourceCustomFlow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomFlowCreate,
		ReadContext:   resourceCustomFlowRead,
		UpdateContext: resourceCustomFlowUpdate,
		DeleteContext: resourceCustomFlowDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceCustomFlowImport},

		Schema: getConfigurationTemplateSchema(CustomFlow),
	}
}

func resourceCustomFlowCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.CustomFlowCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	customFlow, err := apiClient.CustomFlowCreate(payload)
	if err != nil {
		return diag.Errorf("could not create custom flow: %v", err)
	}

	d.SetId(customFlow.Id)

	return nil
}

func resourceCustomFlowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	customFlow, err := apiClient.CustomFlow(d.Id())
	if err != nil {
		return ResourceGetFailure(ctx, "custom flow", d, err)
	}

	if err := writeResourceData(customFlow, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceCustomFlowUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var payload client.CustomFlowCreatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	if _, err := apiClient.CustomFlowUpdate(d.Id(), payload); err != nil {
		return diag.Errorf("could not update custom flow: %v", err)
	}

	return nil
}

func resourceCustomFlowDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.CustomFlowDelete(d.Id()); err != nil {
		return diag.Errorf("could not delete custom flow: %v", err)
	}

	return nil
}

func getCustomFlowByName(name string, meta interface{}) (*client.CustomFlow, error) {
	apiClient := meta.(client.ApiClientInterface)

	customFlows, err := apiClient.CustomFlows(name)
	if err != nil {
		return nil, err
	}

	if len(customFlows) == 0 {
		return nil, fmt.Errorf("custom flow with name %v not found", name)
	}

	if len(customFlows) > 1 {
		return nil, fmt.Errorf("found multiple custom flows with name '%s'. Use id instead or make sure custom flow names are unique %v", name, customFlows)
	}

	return &customFlows[0], nil
}

func getCustomFlow(ctx context.Context, id string, meta interface{}) (*client.CustomFlow, error) {
	if _, err := uuid.Parse(id); err == nil {
		tflog.Info(ctx, "Resolving custom flow by id", map[string]interface{}{"id": id})
		return meta.(client.ApiClientInterface).CustomFlow(id)
	} else {
		tflog.Info(ctx, "Resolving custom flow by name", map[string]interface{}{"name": id})
		return getCustomFlowByName(id, meta)
	}
}

func resourceCustomFlowImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	customFlow, err := getCustomFlow(ctx, d.Id(), meta)
	if err != nil {
		return nil, err
	}

	if err := writeResourceData(customFlow, d); err != nil {
		return nil, fmt.Errorf("schema resource data serialization failed: %v", err)
	}

	return []*schema.ResourceData{d}, nil
}
