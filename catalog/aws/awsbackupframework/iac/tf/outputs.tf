output "framework_arn" {
  description = "The framework's ARN - what report plans reference"
  value       = aws_backup_framework.this.arn
}
