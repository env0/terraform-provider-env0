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

func resourceNotificationProjectAssignment() *schema.Resource {
	allowedEventNames := []string{
		"environmentDeployStarted",
		"environmentDestroyStarted",
		"deploySucceeded",
		"destroySucceeded",
		"deployFailed",
		"destroyFailed",
		"deploymentWaitingForUser",
		"deploymentCancelled",
		"environmentDeployResumed",
		"environmentDestroyResumed",
		"environmentMarkedForAutoDestroy",
		"driftDetected",
		"driftUndetected",
		"driftFailed",
		"budgetExceeded",
	}

	return &schema.Resource{
		CreateContext: resourceNotificationProjectAssignmentCreateOrUpdate,
		ReadContext:   resourceNotificationProjectAssignmentRead,
		UpdateContext: resourceNotificationProjectAssignmentCreateOrUpdate,
		DeleteContext: resourceNotificationProjectAssignmentDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Description: "id of the project",
				Required:    true,
				ForceNew:    true,
			},
			"notification_endpoint_id": {
				Type:        schema.TypeString,
				Description: "id of the notification endpoint",
				Required:    true,
				ForceNew:    true,
			},
			"event_names": {
				Type:        schema.TypeList,
				Description: fmt.Sprintf("a list of notification events (allowed values: %s)", strings.Join(allowedEventNames, ", ")),
				Required:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: NewStringInValidator(allowedEventNames),
				},
			},
		},
	}
}

func resourceNotificationProjectAssignmentCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectId := d.Get("project_id").(string)
	endpointId := d.Get("notification_endpoint_id").(string)

	var payload client.NotificationProjectAssignmentUpdatePayload
	if err := readResourceData(&payload, d); err != nil {
		return diag.Errorf("schema resource data deserialization failed: %v", err)
	}

	apiClient := meta.(client.ApiClientInterface)

	assignment, err := apiClient.NotificationProjectAssignmentUpdate(projectId, endpointId, payload)
	if err != nil {
		return diag.Errorf("could not create or update notification project assignment: %v", err)
	}

	d.SetId(assignment.Id)

	return nil
}

func getNotificationProjectAssignment(d *schema.ResourceData, meta interface{}) (*client.NotificationProjectAssignment, error) {
	projectId := d.Get("project_id").(string)
	endpointId := d.Get("notification_endpoint_id").(string)

	apiClient := meta.(client.ApiClientInterface)

	assignments, err := apiClient.NotificationProjectAssignments(projectId)
	if err != nil {
		return nil, err
	}

	for _, assignment := range assignments {
		if assignment.NotificationEndpointId == endpointId {
			return &assignment, nil
		}
	}

	return nil, nil
}

func resourceNotificationProjectAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	assignment, err := getNotificationProjectAssignment(d, meta)
	if err != nil {
		return ResourceGetFailure(ctx, "notification project assignment", d, err)
	}
	if assignment == nil {
		// Notification endpoint not found.
		tflog.Warn(ctx, "Drift Detected: Terraform will remove id from state", map[string]interface{}{"id": d.Id()})
		d.SetId("")
		return nil
	}

	if err := writeResourceData(assignment, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func resourceNotificationProjectAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectId := d.Get("project_id").(string)
	endpointId := d.Get("notification_endpoint_id").(string)
	// There's no delete API call. When delete is called, remove all events.
	payload := client.NotificationProjectAssignmentUpdatePayload{
		EventNames: []string{},
	}

	apiClient := meta.(client.ApiClientInterface)

	if _, err := apiClient.NotificationProjectAssignmentUpdate(projectId, endpointId, payload); err != nil {
		return diag.Errorf("could not delete notification project assignment: %v", err)
	}

	return nil
}
