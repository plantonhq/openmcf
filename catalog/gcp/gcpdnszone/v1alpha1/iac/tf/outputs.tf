output "zone_id" {
  description = "The managed zone ID (numeric identifier)"
  value       = google_dns_managed_zone.managed_zone.managed_zone_id
}

output "zone_name" {
  description = "The name of the created managed zone (used by GcpDnsRecord managed_zone FK)"
  value       = google_dns_managed_zone.managed_zone.name
}

output "nameservers" {
  description = "The list of nameservers assigned to this managed zone. Configure these at your domain registrar for public zones."
  value       = google_dns_managed_zone.managed_zone.name_servers
}
