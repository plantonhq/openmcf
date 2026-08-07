# Exactly one of the two resources exists (see the count guards in main.tf),
# so each output selects whichever branch was created with one(concat(...)).

# The self-link — the value backend services reference in their
# health_checks list; the composition handle for the load balancing family.
output "self_link" {
  description = "Self-link URI of the health check"
  value       = one(concat(google_compute_health_check.this[*].self_link, google_compute_region_health_check.this[*].self_link))
}

# The name as it exists in GCP.
output "health_check_name" {
  description = "Name of the health check"
  value       = one(concat(google_compute_health_check.this[*].name, google_compute_region_health_check.this[*].name))
}

# The probe protocol GCP computed from the configured block
# (HTTP, HTTPS, HTTP2, TCP, SSL, GRPC, or GRPC_TLS).
output "type" {
  description = "The probe protocol of the health check"
  value       = one(concat(google_compute_health_check.this[*].type, google_compute_region_health_check.this[*].type))
}

# Region of a regional health check; empty for a global one, so downstream
# composition can confirm scope compatibility.
output "region" {
  description = "Region of the health check (empty for global)"
  value       = local.is_regional ? var.spec.region : ""
}
