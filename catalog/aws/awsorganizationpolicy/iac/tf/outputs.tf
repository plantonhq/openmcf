output "policy_id" {
  description = "The policy's AWS-generated ID (p-... - also the provider's import ID; each attachment imports as {target_id}:{policy_id})"
  value       = aws_organizations_policy.this.id
}

output "arn" {
  description = "The policy's ARN"
  value       = aws_organizations_policy.this.arn
}
