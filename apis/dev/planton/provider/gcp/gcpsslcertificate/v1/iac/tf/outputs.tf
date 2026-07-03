# Exactly one of the two resources exists (see the count guards in main.tf),
# so each output selects whichever branch was created with one(concat(...)).
# The private key is write-only in GCP and deliberately never output.

# The self-link — the value a target HTTPS (or SSL) proxy references in its
# ssl_certificates list; the composition handle that terminates TLS at the
# load balancer.
output "self_link" {
  description = "Self-link URI of the SSL certificate"
  value       = one(concat(google_compute_ssl_certificate.this[*].self_link, google_compute_region_ssl_certificate.this[*].self_link))
}

# The name as it exists in GCP.
output "certificate_name" {
  description = "Name of the SSL certificate"
  value       = one(concat(google_compute_ssl_certificate.this[*].name, google_compute_region_ssl_certificate.this[*].name))
}

# Server-assigned numeric ID of the certificate.
output "certificate_id" {
  description = "Server-assigned numeric ID of the certificate"
  value       = tostring(one(concat(google_compute_ssl_certificate.this[*].certificate_id, google_compute_region_ssl_certificate.this[*].certificate_id)))
}

# Expiry parsed by GCP from the uploaded chain. Self-managed certificates do
# NOT renew themselves — plan the create-before-destroy rotation off this.
output "expire_time" {
  description = "Expiry time of the certificate in RFC3339 format"
  value       = one(concat(google_compute_ssl_certificate.this[*].expire_time, google_compute_region_ssl_certificate.this[*].expire_time))
}

# Region of a regional certificate; empty for a global one, so downstream
# composition can confirm scope compatibility.
output "region" {
  description = "Region of the SSL certificate (empty for global)"
  value       = local.is_regional ? var.spec.region : ""
}
