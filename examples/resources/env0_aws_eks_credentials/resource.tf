resource "env0_aws_eks_credentials" "credentials" {
  name           = "example"
  cluster_name   = "my-cluster"
  cluster_region = "us-east-2"
}

data "env0_project" "project" {
  name = "my-project"
}

resource "env0_cloud_credentials_project_assignment" "assignment" {
  credential_id = env0_aws_eks_credentials.credentials.id
  project_id    = data.env0_project.project.id
}
