# AzureMonitorDiagnosticSetting

## Overview

`AzureMonitorDiagnosticSetting` provisions an Azure Monitor diagnostic
setting -- the routing rule that makes a resource observable. It selects
which platform log categories and metrics a target resource emits and sends
them to one or more destinations: a Log Analytics Workspace (queryable,
alertable), a Storage Account (cheap archival), an Event Hub (SIEM/streaming
integration), or an Azure Native ISV partner solution.

Without a diagnostic setting, most Azure resources emit nothing beyond basic
platform metrics. This kind is the glue of the observability story: deploy a
resource, deploy a diagnostic setting on it, and its logs start landing where
alerts and dashboards can see them.

## Key Features

- **Any target** -- the setting is an ARM extension resource; the polymorphic
  `target_resource_id` reference points at any kind's `*_id` output (a vault,
  a cluster, a gateway, a database, a subscription)
- **Category groups** -- `allLogs` and `audit` bundles track new categories
  automatically as Azure adds them; single categories remain available for
  volume control
- **Four destination families** -- workspace (with the modern resource-specific
  table layout), storage, Event Hub (namespace rule + optional named hub),
  partner solution -- any combination in one setting
- **Front-loaded contracts** -- at least one category, at least one
  destination, and category-XOR-category-group are enforced at validation
  time (Azure otherwise accepts an empty setting and then 404s reading it)

## When to Use

- On every production resource whose logs matter -- typically paired with the
  environment's central `AzureLogAnalyticsWorkspace`
- Compliance archival (audit categories to a storage account)
- SIEM integration (audit categories to an Event Hub)
- Up to five settings per target, each routing a different selection

## Spec Highlights

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `setting_name` | string | Yes | Unique among the target's settings (ForceNew) |
| `target_resource_id` | StringValueOrRef | Yes | Any ARM resource (no default kind -- reference explicitly; ForceNew) |
| `enabled_logs[]` | message | At least one log or metric | `category` XOR `category_group` per entry |
| `enabled_metrics[]` | message | At least one log or metric | Metric category (usually "AllMetrics") |
| `log_analytics_workspace_id` | StringValueOrRef | One destination required | Workspace destination |
| `log_analytics_destination_type` | enum | No | DEDICATED (modern tables) / AZURE_DIAGNOSTICS (legacy) |
| `storage_account_id` | StringValueOrRef | One destination required | Archival destination |
| `eventhub_authorization_rule_id` (+ `eventhub_name`) | string | One destination required | Streaming destination (namespace-level rule) |
| `partner_solution_id` | string | One destination required | Azure Native ISV destination |

## Outputs

| Output | Description |
|--------|-------------|
| `diagnostic_setting_id` | The ARM extension-resource ID ({target}/providers/Microsoft.Insights/diagnosticSettings/{name}) |
| `diagnostic_setting_name` | The setting name |
| `target_resource_id` | The resolved target |

## Composition

```yaml
targetResourceId:
  valueFrom:
    kind: AzureKeyVault
    name: my-app-vault
    fieldPath: status.outputs.key_vault_id
logAnalyticsWorkspaceId:
  valueFrom:
    kind: AzureLogAnalyticsWorkspace
    name: my-platform-logs
    fieldPath: status.outputs.workspace_id
```

See `presets/` for workspace-routing, storage-archival, and SIEM-streaming
starting points.
