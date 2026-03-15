package env0

import (
	"context"
	"fmt"
	"strings"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAgentSecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAgentSecretCreate,
		ReadContext:   resourceAgentSecretRead,
		DeleteContext: resourceAgentSecretDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceAgentSecretImport},

		Schema: map[string]*schema.Schema{
			"agent_id": {
				Type:        schema.TypeString,
				Description: "the agent pool id this secret belongs to",
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "the secret description",
				Optional:    true,
				ForceNew:    true,
			},
			"secret": {
				Type:        schema.TypeString,
				Description: "the secret value. Only available on creation. Note: stored in Terraform state",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceAgentSecretCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	agentId := d.Get("agent_id").(string)

	payload := client.AgentSecretCreatePayload{
		Description: d.Get("description").(string),
	}

	agentSecret, err := apiClient.AgentSecretCreate(agentId, payload)
	if err != nil {
		return diag.Errorf("could not create agent secret: %v", err)
	}

	d.SetId(agentSecret.Id)
	d.Set("secret", agentSecret.Secret)

	return nil
}

func resourceAgentSecretRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	agentId := d.Get("agent_id").(string)

	secrets, err := apiClient.AgentSecrets(agentId)
	if err != nil {
		return ResourceGetFailure(ctx, "agent secret", d, err)
	}

	for _, secret := range secrets {
		if secret.Id == d.Id() {
			d.Set("description", secret.Description)
			// Don't overwrite the secret value in state — the GET endpoint doesn't return it.
			return nil
		}
	}

	// Secret not found — drift detected.
	tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]any{"id": d.Id()})
	d.SetId("")

	return nil
}

func resourceAgentSecretDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	agentId := d.Get("agent_id").(string)

	if err := apiClient.AgentSecretDelete(agentId, d.Id()); err != nil {
		return diag.Errorf("could not delete agent secret: %v", err)
	}

	return nil
}

func resourceAgentSecretImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	// Import format: agentId/secretId
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import format, expected agentId/secretId, got: %s", d.Id())
	}

	agentId := parts[0]
	secretId := parts[1]

	apiClient := meta.(client.ApiClientInterface)

	secrets, err := apiClient.AgentSecrets(agentId)
	if err != nil {
		return nil, fmt.Errorf("could not get agent secrets: %w", err)
	}

	for _, secret := range secrets {
		if secret.Id == secretId {
			d.SetId(secretId)
			d.Set("agent_id", agentId)
			d.Set("description", secret.Description)
			// Secret value is not available on import.

			return []*schema.ResourceData{d}, nil
		}
	}

	return nil, fmt.Errorf("agent secret with id %s not found for agent %s", secretId, agentId)
}
