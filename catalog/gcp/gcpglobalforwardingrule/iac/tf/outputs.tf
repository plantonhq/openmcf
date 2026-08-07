# The VIP — the IP address DNS records point at. The API always reports the
# literal IP number here, even when the spec referenced an address resource.
output "ip_address" {
  description = "IP address the forwarding rule accepts traffic on (the load balancer VIP)"
  value       = google_compute_global_forwarding_rule.this.ip_address
}

# Self-link URI of the forwarding rule.
output "self_link" {
  description = "Self-link URI of the global forwarding rule"
  value       = google_compute_global_forwarding_rule.this.self_link
}

# The name as it exists in GCP.
output "forwarding_rule_name" {
  description = "Name of the global forwarding rule in GCP"
  value       = google_compute_global_forwarding_rule.this.name
}

# Server-assigned numeric ID of the forwarding rule.
output "forwarding_rule_id" {
  description = "Server-assigned numeric ID of the global forwarding rule"
  value       = google_compute_global_forwarding_rule.this.forwarding_rule_id
}

# PSC connection id — populated only for Private Service Connect frontends.
output "psc_connection_id" {
  description = "Private Service Connect connection id (PSC frontends only)"
  value       = google_compute_global_forwarding_rule.this.psc_connection_id
}

# PSC connection status — ACCEPTED means the producer admitted the
# connection. Populated only for Private Service Connect frontends.
output "psc_connection_status" {
  description = "Private Service Connect connection status (PSC frontends only)"
  value       = google_compute_global_forwarding_rule.this.psc_connection_status
}
