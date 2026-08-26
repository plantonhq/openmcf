output "framework_arn" {
  description = "The framework's ARN - what report plans reference"
  value       = aws_backup_framework.this.arn
}

output "region" {
  description = "The AWS region the framework lives in (frameworks are addressed by region + name)"
  value       = aws_backup_framework.this.region
}
