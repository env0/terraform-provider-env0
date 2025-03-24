package client

type User struct {
	CreatedAt   string         `json:"created_at"`
	Email       string         `json:"email"`
	FamilyName  string         `json:"family_name"`
	GivenName   string         `json:"given_name"`
	Name        string         `json:"name"`
	Picture     string         `json:"picture"`
	UserId      string         `json:"user_id"`
	LastLogin   string         `json:"last_login"`
	AppMetadata map[string]any `json:"app_metadata"`
}
