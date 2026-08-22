output "location_id" {
  description = "The Cloudflare-assigned UUID of the location"
  value       = cloudflare_zero_trust_dns_location.main.id
}

output "doh_subdomain" {
  description = "The location's unique DNS-over-HTTPS subdomain"
  value       = cloudflare_zero_trust_dns_location.main.doh_subdomain
}

output "ip" {
  description = "The IPv4 destination assigned to the location's plain-DNS endpoint"
  value       = cloudflare_zero_trust_dns_location.main.ip
}
