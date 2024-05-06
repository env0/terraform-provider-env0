resource "random_string" "random" {
  length    = 8
  special   = false
  min_lower = 8
}

resource "env0_aws_credentials" "aws_cred1" {
  name = "Test Role arn1 ${random_string.random.result}"
  arn  = "Role ARN1"
}

resource "env0_aws_credentials" "aws_cred2" {
  name = "Test Role arn2 ${random_string.random.result}"
  arn  = "Role ARN2"
}

resource "env0_gcp_credentials" "gcp_cred" {
  name                = "name ${random_string.random.result}"
  service_account_key = "your_account_key"
  project_id          = "your_project_id"
}

data "env0_cloud_credentials" "all_aws_credentials" {
  depends_on      = [env0_aws_credentials.aws_cred1, env0_aws_credentials.aws_cred2, env0_gcp_credentials.gcp_cred]
  credential_type = "AWS_ASSUMED_ROLE_FOR_DEPLOYMENT"
}

data "env0_aws_credentials" "aws_credentials1" {
  name = data.env0_cloud_credentials.all_aws_credentials.names[index(data.env0_cloud_credentials.all_aws_credentials.names, env0_aws_credentials.aws_cred1.name)]
}

data "env0_aws_credentials" "aws_credentials2" {
  name = data.env0_cloud_credentials.all_aws_credentials.names[index(data.env0_cloud_credentials.all_aws_credentials.names, env0_aws_credentials.aws_cred2.name)]
}

output "credentials_name" {
  value = var.second_run ? replace(data.env0_aws_credentials.aws_credentials1.name, random_string.random.result, "") : ""
}

resource "env0_kubeconfig_credentials" "kubeconfig_credentials" {
  name        = "kubeconfig-${random_string.random.result}"
  kube_config = <<EOT
    apiVersion: v1
    clusters:
    - cluster:
        certificate-authority-data: <ca-data-here>
        server: https://your-k8s-cluster.com
      name: <cluster-name>
    contexts:
    - context:
        cluster:  <cluster-name>
        user:  <cluster-name-user>
      name:  <cluster-name>
    current-context:  <cluster-name>
    kind: Config
    preferences: {}
    users:
    - name:  <cluster-name-user>
      user:
        token: <secret-token-here>
    EOT
}

resource "env0_aws_eks_credentials" "aws_eks_credentials" {
  name           = "aws-eks-${random_string.random.result}"
  cluster_name   = "my-cluster"
  cluster_region = "us-east-2"
}

resource "env0_azure_aks_credentials" "azure_aks_credentials" {
  name           = "azure-aks-${random_string.random.result}"
  cluster_name   = "my-cluster"
  resource_group = "rg1"
}
