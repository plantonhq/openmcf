# Exactly one of the three count-gated endpoint resources materializes,
# so each output coalesces per attribute across the variants.

output "endpoint_id" {
  description = "The Azure Resource Manager ID of the endpoint. Format: {profile_id}/{TYPE}/{name} where {TYPE} is AzureEndpoints, ExternalEndpoints, or NestedEndpoints per the spec's variant"
  value = coalesce(
    try(azurerm_traffic_manager_azure_endpoint.main[0].id, null),
    try(azurerm_traffic_manager_external_endpoint.main[0].id, null),
    try(azurerm_traffic_manager_nested_endpoint.main[0].id, null),
  )
}

output "endpoint_name" {
  description = "The endpoint's name within its profile"
  value = coalesce(
    try(azurerm_traffic_manager_azure_endpoint.main[0].name, null),
    try(azurerm_traffic_manager_external_endpoint.main[0].name, null),
    try(azurerm_traffic_manager_nested_endpoint.main[0].name, null),
  )
}
