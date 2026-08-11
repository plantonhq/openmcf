output "backup_protected_vm_id" {
  description = "The Azure Resource Manager ID of the protected item"
  value       = azurerm_backup_protected_vm.main.id
}
