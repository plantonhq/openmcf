output "script_id" {
  description = "The Cloudflare-assigned ID of the deployed Worker script"
  value       = cloudflare_workers_script.main.id
}

output "script_name" {
  description = "The Worker script name (the target a service binding references)"
  value       = cloudflare_workers_script.main.script_name
}

output "custom_domain_hostnames" {
  description = "The custom-domain hostnames attached to this Worker"
  value       = [for cd in cloudflare_workers_custom_domain.main : cd.hostname]
}

output "route_patterns" {
  description = "The route patterns mapped to this Worker"
  value       = [for r in cloudflare_workers_route.main : r.pattern]
}

# Keyed by hostname (the for_each key) so import can address
# cloudflare_workers_custom_domain as {account_id}/{domain_id}.
output "custom_domain_ids" {
  description = "Cloudflare-assigned custom-domain ids, keyed by hostname"
  value       = { for k, cd in cloudflare_workers_custom_domain.main : k => cd.id }
}

# Keyed by the same index string the module uses for for_each so import can
# address cloudflare_workers_route as {zone_id}/{route_id}.
output "route_ids" {
  description = "Cloudflare-assigned route ids, keyed by list index"
  value       = { for k, r in cloudflare_workers_route.main : k => r.id }
}

output "route_zone_ids" {
  description = "Zone id of each route, keyed the same way as route_ids"
  value       = { for k, r in cloudflare_workers_route.main : k => r.zone_id }
}
