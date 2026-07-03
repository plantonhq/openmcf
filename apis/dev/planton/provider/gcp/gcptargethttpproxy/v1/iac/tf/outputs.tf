# The self-link — the value a global forwarding rule references as its
# target; the composition handle that puts a VIP in front of this proxy.
output "self_link" {
  description = "Self-link URI of the target HTTP proxy"
  value       = google_compute_target_http_proxy.this.self_link
}

# The name as it exists in GCP.
output "proxy_name" {
  description = "Name of the target HTTP proxy in GCP"
  value       = google_compute_target_http_proxy.this.name
}

# Server-assigned numeric ID of the proxy.
output "proxy_id" {
  description = "Server-assigned numeric ID of the target HTTP proxy"
  value       = google_compute_target_http_proxy.this.proxy_id
}

# Server-computed fingerprint for optimistic concurrency control.
output "fingerprint" {
  description = "Fingerprint of the target HTTP proxy"
  value       = google_compute_target_http_proxy.this.fingerprint
}
