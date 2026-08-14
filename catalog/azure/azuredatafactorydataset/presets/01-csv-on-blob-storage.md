# CSV on Blob Storage

This preset creates the most common dataset shape: delimited text (CSV) files in a blob storage container, with a run-time-parameterized folder path so one dataset serves every day's drop.

## When to Use

- Raw file drops landing in blob storage as CSV (the classic ingestion entry point)
- Any feed partitioned by date or run -- the parameterized path covers all partitions with ONE dataset

## Key Configuration Choices

- **`linkedServiceName` by reference** -- wires an AzureDataFactoryLinkedService's name output; the connection must be a blob storage (or Data Lake Gen2) type for this location shape
- **`dynamicPathEnabled: true`** -- makes `path` a Data Factory expression (`@{dataset().runDate}`), evaluated per run from the `parameters` map; pipelines override the value per activity
- **`firstRowAsHeader: true`** -- most exported CSVs carry a header row; omit (false) for headerless files
- **Delimiter defaults left to Azure** -- omit `columnDelimiter`/`quoteCharacter`/`escapeCharacter` for the standard `,` / `"` / `\`; set them only for non-standard files

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-data-factory>` | The Planton name of your `AzureDataFactory` resource | Planton console (or replace `valueFrom` with `value:` and the factory's ARM ID) |
| `<your-blob-linked-service>` | The Planton name of your blob storage `AzureDataFactoryLinkedService` | Planton console (or replace `valueFrom` with `value:` and the connection's name) |

## Related Presets

- **SQL Table** -- the table-shaped sibling for copy activity sinks and sources.
- **Parquet on Data Lake** -- the curated-side format for the same files after processing.
