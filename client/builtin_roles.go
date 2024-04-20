package client

const (
	Admin    string = "Admin"
	Deployer string = "Deployer"
	Planner  string = "Planner"
	Viewer   string = "Viewer"
)

func IsBuiltinRole(role string) bool {
	return role == Admin || role == Deployer || role == Planner || role == Viewer
}
