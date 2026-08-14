output "aggregator_name" {
  description = "The aggregator's name (also the provider's import ID); empty on a grants-only deployment"
  value       = length(aws_config_configuration_aggregator.this) > 0 ? aws_config_configuration_aggregator.this[0].name : ""
}

output "aggregator_arn" {
  description = "The aggregator's ARN; empty on a grants-only deployment"
  value       = length(aws_config_configuration_aggregator.this) > 0 ? aws_config_configuration_aggregator.this[0].arn : ""
}

output "authorization_arns" {
  description = "The grants' ARNs, keyed {account_id}:{authorized_aws_region} (each key is that grant's provider import ID)"
  value       = { for key, grant in aws_config_aggregate_authorization.grants : key => grant.arn }
}
