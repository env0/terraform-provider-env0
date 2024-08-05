package env0

import (
	"context"
	"fmt"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AgentProjectAssignment struct {
	AgentId   string
	ProjectId string
}

func resourceAgentProjectAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAgentProjectAssignmentCreateOrUpdate,
		UpdateContext: resourceAgentProjectAssignmentCreateOrUpdate,
		ReadContext:   resourceAgentProjectAssignmentRead,
		DeleteContext: resourceAgentProjectAssignmentDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceAgentProjectAssignmentImport},

		Description: "assign an agent to a project (multiple self-hosted agents). More details here: https://docs.env0.com/docs/multiple-self-hosted-agents",

		Schema: map[string]*schema.Schema{
			"agent_id": {
				Type:        schema.TypeString,
				Description: "id of the agent",
				Required:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "id of the project",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceAgentProjectAssignmentCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	var assignment AgentProjectAssignment
	if err := readResourceData(&assignment, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	payload := client.AssignProjectsAgentsAssignmentsPayload{
		assignment.ProjectId: assignment.AgentId,
	}

	if _, err := apiClient.AssignAgentsToProjects(payload); err != nil {
		return diag.Errorf("failed to assign project '%s' to agent '%s': %v", assignment.ProjectId, assignment.AgentId, err)
	}

	d.SetId(assignment.ProjectId)

	return nil
}

func resourceAgentProjectAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	assignments, err := apiClient.ProjectsAgentsAssignments()
	if err != nil {
		return diag.Errorf("failed to get project agent assignments: %v", err)
	}

	// If there's no assignment for the project, it should default to the default agent.

	assignment := AgentProjectAssignment{
		AgentId:   assignments.DefaultAgent,
		ProjectId: d.Id(),
	}

	for projectId, agent := range assignments.ProjectsAgents {
		if projectId == assignment.ProjectId {
			assignment.AgentId = agent.(string)
			break
		}
	}

	if err := writeResourceData(&assignment, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceAgentProjectAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	// When deleting an assignment, revert the project assignment to the default agent.

	assignments, err := apiClient.ProjectsAgentsAssignments()
	if err != nil {
		return diag.Errorf("failed to get project agent assignments: %v", err)
	}

	payload := client.AssignProjectsAgentsAssignmentsPayload{
		d.Id(): assignments.DefaultAgent,
	}

	if _, err := apiClient.AssignAgentsToProjects(payload); err != nil {
		return diag.Errorf("failed to assign project '%s' to back to default agent '%s': %v", d.Id(), assignments.DefaultAgent, err)
	}

	return nil
}

func resourceAgentProjectAssignmentImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	apiClient := meta.(client.ApiClientInterface)

	// Validate the project exists.
	if _, err := apiClient.Project(d.Id()); err != nil {
		return nil, fmt.Errorf("unable to get or find a project with id '%s': %w", d.Id(), err)
	}

	assignments, err := apiClient.ProjectsAgentsAssignments()
	if err != nil {
		return nil, fmt.Errorf("failed to get project agent assignments: %w", err)
	}

	// Import using the default agnet if there's no assignment for the project.

	assignment := AgentProjectAssignment{
		AgentId:   assignments.DefaultAgent,
		ProjectId: d.Id(),
	}

	for projectId, agent := range assignments.ProjectsAgents {
		if projectId == assignment.ProjectId {
			assignment.AgentId = agent.(string)
			break
		}
	}
	if err := writeResourceData(&assignment, d); err != nil {
		return nil, fmt.Errorf("schema resource data serialization failed: %v", err)
	}

	return []*schema.ResourceData{d}, nil
}
