# AzureDataFactoryLinkedService Pulumi Module

## Overview

Creates one linked service inside an Azure Data Factory -- a saved connection telling pipelines, datasets, and data flows where an external system lives and how to authenticate to it. The spec's variant block (exactly one of 23) selects the connection type; the module creates the matching classic-SDK resource.

## Resources Created

Exactly one of the 23 `datafactory.LinkedService*` resources (or `datafactory.LinkedCustomService` for the raw-JSON custom form): Azure Blob Storage, Databricks, Azure Files, Azure Function, Azure AI Search, Azure SQL Database, Table Storage, Cosmos DB (SQL API), Cosmos DB for MongoDB, custom, Data Lake Storage Gen2, Key Vault, Kusto (Azure Data Explorer), MySQL, OData, ODBC, PostgreSQL, SFTP, Snowflake, SQL Managed Instance, SQL Server, Synapse, and web.

## Outputs

- `linked_service_id` -- the linked service's ARM resource ID (`{factory_id}/linkedservices/{name}`, the same shape for all 23 types)
- `linked_service_name` -- the linked service's name (what datasets, data flows, and other linked services' Key-Vault-sourced secret references resolve against)

## Behavior Notes

- **PARITY-EXCEPTION (mysql)**: the classic SDK (pulumi-azure v6.38.0) does not carry the MySQL resource's `driver_version` argument -- this engine cannot send it, so Azure applies its own driver default; the Terraform module sends the effective value (default `V2`). Remove when the SDK catches up.
- **ENGINE-SHAPE (sftp)**: the bridged SDK pluralizes the SFTP password's Key Vault block name (`KeyVaultPasswords`, an array the module fills with the spec's single entry) -- a name difference only; both engines write the same ARM object.
- **Secrets never read back**: Azure stores connection strings, passwords, and keys as hidden secure strings; Pulumi state keeps the configured value.
- **Key-Vault-sourced secret blocks** name a KEY VAULT linked service in this same factory plus a secret name -- Data Factory resolves the secret at run time.
- **SQL Managed Instance carries no `additional_properties`** -- the one linked service the provider models without that argument.
- **`use_managed_identity` is sent only when true** on all four variants that carry it (blob storage, Azure SQL Database, Data Lake Gen2, Kusto) -- the provider's constraint checks fire on the argument's PRESENCE, so an explicit false alongside the alternative identity is rejected; unset means false anyway.
- **No tags**: linked services are ARM sub-resources of the factory and expose none.
