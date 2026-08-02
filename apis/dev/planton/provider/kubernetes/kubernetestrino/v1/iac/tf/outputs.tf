output "namespace" {
  description = "Namespace Trino runs in."
  value       = local.namespace
}

output "coordinator_service" {
  description = "Name of the coordinator Service — the handle exposure kinds route to."
  value       = local.coordinator_service
}

output "coordinator_endpoint" {
  description = "In-cluster coordinator endpoint — the URL SQL clients and BI tools connect to (with auth on, pair it with the admin credential)."
  value       = local.coordinator_endpoint
}

output "admin_username" {
  description = "The bootstrap admin username. Empty when auth is disabled or a bring-your-own password file is declared."
  value       = local.admin_username
}

output "admin_password_secret" {
  description = "The Secret and key holding the bootstrap admin password (a NAME reference — the value stays in-cluster). Empty when auth is disabled or bring-your-own."
  value = {
    name = local.admin_password_secret_output_name
    key  = local.admin_password_secret_output_key
  }
}

output "worker_service" {
  description = "Name of the worker Service — internal; exported for network-policy composition."
  value       = local.worker_service
}

output "port_forward_command" {
  description = "Port-forward command for reaching the Trino Web UI from a workstation when no exposure is composed."
  value       = local.port_forward_command
}
