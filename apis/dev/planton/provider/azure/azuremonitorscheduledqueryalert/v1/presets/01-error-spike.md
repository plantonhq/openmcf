# Error Spike Alert (Row Count)

This preset creates the bread-and-butter log alert: count the error rows a KQL query returns and fire when the count crosses a line -- here, more than 10 application exceptions in 10 minutes.

## When to Use

- Every service whose logs land in the workspace (via diagnostic settings or Application Insights)
- Any "N bad events in M minutes" condition -- swap the query and threshold

## Key Configuration Choices

- **COUNT aggregation** -- the row-count evaluation style; the query's result rows ARE the signal
- **Auto-mitigation on** -- the alert resolves itself when errors quiet down, keeping the alert list honest
- **Region matches the workspace** -- Azure requires the rule in the same region as the resource it queries
- **Email subject override** -- the action block can brand notifications per rule

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-observability-rg` | Resource group holding the rule | `AzureResourceGroup` status outputs |
| `my-platform-logs` | The queried workspace | `AzureLogAnalyticsWorkspace` status outputs |
| `my-platform-oncall` | The action group to notify | `AzureMonitorActionGroup` status outputs |
| The KQL query | Your error condition | The workspace's Logs blade -- test it there first |

## Related Presets

- **02-latency-threshold** -- The metric-measurement style over a numeric column
- **03-missing-heartbeat** -- The absence-of-data pattern
