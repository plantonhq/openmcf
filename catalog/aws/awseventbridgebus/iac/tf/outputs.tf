output "bus_name" {
  description = "The name of the EventBridge custom event bus."
  value       = aws_cloudwatch_event_bus.this.name
}

output "bus_arn" {
  description = "The ARN of the EventBridge custom event bus."
  value       = aws_cloudwatch_event_bus.this.arn
}

output "archives" {
  description = "The event archives declared in spec.archives, in spec order -- each entry's name and the ARN AWS assigned it, for replay operations (StartReplay) and IAM policies."
  value       = [for archive in var.spec.archives : { name = archive.name, arn = aws_cloudwatch_event_archive.this[archive.name].arn }]
}
