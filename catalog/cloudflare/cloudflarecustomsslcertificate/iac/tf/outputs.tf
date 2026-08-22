output "certificate_id" {
  description = "The ID of the uploaded certificate"
  value       = cloudflare_custom_ssl.main.id
}

output "zone_id" {
  description = "The zone the certificate belongs to"
  value       = var.spec.zone_id
}

output "expires_on" {
  description = "When the certificate expires (RFC3339)"
  value       = cloudflare_custom_ssl.main.expires_on
}

output "status" {
  description = "The certificate's deployment status (deployment is asynchronous)"
  value       = cloudflare_custom_ssl.main.status
}
