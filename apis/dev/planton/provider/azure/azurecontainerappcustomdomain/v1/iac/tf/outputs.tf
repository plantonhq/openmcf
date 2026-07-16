# Exactly one of the two count-gated variants materializes, so each
# output coalesces PER ATTRIBUTE across them.

output "custom_domain_id" {
  description = "The binding's resource ID -- the providers' synthetic identifier ({container-app-id}/customDomainName/{domain}); Azure models the binding inside the app's ingress configuration."
  value = coalesce(
    try(azurerm_container_app_custom_domain.byo[0].id, null),
    try(azurerm_container_app_custom_domain.managed[0].id, null),
  )
}

output "managed_certificate_id" {
  description = "The ARM ID of the Azure-managed certificate attached to this binding. Empty for bring-your-own bindings, and empty until Azure attaches the managed certificate (asynchronous)."
  value = try(
    coalesce(
      try(azurerm_container_app_custom_domain.byo[0].container_app_environment_managed_certificate_id, null),
      try(azurerm_container_app_custom_domain.managed[0].container_app_environment_managed_certificate_id, null),
    ),
    "",
  )
}
