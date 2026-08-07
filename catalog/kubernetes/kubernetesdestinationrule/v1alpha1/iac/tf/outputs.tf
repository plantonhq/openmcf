output "destination_rule_name" {
  description = "Name of the created DestinationRule (equals metadata.name)."
  value       = var.metadata.name
}

output "namespace" {
  description = "Namespace of the created DestinationRule."
  value       = var.spec.namespace
}
