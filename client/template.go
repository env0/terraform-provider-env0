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

	result, err := self.http.Post("/blueprints", payload)
	if err != nil {
		return Template{}, err
	}
	return result.(Template), nil
}

func (self *ApiClient) Template(id string) (Template, error) {
	result, err := self.http.Get("/blueprints/"+id, nil)
	if err != nil {
		return Template{}, err
	}
	return result.(Template), nil
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

	result, err := self.http.Put("/blueprints/"+id, payload)
	if err != nil {
		return Template{}, err
	}
	return result.(Template), nil
}

func (self *ApiClient) Templates() ([]Template, error) {
	organizationId, err := self.organizationId()
	if err != nil {
		return nil, err
	}

	result, err := self.http.Get("/blueprints", map[string]string{"organizationId": organizationId})
	if err != nil {
		return nil, err
	}

	templates := result.([]Template)
	return templates, err
}
