package client

const (
	AdminRole    string = "Admin"
	DeployerRole string = "Deployer"
	PlannerRole  string = "Planner"
	ViewerRole   string = "Viewer"
	UserRole     string = "User"
)

func IsBuiltinRole(role string) bool {
	return role == AdminRole || role == DeployerRole || role == PlannerRole || role == ViewerRole || role == UserRole
}

func IsCustomRole(role string) bool {
	return !IsBuiltinRole(role)
}

var ProjectBuiltinRoles = []string{AdminRole, DeployerRole, PlannerRole, ViewerRole}
var EnvironmentBuiltinRoles = []string{AdminRole, DeployerRole, PlannerRole, ViewerRole}
var OrganizationBuiltinRoles = []string{UserRole, AdminRole}
