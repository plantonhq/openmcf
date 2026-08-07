output "failover_group_id" {
  description = "The Azure Resource Manager ID of the failover group"
  value       = azurerm_mssql_failover_group.main.id
}

output "failover_group_name" {
  description = "The failover group's name (also its listener DNS label)"
  value       = azurerm_mssql_failover_group.main.name
}

# The listener FQDNs are DNS-derived from the group name (Azure does not
# return them as resource attributes); compose them so downstream connection
# strings reference a single failover-following endpoint.
output "read_write_listener_endpoint" {
  description = "The read-write listener FQDN (always points at the current primary)"
  value       = "${azurerm_mssql_failover_group.main.name}.database.windows.net"
}

output "read_only_listener_endpoint" {
  description = "The read-only listener FQDN for read-only workloads on a secondary"
  value       = "${azurerm_mssql_failover_group.main.name}.secondary.database.windows.net"
}
