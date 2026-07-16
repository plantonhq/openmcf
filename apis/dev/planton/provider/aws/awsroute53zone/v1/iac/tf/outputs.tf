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
