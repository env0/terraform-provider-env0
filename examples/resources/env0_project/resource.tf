resource "env0_project" "example" {
  name        = "example"
  description = "Example project"
}

# Wait for the project's environments to be destroyed or archived before deleting it.
# The wait is bounded by the 'delete' timeout of the 'timeouts' block,
# defaulting to 10 minutes.
resource "env0_project" "example_with_destroy_wait" {
  name        = "project with a bounded destroy wait"
  description = "Example project"

  wait = true

  timeouts {
    delete = "20m"
  }
}
