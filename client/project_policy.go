package client

import "encoding/json"

type Policy struct {
	Id                          string `json:"id"`
	ProjectId                   string `json:"projectId"`
	NumberOfEnvironments        int    `json:"numberOfEnvironments"`
	NumberOfEnvironmentsTotal   int    `json:"numberOfEnvironmentsTotal"`
	RequiresApprovalDefault     bool   `json:"requiresApprovalDefault"`
	IncludeCostEstimation       bool   `json:"includeCostEstimation"`
	SkipApplyWhenPlanIsEmpty    bool   `json:"skipApplyWhenPlanIsEmpty"`
	DisableDestroyEnvironments  bool   `json:"disableDestroyEnvironments"`
	SkipRedundantDeployments    bool   `json:"skipRedundantDeployments"`
	UpdatedBy                   string `json:"updatedBy"`
	RunPullRequestPlanDefault   bool   `json:"runPullRequestPlanDefault"`
	ContinuousDeploymentDefault bool   `json:"continuousDeploymentDefault"`
}

type PolicyUpdatePayload struct {
	ProjectId                   string `json:"projectId"`
	NumberOfEnvironments        int    `json:"numberOfEnvironments"`
	NumberOfEnvironmentsTotal   int    `json:"numberOfEnvironmentsTotal"`
	RequiresApprovalDefault     bool   `json:"requiresApprovalDefault"`
	IncludeCostEstimation       bool   `json:"includeCostEstimation"`
	SkipApplyWhenPlanIsEmpty    bool   `json:"skipApplyWhenPlanIsEmpty"`
	DisableDestroyEnvironments  bool   `json:"disableDestroyEnvironments"`
	SkipRedundantDeployments    bool   `json:"skipRedundantDeployments"`
	RunPullRequestPlanDefault   bool   `json:"runPullRequestPlanDefault"`
	ContinuousDeploymentDefault bool   `json:"continuousDeploymentDefault"`
}

func (p PolicyUpdatePayload) MarshalJSON() ([]byte, error) {
	type serial struct {
		ProjectId                   string `json:"projectId"`
		NumberOfEnvironments        *int   `json:"numberOfEnvironments"`
		NumberOfEnvironmentsTotal   *int   `json:"numberOfEnvironmentsTotal"`
		RequiresApprovalDefault     bool   `json:"requiresApprovalDefault"`
		IncludeCostEstimation       bool   `json:"includeCostEstimation"`
		SkipApplyWhenPlanIsEmpty    bool   `json:"skipApplyWhenPlanIsEmpty"`
		DisableDestroyEnvironments  bool   `json:"disableDestroyEnvironments"`
		SkipRedundantDeployments    bool   `json:"skipRedundantDeployments"`
		RunPullRequestPlanDefault   bool   `json:"runPullRequestPlanDefault"`
		ContinuousDeploymentDefault bool   `json:"continuousDeploymentDefault"`
	}

	s := serial{
		ProjectId:                   p.ProjectId,
		RequiresApprovalDefault:     p.RequiresApprovalDefault,
		IncludeCostEstimation:       p.IncludeCostEstimation,
		SkipApplyWhenPlanIsEmpty:    p.SkipApplyWhenPlanIsEmpty,
		DisableDestroyEnvironments:  p.DisableDestroyEnvironments,
		SkipRedundantDeployments:    p.SkipRedundantDeployments,
		RunPullRequestPlanDefault:   p.RunPullRequestPlanDefault,
		ContinuousDeploymentDefault: p.ContinuousDeploymentDefault,
	}

	if p.NumberOfEnvironments != 0 {
		s.NumberOfEnvironments = &p.NumberOfEnvironments
	}
	if p.NumberOfEnvironmentsTotal != 0 {
		s.NumberOfEnvironmentsTotal = &p.NumberOfEnvironmentsTotal
	}

	return json.Marshal(s)
}

// Policy retrieves a policy from the API
func (self *ApiClient) Policy(projectId string) (Policy, error) {
	u, err := newQueryURL("/policies", parameter{"projectId", projectId})
	if err != nil {
		return Policy{}, err
	}

	var result Policy
	err = self.http.Get(u.String(), nil, &result)
	if err != nil {
		return Policy{}, err
	}
	return result, nil
}

// PolicyUpdate updates a policy through the API
func (self *ApiClient) PolicyUpdate(payload PolicyUpdatePayload) (Policy, error) {
	var result Policy
	err := self.http.Put("/policies", payload, &result)
	if err != nil {
		return Policy{}, err
	}
	return result, nil
}
