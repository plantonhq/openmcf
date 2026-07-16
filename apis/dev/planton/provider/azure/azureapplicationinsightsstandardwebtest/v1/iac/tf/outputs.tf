# The ARM ID is the seam an AzureMonitorMetricAlert references (web_test_id)
# to alert on the test's availability.
output "web_test_id" {
  description = "The Azure Resource Manager ID of the web test"
  value       = azurerm_application_insights_standard_web_test.main.id
}

output "web_test_name" {
  description = "The name of the web test resource"
  value       = azurerm_application_insights_standard_web_test.main.name
}

output "synthetic_monitor_id" {
  description = "The synthetic monitor id Azure assigns the test"
  value       = azurerm_application_insights_standard_web_test.main.synthetic_monitor_id
}
