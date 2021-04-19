package env0apiclient

import (
	"net/url"
)

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

func (self *ApiClient) ProjectCreate(name string) (Project, error) {
	var result Project
	request := map[string]interface{}{"name": name}
	err := self.postJSON("/projects", request, &result)
	if err != nil {
		return Project{}, err
	}
	return result, nil
}

func (self *ApiClient) ProjectDelete(id string) error {
	return self.delete("/projects/" + id)
}
