output "canary_name" {
  description = "The canary's name (the provider's import ID for the canary); empty on groups-only instances"
  value       = var.spec.canary != null ? aws_synthetics_canary.this[0].name : ""
}

output "canary_arn" {
  description = "The canary's ARN; empty on groups-only instances"
  value       = var.spec.canary != null ? aws_synthetics_canary.this[0].arn : ""
}

output "engine_arn" {
  description = "ARN of the Synthetics-managed Lambda behind the canary"
  value       = var.spec.canary != null ? aws_synthetics_canary.this[0].engine_arn : ""
}

output "source_location_arn" {
  description = "ARN of the canary's staged code location"
  value       = var.spec.canary != null ? aws_synthetics_canary.this[0].source_location_arn : ""
}

output "canary_status" {
  description = "The canary's lifecycle status after apply (READY, RUNNING, ...)"
  value       = var.spec.canary != null ? aws_synthetics_canary.this[0].status : ""
}

output "group_arns" {
  description = "Owned group ARNs keyed by group name"
  value       = { for name, group in aws_synthetics_group.this : name => group.arn }
}

output "group_ids" {
  description = "Owned group IDs keyed by group name"
  value       = { for name, group in aws_synthetics_group.this : name => group.group_id }
}
