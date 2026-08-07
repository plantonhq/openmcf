output "service_plan_id" {
  description = "The Azure Resource Manager ID of the Service Plan -- the value downstream app kinds reference as service_plan_id"
  value       = azurerm_service_plan.main.id
}

output "service_plan_name" {
  description = "The name of the Service Plan"
  value       = azurerm_service_plan.main.name
}

output "os_type" {
  description = "The configured operating system type (Linux, Windows, or WindowsContainer)"
  value       = azurerm_service_plan.main.os_type
}

output "sku_name" {
  description = "The configured SKU name in Azure's spelling (e.g. P1v3, EP1, Y1)"
  value       = azurerm_service_plan.main.sku_name
}

output "kind" {
  description = "Azure's computed plan kind (e.g. linux, elastic, functionapp) -- the API's own classification, read back after creation"
  value       = azurerm_service_plan.main.kind
}

output "reserved" {
  description = "Whether the plan runs Linux workers (reserved = true in the Azure API), read back after creation"
  value       = azurerm_service_plan.main.reserved
}
