output "dashboard_name" {
  description = "The dashboard's name (the provider's import ID)"
  value       = aws_cloudwatch_dashboard.this.dashboard_name
}

output "dashboard_arn" {
  description = "The dashboard's ARN"
  value       = aws_cloudwatch_dashboard.this.dashboard_arn
}
