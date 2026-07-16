output "address" {
  description = "The reserved IP address or start of the reserved range"
  value       = google_compute_address.this.address
}

output "self_link" {
  description = "Self-link URL of the regional address resource"
  value       = google_compute_address.this.self_link
}

output "name" {
  description = "Name of the regional address resource in GCP"
  value       = google_compute_address.this.name
}

# Export the plain region NAME from the spec (matching the Pulumi module) —
# the provider's region attribute can carry a self-link, which API callers
# and verification cannot use directly.
output "region" {
  description = "Region of the address reservation"
  value       = var.spec.region
}
