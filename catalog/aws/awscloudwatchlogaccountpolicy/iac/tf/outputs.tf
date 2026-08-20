output "policy_name" {
  description = "The account policy's name"
  value       = aws_cloudwatch_log_account_policy.this.policy_name
}

output "policy_type" {
  description = "The account policy's type (with the name, the provider's import ID \"policy_name:policy_type\")"
  value       = aws_cloudwatch_log_account_policy.this.policy_type
}
