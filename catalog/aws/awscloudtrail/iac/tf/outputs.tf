output "trail_arn" {
  description = "The trail's ARN (also the provider's import ID)"
  value       = aws_cloudtrail.this.arn
}

output "home_region" {
  description = "The trail's home region"
  value       = aws_cloudtrail.this.home_region
}

output "sns_topic_arn" {
  description = "The ARN of the SNS topic notified on log delivery"
  value       = aws_cloudtrail.this.sns_topic_arn
}
