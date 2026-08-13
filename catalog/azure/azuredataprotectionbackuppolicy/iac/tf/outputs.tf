# Exactly one variant resource exists per deployment (the spec's
# exactly-one CEL); the outputs coalesce across the six.
output "backup_policy_id" {
  description = "The Azure Resource Manager ID of the policy -- what backup instances bind their policy by"
  value = try(
    azurerm_data_protection_backup_policy_blob_storage.main[0].id,
    azurerm_data_protection_backup_policy_disk.main[0].id,
    azurerm_data_protection_backup_policy_kubernetes_cluster.main[0].id,
    azurerm_data_protection_backup_policy_mysql_flexible_server.main[0].id,
    azurerm_data_protection_backup_policy_postgresql_flexible_server.main[0].id,
    azurerm_data_protection_backup_policy_data_lake_storage.main[0].id,
    ""
  )
}

output "backup_policy_name" {
  description = "The policy's name, unique on its vault"
  value       = var.spec.name
}
