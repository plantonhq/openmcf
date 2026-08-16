output "report_plan_arn" {
  description = "The report plan's ARN"
  value       = aws_backup_report_plan.this.arn
}
