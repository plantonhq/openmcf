output "subscription_arn" {
  description = "The ARN of the SNS subscription."
  value       = aws_sns_topic_subscription.this.arn
}

output "owner_id" {
  description = "The AWS account ID that owns the subscription."
  value       = aws_sns_topic_subscription.this.owner_id
}

output "pending_confirmation" {
  description = "True while the subscription is awaiting endpoint confirmation."
  value       = aws_sns_topic_subscription.this.pending_confirmation
}

output "confirmation_was_authenticated" {
  description = "True when the confirmation request was authenticated (signed)."
  value       = aws_sns_topic_subscription.this.confirmation_was_authenticated
}
