output "certificate_id" {
  description = "The ID of the uploaded certificate -- what consumers reference"
  value       = cloudflare_mtls_certificate.main.id
}

output "expires_on" {
  description = "When the certificate expires (RFC3339)"
  value       = cloudflare_mtls_certificate.main.expires_on
}

output "serial_number" {
  description = "The certificate's serial number"
  value       = cloudflare_mtls_certificate.main.serial_number
}
