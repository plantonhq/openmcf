output "query_log_config_id" {
  description = "The configuration's id (rqlc-...) - the provider's import ID"
  value       = aws_route53_resolver_query_log_config.this.id
}

output "query_log_config_arn" {
  description = "The configuration's ARN"
  value       = aws_route53_resolver_query_log_config.this.arn
}

output "share_status" {
  description = "Whether the configuration is shared via RAM (NOT_SHARED / SHARED_BY_ME / SHARED_WITH_ME)"
  value       = aws_route53_resolver_query_log_config.this.share_status
}

output "association_ids" {
  description = "AWS-generated association IDs (rqlca-...) keyed by the resolved VPC id"
  value       = { for vpc_id, assoc in aws_route53_resolver_query_log_config_association.this : vpc_id => assoc.id }
}
