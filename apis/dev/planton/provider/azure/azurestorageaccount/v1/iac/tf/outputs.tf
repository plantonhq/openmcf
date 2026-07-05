# The ARM id other resources reference: AzureStorageContainer's parent,
# data-plane role-assignment scopes, private endpoints.
output "storage_account_id" {
  description = "The Azure Resource Manager ID of the storage account"
  value       = azurerm_storage_account.main.id
}

# The DNS prefix of every endpoint, and what app-hosting kinds bind to.
output "storage_account_name" {
  description = "The name of the storage account"
  value       = azurerm_storage_account.main.name
}

output "resource_group_name" {
  description = "The resource group the account lives in"
  value       = azurerm_storage_account.main.resource_group_name
}

output "primary_blob_endpoint" {
  description = "The primary blob service endpoint URL"
  value       = azurerm_storage_account.main.primary_blob_endpoint
}

# The bare hostname a CDN/Front Door origin or custom-domain CNAME
# points at.
output "primary_blob_host" {
  description = "The primary blob service hostname"
  value       = azurerm_storage_account.main.primary_blob_host
}

output "primary_queue_endpoint" {
  description = "The primary queue service endpoint URL"
  value       = azurerm_storage_account.main.primary_queue_endpoint
}

output "primary_table_endpoint" {
  description = "The primary table service endpoint URL"
  value       = azurerm_storage_account.main.primary_table_endpoint
}

output "primary_file_endpoint" {
  description = "The primary file service endpoint URL"
  value       = azurerm_storage_account.main.primary_file_endpoint
}

output "primary_dfs_endpoint" {
  description = "The primary Data Lake Storage Gen2 endpoint URL"
  value       = azurerm_storage_account.main.primary_dfs_endpoint
}

output "primary_web_endpoint" {
  description = "The primary static-website endpoint URL"
  value       = azurerm_storage_account.main.primary_web_endpoint
}

output "primary_web_host" {
  description = "The primary static-website hostname"
  value       = azurerm_storage_account.main.primary_web_host
}

# The secondary endpoints are live only on the read-access replication
# types (RA_GRS / RA_GZRS) -- the read-only mirror in the paired region.
output "secondary_blob_endpoint" {
  description = "The secondary (read-only, paired-region) blob endpoint URL"
  value       = azurerm_storage_account.main.secondary_blob_endpoint
}

output "secondary_queue_endpoint" {
  description = "The secondary queue endpoint URL"
  value       = azurerm_storage_account.main.secondary_queue_endpoint
}

output "secondary_table_endpoint" {
  description = "The secondary table endpoint URL"
  value       = azurerm_storage_account.main.secondary_table_endpoint
}

output "secondary_file_endpoint" {
  description = "The secondary file endpoint URL"
  value       = azurerm_storage_account.main.secondary_file_endpoint
}

output "secondary_dfs_endpoint" {
  description = "The secondary Data Lake Storage Gen2 endpoint URL"
  value       = azurerm_storage_account.main.secondary_dfs_endpoint
}

output "secondary_web_endpoint" {
  description = "The secondary static-website endpoint URL"
  value       = azurerm_storage_account.main.secondary_web_endpoint
}

# Static credential material -- prefer Entra data-plane roles scoped to
# storage_account_id; reference the keys only where a consumer genuinely
# requires key auth.
output "primary_access_key" {
  description = "The account's first shared access key"
  value       = azurerm_storage_account.main.primary_access_key
  sensitive   = true
}

output "secondary_access_key" {
  description = "The account's second shared access key (for zero-downtime rotation)"
  value       = azurerm_storage_account.main.secondary_access_key
  sensitive   = true
}

output "primary_connection_string" {
  description = "Connection string carrying the primary access key"
  value       = azurerm_storage_account.main.primary_connection_string
  sensitive   = true
}

output "secondary_connection_string" {
  description = "Connection string carrying the secondary access key"
  value       = azurerm_storage_account.main.secondary_connection_string
  sensitive   = true
}

output "primary_blob_connection_string" {
  description = "Blob-service-only connection string on the primary key"
  value       = azurerm_storage_account.main.primary_blob_connection_string
  sensitive   = true
}

output "secondary_blob_connection_string" {
  description = "Blob-service-only connection string on the secondary key"
  value       = azurerm_storage_account.main.secondary_blob_connection_string
  sensitive   = true
}

# Populated only when the identity type includes SYSTEM_ASSIGNED; grant
# this principal roles (e.g. Key Vault Crypto User) to let the account
# act on other resources.
output "identity_principal_id" {
  description = "The principal ID of the account's system-assigned identity"
  value       = try(azurerm_storage_account.main.identity[0].principal_id, "")
}
