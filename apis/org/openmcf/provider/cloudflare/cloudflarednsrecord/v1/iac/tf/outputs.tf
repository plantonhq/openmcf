# outputs.tf

output "record_id" {
  description = "The unique identifier of the created DNS record"
  value       = cloudflare_record.main.id
}

output "hostname" {
  description = "The fully qualified hostname of the DNS record"
  value       = cloudflare_record.main.hostname
}

output "record_type" {
  description = "The DNS record type that was created"
  value       = cloudflare_record.main.type
}

output "proxied" {
  description = "Whether the record is proxied through Cloudflare"
  value       = cloudflare_record.main.proxied
}
