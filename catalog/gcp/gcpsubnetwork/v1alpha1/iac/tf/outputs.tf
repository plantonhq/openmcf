# The self-link — the value GKE clusters, compute instances, and other
# subnet consumers reference.
output "subnetwork_self_link" {
  description = "Self-link URI of the subnetwork"
  value       = google_compute_subnetwork.main.self_link
}

# The region the subnetwork lives in.
output "region" {
  description = "The region the subnetwork lives in"
  value       = google_compute_subnetwork.main.region
}

# The primary IPv4 CIDR range (empty for IPV6_ONLY subnets).
output "ip_cidr_range" {
  description = "The primary IPv4 CIDR range of the subnet"
  value       = google_compute_subnetwork.main.ip_cidr_range
}

# Secondary (alias) ranges — GKE clusters select pod/service ranges by name.
output "secondary_ranges" {
  description = "Secondary ranges with their names and CIDRs"
  value = [
    for secondary_range in google_compute_subnetwork.main.secondary_ip_range : {
      range_name    = secondary_range.range_name
      ip_cidr_range = secondary_range.ip_cidr_range
    }
  ]
}

# The name as it exists in GCP — referenced by consumers that address
# subnets by name (e.g. Cloud Run Direct VPC egress).
output "subnetwork_name" {
  description = "Name of the subnetwork"
  value       = google_compute_subnetwork.main.name
}

# IPv4 address of the subnet's default gateway.
output "gateway_address" {
  description = "IPv4 address of the subnet's default gateway"
  value       = google_compute_subnetwork.main.gateway_address
}

# Server-assigned numeric ID, exported as a string for a stable cross-engine
# output shape.
output "subnetwork_id" {
  description = "Server-assigned numeric ID of the subnetwork"
  value       = tostring(google_compute_subnetwork.main.subnetwork_id)
}

# IPv6 prefixes actually allocated (empty when the stack type has no IPv6).
output "internal_ipv6_prefix" {
  description = "Internal IPv6 prefix allocated to the subnet"
  value       = google_compute_subnetwork.main.internal_ipv6_prefix
}

output "external_ipv6_prefix" {
  description = "External IPv6 prefix allocated to the subnet"
  value       = google_compute_subnetwork.main.external_ipv6_prefix
}
