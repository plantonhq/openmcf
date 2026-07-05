output "task_definition_arn" {
  description = "The full ARN of the registered revision (family:revision) -- the handle ECS services reference; each new revision changes it and rolls the referencing service."
  value       = aws_ecs_task_definition.this.arn
}

output "arn_without_revision" {
  description = "The ARN without the revision suffix, for consumers that track the family's latest ACTIVE revision."
  value       = aws_ecs_task_definition.this.arn_without_revision
}

output "family" {
  description = "The family name the revisions are registered under."
  value       = aws_ecs_task_definition.this.family
}

output "revision" {
  description = "The revision number this deployment registered."
  value       = aws_ecs_task_definition.this.revision
}

output "log_group_name" {
  description = "The CloudWatch log group the task's containers log to (auto-created or referenced); empty when logging is disabled."
  value       = local.log_group_name
}

output "log_group_arn" {
  description = "The ARN of the auto-created log group; empty when logging is disabled or an existing group is referenced."
  value       = local.create_log_group ? aws_cloudwatch_log_group.this[0].arn : ""
}
