# Composition handles. All names derive from the fullnameOverride pin
# (= metadata.name); the front-door Service name is pinned by the
# module. Twin of the Pulumi module's outputs.go.

output "namespace" {
  description = "Namespace Harbor is installed into."
  value       = local.namespace
}

output "expose_service" {
  description = "The front-door Service (nginx) — the ONE service clients and composed exposure kinds point at; it proxies UI, API, and OCI traffic."
  value       = local.release_name
}

output "kube_endpoint" {
  description = "In-cluster URL of the front door."
  value       = local.kube_endpoint
}

output "external_url" {
  description = "The declared external URL — the address OCI clients must use (token-service auth is bound to it)."
  value       = var.spec.external_url
}

output "core_service" {
  description = "Core (API) Service name."
  value       = "${local.release_name}-core"
}

output "portal_service" {
  description = "Portal (web UI) Service name."
  value       = "${local.release_name}-portal"
}

output "registry_service" {
  description = "Registry (OCI distribution) Service name."
  value       = "${local.release_name}-registry"
}

output "jobservice_service" {
  description = "Jobservice Service name."
  value       = "${local.release_name}-jobservice"
}

output "trivy_service" {
  description = "Trivy Service name — empty when the scanner is disabled."
  value       = local.trivy_enabled ? "${local.release_name}-trivy" : ""
}

output "database_service" {
  description = "Internal database Service name — set only on the internal database arm."
  value       = local.internal_database ? "${local.release_name}-database" : ""
}

output "redis_service" {
  description = "Internal Redis Service name — set only on the internal cache arm."
  value       = local.internal_redis ? "${local.release_name}-redis" : ""
}

output "admin_username" {
  description = "Admin username — always \"admin\" (a Harbor constant)."
  value       = "admin"
}

output "admin_password_secret" {
  description = "The Secret (name/key) holding the admin password — the module-owned <name>-admin-auth on the generated arm, or the declared Secret/key."
  value = {
    name = local.admin_secret_name
    key  = local.admin_secret_key
  }
}

output "port_forward_command" {
  description = "Copy-paste command to reach the UI/API from a workstation. Pushes/pulls through this tunnel only authenticate when external_url matches the forwarded address."
  value       = local.port_forward_command
}
