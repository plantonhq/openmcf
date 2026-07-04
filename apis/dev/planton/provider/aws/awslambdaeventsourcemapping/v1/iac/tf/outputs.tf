output "uuid" {
  description = "The server-assigned mapping UUID -- the identity AWS APIs key on."
  value       = aws_lambda_event_source_mapping.this.uuid
}

output "mapping_arn" {
  description = "The mapping ARN."
  value       = aws_lambda_event_source_mapping.this.arn
}

output "function_arn" {
  description = "The ARN of the function the mapping invokes, as resolved by AWS."
  value       = aws_lambda_event_source_mapping.this.function_arn
}

output "state" {
  description = "The mapping state as last observed at deploy time (e.g. Enabled, Disabled)."
  value       = aws_lambda_event_source_mapping.this.state
}
