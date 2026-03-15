package env0

import (
	"context"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAgentPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAgentPoolCreate,
		ReadContext:   resourceAgentPoolRead,
		UpdateContext: resourceAgentPoolUpdate,
		DeleteContext: resourceAgentPoolDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceAgentPoolImport},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "the agent pool name",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "the agent pool description",
				Optional:    true,
			},
			"agent_key": {
				Type:        schema.TypeString,
				Description: "the agent key identifier. Deprecated: use id instead",
				Computed:    true,
				Deprecated:  "use id instead",
			},
			"logs": {
				Type:        schema.TypeList,
				Description: "self-hosted logs configuration (dynamo self-hosted)",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeString,
							Description: "the AWS account id for self-hosted logs",
							Required:    true,
						},
						"region": {
							Type:        schema.TypeString,
							Description: "the AWS region for self-hosted logs",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func agentPoolLogsFromSchema(d *schema.ResourceData) *client.AgentPoolLogsConfig {
	logsRaw, ok := d.GetOk("logs")
	if !ok {
		return nil
	}

	logsList := logsRaw.([]any)
	if len(logsList) == 0 {
		return nil
	}

	logsMap := logsList[0].(map[string]any)

	selfHosted := &client.AgentPoolSelfHostedLogs{
		AccountId: logsMap["account_id"].(string),
		Region:    logsMap["region"].(string),
	}

	return &client.AgentPoolLogsConfig{
		Dynamo: &client.AgentPoolDynamoLogs{
			SelfHosted: selfHosted,
		},
	}
}

func agentPoolLogsToSchema(logs *client.AgentPoolLogsConfig) []any {
	if logs == nil || logs.Dynamo == nil || logs.Dynamo.SelfHosted == nil {
		return nil
	}

	sh := logs.Dynamo.SelfHosted

	return []any{map[string]any{
		"account_id": sh.AccountId,
		"region":     sh.Region,
	}}
}

func setAgentPoolSchema(d *schema.ResourceData, agentPool *client.AgentPool) {
	d.Set("name", agentPool.Name)
	d.Set("description", agentPool.Description)
	d.Set("agent_key", agentPool.AgentKey)
	d.Set("logs", agentPoolLogsToSchema(agentPool.Logs))
}

func resourceAgentPoolCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	payload := client.AgentPoolCreatePayload{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	agentPool, err := apiClient.AgentPoolCreate(payload)
	if err != nil {
		return diag.Errorf("could not create agent pool: %v", err)
	}

	d.SetId(agentPool.Id)

	// Logs can only be set via PATCH after creation.
	logs := agentPoolLogsFromSchema(d)
	if logs != nil {
		updatePayload := client.AgentPoolUpdatePayload{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Logs:        logs,
		}

		_, err := apiClient.AgentPoolUpdate(agentPool.Id, updatePayload)
		if err != nil {
			return diag.Errorf("agent pool created successfully but failed to set logs configuration: %v", err)
		}
	}

	return resourceAgentPoolRead(ctx, d, meta)
}

func resourceAgentPoolRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	agentPool, err := apiClient.AgentPool(d.Id())
	if err != nil {
		return ResourceGetFailure(ctx, "agent pool", d, err)
	}

	setAgentPoolSchema(d, agentPool)

	return nil
}

func resourceAgentPoolUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	payload := client.AgentPoolUpdatePayload{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Logs:        agentPoolLogsFromSchema(d),
	}

	_, err := apiClient.AgentPoolUpdate(d.Id(), payload)
	if err != nil {
		return diag.Errorf("could not update agent pool: %v", err)
	}

	return resourceAgentPoolRead(ctx, d, meta)
}

func resourceAgentPoolDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	if err := apiClient.AgentPoolDelete(d.Id()); err != nil {
		return diag.Errorf("could not delete agent pool: %v", err)
	}

	return nil
}

func getAgentPoolByName(name string, meta any) (*client.AgentPool, error) {
	apiClient := meta.(client.ApiClientInterface)

	pools, err := apiClient.AgentPools()
	if err != nil {
		return nil, err
	}

	var foundPools []client.AgentPool

	for _, pool := range pools {
		if pool.Name == name {
			foundPools = append(foundPools, pool)
		}
	}

	if len(foundPools) == 0 {
		return nil, fmt.Errorf("agent pool with name %v not found", name)
	}

	if len(foundPools) > 1 {
		return nil, fmt.Errorf("found multiple agent pools with name: %s. Use id instead or make sure agent pool names are unique", name)
	}

	return &foundPools[0], nil
}

func resourceAgentPoolImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	apiClient := meta.(client.ApiClientInterface)

	// Try by id first; fall back to name lookup.
	agentPool, err := apiClient.AgentPool(d.Id())
	if err != nil {
		tflog.Info(ctx, "Agent pool not found by id, trying by name", map[string]any{"id": d.Id()})

		agentPool, err = getAgentPoolByName(d.Id(), meta)
		if err != nil {
			return nil, fmt.Errorf("could not get agent pool: %w", err)
		}
	}

	d.SetId(agentPool.Id)
	setAgentPoolSchema(d, agentPool)

	return []*schema.ResourceData{d}, nil
}
