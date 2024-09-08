package env0

import (
	"context"
	"fmt"
	"strings"

	"github.com/env0/terraform-provider-env0/client"
	"github.com/env0/terraform-provider-env0/client/http"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataProjectRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "the name of the project",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
				Computed:     true,
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "id of the project",
				Optional:     true,
				ExactlyOneOf: []string{"name", "id"},
				Computed:     true,
			},
			"parent_project_name": {
				Type:          schema.TypeString,
				Description:   "the name of the parent project. Can be used as a filter when there are multiple subprojects with the same name under different parent projects",
				Optional:      true,
				ConflictsWith: []string{"parent_project_path"},
			},
			"parent_project_path": {
				Type:          schema.TypeString,
				Description:   "a path of ancestors projects divided by the prefix '|'. Can be used as a filter when there are multiple subprojects with the same name under different parent projects. For example: 'App|Dev|us-east-1' will search for a project with the hierarchy 'App -> Dev -> us-east-1' ('us-east-1' being the parent)",
				Optional:      true,
				ConflictsWith: []string{"parent_project_name"},
			},
			"parent_project_id": {
				Type:        schema.TypeString,
				Description: "the id of the parent project. Can be used as a filter when there are multiple subprojects with the same name under different parent projects",
				Optional:    true,
				Computed:    true,
			},
			"created_by": {
				Type:        schema.TypeString,
				Description: "textual description of the entity who created the project",
				Computed:    true,
			},
			"role": {
				Type:        schema.TypeString,
				Description: "role of the authenticated user (through api key) in the project",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "textual description of the project",
				Computed:    true,
			},
			"hierarchy": {
				Type:        schema.TypeString,
				Description: "the hierarchy of the project",
				Computed:    true,
			},
		},
	}
}

func dataProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var err error
	var project client.Project

	id, ok := d.GetOk("id")
	if ok {
		project, err = getProjectById(id.(string), meta)
		if err != nil {
			return diag.Errorf("%v", err)
		}
	} else {
		name, ok := d.GetOk("name")
		if !ok {
			return diag.Errorf("either 'name' or 'id' must be specified")
		}

		project, err = getProjectByName(name.(string), d.Get("parent_project_id").(string), d.Get("parent_project_name").(string), d.Get("parent_project_path").(string), meta)
		if err != nil {
			return diag.Errorf("%v", err)
		}
	}

	if err := writeResourceData(&project, d); err != nil {
		return diag.Errorf("schema resource data serialization failed: %v", err)
	}

	return nil
}

func filterByParentProjectId(parentId string, projects []client.Project) []client.Project {
	filteredProjects := make([]client.Project, 0)

	for _, project := range projects {
		if len(project.ParentProjectId) == 0 {
			continue
		}

		if project.ParentProjectId == parentId {
			filteredProjects = append(filteredProjects, project)
		}
	}

	return filteredProjects
}

func filterByParentProjectName(parentName string, projects []client.Project, meta interface{}) ([]client.Project, error) {
	filteredProjects := make([]client.Project, 0)

	for _, project := range projects {
		if len(project.ParentProjectId) == 0 {
			continue
		}

		parentProject, err := getProjectById(project.ParentProjectId, meta)
		if err != nil {
			return nil, err
		}

		if parentProject.Name == parentName {
			filteredProjects = append(filteredProjects, project)
		}
	}

	return filteredProjects, nil
}

func getProjectByName(name string, parentId string, parentName string, parentPath string, meta interface{}) (client.Project, error) {
	apiClient := meta.(client.ApiClientInterface)

	projects, err := apiClient.Projects()
	if err != nil {
		return client.Project{}, fmt.Errorf("could not query project by name: %w", err)
	}

	projectsByName := make([]client.Project, 0)

	for _, candidate := range projects {
		if candidate.Name == name && !candidate.IsArchived {
			projectsByName = append(projectsByName, candidate)
		}
	}

	// Use filters to reduce results.

	switch {
	case len(parentId) > 0:
		projectsByName = filterByParentProjectId(parentId, projectsByName)
	case len(parentName) > 0:
		projectsByName, err = filterByParentProjectName(parentName, projectsByName, meta)
		if err != nil {
			return client.Project{}, err
		}
	case len(parentPath) > 0:
		projectsByName, err = filterByParentProjectPath(parentPath, projectsByName, meta)
		if err != nil {
			return client.Project{}, err
		}
	}

	if len(projectsByName) > 1 {
		return client.Project{}, fmt.Errorf("found multiple projects for name: %s. Use id or one of the filters to make sure only one '%v' is returned", name, projectsByName)
	}

	if len(projectsByName) == 0 {
		return client.Project{}, fmt.Errorf("could not find a project with name: %s", name)
	}

	return projectsByName[0], nil
}

func pathMatches(path, parentIds []string, meta interface{}) (bool, error) {
	if len(path) != len(parentIds) {
		return false, nil
	}

	apiClient := meta.(client.ApiClientInterface)

	for i := range path {
		parentId := parentIds[i]

		parentProject, err := apiClient.Project(parentId)
		if err != nil {
			return false, fmt.Errorf("failed to get a parent project with id '%s': %w", parentId, err)
		}

		if parentProject.Name != path[i] {
			return false, nil
		}
	}

	return true, nil
}

func filterByParentProjectPath(parentPath string, projectsByName []client.Project, meta interface{}) ([]client.Project, error) {
	filteredProjects := make([]client.Project, 0)

	path := strings.Split(parentPath, "|")

	for _, project := range projectsByName {
		parentIds := strings.Split(project.Hierarchy, "|")
		// right most element is the project itself, remove it.
		parentIds = parentIds[:len(parentIds)-1]

		matches, err := pathMatches(path, parentIds, meta)

		if err != nil {
			return nil, err
		}

		if matches {
			filteredProjects = append(filteredProjects, project)
		}
	}

	return filteredProjects, nil
}

func getProjectById(id string, meta interface{}) (client.Project, error) {
	apiClient := meta.(client.ApiClientInterface)

	project, err := apiClient.Project(id)
	if err != nil {
		if frerr, ok := err.(*http.FailedResponseError); ok && frerr.NotFound() {
			return client.Project{}, fmt.Errorf("could not find a project with id: %s", id)
		}

		return client.Project{}, fmt.Errorf("could not query project: %w", err)
	}

	return project, nil
}
