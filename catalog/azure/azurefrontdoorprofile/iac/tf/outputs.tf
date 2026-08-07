output "profile_id" {
  description = "The Azure Resource Manager ID of the Front Door profile (what endpoints and origin groups reference as their parent)"
  value       = azurerm_cdn_frontdoor_profile.main.id
}

output "profile_name" {
  description = "The profile's name -- the ARM namespace every child resource nests under"
  value       = azurerm_cdn_frontdoor_profile.main.name
}

output "resource_guid" {
  description = "The Front Door service's own GUID for this profile (used for traffic-ownership validation, e.g. apex-domain afdverify records)"
  value       = azurerm_cdn_frontdoor_profile.main.resource_guid
}

output "identity_principal_id" {
  description = "The principal ID of the system-assigned managed identity (empty without one) -- the principal to grant Key Vault access to"
  value       = try(azurerm_cdn_frontdoor_profile.main.identity[0].principal_id, "")
}
