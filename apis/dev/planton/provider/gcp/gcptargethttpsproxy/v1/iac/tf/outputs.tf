# The self-link — the value a global forwarding rule references as its
# target; the composition handle that puts a VIP in front of this proxy.
output "self_link" {
  description = "Self-link URI of the target HTTPS proxy"
  value       = google_compute_target_https_proxy.this.self_link
}

# The name as it exists in GCP.
output "proxy_name" {
  description = "Name of the target HTTPS proxy in GCP"
  value       = google_compute_target_https_proxy.this.name
}

# Server-assigned numeric ID of the proxy.
output "proxy_id" {
  description = "Server-assigned numeric ID of the target HTTPS proxy"
  value       = google_compute_target_https_proxy.this.proxy_id
}

# Server-computed fingerprint for optimistic concurrency control.
output "fingerprint" {
  description = "Fingerprint of the target HTTPS proxy"
  value       = google_compute_target_https_proxy.this.fingerprint
}
