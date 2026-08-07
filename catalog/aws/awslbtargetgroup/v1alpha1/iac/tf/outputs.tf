output "target_group_arn" {
  description = "The ARN of the target group (what listeners, rules, ECS services, and ASGs reference)."
  value       = aws_lb_target_group.this.arn
}

output "target_group_name" {
  description = "The friendly name of the target group (metadata.name, truncated to AWS's 32-character limit)."
  value       = aws_lb_target_group.this.name
}

output "arn_suffix" {
  description = "The ARN suffix used as the TargetGroup dimension in CloudWatch metrics."
  value       = aws_lb_target_group.this.arn_suffix
}
