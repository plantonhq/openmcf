output "schedule_arn" {
  description = "The schedule's ARN"
  value       = aws_scheduler_schedule.this.arn
}

output "group_name" {
  description = "The group the schedule lives in (owned, joined, or AWS's 'default') - with metadata.name this forms the provider's '{group}/{name}' import ID"
  value       = aws_scheduler_schedule.this.group_name
}

output "group_arn" {
  description = "The owned group's ARN. Empty when the instance owns no group"
  value       = var.spec.group != null ? aws_scheduler_schedule_group.this[0].arn : ""
}
