# SQL Table

This preset creates an Azure SQL Database table dataset -- the shape copy activities use as a source or sink when moving data into or out of the warehouse.

## When to Use

- Loading files into Azure SQL (this dataset is the SINK of a copy activity)
- Extracting tables to the lake (this dataset is the SOURCE)
- Leave `schema`/`table` undeclared when pipelines supply a query instead of a table name

## Key Configuration Choices

- **`linkedServiceId` by reference, not name** -- the Azure SQL table shape is the one dataset that references its connection by ARM ID; it wires AzureDataFactoryLinkedService's `linked_service_id` output, and Azure requires the connection to live in the SAME factory
- **No `linkedServiceName` at the root** -- this shape carries its own reference; setting both is rejected by validation
- **Declared columns omitted** -- Data Factory infers the table's columns; add `schemaColumn` only when downstream mappings need a stable declared contract

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-data-factory>` | The Planton name of your `AzureDataFactory` resource | Planton console (or replace `valueFrom` with `value:` and the factory's ARM ID) |
| `<your-sql-linked-service>` | The Planton name of your Azure SQL `AzureDataFactoryLinkedService` | Planton console (or replace `valueFrom` with `value:` and the connection's ARM ID) |

## Related Presets

- **CSV on Blob Storage** -- the file-shaped source that typically feeds this table.
- **Parquet on Data Lake** -- the analytical format on the other side of the warehouse.
