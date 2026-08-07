output "namespace" {
  description = "Namespace Airflow runs in."
  value       = local.namespace
}

output "api_server_service" {
  description = "Name of the API server Service (`<name>-api-server`) — the UI + REST API front door."
  value       = local.api_server_service_name
}

output "api_server_endpoint" {
  description = "In-cluster API server endpoint."
  value       = local.api_server_endpoint
}

output "admin_password_secret" {
  description = "The Secret and key holding the bootstrap admin password (a NAME reference — the value stays in-cluster). Empty when admin_user.create is false: no bootstrap credential exists to point at."
  value = {
    name = local.admin_password_secret_output_name
    key  = local.admin_password_secret_output_key
  }
}

output "metadata_connection_secret_name" {
  description = "Name of the module-owned Secret holding the metadata database connection URI (key `connection`)."
  value       = local.metadata_conn_secret_name
}

output "broker_url_secret_name" {
  description = "Name of the Secret holding the Celery broker URL (key `connection`); empty when no Celery executor is declared."
  value       = local.broker_url_secret_name
}

output "fernet_key_secret_name" {
  description = "Name of the Secret holding the Fernet key — back it up: losing it orphans every credential Airflow has stored."
  value       = local.fernet_key_secret_name
}

output "port_forward_command" {
  description = "Port-forward command for reaching the Airflow UI from a workstation when no exposure is composed."
  value       = local.port_forward_command
}
