# All Logs and Metrics to a Workspace

This preset routes everything a resource emits -- every log category (via the `allLogs` group) and all metrics -- to a Log Analytics Workspace in the modern resource-specific table layout. It is the wiring that makes a resource observable: queryable with KQL, alertable with scheduled query rules.

## When to Use

- The default diagnostic wiring for any production resource (the example targets a Key Vault; swap the target reference for any kind)
- Environments standardizing on "everything lands in the workspace"

## Key Configuration Choices

- **`allLogs` category group** -- tracks new log categories automatically as Azure adds them; pick single categories instead when volume costs demand selectivity
- **DEDICATED destination type** -- resource-specific tables with typed columns; prefer it wherever the resource type supports it
- **Polymorphic target** -- `targetResourceId` carries no default kind; reference any resource's `*_id` output explicitly with valueFrom

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-app-vault` | The resource whose telemetry is routed | Any kind's status outputs (`*_id`) |
| `my-platform-logs` | The destination workspace | `AzureLogAnalyticsWorkspace` status outputs |

## Related Presets

- **02-archive-to-storage** -- Cheap long-term archival instead of (or alongside) the workspace
- **03-stream-to-siem** -- Stream to an Event Hub for an external SIEM
