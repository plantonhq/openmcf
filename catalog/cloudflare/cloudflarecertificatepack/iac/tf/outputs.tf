output "certificate_pack_id" {
  description = "The certificate pack identifier"
  value       = cloudflare_certificate_pack.main.id
}

# No status output: pack issuance is asynchronous (initializing ->
# pending_validation -> active), and a point-in-time phase is never a stable
# stack output -- it flips on the first refresh after the transition and
# re-plans forever. Read issuance status from the Cloudflare API instead.

output "primary_certificate" {
  description = "The identifier of the primary certificate in the pack"
  value       = cloudflare_certificate_pack.main.primary_certificate
}

output "zone_id" {
  description = "The Cloudflare zone the pack was ordered in (a pack's API identity is zone_id + certificate_pack_id)"
  value       = cloudflare_certificate_pack.main.zone_id
}
