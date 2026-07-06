# Fully qualified policy resource path — the canonical identifier for
# API calls and audit tooling.
output "policy_id" {
  description = "Fully qualified resource path of the service connection policy"
  value       = google_network_connectivity_service_connection_policy.this.id
}

# Short policy name — what an operator sees in the console's service
# connection policy list for the network.
output "name" {
  description = "Short name of the service connection policy"
  value       = local.policy_name
}

# The underlying connectivity mechanism the automation uses for this
# policy (PSC) — confirms the policy is PSC-backed without inspecting
# individual connections.
output "infrastructure" {
  description = "Type of underlying resources the automation creates (PSC)"
  value       = google_network_connectivity_service_connection_policy.this.infrastructure
}

# Server-computed etag — changes on every policy mutation, useful for
# change detection when auditing shared networks.
output "etag" {
  description = "Server-computed etag of the policy"
  value       = google_network_connectivity_service_connection_policy.this.etag
}
