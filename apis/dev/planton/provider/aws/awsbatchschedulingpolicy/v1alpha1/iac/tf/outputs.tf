output "scheduling_policy_arn" {
  description = "The ARN of the scheduling policy -- what job queues reference through their scheduling_policy field."
  value       = aws_batch_scheduling_policy.this.arn
}

output "scheduling_policy_name" {
  description = "The scheduling policy's name (derived from metadata.name)."
  value       = aws_batch_scheduling_policy.this.name
}
