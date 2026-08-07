# Static Threshold Alert

This preset creates the classic metric alert: a fixed threshold on a platform metric (storage availability below 99.9%, averaged over 15 minutes), evaluated every minute, paging the on-call action group at severity 1.

## When to Use

- SLO-backed conditions with a known line (availability floors, latency ceilings, saturation limits)
- Any resource kind -- swap the scope reference and the metric namespace/name

## Key Configuration Choices

- **AVERAGE over PT15M, evaluated PT1M** -- smooths blips while alerting within a minute of a sustained breach
- **Severity 1** -- availability is a paging condition; the default severity 3 is for informational rules
- **Auto-mitigate (default)** -- the alert resolves itself when availability recovers; no manual closing
- **Description is the runbook seam** -- it is delivered with every notification; write it for the person paged at 3 AM

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-observability-rg` | Resource group holding the rule | `AzureResourceGroup` status outputs |
| `my-app-storage` | The monitored resource | Any kind's status outputs (`*_id`) |
| `my-platform-oncall` | The action group to page | `AzureMonitorActionGroup` status outputs |

## Related Presets

- **02-dynamic-anomaly** -- Machine-learned thresholds when no fixed line exists
- **03-webtest-availability** -- Availability-test-driven alerting for user-facing endpoints
