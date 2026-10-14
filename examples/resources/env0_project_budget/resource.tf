resource "env0_project" "project" {
  name = "example"
}

resource "env0_project_budget" "project_budget" {
  project_id = env0_project.project.id
  amount     = 1000
  timeframe  = "MONTHLY"
  thresholds = [50, 80, 100] # percentages of amount
}
