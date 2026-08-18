output "registry_id" {
  description = "The registry id (the account's 12-digit id) - the import ID for the policy/scanning/replication singletons"
  value       = data.aws_caller_identity.this.account_id
}

output "registry_url" {
  description = "The registry's pull URL base - the prefix cached upstream images pull through"
  value       = "${data.aws_caller_identity.this.account_id}.dkr.ecr.${var.spec.region}.amazonaws.com"
}

output "pull_through_cache_rule_registry_ids" {
  description = "Cache-rule registry ids keyed by ecr_repository_prefix (each rule's import ID is its prefix)"
  value       = { for prefix, cache_rule in aws_ecr_pull_through_cache_rule.this : prefix => cache_rule.registry_id }
}

output "repository_creation_template_registry_ids" {
  description = "Creation-template registry ids keyed by prefix (each template's import ID is its prefix)"
  value       = { for prefix, template in aws_ecr_repository_creation_template.this : prefix => template.registry_id }
}

output "pull_time_update_exclusion_arns" {
  description = "Pull-time update exclusions keyed by the resolved principal ARN (each exclusion's import ID is the ARN itself)"
  value       = { for arn, exclusion in aws_ecr_pull_time_update_exclusion.this : arn => exclusion.principal_arn }
}
