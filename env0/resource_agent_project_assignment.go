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

const ENV0_DEFAULT = "ENV0_DEFAULT"

func getAgentProjectAssignment(d *schema.ResourceData, meta interface{}) (*AgentProjectAssignment, error) {
	apiClient := meta.(client.ApiClientInterface)

	assignments, err := apiClient.ProjectsAgentsAssignments()
	if err != nil {
		return nil, fmt.Errorf("failed to get project agent assignments: %w", err)
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

	return &assignment, nil
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
	assignment, err := getAgentProjectAssignment(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := writeResourceData(assignment, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceAgentProjectAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	// When deleting an assignment, revert the project assignment to the default.

	payload := client.AssignProjectsAgentsAssignmentsPayload{
		d.Id(): "ENV0_DEFAULT",
	}

	if _, err := apiClient.AssignAgentsToProjects(payload); err != nil {
		return diag.Errorf("failed to assign project '%s' to back to default agent: %v", d.Id(), err)
	}

	return nil
}

func resourceAgentProjectAssignmentImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	apiClient := meta.(client.ApiClientInterface)

	// Validate the project exists.
	if _, err := apiClient.Project(d.Id()); err != nil {
		return nil, fmt.Errorf("unable to get or find a project with id '%s': %w", d.Id(), err)
	}

	assignment, err := getAgentProjectAssignment(d, meta)
	if err != nil {
		return nil, err
	}

	if err := writeResourceData(assignment, d); err != nil {
		return nil, fmt.Errorf("schema resource data serialization failed: %v", err)
	}

	return []*schema.ResourceData{d}, nil
}
