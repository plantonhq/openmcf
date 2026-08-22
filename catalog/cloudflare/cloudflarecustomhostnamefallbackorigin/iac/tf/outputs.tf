output "status" {
  description = "The status of the fallback origin"
  value       = cloudflare_custom_hostname_fallback_origin.main.status
}

output "created_at" {
  description = "RFC3339 timestamp of when the fallback origin was created"
  value       = cloudflare_custom_hostname_fallback_origin.main.created_at
}

output "updated_at" {
  description = "RFC3339 timestamp of when the fallback origin was last updated"
  value       = cloudflare_custom_hostname_fallback_origin.main.updated_at
}

output "errors" {
  description = "Any errors reported while deploying the fallback origin"
  value       = try(cloudflare_custom_hostname_fallback_origin.main.errors, [])
}

output "zone_id" {
  description = "The Cloudflare zone this singleton belongs to (the fallback origin has no resource id; its API identity IS the zone)"
  value       = cloudflare_custom_hostname_fallback_origin.main.zone_id
}
