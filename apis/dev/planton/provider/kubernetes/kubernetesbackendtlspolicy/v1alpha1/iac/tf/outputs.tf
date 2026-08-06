# Stack outputs — identical names and derivations in the Pulumi module's
# outputs.go / main.go exports.

output "policy_name" {
  description = "Name of the created BackendTLSPolicy (equals metadata.name)."
  value       = var.metadata.name
}

output "namespace" {
  description = "Namespace of the created BackendTLSPolicy."
  value       = var.spec.namespace
}
