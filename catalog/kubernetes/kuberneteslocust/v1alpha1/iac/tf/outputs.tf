output "namespace" {
  description = "Namespace Locust runs in."
  value       = local.namespace
}

output "master_service" {
  description = "Name of the master Service — the handle exposure kinds route to. Serves the web UI (8089) and the worker-connect ports (5557/5558)."
  value       = local.master_service
}

output "web_endpoint" {
  description = "In-cluster web endpoint — the web UI and REST API (with login on, pair it with the credential handle)."
  value       = local.web_endpoint
}

output "master_bind_endpoint" {
  description = "In-cluster worker-connect endpoint — where additional workers register with the master."
  value       = local.master_bind_endpoint
}

output "web_ui_username" {
  description = "The web-UI login username. Empty when the login is disabled or the run is headless."
  value       = local.web_username
}

output "web_ui_password_secret" {
  description = "The Secret and key holding the web-UI login password (a NAME reference — the value stays in-cluster). Empty when the login is disabled or the run is headless."
  value = {
    name = local.web_ui_password_secret_output_name
    key  = local.web_ui_password_secret_output_key
  }
}

output "port_forward_command" {
  description = "Port-forward command for reaching the web UI from a workstation when no exposure is composed."
  value       = local.port_forward_command
}
