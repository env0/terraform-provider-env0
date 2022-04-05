package client

import "errors"

type TeamCreatePayload struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	OrganizationId string `json:"organizationId"`
}

type TeamUpdatePayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Team struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	OrganizationId string `json:"organizationId"`
}

func (self *ApiClient) TeamCreate(payload TeamCreatePayload) (Team, error) {
	if payload.Name == "" {
		return Team{}, errors.New("Must specify team name on creation")
	}
	if payload.OrganizationId != "" {
		return Team{}, errors.New("Must not specify organizationId")
	}
	organizationId, err := self.organizationId()
	if err != nil {
		return Team{}, err
	}
	payload.OrganizationId = organizationId

	var result Team
	err = self.http.Post("/teams", payload, &result)
	if err != nil {
		return Team{}, err
	}
	return result, nil
}

func (self *ApiClient) Team(id string) (Team, error) {
	var result Team
	err := self.http.Get("/teams/"+id, nil, &result)
	if err != nil {
		return Team{}, err
	}
	return result, nil
}

func (self *ApiClient) TeamDelete(id string) error {
	return self.http.Delete("/teams/" + id)
}

func (self *ApiClient) TeamUpdate(id string, payload TeamUpdatePayload) (Team, error) {
	if payload.Name == "" {
		return Team{}, errors.New("Must specify team name on update")
	}

	var result Team
	err := self.http.Put("/teams/"+id, payload, &result)
	if err != nil {
		return Team{}, err
	}
	return result, nil
}

func (self *ApiClient) Teams() ([]Team, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return nil, err
	}
	var result []Team
	err = self.http.Get("/teams/organizations/"+organizationId, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}
