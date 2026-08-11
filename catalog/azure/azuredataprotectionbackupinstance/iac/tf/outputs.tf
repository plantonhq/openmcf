# Exactly one variant resource exists per deployment (the spec's
# exactly-one CEL); the outputs coalesce across the six.
output "backup_instance_id" {
  description = "The Azure Resource Manager ID of the backup instance"
  value = try(
    azurerm_data_protection_backup_instance_blob_storage.main[0].id,
    azurerm_data_protection_backup_instance_disk.main[0].id,
    azurerm_data_protection_backup_instance_kubernetes_cluster.main[0].id,
    azurerm_data_protection_backup_instance_mysql_flexible_server.main[0].id,
    azurerm_data_protection_backup_instance_postgresql_flexible_server.main[0].id,
    azurerm_data_protection_backup_instance_data_lake_storage.main[0].id,
    ""
  )
}

output "backup_instance_name" {
  description = "The instance's name, unique on its vault"
  value       = var.spec.name
}
