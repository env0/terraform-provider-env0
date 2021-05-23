package client

//templates are actually called "blueprints" in some parts of the API, this layer
//attempts to abstract this detail away - all the users of api client should
//only use "template", no mention of blueprint

import (
	"errors"
	"net/url"
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
	err = self.postJSON("/blueprints", payload, &result)
	if err != nil {
		return Template{}, err
	}
	return result, nil
}

func (self *ApiClient) Template(id string) (Template, error) {
	var result Template
	err := self.getJSON("/blueprints/"+id, nil, &result)
	if err != nil {
		return Template{}, err
	}
	return result, nil
}

func (self *ApiClient) TemplateDelete(id string) error {
	return self.delete("/blueprints/" + id)
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
	err = self.putJSON("/blueprints/"+id, payload, &result)
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
	params := url.Values{}
	params.Add("organizationId", organizationId)
	err = self.getJSON("/blueprints", params, &result)
	if err != nil {
		return nil, err
	}
	return result, err
}
