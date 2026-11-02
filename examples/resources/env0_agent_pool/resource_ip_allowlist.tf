# Restricting the IP addresses an agent may be used from
# Applies on top of the organization's IP allowlist. Terraform owns the list, so one set in the env0 UI is cleared on the next apply.

resource "env0_agent_pool" "ip_restricted_agent" {
  name = "ip restricted agent"
  ip_allowlist = [
    "198.51.100.7",
    "198.51.100.0/28",
    "198.51.100.40-198.51.100.60",
    "2001:db8::/32",
  ]
}
