resource "env0_aws_eks_credentials" "example" {
  name           = "example"
  cluster_name   = "my-cluster"
  cluster_region = "us-east-2"
}

data "env0_aws_eks_credentials" "by_id" {
  id = env0_aws_eks_credentials.example.id
}

data "env0_aws_eks_credentials" "by_name" {
  name = env0_aws_eks_credentials.example.name
}
