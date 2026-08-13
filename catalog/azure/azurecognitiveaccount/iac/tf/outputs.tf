output "cognitive_account_id" {
  description = "The Azure Resource Manager ID of the account -- what model deployments and projects reference as their cognitive_account_id"
  value       = azurerm_cognitive_account.main.id
}

output "cognitive_account_name" {
  description = "The name of the account"
  value       = azurerm_cognitive_account.main.name
}

output "endpoint" {
  description = "The account's endpoint URL -- what applications call with a key or an Entra ID token"
  value       = azurerm_cognitive_account.main.endpoint
}

output "primary_access_key" {
  description = "The account's primary access key (empty when local_auth_enabled is false)"
  value       = azurerm_cognitive_account.main.primary_access_key
  sensitive   = true
}

output "secondary_access_key" {
  description = "The account's secondary access key, for zero-downtime rotation (empty when local_auth_enabled is false)"
  value       = azurerm_cognitive_account.main.secondary_access_key
  sensitive   = true
}

output "system_assigned_identity_principal_id" {
  description = "The principal (object) ID of the account's system-assigned identity, when one is enabled"
  value       = try(azurerm_cognitive_account.main.identity[0].principal_id, "")
}

output "rai_blocklist_ids" {
  description = "The ARM ID of each responsible-AI blocklist on the account, keyed by the blocklist's name from the spec"
  value       = { for name, blocklist in azurerm_cognitive_account_rai_blocklist.rai_blocklists : name => blocklist.id }
}

output "rai_policy_ids" {
  description = "The ARM ID of each responsible-AI policy on the account, keyed by the policy's name from the spec"
  value       = { for name, policy in azurerm_cognitive_account_rai_policy.rai_policies : name => policy.id }
}
