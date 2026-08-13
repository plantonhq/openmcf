# AzureDataFactoryDataset Pulumi Module

## Overview

Creates one dataset inside an Azure Data Factory -- a named view of data telling pipelines and data flows what the data looks like and where it sits within a system a linked service already connects to. The spec's variant block (exactly one of 13) selects the data shape, and the module dispatches to the matching Pulumi resource.

## Resources Created

Exactly one of the 13 `datafactory.Dataset*` resources (or `datafactory.CustomDataset` for the raw-JSON custom form): Azure Blob, Azure SQL Table, binary, Cosmos DB (SQL API), custom, delimited text (CSV), HTTP, JSON, MySQL, Parquet, PostgreSQL, Snowflake, and SQL Server Table.

## Module Structure

- `module/main.go` -- provider setup and the 13-way variant dispatch
- `module/shared.go` -- the shared optional fields every dataset resource carries (sent only when set)
- `module/variants_files.go` -- the file-format builders (azure blob, binary, delimited text, HTTP, JSON, Parquet)
- `module/variants_tables.go` -- the table-form builders (Azure SQL, Cosmos DB, MySQL, PostgreSQL, Snowflake, SQL Server) and the custom form
- `module/outputs.go` -- exported output names

## Outputs

- `dataset_id` -- the dataset's ARM resource ID (`{factory_id}/datasets/{name}`, the same shape for all 13 types)
- `dataset_name` -- the dataset's name (what pipelines and data flows resolve against)

## Behavior Notes

- **All 13 shapes share one name namespace** inside the factory, so switching variant blocks replaces the dataset (a different resource is created at the same ARM address).
- **The linked service reference has three wire forms**: 11 variants send the shared `linked_service_name`; `azure_sql_table` sends its own `linked_service_id` (must belong to the same factory -- Azure rejects cross-factory references); `custom` sends a `linked_service` block that can also carry per-use parameter values.
- **Omitted parse settings fall back to the provider's own defaults** on the delimited text form: `,` column delimiter, `"` quote character, `\` escape character, first row not a header.
- **The `dynamic_*_enabled` flags mark the paired value as a Data Factory expression** (evaluated at run time, e.g. `@{dataset().runDate}` against dataset parameters) instead of a literal string.
- **Parquet's `compression_level` argument is deliberately not modeled**: it exists in the provider's schema but the provider's create/update code never reads it, so setting it would silently do nothing (recorded in `iac/provider-parity.yaml`).
- **No tags**: datasets are ARM sub-resources of the factory and expose none.
