output "rule_arn" {
  description = "The rule's ARN (organization rules report the organization rule ARN)"
  value = coalesce(
    length(aws_config_config_rule.this) > 0 ? aws_config_config_rule.this[0].arn : null,
    length(aws_config_organization_managed_rule.this) > 0 ? aws_config_organization_managed_rule.this[0].arn : null,
    length(aws_config_organization_custom_rule.this) > 0 ? aws_config_organization_custom_rule.this[0].arn : null,
    length(aws_config_organization_custom_policy_rule.this) > 0 ? aws_config_organization_custom_policy_rule.this[0].arn : null,
  )
}

output "rule_name" {
  description = "The rule's name (metadata.name echoed) - the key remediations and aggregator queries address rules by"
  value       = var.metadata.name
}

output "rule_id" {
  description = "The rule's AWS-assigned ID (account-scoped rules; empty for organization rules)"
  value       = length(aws_config_config_rule.this) > 0 ? aws_config_config_rule.this[0].rule_id : ""
}

output "remediation_arn" {
  description = "The remediation configuration's ARN (set only when spec.remediation is configured)"
  value       = length(aws_config_remediation_configuration.this) > 0 ? aws_config_remediation_configuration.this[0].arn : ""
}
