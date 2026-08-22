output "snippet_name" {
  description = "The snippet's name -- its identity in the zone, referenced by snippet rules"
  value       = cloudflare_snippet.main.snippet_name
}

output "zone_id" {
  description = "The zone the snippet is deployed to"
  value       = var.spec.zone_id
}
