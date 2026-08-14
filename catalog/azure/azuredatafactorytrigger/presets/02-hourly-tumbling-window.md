# Hourly Tumbling Window

This preset creates the standard incremental-load trigger: one contiguous hour per run, each run receiving its window's bounds, with retry and a deliberate concurrency cap.

## When to Use

- Incremental data loads where the pipeline's work is a function of a time range
- Backfills: point `startTime` at history and every window since then runs, rate-limited

## Key Configuration Choices

- **`activated: false`** -- deploy stopped; a started trigger with a past `startTime` begins BACKFILLING immediately
- **`maxConcurrency: 4`** -- Azure's default is 50 parallel windows; cap it at what the sink tables tolerate
- **`retry`** -- per-window retry so transient failures do not strand single windows in a sea of green ones
- **Window bounds as parameters** -- `@trigger().outputs.windowStartTime/windowEndTime` make the pipeline's work a pure function of its window

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-data-factory>` | The Planton name of your `AzureDataFactory` resource | Planton console (or replace `valueFrom` with `value:` and the factory's ARM ID) |
| `<your-pipeline>` | The Planton name of the `AzureDataFactoryPipeline` to run per window | Planton console (or replace `valueFrom` with `value:` and the pipeline's name) |
| `startTime` | The FIRST window's start (UTC) -- history backfills from here | Your data's start date, chosen deliberately |

## Related Presets

- **Daily Schedule** -- clock-shaped work without windows.
- **Blob Landing Event** -- file-arrival processing.
