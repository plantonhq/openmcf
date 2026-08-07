output "autoscaling_group_name" {
  description = "The name of the auto-scaling group (the handle CloudWatch dimensions and ECS capacity providers reference)."
  value       = aws_autoscaling_group.this.name
}

output "autoscaling_group_arn" {
  description = "The ARN of the auto-scaling group, for IAM policies and EventBridge rules scoped to this group."
  value       = aws_autoscaling_group.this.arn
}
