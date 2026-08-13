# Exactly one of the 13 resources exists (the spec's variant block
# decides); all share the {factory_id}/datasets/{name} ID shape.
output "dataset_id" {
  description = "The Azure Resource Manager ID of the dataset ({factory_id}/datasets/{name})"
  value = try(
    azurerm_data_factory_dataset_azure_blob.main[0].id,
    azurerm_data_factory_dataset_azure_sql_table.main[0].id,
    azurerm_data_factory_dataset_binary.main[0].id,
    azurerm_data_factory_dataset_cosmosdb_sqlapi.main[0].id,
    azurerm_data_factory_custom_dataset.main[0].id,
    azurerm_data_factory_dataset_delimited_text.main[0].id,
    azurerm_data_factory_dataset_http.main[0].id,
    azurerm_data_factory_dataset_json.main[0].id,
    azurerm_data_factory_dataset_mysql.main[0].id,
    azurerm_data_factory_dataset_parquet.main[0].id,
    azurerm_data_factory_dataset_postgresql.main[0].id,
    azurerm_data_factory_dataset_snowflake.main[0].id,
    azurerm_data_factory_dataset_sql_server_table.main[0].id,
    null
  )
}

output "dataset_name" {
  description = "The dataset's name -- what pipelines and data flows resolve against"
  value = try(
    azurerm_data_factory_dataset_azure_blob.main[0].name,
    azurerm_data_factory_dataset_azure_sql_table.main[0].name,
    azurerm_data_factory_dataset_binary.main[0].name,
    azurerm_data_factory_dataset_cosmosdb_sqlapi.main[0].name,
    azurerm_data_factory_custom_dataset.main[0].name,
    azurerm_data_factory_dataset_delimited_text.main[0].name,
    azurerm_data_factory_dataset_http.main[0].name,
    azurerm_data_factory_dataset_json.main[0].name,
    azurerm_data_factory_dataset_mysql.main[0].name,
    azurerm_data_factory_dataset_parquet.main[0].name,
    azurerm_data_factory_dataset_postgresql.main[0].name,
    azurerm_data_factory_dataset_snowflake.main[0].name,
    azurerm_data_factory_dataset_sql_server_table.main[0].name,
    null
  )
}
