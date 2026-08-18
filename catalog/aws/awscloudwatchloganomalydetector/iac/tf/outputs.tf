output "anomaly_detector_arn" {
  description = "The detector's ARN (its identity and the provider's import ID)"
  value       = aws_cloudwatch_log_anomaly_detector.this.arn
}
