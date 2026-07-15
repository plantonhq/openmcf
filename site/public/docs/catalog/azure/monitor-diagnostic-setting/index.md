---
title: "Monitor Diagnostic Setting"
description: "Monitor Diagnostic Setting deployment documentation"
icon: "package"
order: 100
componentName: "azuremonitordiagnosticsetting"
---

# Azure Monitor Diagnostic Setting

Deploys an Azure Monitor diagnostic setting -- the routing rule that selects which platform logs and metrics a resource emits and where they land: a Log Analytics Workspace, a Storage Account, an Event Hub, or a partner monitoring solution. Without one, most resources emit nothing beyond basic metrics.

## What Gets Created

When you deploy an AzureMonitorDiagnosticSetting resource, Planton provisions:

- **Diagnostic setting** -- an `azurerm_monitor_diagnostic_setting` on the referenced target resource, carrying the log/metric selection and the destination wiring

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A target resource** whose telemetry you are routing (any kind -- reference its `*_id` output)
- **At least one destination** -- typically an AzureLogAnalyticsWorkspace

## Quick Start

Create a file `diagnostics.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureMonitorDiagnosticSetting
metadata:
  name: vault-diagnostics
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureMonitorDiagnosticSetting.vault-diagnostics
spec:
  settingName: route-to-workspace
  targetResourceId:
    valueFrom:
      kind: AzureKeyVault
      name: app-vault
      fieldPath: status.outputs.key_vault_id
  enabledLogs:
    - categoryGroup: allLogs
  enabledMetrics:
    - category: AllMetrics
  logAnalyticsWorkspaceId:
    valueFrom:
      kind: AzureLogAnalyticsWorkspace
      name: platform-logs
      fieldPath: status.outputs.workspace_id
```

Deploy it:

```bash
planton pulumi up --manifest diagnostics.yaml
```

## Common Configurations

- **Modern table layout**: `logAnalyticsDestinationType: DEDICATED` -- resource-specific tables with typed columns (prefer where supported)
- **Compliance archive**: route the `audit` category group to a `storageAccountId` -- pennies-per-GB retention for years
- **SIEM stream**: set `eventhubAuthorizationRuleId` (a namespace-level rule) + `eventhubName`
- **Selective categories**: replace the group with `category` entries when volume costs demand it (discover names with `az monitor diagnostic-settings categories list --resource <id>`)

## Notes

- A target can carry up to five settings, each routing a different selection to different destinations.
- Category availability is defined per resource type by Azure -- the portal's "Diagnostic settings" blade lists them.
