package client

//templates are actually called "blueprints" in some parts of the API, this layer
//attempts to abstract this detail away - all the users of api client should
//only use "template", no mention of blueprint

import (
	"errors"
)

func (self *ApiClient) TemplateCreate(payload TemplateCreatePayload) (Template, error) {
	if payload.Name == "" {
		return Template{}, errors.New("Must specify template name on creation")
	}
	if payload.OrganizationId != "" {
		return Template{}, errors.New("Must not specify organizationId")
	}
	organizationId, err := self.organizationId()
	if err != nil {
		return Template{}, nil
	}
	payload.OrganizationId = organizationId

	var result Template
	err = self.http.Post("/blueprints", payload, &result)
	if err != nil {
		return Template{}, err
	}
	return result, nil
}

func (self *ApiClient) Template(id string) (Template, error) {
	var result Template
	err := self.http.Get("/blueprints/"+id, nil, &result)
	if err != nil {
		return Template{}, err
	}
	return result, nil
}

func (self *ApiClient) TemplateDelete(id string) error {
	return self.http.Delete("/blueprints/" + id)
}

func (self *ApiClient) TemplateUpdate(id string, payload TemplateCreatePayload) (Template, error) {
	if payload.Name == "" {
		return Template{}, errors.New("Must specify template name on creation")
	}
	if payload.OrganizationId != "" {
		return Template{}, errors.New("Must not specify organizationId")
	}
	organizationId, err := self.organizationId()
	if err != nil {
		return Template{}, err
	}
	payload.OrganizationId = organizationId

	var result Template
	err = self.http.Put("/blueprints/"+id, payload, &result)
	if err != nil {
		return Template{}, err
	}
	return result, nil
}

func (self *ApiClient) Templates() ([]Template, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return nil, err
	}
	var result []Template
	err = self.http.Get("/blueprints", map[string]string{"organizationId": organizationId}, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}


func (self *ApiClient) AssignTemplateToProject(id string, payload TemplateAssignmentToProjectPayload) (Template, error) {
	if payload.ProjectId == "" {
		return Template{}, errors.New("Must specify projectId on assignment to a template")
	}
	
	var result Template
	err := self.http.Patch("/blueprints/"+id+"/projects", payload, &result)
	if err != nil {
		return Template{}, err
	}
	return result, nil
}