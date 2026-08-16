output "policy_id" {
  description = "The policy's id (policy-...) - the provider's import ID"
  value       = aws_dlm_lifecycle_policy.this.id
}

output "policy_arn" {
  description = "The policy's ARN"
  value       = aws_dlm_lifecycle_policy.this.arn
}
