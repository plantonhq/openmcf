---
title: "Audit Trail to Storage Archival"
description: "This preset routes a resource's audit-class logs (via the `audit` category group) to a storage account -- the pennies-per-GB destination for compliance retention measured in years, where the..."
type: "preset"
rank: "02"
presetSlug: "02-archive-to-storage"
componentSlug: "monitor-diagnostic-setting"
componentTitle: "Monitor Diagnostic Setting"
provider: "azure"
icon: "package"
order: 2
---

# Audit Trail to Storage Archival

This preset routes a resource's audit-class logs (via the `audit` category group) to a storage account -- the pennies-per-GB destination for compliance retention measured in years, where the workspace's per-GB pricing would be wasteful.

## When to Use

- Multi-year compliance retention of audit logs (the workspace keeps at most 730 days)
- Cold telemetry nobody queries day to day but regulators may demand
- Alongside a workspace setting: a resource can carry up to five settings, each routing a different selection

## Key Configuration Choices

- **`audit` category group** -- just the audit-relevant categories, tracked automatically as Azure adds them
- **Storage destination only** -- archived logs are not KQL-queryable; add a workspace destination (or a second setting) for the operational window
- **Pair with lifecycle rules** -- the storage account's lifecycle management tiers the archive to Cool/Archive automatically

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-app-vault` | The resource whose audit trail is archived | Any kind's status outputs (`*_id`) |
| `my-archive-account` | The archival storage account | `AzureStorageAccount` status outputs |

## Related Presets

- **01-logs-to-workspace** -- The queryable operational window
- **03-stream-to-siem** -- Real-time streaming to external systems
