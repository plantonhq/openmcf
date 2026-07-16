output "job_id" {
  description = "The Azure Resource Manager ID of the Container App Job"
  value       = azurerm_container_app_job.main.id
}

output "job_name" {
  description = "The name of the Container App Job"
  value       = azurerm_container_app_job.main.name
}

output "event_stream_endpoint" {
  description = "The endpoint streaming the job's execution events"
  value       = azurerm_container_app_job.main.event_stream_endpoint
}

output "outbound_ip_addresses" {
  description = "Outbound IP addresses used by the job's replicas for egress traffic"
  value       = azurerm_container_app_job.main.outbound_ip_addresses
}

output "identity_principal_id" {
  description = "The principal ID of the job's system-assigned managed identity (empty unless SYSTEM_ASSIGNED is enabled)"
  value       = try(azurerm_container_app_job.main.identity[0].principal_id, "")
}
