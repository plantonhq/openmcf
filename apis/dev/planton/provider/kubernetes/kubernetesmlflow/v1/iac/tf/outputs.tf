output "namespace" {
  description = "Namespace MLflow runs in."
  value       = local.namespace
}

output "service" {
  description = "Name of the MLflow Service — the handle exposure kinds route to."
  value       = local.name
}

output "tracking_endpoint" {
  description = "In-cluster tracking endpoint — set it as MLFLOW_TRACKING_URI in training jobs and pipelines (with auth on, pair it with MLFLOW_TRACKING_USERNAME/PASSWORD)."
  value       = local.tracking_endpoint
}

output "admin_password_secret" {
  description = "The Secret and key holding the bootstrap admin password (a NAME reference — the value stays in-cluster). Empty when auth is disabled."
  value = {
    name = local.admin_password_secret_output_name
    key  = local.admin_password_secret_output_key
  }
}

output "backend_store_uri_secret_name" {
  description = "Name of the module-owned Secret holding the backend-store connection URI (key `uri`) — empty on the sqlite arm."
  value       = local.backend_uri_secret_output_name
}

output "port_forward_command" {
  description = "Port-forward command for reaching the MLflow UI from a workstation when no exposure is composed."
  value       = local.port_forward_command
}
