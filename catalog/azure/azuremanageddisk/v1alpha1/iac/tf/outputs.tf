output "disk_id" {
  description = "The Azure Resource Manager ID of the disk"
  value       = azurerm_managed_disk.main.id
}

output "disk_name" {
  description = "The name of the disk"
  value       = azurerm_managed_disk.main.name
}

output "disk_size_gb" {
  description = "The disk's actual size in GiB (inherited from the source when the spec omitted it)"
  value       = azurerm_managed_disk.main.disk_size_gb
}
