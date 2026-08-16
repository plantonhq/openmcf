output "event_data_store_arn" {
  description = "The event data store's ARN (also the provider's import ID)"
  value       = aws_cloudtrail_event_data_store.this.arn
}
