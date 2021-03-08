package env0apiclient

import (
	"errors"
	"net/http"
	"net/url"
)

type ApiClient struct {
	Endpoint             string
	ApiKey               string
	ApiSecret            string
	client               *http.Client
	cachedOrganizationId string
}

func (self *ApiClient) Organization() (Organization, error) {
	var result []Organization
	err := self.getJSON("/organizations", nil, &result)
	if err != nil {
		return Organization{}, err
	}
	if len(result) != 1 {
		return Organization{}, errors.New("Server responded with a too many organizations")
	}
	return result[0], nil
}

func (self *ApiClient) Projects() ([]Project, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return nil, err
	}
	var result []Project
	params := url.Values{}
	params.Add("organizationId", organizationId)
	err = self.getJSON("/projects", params, &result)
	if err != nil {
		return []Project{}, err
	}
	return result, nil
}

func (self *ApiClient) Project(id string) (Project, error) {
	var result Project
	err := self.getJSON("/projects/"+id, nil, &result)
	if err != nil {
		return Project{}, err
	}
	return result, nil
}

func (self *ApiClient) ProjectByName(name string) (Project, error) {
	projects, err := self.Projects()
	if err != nil {
		return Project{}, nil
	}
	for _, project := range projects {
		if project.Name == name {
			return project, nil
		}
	}
	return Project{}, errors.New("No project matched the project name provided")
}
