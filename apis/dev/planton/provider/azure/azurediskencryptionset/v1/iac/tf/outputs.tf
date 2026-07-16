# The ARM ID is the composition seam: managed disks, VMs, and scale sets
# reference it to encrypt their disks with the customer-managed key.
output "disk_encryption_set_id" {
  description = "The Azure Resource Manager ID of the disk encryption set"
  value       = azurerm_disk_encryption_set.main.id
}

output "disk_encryption_set_name" {
  description = "The name of the disk encryption set resource"
  value       = azurerm_disk_encryption_set.main.name
}

# The principal to grant Key Vault crypto access so the set can unwrap the
# key (empty for user-assigned-only sets, which grant their own identities).
output "identity_principal_id" {
  description = "The principal ID of the set's system-assigned identity"
  value       = try(azurerm_disk_encryption_set.main.identity[0].principal_id, "")
}

output "identity_tenant_id" {
  description = "The tenant ID of the set's system-assigned identity"
  value       = try(azurerm_disk_encryption_set.main.identity[0].tenant_id, "")
}
