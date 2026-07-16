output "health_check_id" {
  description = "The health check ID — the value DNS records reference via health_check_id and calculated parents aggregate as children."
  value       = aws_route53_health_check.this.id
}

output "health_check_arn" {
  description = "The ARN of the health check, for IAM policies and the CloudWatch HealthCheckStatus metric dimension."
  value       = aws_route53_health_check.this.arn
}
