package client

import "fmt"

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
