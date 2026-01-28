output "script_id" {
  description = "The ID of the Cloudflare Worker script"
  value       = cloudflare_workers_script.main.id
}

output "script_name" {
  description = "The name of the Worker script"
  value       = cloudflare_workers_script.main.name
}

output "route_urls" {
  description = "List of route URLs where the Worker is accessible"
  value = local.dns_enabled ? [
    "https://${local.dns_hostname}"
  ] : []
}

output "route_pattern" {
  description = "The route pattern for the Worker"
  value       = local.dns_enabled ? local.route_pattern : ""
}

