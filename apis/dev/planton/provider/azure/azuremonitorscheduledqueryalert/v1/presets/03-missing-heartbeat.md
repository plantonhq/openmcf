# Missing Heartbeat Alert (Absence of Data)

This preset creates the absence-of-data alert: fire when heartbeat rows drop BELOW a floor -- the pattern that catches dead agents, stopped VMs, and silently broken ingestion pipelines, which no error-based alert can see.

## When to Use

- Fleets of VMs or agents reporting Heartbeat to the workspace
- Any "silence is the failure" condition -- swap Heartbeat for your own liveness table

## Key Configuration Choices

- **COUNT + LESS_THAN** -- the inversion that makes absence alertable
- **Mute instead of auto-mitigate** (`muteActionsAfterAlertDuration: PT1H`) -- a dead agent cannot report recovery, so auto-resolve would silently close real incidents; muting caps notification noise instead (the two are mutually exclusive)
- **Window wider than frequency** (15m window, 5m cadence) -- tolerates one late heartbeat without paging

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-observability-rg` | Resource group holding the rule | `AzureResourceGroup` status outputs |
| `my-platform-logs` | The queried workspace | `AzureLogAnalyticsWorkspace` status outputs |
| `my-platform-oncall` | The action group to page | `AzureMonitorActionGroup` status outputs |
| `threshold: 1` | The heartbeat floor | Your fleet size and reporting interval |

## Related Presets

- **01-error-spike** -- Presence-of-errors alerting
- **02-latency-threshold** -- Metric-measurement alerting
