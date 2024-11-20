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

type PaginatedTeamsResponse struct {
	Teams       []Team `json:"teams"`
	NextPageKey string `json:"nextPageKey"`
}

func (client *ApiClient) TeamCreate(payload TeamCreatePayload) (Team, error) {
	if payload.Name == "" {
		return Team{}, errors.New("must specify team name on creation")
	}

	if payload.OrganizationId != "" {
		return Team{}, errors.New("must not specify organizationId")
	}

	organizationId, err := client.OrganizationId()
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
	return client.http.Delete("/teams/"+id, nil)
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

func (client *ApiClient) GetTeams(params map[string]string) ([]Team, error) {
	organizationId, err := client.OrganizationId()
	if err != nil {
		return nil, err
	}

	var teams []Team

	for {
		var res PaginatedTeamsResponse

		if err := client.http.Get("/teams/organizations/"+organizationId, params, &res); err != nil {
			return nil, err
		}

		teams = append(teams, res.Teams...)

		nextPageKey := res.NextPageKey
		if nextPageKey == "" {
			break
		}

		params["offset"] = nextPageKey
	}

	return teams, nil
}

func (client *ApiClient) Teams() ([]Team, error) {
	return client.GetTeams(map[string]string{"limit": "100"})
}

func (client *ApiClient) TeamsByName(name string) ([]Team, error) {
	return client.GetTeams(map[string]string{"name": name, "limit": "100"})
}
