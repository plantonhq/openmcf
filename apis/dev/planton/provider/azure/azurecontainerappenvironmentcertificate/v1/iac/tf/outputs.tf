output "certificate_id" {
  description = "The Azure Resource Manager ID of the certificate -- the binding seam AzureContainerAppCustomDomain consumes."
  value       = azurerm_container_app_environment_certificate.main.id
}

output "subject_name" {
  description = "The certificate's subject name as read back from the uploaded material."
  value       = azurerm_container_app_environment_certificate.main.subject_name
}

output "issuer" {
  description = "The certificate's issuer."
  value       = azurerm_container_app_environment_certificate.main.issuer
}

output "issue_date" {
  description = "When the certificate was issued."
  value       = azurerm_container_app_environment_certificate.main.issue_date
}

output "expiration_date" {
  description = "When the certificate expires -- the value to alarm on for inline-PFX certificates, whose rotation is manual."
  value       = azurerm_container_app_environment_certificate.main.expiration_date
}

output "thumbprint" {
  description = "The certificate's SHA-1 thumbprint."
  value       = azurerm_container_app_environment_certificate.main.thumbprint
}
