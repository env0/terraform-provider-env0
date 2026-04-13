# Create agent pool and install agent using helm chart

variable "agent_custom_configuration" {
  type        = map(any)
  description = "See https://docs.envzero.com/guides/admin-guide/self-hosted-kubernetes-agent/custom-optional-configuration"
}

resource "env0_agent_pool" "default" {
  name = "default sefl-hosted agent"
}

resource "env0_agent_secret" "first" {
  agent_id = env0_agent_pool.default.id
}

resource "helm_release" "this" {
  repository = "https://env0.github.io/self-hosted"
  chart      = "env0-agent"

  namespace = "env0-agent-default"
  name      = "env0-agent-default"

  create_namespace = true
  timeout          = 600

  set_sensitive {
    name  = "agentAccessToken"
    value = env0_agent_secret.first.secret
  }

  values = [
    yamlencode(var.agent_custom_configuration)
  ]
}
