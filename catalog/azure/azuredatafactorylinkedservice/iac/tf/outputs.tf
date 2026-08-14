# Exactly one of the 23 resources exists (the spec's variant block
# decides); all share the {factory_id}/linkedservices/{name} ID shape.
output "linked_service_id" {
  description = "The Azure Resource Manager ID of the linked service ({factory_id}/linkedservices/{name})"
  value = try(
    azurerm_data_factory_linked_service_azure_blob_storage.main[0].id,
    azurerm_data_factory_linked_service_azure_databricks.main[0].id,
    azurerm_data_factory_linked_service_azure_file_storage.main[0].id,
    azurerm_data_factory_linked_service_azure_function.main[0].id,
    azurerm_data_factory_linked_service_azure_search.main[0].id,
    azurerm_data_factory_linked_service_azure_sql_database.main[0].id,
    azurerm_data_factory_linked_service_azure_table_storage.main[0].id,
    azurerm_data_factory_linked_service_cosmosdb.main[0].id,
    azurerm_data_factory_linked_service_cosmosdb_mongoapi.main[0].id,
    azurerm_data_factory_linked_custom_service.main[0].id,
    azurerm_data_factory_linked_service_data_lake_storage_gen2.main[0].id,
    azurerm_data_factory_linked_service_key_vault.main[0].id,
    azurerm_data_factory_linked_service_kusto.main[0].id,
    azurerm_data_factory_linked_service_mysql.main[0].id,
    azurerm_data_factory_linked_service_odata.main[0].id,
    azurerm_data_factory_linked_service_odbc.main[0].id,
    azurerm_data_factory_linked_service_postgresql.main[0].id,
    azurerm_data_factory_linked_service_sftp.main[0].id,
    azurerm_data_factory_linked_service_snowflake.main[0].id,
    azurerm_data_factory_linked_service_sql_managed_instance.main[0].id,
    azurerm_data_factory_linked_service_sql_server.main[0].id,
    azurerm_data_factory_linked_service_synapse.main[0].id,
    azurerm_data_factory_linked_service_web.main[0].id,
    null
  )
}

output "linked_service_name" {
  description = "The linked service's name -- what datasets, data flows, and other linked services' Key-Vault-sourced secret references resolve against"
  value = try(
    azurerm_data_factory_linked_service_azure_blob_storage.main[0].name,
    azurerm_data_factory_linked_service_azure_databricks.main[0].name,
    azurerm_data_factory_linked_service_azure_file_storage.main[0].name,
    azurerm_data_factory_linked_service_azure_function.main[0].name,
    azurerm_data_factory_linked_service_azure_search.main[0].name,
    azurerm_data_factory_linked_service_azure_sql_database.main[0].name,
    azurerm_data_factory_linked_service_azure_table_storage.main[0].name,
    azurerm_data_factory_linked_service_cosmosdb.main[0].name,
    azurerm_data_factory_linked_service_cosmosdb_mongoapi.main[0].name,
    azurerm_data_factory_linked_custom_service.main[0].name,
    azurerm_data_factory_linked_service_data_lake_storage_gen2.main[0].name,
    azurerm_data_factory_linked_service_key_vault.main[0].name,
    azurerm_data_factory_linked_service_kusto.main[0].name,
    azurerm_data_factory_linked_service_mysql.main[0].name,
    azurerm_data_factory_linked_service_odata.main[0].name,
    azurerm_data_factory_linked_service_odbc.main[0].name,
    azurerm_data_factory_linked_service_postgresql.main[0].name,
    azurerm_data_factory_linked_service_sftp.main[0].name,
    azurerm_data_factory_linked_service_snowflake.main[0].name,
    azurerm_data_factory_linked_service_sql_managed_instance.main[0].name,
    azurerm_data_factory_linked_service_sql_server.main[0].name,
    azurerm_data_factory_linked_service_synapse.main[0].name,
    azurerm_data_factory_linked_service_web.main[0].name,
    null
  )
}
