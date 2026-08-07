# The ARM ID is the composition seam: rule collection groups nest under
# it, firewalls attach it, and child policies inherit from it.
output "firewall_policy_id" {
  description = "The Azure Resource Manager ID of the firewall policy"
  value       = azurerm_firewall_policy.main.id
}

output "firewall_policy_name" {
  description = "The name of the firewall policy resource"
  value       = azurerm_firewall_policy.main.name
}

# Empty when the policy carries no system-assigned identity. Grant it Key
# Vault secret read access when TLS inspection rides the system identity.
output "identity_principal_id" {
  description = "The principal id of the policy's system-assigned managed identity"
  value = try(
    azurerm_firewall_policy.main.identity[0].principal_id,
    ""
  )
}
