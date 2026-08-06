output "namespace" {
  description = "Namespace Superset runs in."
  value       = local.namespace
}

output "service" {
  description = "Name of the Superset web Service — the handle exposure kinds route to (port 8088)."
  value       = local.service
}

output "endpoint" {
  description = "In-cluster endpoint — the URL browsers and API clients reach Superset at (behind composed exposure)."
  value       = local.endpoint
}

output "admin_username" {
  description = "The bootstrap admin username."
  value       = local.admin_username
}

output "admin_password_secret" {
  description = "The Secret and key holding the bootstrap admin password (a NAME reference — the value stays in-cluster)."
  value = {
    name = local.admin_password_secret_name
    key  = local.admin_password_secret_key
  }
}

output "env_secret_name" {
  description = "Name of the module-owned environment Secret — the chart's runtime-credential contract, exported for audit and advanced composition."
  value       = local.env_secret_name
}

output "secret_key_secret_name" {
  description = "Name of the Secret holding the session-signing SECRET_KEY — needed for `superset re-encrypt-secrets` rotation procedures."
  value       = local.secret_key_secret_name
}

output "port_forward_command" {
  description = "Port-forward command for reaching the Superset UI from a workstation when no exposure is composed."
  value       = local.port_forward_command
}
