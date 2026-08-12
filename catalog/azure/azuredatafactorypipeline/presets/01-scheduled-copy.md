# Scheduled Copy

This preset creates the standard incremental-ingestion pipeline shell: window-parameterized, single-concurrency so runs never overlap, with a placeholder activity to replace with your Studio-authored copy.

## When to Use

- Daily/hourly loads from a source system into the lakehouse or warehouse
- Any pipeline a schedule trigger will run with a per-run window parameter

## Key Configuration Choices

- **Concurrency 1** -- incremental loads must never overlap; queued runs wait rather than fail
- **`windowStart` parameter** -- the trigger supplies the window per run; the pipeline definition stays environment-agnostic
- **Placeholder Wait activity** -- author the real copy in the Data Factory Studio and paste its Code view's "activities" array over the placeholder (the catalog deliberately does not re-model Azure's activity schema)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-data-factory>` | The Planton name of your `AzureDataFactory` resource | Planton console (or replace `valueFrom` with `value:` and the factory's ARM ID) |
| `activitiesJson` | Your pipeline's real activities | Data Factory Studio -> your pipeline -> Code view -> the "activities" array |

## Related Presets

- None yet -- triggers and datasets arrive as their own kinds.
