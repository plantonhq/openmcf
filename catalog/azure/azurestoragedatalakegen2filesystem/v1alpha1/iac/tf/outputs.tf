# ADLS filesystems surface in ARM as blob containers, so the
# filesystem's management/RBAC identity is the container-proxy ID --
# what data-plane role assignments (Storage Blob Data
# Reader/Contributor/Owner) scope to. Constructed from the account ID +
# name (identically on both engines) because the provider's own
# resource id is a data-plane dfs URL nothing management-grain can
# consume.
output "filesystem_id" {
  description = "The Azure Resource Manager (container-proxy) ID of the filesystem"
  value       = local.filesystem_arm_id
}

# The container segment of every abfss:// and dfs URL.
output "filesystem_name" {
  description = "The name of the filesystem"
  value       = azurerm_storage_data_lake_gen2_filesystem.main.name
}

# Parsed from the account ARM ID -- consumers frequently need the
# account/filesystem pair (abfss URLs, Spark configs, mounts).
output "storage_account_name" {
  description = "The name of the storage account the filesystem lives in"
  value       = local.storage_account_name
}
