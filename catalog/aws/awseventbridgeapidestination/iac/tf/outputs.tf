output "connection_arn" {
  description = "The owned connection's ARN - what other instances' destinations and pipe/rule targets reference. Empty when the instance has no connection arm"
  value       = var.spec.connection != null ? aws_cloudwatch_event_connection.this[0].arn : ""
}

output "connection_secret_arn" {
  description = "The Secrets Manager secret AWS created for the connection's credentials (AWS owns its lifecycle). Empty when the instance has no connection arm"
  value       = var.spec.connection != null ? aws_cloudwatch_event_connection.this[0].secret_arn : ""
}

output "api_destination_arn" {
  description = "The owned API destination's ARN - what EventBridge rule targets, pipes, and schedules invoke. Empty when the instance has no destination arm"
  value       = var.spec.destination != null ? aws_cloudwatch_event_api_destination.this[0].arn : ""
}
