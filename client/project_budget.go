package client

type ProjectBudget struct {
	Id         string `json:"id" tfschema:"-"`
	ProjectId  string `json:"projectId"`
	Amount     int    `json:"amount"`
	Timeframe  string `json:"timeframe"`
	Thresholds []int  `json:"thresholds"`
}

type ProjectBudgetUpdatePayload struct {
	Amount     int    `json:"amount"`
	Timeframe  string `json:"timeframe"`
	Thresholds []int  `json:"thresholds"`
}

func (client *ApiClient) ProjectBudget(projectId string) (*ProjectBudget, error) {
	var result ProjectBudget

	if err := client.http.Get("/costs/project/"+projectId+"/budget", nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) ProjectBudgetUpdate(projectId string, payload *ProjectBudgetUpdatePayload) (*ProjectBudget, error) {
	var result ProjectBudget

	if payload.Thresholds == nil {
		payload.Thresholds = []int{}
	}

	err := client.http.Put("/costs/project/"+projectId+"/budget", payload, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *ApiClient) ProjectBudgetDelete(projectId string) error {
	return client.http.Delete("/costs/project/"+projectId+"/budget", nil)
}
