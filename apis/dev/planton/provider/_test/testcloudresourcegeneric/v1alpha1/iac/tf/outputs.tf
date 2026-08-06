# Stack outputs -- names and meanings mirror stack_outputs.proto exactly.

output "id" {
  description = "Synthetic resource id (id_prefix + metadata.name)"
  value       = local.resource_id
}

output "name" {
  description = "The resource name"
  value       = var.metadata.name
}

output "endpoint" {
  description = "Deterministic endpoint derived from inputs"
  value       = local.endpoint
}

output "tags" {
  description = "Echo of spec.commands -- proves repeated fields round-trip"
  value       = var.spec.commands
}
