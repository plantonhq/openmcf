output "load_balancer_id" {
  description = "The ID of the Cloudflare Load Balancer"
  value       = cloudflare_load_balancer.main.id
}

output "load_balancer_dns_record_name" {
  description = "The DNS record name (hostname) of the load balancer"
  value       = cloudflare_load_balancer.main.name
}

output "load_balancer_cname_target" {
  description = "The CNAME target returned by Cloudflare for the load balancer"
  value       = cloudflare_load_balancer.main.id
}

output "pool_id" {
  description = "The ID of the load balancer pool"
  value       = cloudflare_load_balancer_pool.main.id
}

output "monitor_id" {
  description = "The ID of the health check monitor"
  value       = cloudflare_load_balancer_monitor.health_check.id
}

