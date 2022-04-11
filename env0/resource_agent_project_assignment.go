package env0

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Updating an agent project assignment overrides all project organization assignments.
// Therefore, must extract all existing assignments and append/remove the created/deleted assignment.
// Since Terraform may run the assignments in parallel a mutex is required.
var apaLock sync.Mutex

// id is <agent_id>_<project_id>

type AgentProjectAssignment struct {
	AgentId   string `json:"agent_id"`
	ProjectId string `json:"project_id"`
}

func GetAgentProjectAssignmentId(agentId string, projectId string) string {
	return agentId + "_" + projectId
}

func (a *AgentProjectAssignment) GetId() string {
	return GetAgentProjectAssignmentId(a.AgentId, a.ProjectId)
}

func resourceAgentProjectAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAgentProjectAssignmentCreate,
		ReadContext:   resourceAgentProjectAssignmentRead,
		DeleteContext: resourceAgentProjectAssignmentDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceAgentProjectAssignmentImport},

		Schema: map[string]*schema.Schema{
			"agent_id": {
				Type:        schema.TypeString,
				Description: "id of the agent",
				Required:    true,
				ForceNew:    true,
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

func resourceAgentProjectAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var newAssignment AgentProjectAssignment
	if err := readResourceData(&newAssignment, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	apaLock.Lock()
	defer apaLock.Unlock()

	apiClient := meta.(client.ApiClientInterface)
	assignments, err := apiClient.ProjectsAgentsAssignments()
	if err != nil {
		return diag.Errorf("could not get project-agent assignments: %v", err)
	}

	for projectId, agentId := range assignments.ProjectsAgents {
		if projectId == newAssignment.ProjectId && agentId == newAssignment.AgentId {
			return diag.Errorf("assignment for project id %v and agent id %v already exist", projectId, agentId)
		}
	}

	assignments.ProjectsAgents[newAssignment.ProjectId] = newAssignment.AgentId

	if _, err := apiClient.AssignAgentsToProjects(assignments.ProjectsAgents); err != nil {
		return diag.Errorf("could not update project-agent assignments: %v", err)
	}

	d.SetId(newAssignment.GetId())

	return nil
}

func resourceAgentProjectAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(client.ApiClientInterface)

	apaLock.Lock()
	defer apaLock.Unlock()

	assignments, err := apiClient.ProjectsAgentsAssignments()
	if err != nil {
		return diag.Errorf("could not get project-agent assignments: %v", err)
	}

	var assignment *AgentProjectAssignment
	for projectId, agentId := range assignments.ProjectsAgents {
		if d.Id() == GetAgentProjectAssignmentId(agentId.(string), projectId) {
			assignment = &AgentProjectAssignment{
				AgentId:   agentId.(string),
				ProjectId: projectId,
			}
			break
		}
	}

	if assignment == nil {
		log.Printf("[WARN] Drift Detected: Terraform will remove %s from state", d.Id())
		d.SetId("")
		return nil
	}

	if err := writeResourceData(assignment, d); err != nil {
		diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceAgentProjectAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectId := d.Get("project_id").(string)
	agentId := d.Get("agent_id").(string)

	apaLock.Lock()
	defer apaLock.Unlock()

	apiClient := meta.(client.ApiClientInterface)

	assignments, err := apiClient.ProjectsAgentsAssignments()
	if err != nil {
		return diag.Errorf("could not get project-agent assignments: %v", err)
	}

	newAssignments := make(map[string]interface{})

	// Remove from the assignments the deleted assignment.
	for projectIdOther, agentIdOther := range assignments.ProjectsAgents {
		if projectId != projectIdOther || agentId != agentIdOther {
			newAssignments[projectIdOther] = agentIdOther
		}
	}

	if _, err := apiClient.AssignAgentsToProjects(newAssignments); err != nil {
		return diag.Errorf("could not update project-agent assignments: %v", err)
	}

	return nil
}

func resourceAgentProjectAssignmentImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	splitAgentProject := strings.Split(d.Id(), "_")
	if len(splitAgentProject) != 2 {
		return nil, fmt.Errorf("the id %v is invalid must be <agent_id>_<project_id>", d.Id())
	}

	agentId := splitAgentProject[0]
	projectId := splitAgentProject[1]

	apaLock.Lock()
	defer apaLock.Unlock()

	apiClient := meta.(client.ApiClientInterface)

	assignments, err := apiClient.ProjectsAgentsAssignments()
	if err != nil {
		return nil, err
	}

	var assignment *AgentProjectAssignment
	for projectIdOther, agentIdOther := range assignments.ProjectsAgents {
		if projectIdOther == projectId && agentIdOther == agentId {
			assignment = &AgentProjectAssignment{
				AgentId:   agentId,
				ProjectId: projectId,
			}
			break
		}
	}

	if assignment == nil {
		return nil, fmt.Errorf("assignment with id %v not found", d.Id())
	}

	if err := writeResourceData(assignment, d); err != nil {
		diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return []*schema.ResourceData{d}, nil
}
