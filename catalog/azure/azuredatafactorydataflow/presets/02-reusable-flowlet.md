# Reusable Flowlet

This preset creates a flowlet -- the reusable transformation snippet other data flows embed by reference -- with a placeholder scrubbing script to replace with your shared logic.

## When to Use

- Organization-wide cleanup/conformance rules (PII scrubbing, column standardization) that every ingest flow should apply identically
- Any logic you would otherwise copy-paste between data flow scripts

## Key Configuration Choices

- **`flowlet: true`** -- creates the flowlet form; sources and sinks stay omitted because the embedding flow supplies them
- **The name is the contract** -- embedding flows reference this flowlet by name; treat it like a package name (stable; version by suffix when logic changes shape)
- **Placeholder script** -- author the real flowlet in the Data Factory Studio and paste its Script view over the placeholder

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-data-factory>` | The Planton name of your `AzureDataFactory` resource | Planton console (or replace `valueFrom` with `value:` and the factory's ARM ID) |
| `script` | Your flowlet's real script | Data Factory Studio -> your flowlet -> Script view |

## Related Presets

- **Mapping Transformation** -- the runnable flow that embeds this flowlet at its source.
