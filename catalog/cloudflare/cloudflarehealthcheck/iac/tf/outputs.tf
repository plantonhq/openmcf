output "healthcheck_id" {
  description = "The ID of the created health check"
  value       = cloudflare_healthcheck.main.id
}

output "zone_id" {
  description = "The zone the health check belongs to"
  value       = var.spec.zone_id
}
