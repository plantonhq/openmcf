output "flow_log_id" {
  description = "The Azure Resource Manager ID of the flow log"
  value       = azurerm_network_watcher_flow_log.main.id
}

output "flow_log_name" {
  description = "The name of the flow log resource"
  value       = azurerm_network_watcher_flow_log.main.name
}

output "network_watcher_name" {
  description = "The Network Watcher the flow log attached to (the auto-created regional singleton unless the spec addressed a self-managed one)"
  value       = azurerm_network_watcher_flow_log.main.network_watcher_name
}
