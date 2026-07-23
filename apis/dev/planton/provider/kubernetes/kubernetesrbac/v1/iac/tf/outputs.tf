# Terraform outputs for the Kubernetes RBAC grant.
# Keys match stack_outputs.proto field names and the Pulumi module's exports.

output "role_name" {
  description = "The name of the role in the grant: the created Role/ClusterRole, or the existing role that was bound to"
  value       = local.role_name
}

output "role_kind" {
  description = "The Kubernetes kind of the role in the grant: Role or ClusterRole"
  value       = local.role_kind
}

output "binding_name" {
  description = "The name of the created binding; empty when the grant has no subjects (no binding is created)"
  value       = local.binding_name
}

output "binding_kind" {
  description = "The Kubernetes kind of the created binding: RoleBinding or ClusterRoleBinding; empty when no binding is created"
  value       = local.binding_kind
}

output "namespace" {
  description = "The namespace the grant applies to; empty for cluster-scoped grants"
  value       = local.namespace
}
