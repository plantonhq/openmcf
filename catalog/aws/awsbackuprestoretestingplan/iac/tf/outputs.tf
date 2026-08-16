output "restore_testing_plan_arn" {
  description = "The restore testing plan's ARN (the plan and its selections import by name)"
  value       = aws_backup_restore_testing_plan.this.arn
}
