# outputs.tf

output "zone_id" {
  description = "The unique identifier of the created Cloudflare zone"
  value       = cloudflare_zone.main.id
}

output "nameservers" {
  description = "The Cloudflare nameservers assigned to this zone"
  value       = cloudflare_zone.main.name_servers
}

output "zone_name" {
  description = "The zone name (same as input)"
  value       = cloudflare_zone.main.zone
}

output "status" {
  description = "The current status of the zone"
  value       = cloudflare_zone.main.status
}

