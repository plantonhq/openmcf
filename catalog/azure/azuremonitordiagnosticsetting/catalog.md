# Azure Monitor Diagnostic Setting

Deploys an Azure Monitor diagnostic setting -- how a resource's platform telemetry LEAVES the resource. It selects which log categories and metrics the target emits and routes them to destinations: a Log Analytics Workspace (queryable with KQL, alertable with scheduled query rules), a Storage Account (cheap archival), an Event Hub (streaming to SIEMs), or an Azure Native partner solution. Without a diagnostic setting, most Azure resources emit nothing beyond basic platform metrics -- this kind closes the logging pipeline that scheduled query alerts watch. The component integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring to targets and destinations.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Diagnostic Setting** -- a `Microsoft.Insights/diagnosticSettings` EXTENSION resource living ON the target (any ARM resource -- a Key Vault, an AKS cluster, a gateway, a subscription), carrying the category selections and destination routing. A target can carry up to five settings

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### Azure Subscription

- **A target resource** that emits diagnostics -- reference its `*_id` output or paste its ARM ID. Category names are defined per resource type: `az monitor diagnostic-settings categories list --resource <id>` lists them.
- **At least one destination** -- a Log Analytics Workspace, a Storage Account, an Event Hub namespace authorization rule (with send rights), or a partner solution.

## Deploy

### Console

Open the deployment store, find **Azure Monitor Diagnostic Setting**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **logs-to-workspace** preset in the [Presets](#presets) tab to pre-populate the everyday observability wiring.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureMonitorDiagnosticSetting
metadata:
  name: vault-diagnostics
  org: acme-corp
  env: prod
spec:
  settingName: route-to-workspace
  targetResourceId:
    valueFrom:
      kind: AzureKeyVault
      name: my-app-vault
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
  logAnalyticsDestinationType: DEDICATED
```

```shell
planton apply -f diagnostic-setting.yaml
```

This routes every current AND future log category of the vault (the allLogs group tracks new ones automatically) plus its metrics into the workspace, landing in modern resource-specific tables.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the setting between the resource it drains and the workspace it fills -- the InfraPipeline resolves the dependency graph and deploys target, workspace, then setting.

## Key Configuration

These are the most important decisions when configuring a diagnostic setting. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Telemetry selection** -- each `enabledLogs` entry names exactly one of a single `category` (AuditEvent on a vault, kube-audit on AKS) or a `categoryGroup` (Azure's curated allLogs/audit bundles, which track new categories automatically). At least one log or metric selection is required -- Azure rejects a setting that routes nothing. Volume note: allLogs on a chatty resource can dominate workspace ingestion costs.

**Destinations** -- at least one of: `logAnalyticsWorkspaceId` (queryable + alertable -- pair with `logAnalyticsDestinationType: DEDICATED` for modern typed tables), `storageAccountId` (cheap long-term archival), the Event Hub pair (`eventhubAuthorizationRuleId` -- a NAMESPACE-scoped send rule -- with the optional `eventhubName` picking one hub), or `partnerSolutionId` (Elastic/Datadog via Marketplace). One setting can feed several at once.

**The extension model** -- `settingName` and `targetResourceId` are fixed at creation (the setting lives ON the target); selections and destinations update in place -- re-routing telemetry is the designed operation.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **Any diagnostics-emitting kind** | `targetResourceId` | the target's `*_id` output (explicit valueFrom) |
| **AzureLogAnalyticsWorkspace** | `logAnalyticsWorkspaceId` | `status.outputs.workspace_id` |
| **AzureStorageAccount** | `storageAccountId` | `status.outputs.storage_account_id` |
| **AzureEventHubAuthorizationRule** | `eventhubAuthorizationRuleId` | `status.outputs.authorization_rule_id` |
| **AzureEventHub** | `eventhubName` | `status.outputs.event_hub_name` |

### What This Component Provides

After provisioning, `status.outputs` contains (the setting is a leaf -- operator identifiers):

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `diagnostic_setting_id` | Azure Resource Manager ID of the setting | Portal navigation, CLI references |
| `diagnostic_setting_name` | Name of the setting | Azure CLI references |
| `target_resource_id` | The as-deployed target | The routing audit trail -- what this setting actually drains |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Logs to workspace** -- allLogs + AllMetrics into the platform workspace: the everyday observability wiring, and what makes scheduled query alerts see the resource. Start from the **logs-to-workspace** preset.

**Archive to storage** -- audit categories into blob storage for long-horizon compliance retention. Start from the **archive-to-storage** preset.

**Stream to SIEM** -- categories into an Event Hub an external pipeline consumes. Start from the **stream-to-siem** preset.

## Works With

- [**Azure Log Analytics Workspace**](/cloud-catalog/azure-log-analytics-workspace) -- the queryable, alertable destination
- [**Azure Monitor Scheduled Query Alert**](/cloud-catalog/azure-monitor-scheduled-query-alert) -- watches the workspace this setting fills
- [**Azure Storage Account**](/cloud-catalog/azure-storage-account) -- the archival destination
- [**Azure Key Vault**](/cloud-catalog/azure-key-vault) -- the classic audit-log target
