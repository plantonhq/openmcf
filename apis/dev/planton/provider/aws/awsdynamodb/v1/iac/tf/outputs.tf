output "table_name" {
  description = "The table name -- what SDK calls, IAM policy resources, and application configuration reference."
  value       = aws_dynamodb_table.this.name
}

output "table_arn" {
  description = "The table ARN -- the join key for IAM policies, resource policies, and cross-service integrations."
  value       = aws_dynamodb_table.this.arn
}

output "table_id" {
  description = "The provider-assigned table identifier."
  value       = aws_dynamodb_table.this.id
}

output "stream_arn" {
  description = "The DynamoDB Streams ARN -- what Lambda event-source mappings attach to. Empty when streams are disabled."
  value       = aws_dynamodb_table.this.stream_arn
}

output "stream_label" {
  description = "The stream label; combined with the account and table name it uniquely identifies the stream. Empty when streams are disabled."
  value       = aws_dynamodb_table.this.stream_label
}
