# Overview

The **AzureDataFactoryLinkedService** component deploys one linked service inside an Azure Data Factory (AzureDataFactory) -- a saved connection in the factory's address book: where an external system lives and how pipelines, datasets, and data flows authenticate to it. One kind covers all 23 connection types azurerm models as first-class resources, plus the raw-JSON custom form for every other Data Factory connector.

## Purpose

- **Connections as configuration**: the credentials posture that production depends on lives in the manifest, reviewed and versioned like everything else -- not clicked together in a portal.
- **One kind, 23 types**: the variant block declares the connection type; Azure stores every type in one factory-scoped namespace, and so does the catalog.
- **Secrets steered to the safe patterns**: every variant that offers managed identity or Key-Vault-sourced secrets models them first-class -- a connection can carry ZERO secret material and still authenticate.

## Key Features

- Full azurerm v5 surface across ALL TWENTY-THREE provider linked service resources: storage (blob, file share, table, Data Lake Gen2), databases (Azure SQL, SQL Server, SQL Managed Instance, Synapse, PostgreSQL, MySQL, two Cosmos DB forms, Snowflake, Kusto, ODBC), services (Key Vault, AI Search, Azure Function, Databricks), protocol endpoints (web, OData, SFTP), and the raw-JSON custom escape hatch.
- Chart-ready: `data_factory_id` defaults its reference to AzureDataFactory's ID output; the key vault variant's `key_vault_id` to AzureKeyVault's; blob's `service_endpoint` and Gen2's `url` to AzureStorageAccount's endpoint outputs; search's `url` to AzureSearchService's endpoint output; and every Key-Vault-sourced secret block defaults to ANOTHER AzureDataFactoryLinkedService's name output (the Key Vault connection).
- Secure by default: every secret-bearing field is marked sensitive; per-variant validation mirrors the provider's own authentication matrices exactly (one connection form, one identity, complete service principals).

## Use Cases

- **The factory's secrets backbone**: one Key Vault connection, then every other connection references its secrets by name -- rotation happens in the vault, never in manifests.
- **Lakehouse plumbing**: blob / Data Lake Gen2 connections on managed identity feed datasets and data flows without a single stored key.
- **Warehouse loads**: Azure SQL / Synapse / Snowflake connections with Key-Vault-held connection strings behind ingestion pipelines.
- **Partner exchange**: an SFTP connection with a Key-Vault-held private key lands partner files into the lake on a schedule.
- **Anything else Data Factory speaks**: the custom form carries any connector type (Salesforce, SAP, Oracle, REST, ...) as typed JSON.

## Future Enhancements

- Datasets and integration runtimes arrive as their own kinds, completing the Data Factory family (the shared `integration_runtime_name` reference upgrades in place when the runtime kind ships).
