output "topic_arn" {
  description = "The ARN of the SNS topic."
  value       = aws_sns_topic.this.arn
}

output "topic_name" {
  description = "The name of the SNS topic (includes .fifo suffix for FIFO topics)."
  value       = aws_sns_topic.this.name
}

output "owner" {
  description = "The AWS account ID that owns the topic."
  value       = aws_sns_topic.this.owner
}

output "beginning_archive_time" {
  description = "Timestamp from which archived messages are replayable (empty unless a FIFO archive policy is active)."
  value       = aws_sns_topic.this.beginning_archive_time
}
