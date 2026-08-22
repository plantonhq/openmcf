output "policy_id" {
  description = "The policy's identity: its name (account scope) or the target resource ARN (resource scope) - also the provider's import ID"
  value       = aws_cloudwatch_log_resource_policy.this.id
}

output "policy_scope" {
  description = "The scope AWS recorded (ACCOUNT or RESOURCE)"
  value       = aws_cloudwatch_log_resource_policy.this.policy_scope
}

output "revision_id" {
  description = "AWS's revision ID after the last apply (the optimistic-concurrency token)"
  value       = aws_cloudwatch_log_resource_policy.this.revision_id
}
