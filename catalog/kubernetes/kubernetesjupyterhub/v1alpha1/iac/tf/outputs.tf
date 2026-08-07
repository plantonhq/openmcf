output "namespace" {
  description = "Namespace JupyterHub runs in."
  value       = local.namespace
}

output "proxy_public_service" {
  description = "Name of the proxy-public Service — the instance's front door (chart-fixed name)."
  value       = local.proxy_public_service_name
}

output "endpoint" {
  description = "In-cluster endpoint of the front door."
  value       = local.proxy_public_endpoint
}

output "hub_service" {
  description = "Name of the hub's own Service (chart-fixed name; the hub REST API handle for in-cluster automation)."
  value       = local.hub_service_name
}

output "shared_password_secret" {
  description = "The Secret and key holding the shared sign-in password (a NAME reference — the value stays in-cluster). Empty for OAuth/OIDC/native sign-in: those identities live with the provider."
  value = {
    name = local.shared_password_output_name
    key  = local.shared_password_output_key
  }
}

output "port_forward_command" {
  description = "Port-forward command for reaching JupyterHub from a workstation when no exposure is composed."
  value       = local.port_forward_command
}
