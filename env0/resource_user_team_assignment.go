package env0

import (
	"context"
	"fmt"
	"sync"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Updating a user team assignment overrides all user team assignments
// Therefore, must extract all existing users and append/remove the created/deleted user.
// Since Terraform may run the assignments in parallel a mutex is required.
var utaLock sync.Mutex

// id is <user_id>_<team_id>

type UserTeamAssignment struct {
	UserId string `json:"user_id"`
	TeamId string `json:"team_id"`
}

func GetUserTeamAssignmentId(userId string, teamId string) string {
	return userId + "_" + teamId
}

func (a *UserTeamAssignment) GetId() string {
	return GetUserTeamAssignmentId(a.UserId, a.TeamId)
}

func GetUserTeamAssignmentFromId(id string) (*UserTeamAssignment, error) {
	// lastSplit is used to avoid issues where the user_id has underscores in it.
	splitUserTeam := lastSplit(id, "_")
	if len(splitUserTeam) != 2 {
		return nil, fmt.Errorf("the id %v is invalid must be <user_id>_<team_id>", id)
	}
	return &UserTeamAssignment{
		UserId: splitUserTeam[0],
		TeamId: splitUserTeam[1],
	}, nil
}

func resourceUserTeamAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserTeamAssignmentCreate,
		ReadContext:   resourceUserTeamAssignmentRead,
		DeleteContext: resourceUserTeamAssignmentDelete,

		Importer: &schema.ResourceImporter{StateContext: resourceUserTeamAssignmentImport},

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Description: `id of the user. Note: can also be an id of a "User" API key`,
				Required:    true,
				ForceNew:    true,
			},
			"team_id": {
				Type:        schema.TypeString,
				Description: "id of the team",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceUserTeamAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var newAssignment UserTeamAssignment
	if err := readResourceData(&newAssignment, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	apiClient := meta.(client.ApiClientInterface)

	utaLock.Lock()
	defer utaLock.Unlock()

	team, err := apiClient.Team(newAssignment.TeamId)
	if err != nil {
		return diag.Errorf("could not get team: %v", err)
	}

	userIds := []string{newAssignment.UserId}

	for _, user := range team.Users {
		if user.UserId == newAssignment.UserId {
			return diag.Errorf("assignment for user id %v and team id %v already exist", newAssignment.UserId, newAssignment.TeamId)
		}
		userIds = append(userIds, user.UserId)
	}

	if _, err := apiClient.TeamUpdate(team.Id, client.TeamUpdatePayload{
		Name:        team.Name,
		Description: team.Description,
		UserIds:     userIds,
	}); err != nil {
		return diag.Errorf("could not update team with new assignment: %v", err)
	}

	d.SetId(newAssignment.GetId())

	return nil
}

func resourceUserTeamAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	assignment, err := GetUserTeamAssignmentFromId(d.Id())
	if err != nil {
		return diag.Errorf("%v", err)
	}

	apiClient := meta.(client.ApiClientInterface)

	utaLock.Lock()
	defer utaLock.Unlock()

	team, err := apiClient.Team(assignment.TeamId)
	if err != nil {
		return ResourceGetFailure(ctx, "team", d, err)
	}

	found := false
	for _, user := range team.Users {
		if user.UserId == assignment.UserId {
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]interface{}{"id": d.Id()})
		d.SetId("")
		return nil
	}

	if err := writeResourceData(assignment, d); err != nil {
		diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceUserTeamAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var assignment UserTeamAssignment
	if err := readResourceData(&assignment, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	apiClient := meta.(client.ApiClientInterface)

	utaLock.Lock()
	defer utaLock.Unlock()

	team, err := apiClient.Team(assignment.TeamId)
	if err != nil {
		return diag.Errorf("could not get team: %v", err)
	}

	userIds := []string{}

	for _, user := range team.Users {
		if user.UserId == assignment.UserId {
			continue
		}
		userIds = append(userIds, user.UserId)
	}

	if _, err := apiClient.TeamUpdate(team.Id, client.TeamUpdatePayload{
		Name:        team.Name,
		Description: team.Description,
		UserIds:     userIds,
	}); err != nil {
		return diag.Errorf("could not update team with removed assignment: %v", err)
	}

	return nil
}

func resourceUserTeamAssignmentImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	assignment, err := GetUserTeamAssignmentFromId(d.Id())
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	apiClient := meta.(client.ApiClientInterface)

	utaLock.Lock()
	defer utaLock.Unlock()

	team, err := apiClient.Team(assignment.TeamId)
	if err != nil {
		if frerr, ok := err.(*http.FailedResponseError); ok && frerr.NotFound() {
			return nil, fmt.Errorf("team %v not found", assignment.TeamId)
		}
		return nil, err
	}

	found := false
	for _, user := range team.Users {
		if user.UserId == assignment.UserId {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("user %v not assigned to team %v", assignment.UserId, assignment.TeamId)
	}

	if err := writeResourceData(assignment, d); err != nil {
		diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return []*schema.ResourceData{d}, nil
}
