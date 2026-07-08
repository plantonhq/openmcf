output "delivery_stream_arn" {
  description = "The ARN of the Kinesis Firehose delivery stream."
  value       = aws_kinesis_firehose_delivery_stream.this.arn
}

output "delivery_stream_name" {
  description = "The name of the Kinesis Firehose delivery stream."
  value       = aws_kinesis_firehose_delivery_stream.this.name
}

output "destination_id" {
  description = "Identifier of the destination configuration within the delivery stream (required by the UpdateDestination API)."
  value       = aws_kinesis_firehose_delivery_stream.this.destination_id
}

output "version_id" {
  description = "Version of the delivery stream configuration (incremented by AWS on every update; the UpdateDestination optimistic-concurrency token)."
  value       = aws_kinesis_firehose_delivery_stream.this.version_id
}
