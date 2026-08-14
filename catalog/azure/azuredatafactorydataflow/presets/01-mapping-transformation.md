# Mapping Transformation

This preset creates the standard transformation shell: one source, one sink, and a pass-through script to replace with your Studio-authored flow.

## When to Use

- Cleaning/conforming raw landing data into curated tables
- Any transformation a pipeline's Execute Data Flow activity will run

## Key Configuration Choices

- **Stream names line up** -- the `sources`/`sinks` entries name the same streams the script names (`rawOrders`, `curatedOrders`); a mismatch fails at deploy time
- **Linked-service bindings** -- the endpoints bind to linked services by name; datasets work too (`dataset:` instead of `linkedService:`)
- **Placeholder pass-through script** -- author the real flow in the Data Factory Studio and paste its "Script" view over the placeholder (the catalog deliberately does not re-model the script language)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-data-factory>` | The Planton name of your `AzureDataFactory` resource | Planton console (or replace `valueFrom` with `value:` and the factory's ARM ID) |
| `<your-source-linked-service>` | The factory-scoped linked service the source reads through | Data Factory Studio -> Manage -> Linked services |
| `<your-sink-linked-service>` | The factory-scoped linked service the sink writes through | Data Factory Studio -> Manage -> Linked services |
| `script` | Your flow's real script | Data Factory Studio -> your data flow -> Script view |

## Related Presets

- **Reusable Flowlet** -- shared cleanup logic embedded by every ingest flow.
