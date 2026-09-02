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

# status is deliberately NOT an output: deployment is asynchronous (pending
# before active), so a point-in-time phase flips on the first refresh after
# the transition and re-plans forever (the class was measured live
# 2026-08-28 on the sibling Authenticated Origin Pulls certificate).
