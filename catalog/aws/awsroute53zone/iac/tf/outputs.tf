output "zone_id" {
  description = "The hosted zone ID — the join key every DNS-composing resource references."
  value       = aws_route53_zone.this.zone_id
}

output "zone_name" {
  description = "The zone's domain name as normalized by Route 53 (trailing dot removed)."
  value       = aws_route53_zone.this.name
}

output "nameservers" {
  description = "The four authoritative name servers assigned to the zone (registrar NS delegation values for public zones)."
  value       = aws_route53_zone.this.name_servers
}

output "primary_name_server" {
  description = "The first (primary) name server of the zone's delegation set (the SOA MNAME)."
  value       = aws_route53_zone.this.primary_name_server
}

output "zone_arn" {
  description = "The ARN of the hosted zone, for IAM policies scoping route53:ChangeResourceRecordSets."
  value       = aws_route53_zone.this.arn
}

output "ds_record" {
  description = "The key-signing key's DS record — register it with the parent zone (the registrar) to complete the DNSSEC chain of trust. Empty when DNSSEC is not enabled."
  value       = local.dnssec_enabled ? aws_route53_key_signing_key.this[0].ds_record : ""
}

output "dnskey_record" {
  description = "The key-signing key's DNSKEY record (the public key in DNS record form). Empty when DNSSEC is not enabled."
  value       = local.dnssec_enabled ? aws_route53_key_signing_key.this[0].dnskey_record : ""
}

output "key_signing_key_tag" {
  description = "The key-signing key's key tag — the short identifier registrars display next to a DS record. Empty when DNSSEC is not enabled."
  value       = local.dnssec_enabled ? tostring(aws_route53_key_signing_key.this[0].key_tag) : ""
}
