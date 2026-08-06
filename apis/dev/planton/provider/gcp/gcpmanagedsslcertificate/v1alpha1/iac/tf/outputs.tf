output "self_link" {
  description = "Self-link URI of the SSL certificate — the value a target HTTPS proxy references in ssl_certificates."
  value       = google_compute_managed_ssl_certificate.this.self_link
}

output "certificate_name" {
  description = "Name of the SSL certificate as it exists in GCP."
  value       = google_compute_managed_ssl_certificate.this.name
}

output "certificate_id" {
  description = "Server-assigned numeric ID of the certificate."
  value       = tostring(google_compute_managed_ssl_certificate.this.certificate_id)
}

output "expire_time" {
  description = "Expiry time of the certificate in RFC3339 format. Empty until provisioning completes."
  value       = google_compute_managed_ssl_certificate.this.expire_time
}
