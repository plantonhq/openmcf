output "window_id" {
  description = "The window's AWS-generated ID (\"mw-...\" - also the provider's import ID)"
  value       = aws_ssm_maintenance_window.this.id
}

output "target_ids" {
  description = "AWS-generated target registration IDs keyed by target name (each imports as \"{window_id}/{target_id}\")"
  value       = { for name, t in aws_ssm_maintenance_window_target.this : name => t.id }
}

output "task_ids" {
  description = "AWS-generated task registration IDs keyed by task name (each imports as \"{window_id}/{task_id}\")"
  value       = { for name, t in aws_ssm_maintenance_window_task.this : name => t.window_task_id }
}
