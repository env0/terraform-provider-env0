package env0apiclient

type Organization struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	CreatedBy    string `json:"createdBy"`
	Role         string `json:"role"`
	IsSelfHosted bool   `json:"isSelfHosted"`
}
