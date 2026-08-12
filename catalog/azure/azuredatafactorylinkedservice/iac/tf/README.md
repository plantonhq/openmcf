# AzureDataFactoryLinkedService Terraform Module

## Overview

Creates one linked service inside an Azure Data Factory -- a saved connection telling pipelines, datasets, and data flows where an external system lives and how to authenticate to it. The spec's variant block (exactly one of 23) selects the connection type.

## Resources Created

Exactly one of the 23 `azurerm_data_factory_linked_service_*` resources (or `azurerm_data_factory_linked_custom_service` for the raw-JSON custom form): Azure Blob Storage, Databricks, Azure Files, Azure Function, Azure AI Search, Azure SQL Database, Table Storage, Cosmos DB (SQL API), Cosmos DB for MongoDB, custom, Data Lake Storage Gen2, Key Vault, Kusto (Azure Data Explorer), MySQL, OData, ODBC, PostgreSQL, SFTP, Snowflake, SQL Managed Instance, SQL Server, Synapse, and web.

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureDataFactoryLinkedServiceSpec fields; the factory, Key Vault, storage endpoint, search endpoint, and Key-Vault-linked-service-name references arrive as resolved literals

## Outputs

- `linked_service_id` -- the linked service's ARM resource ID (`{factory_id}/linkedservices/{name}`, the same shape for all 23 types)
- `linked_service_name` -- the linked service's name (what datasets, data flows, and other linked services' Key-Vault-sourced secret references resolve against)

## Behavior Notes

- **All 23 types share one name namespace** inside the factory, so switching variant blocks replaces the linked service (a different provider resource is created at the same ARM address).
- **Secrets never read back**: Azure stores connection strings, passwords, and keys as hidden secure strings and returns them masked or not at all -- Terraform state keeps the configured value; drift on secret fields is suppressed by the provider where the API supports comparison.
- **Key-Vault-sourced secret blocks** (`key_vault_password`, `key_vault_connection_string`, and friends) name a KEY VAULT linked service in this same factory plus a secret name -- Data Factory resolves the secret at run time, so nothing secret sits in configuration.
- **The custom form's integration runtime travels as a block**: its `name` comes from the same root `integration_runtime_name` every other variant sends as a plain argument; only the per-use `parameters` are variant-scoped.
- **SQL Managed Instance carries no `additional_properties`** -- the one linked service the provider models without that argument.
- **`use_managed_identity` is sent only when true** on all four variants that carry it (blob storage, Azure SQL Database, Data Lake Gen2, Kusto) -- the provider's constraint checks fire on the argument's PRESENCE, so an explicit false alongside the alternative identity is rejected; unset means false anyway.
- **No tags**: linked services are ARM sub-resources of the factory and expose none.
