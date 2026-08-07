output "certificate_id" {
  description = "The Azure Resource Manager ID of the managed certificate."
  value       = azurerm_container_app_environment_managed_certificate.main.id
}

output "validation_token" {
  description = "The domain-validation token Azure generated for this certificate. Informational once the certificate issues (the deployment waits for validation)."
  value       = azurerm_container_app_environment_managed_certificate.main.validation_token
}
