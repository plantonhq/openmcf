output "source_name" {
  description = "The vended source's name (the provider's import ID for the source); empty without a source"
  value       = var.spec.vended != null && var.spec.vended.source != null ? aws_cloudwatch_log_delivery_source.this[0].name : ""
}

output "source_arn" {
  description = "The vended source's ARN; empty without a source"
  value       = var.spec.vended != null && var.spec.vended.source != null ? aws_cloudwatch_log_delivery_source.this[0].arn : ""
}

output "source_service" {
  description = "The service AWS recorded as the source's producer"
  value       = var.spec.vended != null && var.spec.vended.source != null ? aws_cloudwatch_log_delivery_source.this[0].service : ""
}

output "destination_arns" {
  description = "Owned destination ARNs keyed by destination name - what other instances' deliveries reference"
  value       = { for name, destination in aws_cloudwatch_log_delivery_destination.this : name => destination.arn }
}

output "delivery_ids" {
  description = "AWS-generated delivery IDs keyed by delivery name (each delivery imports by this ID)"
  value       = { for name, delivery in aws_cloudwatch_log_delivery.this : name => delivery.id }
}

output "delivery_arns" {
  description = "Delivery ARNs keyed by delivery name"
  value       = { for name, delivery in aws_cloudwatch_log_delivery.this : name => delivery.arn }
}

output "cross_account_destination_name" {
  description = "The cross-account destination's name (the provider's import ID for it); empty without that arm"
  value       = var.spec.cross_account_destination != null ? aws_cloudwatch_log_destination.this[0].name : ""
}

output "cross_account_destination_arn" {
  description = "The cross-account destination's ARN - what other accounts' subscription filters target; empty without that arm"
  value       = var.spec.cross_account_destination != null ? aws_cloudwatch_log_destination.this[0].arn : ""
}
