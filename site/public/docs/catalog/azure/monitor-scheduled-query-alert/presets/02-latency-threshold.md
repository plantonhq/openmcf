---
title: "Latency Threshold Alert (Metric Measurement, Per-Dimension)"
description: "This preset creates a metric-measurement log alert: the query computes average request duration per service role, and each role alerts independently when its average exceeds 500ms in 3 of the last 5..."
type: "preset"
rank: "02"
presetSlug: "02-latency-threshold"
componentSlug: "monitor-scheduled-query-alert"
componentTitle: "Monitor Scheduled Query Alert"
provider: "azure"
icon: "package"
order: 2
---

# Latency Threshold Alert (Metric Measurement, Per-Dimension)

This preset creates a metric-measurement log alert: the query computes average request duration per service role, and each role alerts independently when its average exceeds 500ms in 3 of the last 5 evaluation periods.

## When to Use

- Latency/percentile conditions computed from request logs
- One rule covering many services -- the dimension split alerts per `AppRoleName` value

## Key Configuration Choices

- **AVERAGE over `metricMeasureColumn`** -- non-COUNT aggregations compute over the named numeric column; the query must project it (`avg_DurationMs` is KQL's auto-name for `avg(DurationMs)`)
- **`"*"` INCLUDE dimension** -- splits the evaluation across every role value; each alerts (and resolves) independently
- **3 of 5 failing periods** -- the flap damper: sustained degradation pages, a single slow window does not

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-observability-rg` | Resource group holding the rule | `AzureResourceGroup` status outputs |
| `my-platform-logs` | The queried workspace | `AzureLogAnalyticsWorkspace` status outputs |
| `my-platform-oncall` | The action group to notify | `AzureMonitorActionGroup` status outputs |

## Related Presets

- **01-error-spike** -- The row-count style
- **03-missing-heartbeat** -- The absence-of-data pattern
