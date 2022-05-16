package client

import "errors"

type TeamCreatePayload struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	OrganizationId string `json:"organizationId"`
}

type TeamUpdatePayload struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	UserIds     []string `json:"userIds,omitempty"`
}

type Team struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	OrganizationId string `json:"organizationId"`
	Users          []User `json:"users"`
}

func (client *ApiClient) TeamCreate(payload TeamCreatePayload) (Team, error) {
	if payload.Name == "" {
		return Team{}, errors.New("must specify team name on creation")
	}
	if payload.OrganizationId != "" {
		return Team{}, errors.New("must not specify organizationId")
	}
	organizationId, err := client.organizationId()
	if err != nil {
		return Team{}, err
	}
	payload.OrganizationId = organizationId

	var result Team
	err = client.http.Post("/teams", payload, &result)
	if err != nil {
		return Team{}, err
	}
	return result, nil
}

func (client *ApiClient) Team(id string) (Team, error) {
	var result Team
	err := client.http.Get("/teams/"+id, nil, &result)
	if err != nil {
		return Team{}, err
	}
	return result, nil
}

func (client *ApiClient) TeamDelete(id string) error {
	return client.http.Delete("/teams/" + id)
}

func (client *ApiClient) TeamUpdate(id string, payload TeamUpdatePayload) (Team, error) {
	if payload.Name == "" {
		return Team{}, errors.New("must specify team name on update")
	}

	var result Team
	err := client.http.Put("/teams/"+id, payload, &result)
	if err != nil {
		return Team{}, err
	}
	return result, nil
}

func (client *ApiClient) Teams() ([]Team, error) {
	organizationId, err := client.organizationId()
	if err != nil {
		return nil, err
	}
	var result []Team
	err = client.http.Get("/teams/organizations/"+organizationId, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}
