# Exactly one variant resource exists (spec CEL enforces exactly one
# variant block), so each try() resolves against the one that was
# created.

output "datastore_id" {
  description = "The Azure Resource Manager ID of the datastore"
  value = try(
    azurerm_machine_learning_datastore_blobstorage.main[0].id,
    azurerm_machine_learning_datastore_datalake_gen2.main[0].id,
    azurerm_machine_learning_datastore_fileshare.main[0].id,
    ""
  )
}

output "datastore_name" {
  description = "The datastore's name -- what jobs and data assets reference within the workspace"
  value = try(
    azurerm_machine_learning_datastore_blobstorage.main[0].name,
    azurerm_machine_learning_datastore_datalake_gen2.main[0].name,
    azurerm_machine_learning_datastore_fileshare.main[0].name,
    ""
  )
}

output "is_default" {
  description = "Whether this datastore is the workspace's default (settable only on the blob variant; read back on the others)"
  value = try(
    azurerm_machine_learning_datastore_blobstorage.main[0].is_default,
    azurerm_machine_learning_datastore_datalake_gen2.main[0].is_default,
    azurerm_machine_learning_datastore_fileshare.main[0].is_default,
    false
  )
}
