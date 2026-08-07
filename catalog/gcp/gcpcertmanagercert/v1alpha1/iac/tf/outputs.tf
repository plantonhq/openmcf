output "certificate_id" {
  description = "Fully-qualified resource ID (projects/{project}/locations/{location}/certificates/{name})"
  value       = google_certificate_manager_certificate.certificate.id
}

output "certificate_name" {
  description = "Name of the certificate as it exists in GCP — the value a target HTTPS proxy's certificate_manager_certificates list consumes"
  value       = google_certificate_manager_certificate.certificate.name
}

output "san_dnsnames" {
  description = "Subject Alternative Names (dnsName type) in the issued certificate"
  value       = google_certificate_manager_certificate.certificate.san_dnsnames
}

output "location" {
  description = "The Certificate Manager location the certificate lives in"
  value       = var.spec.location != "" ? var.spec.location : "global"
}

output "managed_state" {
  description = "State of a managed certificate (PROVISIONING/FAILED/ACTIVE); empty for self-managed"
  # try() guards the self-managed arm, where the managed block is absent.
  value = try(google_certificate_manager_certificate.certificate.managed[0].state, "")
}
