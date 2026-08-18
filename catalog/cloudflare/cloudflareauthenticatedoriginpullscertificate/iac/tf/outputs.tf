output "certificate_id" {
  description = "The ID of the uploaded certificate -- what per-hostname associations reference"
  value       = local.is_zone_scope ? cloudflare_authenticated_origin_pulls_certificate.zone[0].id : cloudflare_authenticated_origin_pulls_hostname_certificate.hostname[0].id
}

output "zone_id" {
  description = "The zone the certificate belongs to"
  value       = var.spec.zone_id
}

output "expires_on" {
  description = "When the certificate expires (RFC3339)"
  value       = local.is_zone_scope ? cloudflare_authenticated_origin_pulls_certificate.zone[0].expires_on : cloudflare_authenticated_origin_pulls_hostname_certificate.hostname[0].expires_on
}

output "status" {
  description = "The certificate's deployment status (deployment and deletion are asynchronous)"
  value       = local.is_zone_scope ? cloudflare_authenticated_origin_pulls_certificate.zone[0].status : cloudflare_authenticated_origin_pulls_hostname_certificate.hostname[0].status
}
