# Example 1: Project-level Permissions
# Using project_permissions block within the API key resource

# Example: Admin API Key (default)
resource "env0_api_key" "admin_key" {
  name = "admin-api-key"
}

# Example: User API Key with project permissions
resource "env0_project" "project" {
  name = "demo-project"
}

resource "env0_api_key" "user_key_with_project" {
  name              = "user-api-key-project"
  organization_role = "User"

  project_permissions {
    project_id   = env0_project.project.id
    project_role = "Deployer"
  }
}

# Example 2: Team Assignment
# Assign API key to teams for team-based access control

resource "env0_api_key" "team_api_key" {
  name = "team-api-key"
}

resource "env0_team" "team_resource" {
  name = "team-resource"
}

resource "env0_user_team_assignment" "api_key_team_assignment" {
  user_id = env0_api_key.team_api_key.id
  team_id = env0_team.team_resource.id
}

# Example 3: Combined Approach
# Using both project permissions and team assignment

resource "env0_api_key" "combined_key" {
  name              = "combined-access-key"
  organization_role = "User"

  project_permissions {
    project_id   = env0_project.project.id
    project_role = "Viewer"
  }
}

resource "env0_user_team_assignment" "combined_team_assignment" {
  user_id = env0_api_key.combined_key.id
  team_id = env0_team.team_resource.id
}

# Example 4: API Key with Custom Role
# Create a custom role and assign it to an API key

resource "env0_custom_role" "deployer_role" {
  name = "custom-deployer"
  permissions = [
    "VIEW_ORGANIZATION",
    "VIEW_PROJECT",
    "VIEW_ENVIRONMENT",
    "RUN_PLAN",
    "RUN_APPLY"
  ]
}

resource "env0_api_key" "custom_role_key" {
  name              = "custom-role-api-key"
  organization_role = env0_custom_role.deployer_role.id

  project_permissions {
    project_id   = env0_project.project.id
    project_role = "Deployer"
  }
}
