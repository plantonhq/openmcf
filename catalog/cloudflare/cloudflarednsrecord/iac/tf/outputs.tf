# outputs.tf

output "record_id" {
  description = "The unique identifier of the created DNS record"
  value       = cloudflare_dns_record.main.id
}

output "record_name" {
  description = "The DNS record name as stored by Cloudflare"
  value       = cloudflare_dns_record.main.name
}

output "record_type" {
  description = "The DNS record type that was created"
  value       = cloudflare_dns_record.main.type
}

output "proxied" {
  description = "Whether the record is proxied through Cloudflare"
  value       = cloudflare_dns_record.main.proxied
}

output "zone_id" {
  description = "The Cloudflare zone ID the record lives in (a record's API identity is zone_id + record_id)"
  value       = cloudflare_dns_record.main.zone_id
}
