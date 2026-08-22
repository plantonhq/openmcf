output "certificate_pack_id" {
  description = "The certificate pack identifier"
  value       = cloudflare_certificate_pack.main.id
}

output "status" {
  description = "The order/issuance status"
  value       = cloudflare_certificate_pack.main.status
}

output "primary_certificate" {
  description = "The identifier of the primary certificate in the pack"
  value       = cloudflare_certificate_pack.main.primary_certificate
}

output "zone_id" {
  description = "The Cloudflare zone the pack was ordered in (a pack's API identity is zone_id + certificate_pack_id)"
  value       = cloudflare_certificate_pack.main.zone_id
}
