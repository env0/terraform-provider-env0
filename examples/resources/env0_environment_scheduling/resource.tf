data "env0_environment" "example" {
  name = "Environment Name"
}

resource "env0_environment_schedling" "example" {
  environment_id = data.env0_environment.example.id
  deploy_cron    = "5 * * * *"
  destroy_cron   = "10 * * * *"
}
