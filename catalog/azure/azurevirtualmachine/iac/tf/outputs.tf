output "vm_id" {
  description = "The Azure Resource Manager ID of the VM"
  value       = local.is_linux ? azurerm_linux_virtual_machine.main[0].id : azurerm_windows_virtual_machine.main[0].id
}

output "vm_name" {
  description = "The name of the VM"
  value       = local.is_linux ? azurerm_linux_virtual_machine.main[0].name : azurerm_windows_virtual_machine.main[0].name
}

output "virtual_machine_guid" {
  description = "The 128-bit unique GUID Azure assigns the VM"
  value       = local.is_linux ? azurerm_linux_virtual_machine.main[0].virtual_machine_id : azurerm_windows_virtual_machine.main[0].virtual_machine_id
}

output "private_ip_address" {
  description = "The primary private IP address across the VM's attached NICs"
  value       = local.is_linux ? azurerm_linux_virtual_machine.main[0].private_ip_address : azurerm_windows_virtual_machine.main[0].private_ip_address
}

output "public_ip_address" {
  description = "The primary public IP address across the VM's attached NICs (empty for private-only VMs)"
  value       = local.is_linux ? azurerm_linux_virtual_machine.main[0].public_ip_address : azurerm_windows_virtual_machine.main[0].public_ip_address
}

output "computer_name" {
  description = "The OS hostname the VM booted with"
  value       = local.is_linux ? azurerm_linux_virtual_machine.main[0].computer_name : azurerm_windows_virtual_machine.main[0].computer_name
}

output "system_assigned_identity_principal_id" {
  description = "The principal ID of the VM's system-assigned identity (populated only when the identity type includes SYSTEM_ASSIGNED)"
  value = local.is_linux ? (
    try(azurerm_linux_virtual_machine.main[0].identity[0].principal_id, "")
    ) : (
    try(azurerm_windows_virtual_machine.main[0].identity[0].principal_id, "")
  )
}
