# Custom Event

This preset creates the standard event-driven trigger: fire the load pipeline when an upstream system publishes a readiness event to an Event Grid custom topic.

## When to Use

- Orchestration driven by upstream systems' own signals ("batch ready", "export complete") rather than clocks or file arrivals
- Decoupling producers from pipelines: upstream publishes to the topic; this trigger subscribes

## Key Configuration Choices

- **`activated: false`** -- deploy stopped, verify the pipeline's Debug run, then flip to true
- **`events` is your topic's own vocabulary** -- free-form strings matched against the published events' eventType field
- **`subjectBeginsWith`** -- narrows firings to one subject subtree when many producers share the topic
- **Event payload as parameters** -- `@triggerBody().event.data.*` passes the event's own fields to the run

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-data-factory>` | The Planton name of your `AzureDataFactory` resource | Planton console (or replace `valueFrom` with `value:` and the factory's ARM ID) |
| `<your-eventgrid-topic>` | The Planton name of the `AzureEventgridTopic` publishing the events | Planton console (or replace `valueFrom` with `value:` and the topic's ARM ID) |
| `<your-pipeline>` | The Planton name of the `AzureDataFactoryPipeline` to run per event | Planton console (or replace `valueFrom` with `value:` and the pipeline's name) |
| `events` | Your topic's event-type vocabulary | The publishing system's documentation |

## Related Presets

- **Blob Landing Event** -- file-arrival processing.
- **Hourly Tumbling Window** -- windowed incremental loads.
