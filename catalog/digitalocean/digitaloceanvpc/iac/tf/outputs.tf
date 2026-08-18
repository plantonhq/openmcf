output "vpc_id" {
  description = "The unique identifier (UUID) of the created DigitalOcean VPC."
  value       = digitalocean_vpc.vpc.id
}

output "ip_range" {
  description = "The VPC's IP range in CIDR notation, as reported by DigitalOcean (covers the auto-assigned case)."
  value       = digitalocean_vpc.vpc.ip_range
}

output "urn" {
  description = "The uniform resource name (URN) of the VPC."
  value       = digitalocean_vpc.vpc.urn
}
