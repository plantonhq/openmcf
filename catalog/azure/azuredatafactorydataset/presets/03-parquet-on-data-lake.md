# Parquet on Data Lake

This preset creates a Parquet dataset over Azure Data Lake Storage Gen2 -- the curated-side shape of a lakehouse, where processed data lands in an analytical format.

## When to Use

- The output side of data flows and copy activities that transform raw files into columnar data
- Any feed analytical engines (Synapse, Databricks, Fabric) read directly from the lake

## Key Configuration Choices

- **`azureBlobFsLocation`** -- the Data Lake Gen2 location shape (file system + path), for storage accounts with the hierarchical namespace; the linked service must be a Gen2 (or blob) connection type
- **`compressionCodec: snappy`** -- Parquet's native pairing (fast, splittable); note the compression LEVEL argument is deliberately absent -- Azure ignores it for Parquet
- **`folder: curated`** -- a Studio display folder only; it groups datasets in the authoring UI and has no effect on the data

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-data-factory>` | The Planton name of your `AzureDataFactory` resource | Planton console (or replace `valueFrom` with `value:` and the factory's ARM ID) |
| `<your-gen2-linked-service>` | The Planton name of your Data Lake Gen2 `AzureDataFactoryLinkedService` | Planton console (or replace `valueFrom` with `value:` and the connection's name) |

## Related Presets

- **CSV on Blob Storage** -- the raw-side format these files typically start as.
- **SQL Table** -- the warehouse shape for serving the same data relationally.
