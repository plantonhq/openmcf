output "fqdn" {
  description = "Fully qualified domain name of the created DNS record."
  value       = aws_route53_record.this.fqdn
}

output "record_type" {
  description = "DNS record type (A, AAAA, CNAME, ...)."
  value       = aws_route53_record.this.type
}

output "zone_id" {
  description = "Route 53 hosted zone ID the record lives in."
  value       = var.spec.zone_id
}

output "is_alias" {
  description = "Whether this is an alias record."
  value       = local.is_alias
}

output "set_identifier" {
  description = "Set identifier distinguishing this record within its routing group."
  value       = var.spec.set_identifier
}
