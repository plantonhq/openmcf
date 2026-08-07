# Stack outputs — identical names and derivations in the Pulumi module's
# outputs.go / main.go exports.

output "route_name" {
  description = "Name of the created UDPRoute (equals metadata.name)."
  value       = var.metadata.name
}

output "namespace" {
  description = "Namespace of the created UDPRoute."
  value       = var.spec.namespace
}
