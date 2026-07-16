# Exactly one of the three mode/OS resources exists; every output
# coalesces across them.

output "scale_set_id" {
  description = "The Azure Resource Manager ID of the scale set"
  value = coalesce(
    one(azurerm_linux_virtual_machine_scale_set.main[*].id),
    one(azurerm_windows_virtual_machine_scale_set.main[*].id),
    one(azurerm_orchestrated_virtual_machine_scale_set.main[*].id),
  )
}

output "scale_set_name" {
  description = "The name of the scale set"
  value = coalesce(
    one(azurerm_linux_virtual_machine_scale_set.main[*].name),
    one(azurerm_windows_virtual_machine_scale_set.main[*].name),
    one(azurerm_orchestrated_virtual_machine_scale_set.main[*].name),
  )
}

output "unique_id" {
  description = "The scale set's globally unique ARM-assigned identifier"
  value = coalesce(
    one(azurerm_linux_virtual_machine_scale_set.main[*].unique_id),
    one(azurerm_windows_virtual_machine_scale_set.main[*].unique_id),
    one(azurerm_orchestrated_virtual_machine_scale_set.main[*].unique_id),
  )
}

output "system_assigned_identity_principal_id" {
  description = "The system-assigned managed identity's principal ID (empty unless the identity type includes SYSTEM_ASSIGNED -- UNIFORM sets only)"
  value = try(coalesce(
    try(one(azurerm_linux_virtual_machine_scale_set.main[*].identity[0].principal_id), null),
    try(one(azurerm_windows_virtual_machine_scale_set.main[*].identity[0].principal_id), null),
  ), "")
}
