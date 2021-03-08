package env0apiclient

type Organization struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	CreatedBy    string `json:"createdBy"`
	Role         string `json:"role"`
	IsSelfHosted bool   `json:"isSelfHosted"`
}

type User struct {
	CreatedAt   string                 `json:"created_at"`
	Email       string                 `json:"email"`
	FamilyName  string                 `json:"family_name"`
	GivenName   string                 `json:"given_name"`
	Name        string                 `json:"name"`
	Picture     string                 `json:"picture"`
	UserId      string                 `json:"user_id"`
	LastLogin   string                 `json:"last_login"`
	AppMetadata map[string]interface{} `json:"app_metadata"`
}

type Project struct {
	IsArchived     bool   `json:"isArchived"`
	OrganizationId string `json:"organizationId"`
	UpdatedAt      string `json:"updatedAt"`
	CreatedAt      string `json:"createdAt"`
	Id             string `json:"id"`
	Name           string `json:"name"`
	CreatedBy      string `json:"createdBy"`
	Role           string `json:"role"`
	CreatedByUser  User   `json:"createdByUser"`
}
