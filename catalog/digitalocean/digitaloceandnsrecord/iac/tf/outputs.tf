output "record_id" {
  description = "The unique numeric identifier of the created DNS record (string form)"
  value       = digitalocean_record.dns_record.id
}

output "hostname" {
  description = "The fully-qualified hostname of the record (the provider's computed fqdn)"
  value       = digitalocean_record.dns_record.fqdn
}

output "record_type" {
  description = "The DNS record type that was created"
  value       = digitalocean_record.dns_record.type
}

output "domain" {
  description = "The domain name (DNS zone) the record was created in"
  value       = digitalocean_record.dns_record.domain
}

output "ttl_seconds" {
  description = "The TTL applied to the record in seconds (the API default when the spec left it unset)"
  value       = digitalocean_record.dns_record.ttl
}
