# Daily Schedule

This preset creates the standard nightly trigger: daily at 02:00 UTC, deployed STOPPED so nothing runs until you flip it on.

## When to Use

- Nightly/daily pipeline runs at a fixed clock time
- Any genuinely clock-shaped work (reports, cache warms) -- for windowed incremental loads, prefer the Hourly Tumbling Window preset

## Key Configuration Choices

- **`activated: false`** -- deploy stopped, verify the pipeline's Debug run, then flip to true and re-apply; a started trigger runs real, billed pipeline runs
- **Explicit `startTime`** -- an omitted start time starts the recurrence from the moment of deployment
- **`recurrenceSchedule.hours/minutes`** -- narrows the Day frequency to 02:00; the fire time passes to the pipeline via `@trigger().scheduledTime`

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-data-factory>` | The Planton name of your `AzureDataFactory` resource | Planton console (or replace `valueFrom` with `value:` and the factory's ARM ID) |
| `<your-pipeline>` | The Planton name of the `AzureDataFactoryPipeline` to run | Planton console (or replace `valueFrom` with `value:` and the pipeline's name) |
| `startTime` | When the recurrence begins (UTC) | Your go-live date |

## Related Presets

- **Hourly Tumbling Window** -- windowed incremental loads with backfill.
- **Blob Landing Event** -- file-arrival processing.
- **Custom Event** -- event-driven orchestration from an Event Grid topic.
