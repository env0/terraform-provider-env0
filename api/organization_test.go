package api

import (
	"testing"
)

func TestOrganization(t *testing.T) {
	client, err := NewClientFromEnv()
	if err != nil {
		t.Error("Unable to init api client:", err)
		return
	}

	organization, err := client.Organization()
	if err != nil {
		t.Error("Unable to get organization:", err)
		return
	}
	if organization.IsSelfHosted {
		t.Error("Expected not self hosted")
	}
	if organization.Id == "" {
		t.Error("Expected non empty id")
	}
	if organization.Name == "" {
		t.Error("Expected non empty name")
	}
}
