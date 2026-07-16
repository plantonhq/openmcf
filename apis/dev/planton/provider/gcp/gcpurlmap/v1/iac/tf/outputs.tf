# The self-link — the value target HTTP(S) proxies reference as their url_map;
# the composition handle that puts this routing brain behind a load-balancer
# frontend.
output "self_link" {
  description = "Self-link URI of the URL map"
  value       = google_compute_url_map.this.self_link
}

# The name as it exists in GCP.
output "url_map_name" {
  description = "Name of the URL map in GCP"
  value       = google_compute_url_map.this.name
}

# Server-assigned numeric ID of the URL map.
output "map_id" {
  description = "Server-assigned numeric ID of the URL map"
  value       = google_compute_url_map.this.map_id
}

# Server-computed fingerprint for optimistic concurrency control.
output "fingerprint" {
  description = "Fingerprint of the URL map"
  value       = google_compute_url_map.this.fingerprint
}
