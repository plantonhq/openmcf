output "application_insights_id" {
  description = "The Azure Resource Manager ID of the Application Insights resource"
  value       = azurerm_application_insights.main.id
}

output "application_insights_name" {
  description = "The name of the Application Insights resource"
  value       = azurerm_application_insights.main.name
}

output "instrumentation_key" {
  description = "The instrumentation key for classic SDK configuration (prefer the connection string)"
  value       = azurerm_application_insights.main.instrumentation_key
  sensitive   = true
}

output "connection_string" {
  description = "The connection string SDKs are configured with -- the composition seam app kinds reference"
  value       = azurerm_application_insights.main.connection_string
  sensitive   = true
}

output "app_id" {
  description = "The Application ID used when querying telemetry via the REST API"
  value       = azurerm_application_insights.main.app_id
}
