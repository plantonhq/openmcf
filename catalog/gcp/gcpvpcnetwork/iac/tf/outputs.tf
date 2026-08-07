output "network_self_link" {
  description = "The full self-link URL of the created VPC network"
  value       = google_compute_network.vpc.self_link
}

output "network_name" {
  description = "Name of the VPC network"
  value       = google_compute_network.vpc.name
}

output "network_id" {
  description = "Self-link identifier of the VPC network"
  value       = google_compute_network.vpc.id
}

output "gateway_ipv4" {
  description = "IPv4 address of the default internet gateway for this network"
  value       = google_compute_network.vpc.gateway_ipv4
}

output "internal_ipv6_range" {
  description = "ULA internal IPv6 range assigned when ULA IPv6 is enabled"
  value       = google_compute_network.vpc.internal_ipv6_range
}
