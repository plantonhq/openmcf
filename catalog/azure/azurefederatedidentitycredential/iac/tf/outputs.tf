output "federated_identity_credential_id" {
  description = "The full ARM ID of the federated identity credential"
  value       = azurerm_federated_identity_credential.main.id
}

output "name" {
  description = "The credential's ARM resource name under its parent identity"
  value       = azurerm_federated_identity_credential.main.name
}

output "user_assigned_identity_id" {
  description = "The ARM ID of the parent user-assigned identity"
  value       = azurerm_federated_identity_credential.main.user_assigned_identity_id
}

output "issuer" {
  description = "The OIDC issuer URL the trust matches against the token's iss claim"
  value       = azurerm_federated_identity_credential.main.issuer
}

output "subject" {
  description = "The workload identifier the trust matches against the token's sub claim"
  value       = azurerm_federated_identity_credential.main.subject
}

output "audience" {
  # azurerm exposes the audience as a single-element list (ARM's wire shape);
  # the platform contract is one audience string, so the module exports the
  # sole element -- the same value the Pulumi engine exports from its
  # flattened attribute.
  description = "The audience the trust matches against the token's aud claim"
  value       = one(azurerm_federated_identity_credential.main.audience)
}
