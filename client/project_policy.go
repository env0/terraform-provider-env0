package client

type Policy struct {
	Id                          string  `json:"id" tfschema:"-"`
	ProjectId                   string  `json:"projectId"`
	NumberOfEnvironments        *int    `json:"numberOfEnvironments,omitempty" tfschema:",omitempty"`
	NumberOfEnvironmentsTotal   *int    `json:"numberOfEnvironmentsTotal,omitempty" tfschema:",omitempty"`
	RequiresApprovalDefault     bool    `json:"requiresApprovalDefault"`
	IncludeCostEstimation       bool    `json:"includeCostEstimation"`
	SkipApplyWhenPlanIsEmpty    bool    `json:"skipApplyWhenPlanIsEmpty"`
	DisableDestroyEnvironments  bool    `json:"disableDestroyEnvironments"`
	SkipRedundantDeployments    bool    `json:"skipRedundantDeployments"`
	UpdatedBy                   string  `json:"updatedBy"`
	RunPullRequestPlanDefault   bool    `json:"runPullRequestPlanDefault"`
	ContinuousDeploymentDefault bool    `json:"continuousDeploymentDefault"`
	MaxTtl                      *string `json:"maxTtl,omitempty" tfschema:",omitempty"`
	DefaultTtl                  *string `json:"defaultTtl,omitempty" tfschema:",omitempty"`
	ForceRemoteBackend          bool    `json:"forceRemoteBackend"`
	DriftDetectionCron          string  `json:"driftDetectionCron"`
	DriftDetectionEnabled       bool    `json:"driftDetectionEnabled"`
	VcsPrCommentsEnabledDefault bool    `json:"vcsPrCommentsEnabledDefault"`
}

type PolicyUpdatePayload struct {
	ProjectId                   string `json:"projectId"`
	NumberOfEnvironments        int    `json:"numberOfEnvironments,omitempty"`
	NumberOfEnvironmentsTotal   int    `json:"numberOfEnvironmentsTotal,omitempty"`
	RequiresApprovalDefault     bool   `json:"requiresApprovalDefault"`
	IncludeCostEstimation       bool   `json:"includeCostEstimation"`
	SkipApplyWhenPlanIsEmpty    bool   `json:"skipApplyWhenPlanIsEmpty"`
	DisableDestroyEnvironments  bool   `json:"disableDestroyEnvironments"`
	SkipRedundantDeployments    bool   `json:"skipRedundantDeployments"`
	RunPullRequestPlanDefault   bool   `json:"runPullRequestPlanDefault"`
	ContinuousDeploymentDefault bool   `json:"continuousDeploymentDefault"`
	MaxTtl                      string `json:"maxTtl,omitempty"`
	DefaultTtl                  string `json:"defaultTtl,omitempty"`
	ForceRemoteBackend          bool   `json:"forceRemoteBackend"`
	DriftDetectionCron          string `json:"driftDetectionCron"`
	DriftDetectionEnabled       bool   `json:"driftDetectionEnabled"`
	VcsPrCommentsEnabledDefault bool   `json:"vcsPrCommentsEnabledDefault"`
}

// Policy retrieves a policy from the API
func (client *ApiClient) Policy(projectId string) (Policy, error) {
	u, err := newQueryURL("/policies", parameter{"projectId", projectId})
	if err != nil {
		return Policy{}, err
	}

	var result Policy
	err = client.http.Get(u.String(), nil, &result)
	if err != nil {
		return Policy{}, err
	}

	return result, nil
}

// PolicyUpdate updates a policy through the API
func (client *ApiClient) PolicyUpdate(payload PolicyUpdatePayload) (Policy, error) {
	var result Policy
	err := client.http.Put("/policies", payload, &result)
	if err != nil {
		return Policy{}, err
	}
	return result, nil
}
