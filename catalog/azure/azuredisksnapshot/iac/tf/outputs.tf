output "snapshot_id" {
  description = "The Azure Resource Manager ID of the snapshot -- what disks restore from and gallery image versions build from"
  value       = azurerm_snapshot.main.id
}

output "snapshot_name" {
  description = "The snapshot's name"
  value       = azurerm_snapshot.main.name
}
