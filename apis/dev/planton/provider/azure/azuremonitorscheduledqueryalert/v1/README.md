# AzureMonitorScheduledQueryAlert

## Overview

`AzureMonitorScheduledQueryAlert` provisions an Azure Monitor scheduled query
alert rule (`Microsoft.Insights/scheduledQueryRules`) -- the log-search alert.
It runs a KQL query against a Log Analytics Workspace (or an Application
Insights resource) on a schedule, compares the result to a threshold, and
fires action groups when the condition holds.

This is the alerting half of the logging pipeline: diagnostic settings route
logs into the workspace, and query alerts watch what arrives -- error spikes,
security events, missing heartbeats, anything KQL can express.

## Key Features

- **Both evaluation styles** -- row count (COUNT over the query's result
  rows) and metric measurement (AVERAGE/MIN/MAX/TOTAL over a projected
  numeric column)
- **Dimension splits** -- each combination of the named columns' values
  evaluates and alerts independently (one rule covering many services)
- **Failing periods** -- the flap damper: require N of the last M evaluations
  to breach before firing
- **Noise controls** -- auto-mitigation (stateful resolve) XOR a mute
  duration for conditions that cannot report recovery
- **Managed identity** -- run the query as a system- or user-assigned
  identity when the workspace enforces Entra-only query access

## When to Use

- Every service whose logs land in the workspace -- the bread-and-butter
  "N bad events in M minutes" alert
- Absence-of-data conditions (dead agents, silent pipelines) no error-based
  alert can see
- Per-resource alerting via a projected resource-id column

## Spec Highlights

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `region` | string | Yes | - | Must match the queried resource's region |
| `resource_group` | StringValueOrRef | Yes | - | Resource group |
| `alert_name` | string | Yes | - | Rule name (ForceNew) |
| `scope` | StringValueOrRef | Yes | - | The queried workspace (default) or App Insights resource (ForceNew) |
| `severity` | int32 | No | 3 | 0 (critical) - 4 (verbose) |
| `evaluation_frequency` / `window_duration` | string | No | PT5M / PT5M | ISO 8601 vocabularies |
| `criteria[]` | message | Yes (min 1) | - | query + aggregation + operator + threshold (+ dimensions, failing periods) |
| `auto_mitigation_enabled` XOR `mute_actions_after_alert_duration` | bool / string | No | false / - | Mutually exclusive noise controls |
| `identity` | message | No | - | SYSTEM_ASSIGNED or USER_ASSIGNED |
| `action` | message | No | - | Action-group FKs + custom properties + email subject |
| `tags` | map | No | - | User tags |

## Outputs

| Output | Description |
|--------|-------------|
| `scheduled_query_alert_id` | ARM resource ID |
| `scheduled_query_alert_name` | Rule name |
| `identity_principal_id` | System-assigned identity principal (empty unless enabled) -- grant it workspace read access for Entra-only workspaces |

See `presets/` for error-spike, latency-threshold, and missing-heartbeat
starting points.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
