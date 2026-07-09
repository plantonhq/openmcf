---
title: "Log Analytics Workspace"
description: "Log Analytics Workspace deployment documentation"
icon: "package"
order: 100
componentName: "azureloganalyticsworkspace"
---

# Azure Log Analytics Workspace

Deploys an Azure Log Analytics Workspace -- the central data platform for Azure Monitor. Diagnostic settings route platform logs into it, Application Insights stores telemetry in it, scheduled query alerts watch it, and Container Insights and Microsoft Sentinel build on it.

## What Gets Created

When you deploy an AzureLogAnalyticsWorkspace resource, Planton provisions:

- **Log Analytics Workspace** -- an `azurerm_log_analytics_workspace` in the chosen region and resource group, carrying the pricing tier (pay-as-you-go or a commitment tier), retention, daily quota, security/network posture, and optional managed identity
- **Azure Tags** -- Planton-derived governance tags merged with your own (your values win)

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An Azure Resource Group** where the workspace will be created (can reference an AzureResourceGroup resource)
- **Workspace naming plan** -- 4-63 letters, digits, and hyphens; must start and end with a letter or digit; renaming recreates the workspace and its data

## Quick Start

Create a file `log-analytics.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureLogAnalyticsWorkspace
metadata:
  name: platform-logs
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureLogAnalyticsWorkspace.platform-logs
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: observability-rg
      fieldPath: status.outputs.resource_group_name
  workspaceName: platform-logs
  retentionInDays: 90
  dailyQuotaGb: 25
```

Deploy it:

```bash
planton pulumi up --manifest log-analytics.yaml
```

## Common Configurations

- **Commitment tier**: `sku: CAPACITY_RESERVATION` + `reservationCapacityInGbPerDay: 100` when sustained ingestion crosses ~100 GB/day (fixed tiers up to 50000; a 31-day commitment applies)
- **Keyless posture**: `localAuthenticationEnabled: false` -- agents authenticate with Entra identities; the shared-key outputs stop being usable credentials
- **Private-link only**: `internetIngestionEnabled: false` + `internetQueryEnabled: false` -- both paths then require Azure Monitor Private Link Scope endpoints
- **Compliance retention**: `retentionInDays: 365` (or up to 730) plus `immediateDataPurgeOn30DaysEnabled: true` to drop Azure's post-retention grace store

## Key Outputs

| Output | Use |
| --- | --- |
| `workspace_id` | The ARM ID downstream kinds reference (App Insights, diagnostic settings, query alerts, AKS, Container Apps) |
| `workspace_customer_id` | The GUID agents authenticate against |
| `primary_shared_key` / `secondary_shared_key` | Agent keys (secret-bearing; rotate via the primary/secondary swap) |
