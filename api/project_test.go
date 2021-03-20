package api

import (
	"testing"
)

func TestProject(t *testing.T) {
	client, err := NewClientFromEnv()
	if err != nil {
		t.Error("Unable to init api client:", err)
		return
	}

	projects, err := client.Projects()
	if err != nil {
		t.Error("Unable to get projects:", err)
		return
	}
	if len(projects) == 0 {
		t.Error("Expected at least one project")
		return
	}
	var defaultProject Project
	for _, project := range projects {
		if project.Name == "Default Organization Project" {
			defaultProject = project
		}
	}
	if defaultProject.Name == "" {
		t.Error("Default project not found")
		return
	}
	defaultProject2, err := client.Project(defaultProject.Id)
	if err != nil {
		t.Error("Unable to fetch default project by id:", err)
		return
	}
	if defaultProject2.Name != "Default Organization Project" {
		t.Error("Fetching by id returned incorrect name:", defaultProject2.Name)
		return
	}
}