# AzureMonitorMetricAlert

## Overview

`AzureMonitorMetricAlert` provisions an Azure Monitor metric alert rule --
a condition evaluated over platform metrics (CPU, latency, transactions,
queue depth) on a rolling window, firing action groups when it holds.

Three condition families cover the real alerting space:

- **Static criteria** -- classic thresholds ("availability < 99.9%");
  multiple metrics AND together in one rule
- **Dynamic criteria** -- machine-learning thresholds that learn the
  metric's normal band (including seasonality) and alert on deviation
- **Web-test availability** -- fires when an Application Insights
  availability test fails from N locations (the outside-in "site is down")

## Key Features

- **Polymorphic scopes** -- any metric-emitting resource, a resource group,
  or a whole subscription (multi-scope rules take `target_resource_type` +
  `target_resource_location`)
- **Dimension filters** -- restrict to specific dimension values or split
  the alert across them (each value alerts independently)
- **Stateful by default** -- auto-mitigation resolves the alert when the
  condition stops holding
- **Closed vocabularies** -- aggregations, operators (with the
  static/dynamic direction split enforced), sensitivities, frequencies,
  and windows all validate before deploy

## When to Use

- SLO-backed thresholds on any resource the catalog deploys
- Anomaly detection on rhythmic metrics where fixed lines are wrong
- Availability paging driven by Application Insights web tests

## Spec Highlights

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `resource_group` | StringValueOrRef | Yes | - | Resource group (rules are global -- no region) |
| `alert_name` | string | Yes | - | Rule name (ForceNew) |
| `scopes[]` | StringValueOrRef | Yes (min 1) | - | The evaluated resources (no default kind -- reference explicitly) |
| `severity` | int32 | No | 3 | 0 (critical) - 4 (verbose) |
| `frequency` / `window_size` | string | No | PT1M / PT5M | ISO 8601 vocabularies |
| `static_criteria[]` XOR `dynamic_criteria` XOR `web_test_availability_criteria` | messages | Exactly one family | - | The condition |
| `auto_mitigate` | bool | No | true | Stateful auto-resolve |
| `actions[]` | message | No | - | Action-group FKs + webhook properties |
| `tags` | map | No | - | User tags |

## Outputs

| Output | Description |
|--------|-------------|
| `metric_alert_id` | ARM resource ID |
| `metric_alert_name` | Rule name |

See `presets/` for static-threshold, dynamic-anomaly, and web-test
availability starting points.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
