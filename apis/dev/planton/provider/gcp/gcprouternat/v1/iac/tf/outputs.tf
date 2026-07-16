output "name" {
  description = "Name of the Cloud NAT gateway"
  value       = google_compute_router_nat.nat.name
}

output "router_self_link" {
  description = "Self-link URL of the Cloud Router carrying this NAT"
  value       = google_compute_router.router.self_link
}

output "nat_ips" {
  description = "Self links of the static NAT IPs (manual allocation only; empty for auto-allocation)"
  value       = google_compute_router_nat.nat.nat_ips
}
