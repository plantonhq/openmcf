output "monitor_arn" {
  description = "The monitor's ARN (also the provider's import ID)"
  value       = aws_ce_anomaly_monitor.this.arn
}

output "subscription_arns" {
  description = "Subscription ARNs keyed by subscription name (each subscription imports by its ARN)"
  value       = { for name, subscription in aws_ce_anomaly_subscription.this : name => subscription.arn }
}
