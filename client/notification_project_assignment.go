package client

import "fmt"

type NotificationProjectAssignment struct {
	Id                     string   `json:"id"`
	NotificationEndpointId string   `json:"notificationEndpointId"`
	EventNames             []string `json:"eventNames"`
	CreatedBy              string   `json:"createdBy"`
	CreatedByUser          User     `json:"createdByUser"`
}

type NotificationProjectAssignmentUpdatePayload struct {
	EventNames []string `json:"eventNames"`
}

func (ac *ApiClient) NotificationProjectAssignments(projectId string) ([]NotificationProjectAssignment, error) {
	var result []NotificationProjectAssignment
	if err := ac.http.Get("/notifications/projects/"+projectId, nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (ac *ApiClient) NotificationProjectAssignmentUpdate(projectId string, endpointId string, payload NotificationProjectAssignmentUpdatePayload) (*NotificationProjectAssignment, error) {
	var result NotificationProjectAssignment
	url := fmt.Sprintf("/notifications/projects/%s/endpoints/%s", projectId, endpointId)
	if err := ac.http.Put(url, payload, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
