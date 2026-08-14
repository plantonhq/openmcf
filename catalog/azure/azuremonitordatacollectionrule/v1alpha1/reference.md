# AzureMonitorDataCollectionRule

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Deep-shape example for docs and offline validation: a Linux rule
# collecting filtered syslog and performance counters, landing logs in
# a Log Analytics workspace and metrics in Azure Monitor metrics --
# with an ingestion-time KQL filter on the syslog flow. References are
# literal ARM ids so the manifest validates standalone.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureMonitorDataCollectionRule
metadata:
  name: test-dcr
  id: test-dcr
  org: test-org
  env: test
spec:
  resourceGroup:
    value: test-rg
  name: linux-baseline
  region: eastus
  kind: Linux
  description: Linux fleet baseline - filtered syslog and perf counters
  dataSources:
    syslogs:
      - name: linux-syslog
        facilityNames:
          - auth
          - authpriv
          - daemon
        logLevels:
          - Warning
          - Error
          - Critical
          - Alert
          - Emergency
        streams:
          - Microsoft-Syslog
    performanceCounters:
      - name: linux-perf
        samplingFrequencyInSeconds: 60
        counterSpecifiers:
          - Processor(*)\% Processor Time
          - Memory(*)\% Used Memory
        streams:
          - Microsoft-Perf
      - name: linux-insights-metrics
        # Streams targeting Microsoft-InsightsMetrics require EXACTLY
        # 60-second sampling (Azure enforces at deploy time).
        samplingFrequencyInSeconds: 60
        counterSpecifiers:
          - Processor(*)\% Processor Time
        streams:
          - Microsoft-InsightsMetrics
  destinations:
    logAnalytics:
      - name: ops-workspace
        workspaceResourceId:
          value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.OperationalInsights/workspaces/test-law
    azureMonitorMetrics:
      name: azure-metrics
  dataFlows:
    - streams:
        - Microsoft-Syslog
      destinations:
        - ops-workspace
      # Drop noisy repeated daemon chatter at ingestion, before it bills.
      transformKql: source | where SyslogMessage !has 'systemd' or SeverityLevel != 'warning'
      outputStream: Microsoft-Syslog
    - streams:
        - Microsoft-Perf
      destinations:
        - ops-workspace
    - streams:
        - Microsoft-InsightsMetrics
      destinations:
        - azure-metrics
  tags:
    cost-center: platform-observability
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.region` | `string` | yes |  |  |
| `spec.kind` | `string` |  |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.dataCollectionEndpointId` | `string \| valueFrom` |  |  |  |
| `spec.identity` | `AzureMonitorDataCollectionRuleIdentity` |  |  |  |
| `spec.identity.type` | `enum` | yes |  |  |
| `spec.identity.identityIds` | `[]string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.dataSources` | `AzureMonitorDataCollectionRuleDataSources` |  |  |  |
| `spec.dataSources.syslogs` | `[]AzureMonitorDataCollectionRuleSyslog` |  |  |  |
| `spec.dataSources.syslogs[].name` | `string` | yes |  |  |
| `spec.dataSources.syslogs[].facilityNames` | `[]string` | yes |  |  |
| `spec.dataSources.syslogs[].logLevels` | `[]string` | yes |  |  |
| `spec.dataSources.syslogs[].streams` | `[]string` | yes |  |  |
| `spec.dataSources.performanceCounters` | `[]AzureMonitorDataCollectionRulePerformanceCounter` |  |  |  |
| `spec.dataSources.performanceCounters[].name` | `string` | yes |  |  |
| `spec.dataSources.performanceCounters[].samplingFrequencyInSeconds` | `int32` | yes |  |  |
| `spec.dataSources.performanceCounters[].counterSpecifiers` | `[]string` | yes |  |  |
| `spec.dataSources.performanceCounters[].streams` | `[]string` | yes |  |  |
| `spec.dataSources.windowsEventLogs` | `[]AzureMonitorDataCollectionRuleWindowsEventLog` |  |  |  |
| `spec.dataSources.windowsEventLogs[].name` | `string` | yes |  |  |
| `spec.dataSources.windowsEventLogs[].xPathQueries` | `[]string` | yes |  |  |
| `spec.dataSources.windowsEventLogs[].streams` | `[]string` | yes |  |  |
| `spec.dataSources.extensions` | `[]AzureMonitorDataCollectionRuleExtension` |  |  |  |
| `spec.dataSources.extensions[].name` | `string` | yes |  |  |
| `spec.dataSources.extensions[].extensionName` | `string` | yes |  |  |
| `spec.dataSources.extensions[].extensionJson` | `string` |  |  |  |
| `spec.dataSources.extensions[].inputDataSources` | `[]string` |  |  |  |
| `spec.dataSources.extensions[].streams` | `[]string` | yes |  |  |
| `spec.dataSources.iisLogs` | `[]AzureMonitorDataCollectionRuleIisLog` |  |  |  |
| `spec.dataSources.iisLogs[].name` | `string` | yes |  |  |
| `spec.dataSources.iisLogs[].logDirectories` | `[]string` |  |  |  |
| `spec.dataSources.iisLogs[].streams` | `[]string` | yes |  |  |
| `spec.dataSources.logFiles` | `[]AzureMonitorDataCollectionRuleLogFile` |  |  |  |
| `spec.dataSources.logFiles[].name` | `string` | yes |  |  |
| `spec.dataSources.logFiles[].filePatterns` | `[]string` | yes |  |  |
| `spec.dataSources.logFiles[].format` | `string` | yes |  |  |
| `spec.dataSources.logFiles[].settings` | `AzureMonitorDataCollectionRuleLogFileSettings` |  |  |  |
| `spec.dataSources.logFiles[].settings.text` | `AzureMonitorDataCollectionRuleLogFileSettingsText` | yes |  |  |
| `spec.dataSources.logFiles[].settings.text.recordStartTimestampFormat` | `string` | yes |  |  |
| `spec.dataSources.logFiles[].streams` | `[]string` | yes |  |  |
| `spec.dataSources.prometheusForwarders` | `[]AzureMonitorDataCollectionRulePrometheusForwarder` |  |  |  |
| `spec.dataSources.prometheusForwarders[].name` | `string` | yes |  |  |
| `spec.dataSources.prometheusForwarders[].labelIncludeFilters` | `[]AzureMonitorDataCollectionRuleLabelIncludeFilter` |  |  |  |
| `spec.dataSources.prometheusForwarders[].labelIncludeFilters[].label` | `string` | yes |  |  |
| `spec.dataSources.prometheusForwarders[].labelIncludeFilters[].value` | `string` | yes |  |  |
| `spec.dataSources.prometheusForwarders[].streams` | `[]string` | yes |  |  |
| `spec.dataSources.windowsFirewallLogs` | `[]AzureMonitorDataCollectionRuleWindowsFirewallLog` |  |  |  |
| `spec.dataSources.windowsFirewallLogs[].name` | `string` | yes |  |  |
| `spec.dataSources.windowsFirewallLogs[].streams` | `[]string` | yes |  |  |
| `spec.dataSources.platformTelemetries` | `[]AzureMonitorDataCollectionRulePlatformTelemetry` |  |  |  |
| `spec.dataSources.platformTelemetries[].name` | `string` | yes |  |  |
| `spec.dataSources.platformTelemetries[].streams` | `[]string` | yes |  |  |
| `spec.dataSources.dataImport` | `AzureMonitorDataCollectionRuleDataImport` |  |  |  |
| `spec.dataSources.dataImport.eventHubDataSource` | `AzureMonitorDataCollectionRuleDataImportEventHub` | yes |  |  |
| `spec.dataSources.dataImport.eventHubDataSource.name` | `string` | yes |  |  |
| `spec.dataSources.dataImport.eventHubDataSource.stream` | `string` | yes |  |  |
| `spec.dataSources.dataImport.eventHubDataSource.consumerGroup` | `string` |  |  |  |
| `spec.destinations` | `AzureMonitorDataCollectionRuleDestinations` | yes |  |  |
| `spec.destinations.logAnalytics` | `[]AzureMonitorDataCollectionRuleLogAnalytics` |  |  |  |
| `spec.destinations.logAnalytics[].name` | `string` | yes |  |  |
| `spec.destinations.logAnalytics[].workspaceResourceId` | `string \| valueFrom` | yes |  | AzureLogAnalyticsWorkspace (`status.outputs.workspace_id`) |
| `spec.destinations.azureMonitorMetrics` | `AzureMonitorDataCollectionRuleAzureMonitorMetrics` |  |  |  |
| `spec.destinations.azureMonitorMetrics.name` | `string` | yes |  |  |
| `spec.destinations.eventHub` | `AzureMonitorDataCollectionRuleEventHubDestination` |  |  |  |
| `spec.destinations.eventHub.name` | `string` | yes |  |  |
| `spec.destinations.eventHub.eventHubId` | `string \| valueFrom` | yes |  | AzureEventHub (`status.outputs.event_hub_id`) |
| `spec.destinations.eventHubDirect` | `AzureMonitorDataCollectionRuleEventHubDestination` |  |  |  |
| `spec.destinations.eventHubDirect.name` | `string` | yes |  |  |
| `spec.destinations.eventHubDirect.eventHubId` | `string \| valueFrom` | yes |  | AzureEventHub (`status.outputs.event_hub_id`) |
| `spec.destinations.monitorAccounts` | `[]AzureMonitorDataCollectionRuleMonitorAccount` |  |  |  |
| `spec.destinations.monitorAccounts[].name` | `string` | yes |  |  |
| `spec.destinations.monitorAccounts[].monitorAccountId` | `string \| valueFrom` | yes |  |  |
| `spec.destinations.storageBlobs` | `[]AzureMonitorDataCollectionRuleStorageBlobDestination` |  |  |  |
| `spec.destinations.storageBlobs[].name` | `string` | yes |  |  |
| `spec.destinations.storageBlobs[].containerName` | `string` | yes |  |  |
| `spec.destinations.storageBlobs[].storageAccountId` | `string \| valueFrom` | yes |  | AzureStorageAccount (`status.outputs.storage_account_id`) |
| `spec.destinations.storageBlobDirects` | `[]AzureMonitorDataCollectionRuleStorageBlobDestination` |  |  |  |
| `spec.destinations.storageBlobDirects[].name` | `string` | yes |  |  |
| `spec.destinations.storageBlobDirects[].containerName` | `string` | yes |  |  |
| `spec.destinations.storageBlobDirects[].storageAccountId` | `string \| valueFrom` | yes |  | AzureStorageAccount (`status.outputs.storage_account_id`) |
| `spec.destinations.storageTableDirects` | `[]AzureMonitorDataCollectionRuleStorageTableDirect` |  |  |  |
| `spec.destinations.storageTableDirects[].name` | `string` | yes |  |  |
| `spec.destinations.storageTableDirects[].tableName` | `string` | yes |  |  |
| `spec.destinations.storageTableDirects[].storageAccountId` | `string \| valueFrom` | yes |  | AzureStorageAccount (`status.outputs.storage_account_id`) |
| `spec.dataFlows` | `[]AzureMonitorDataCollectionRuleDataFlow` | yes |  |  |
| `spec.dataFlows[].streams` | `[]string` | yes |  |  |
| `spec.dataFlows[].destinations` | `[]string` | yes |  |  |
| `spec.dataFlows[].builtInTransform` | `string` |  |  |  |
| `spec.dataFlows[].outputStream` | `string` |  |  |  |
| `spec.dataFlows[].transformKql` | `string` |  |  |  |
| `spec.streamDeclarations` | `[]AzureMonitorDataCollectionRuleStreamDeclaration` |  |  |  |
| `spec.streamDeclarations[].streamName` | `string` | yes |  |  |
| `spec.streamDeclarations[].columns` | `[]AzureMonitorDataCollectionRuleStreamDeclarationColumn` | yes |  |  |
| `spec.streamDeclarations[].columns[].name` | `string` | yes |  |  |
| `spec.streamDeclarations[].columns[].type` | `string` | yes |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.resourceGroup

`string | valueFrom` · required

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.region

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.kind

`string` · optional (explicit presence)

- rule: {"string":{"in":["Linux","Windows","AgentDirectToStore","WorkspaceTransforms"]}}

### spec.description

`string`

### spec.dataCollectionEndpointId

`string | valueFrom`

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.identity

`AzureMonitorDataCollectionRuleIdentity`

- rule: identity_ids is required for USER_ASSIGNED and must be empty for SYSTEM_ASSIGNED

### spec.identity.type

`enum` · required

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_monitor_data_collection_rule_identity_type_unspecified`
- `SYSTEM_ASSIGNED`
- `USER_ASSIGNED`

### spec.identity.identityIds

`[]string | valueFrom`

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.dataSources

`AzureMonitorDataCollectionRuleDataSources`

### spec.dataSources.syslogs

`[]AzureMonitorDataCollectionRuleSyslog`

### spec.dataSources.syslogs[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1","maxLen":"32"}}

### spec.dataSources.syslogs[].facilityNames

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"in":["*","alert","audit","auth","authpriv","clock","cron","daemon","ftp","kern","local0","local1","local2","local3","local4","local5","local6","local7","lpr","mail","mark","news","nopri","ntp","syslog","user","uucp"]}}}}

### spec.dataSources.syslogs[].logLevels

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"in":["*","Alert","Critical","Debug","Emergency","Error","Info","Notice","Warning"]}}}}

### spec.dataSources.syslogs[].streams

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.dataSources.performanceCounters

`[]AzureMonitorDataCollectionRulePerformanceCounter`

### spec.dataSources.performanceCounters[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1","maxLen":"32"}}

### spec.dataSources.performanceCounters[].samplingFrequencyInSeconds

`int32` · required · optional (explicit presence)

- rule: {"required":true,"int32":{"lte":1800,"gte":1}}

### spec.dataSources.performanceCounters[].counterSpecifiers

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.dataSources.performanceCounters[].streams

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.dataSources.windowsEventLogs

`[]AzureMonitorDataCollectionRuleWindowsEventLog`

### spec.dataSources.windowsEventLogs[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1","maxLen":"32"}}

### spec.dataSources.windowsEventLogs[].xPathQueries

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.dataSources.windowsEventLogs[].streams

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.dataSources.extensions

`[]AzureMonitorDataCollectionRuleExtension`

### spec.dataSources.extensions[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1","maxLen":"32"}}

### spec.dataSources.extensions[].extensionName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.dataSources.extensions[].extensionJson

`string`

### spec.dataSources.extensions[].inputDataSources

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.dataSources.extensions[].streams

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.dataSources.iisLogs

`[]AzureMonitorDataCollectionRuleIisLog`

### spec.dataSources.iisLogs[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1","maxLen":"32"}}

### spec.dataSources.iisLogs[].logDirectories

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.dataSources.iisLogs[].streams

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.dataSources.logFiles

`[]AzureMonitorDataCollectionRuleLogFile`

### spec.dataSources.logFiles[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1","maxLen":"32"}}

### spec.dataSources.logFiles[].filePatterns

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.dataSources.logFiles[].format

`string` · required

- rule: {"required":true,"string":{"in":["json","text"]}}

### spec.dataSources.logFiles[].settings

`AzureMonitorDataCollectionRuleLogFileSettings`

### spec.dataSources.logFiles[].settings.text

`AzureMonitorDataCollectionRuleLogFileSettingsText` · required

- rule: {"required":true}

### spec.dataSources.logFiles[].settings.text.recordStartTimestampFormat

`string` · required

- rule: {"required":true,"string":{"in":["ISO 8601","YYYY-MM-DD HH:MM:SS","M/D/YYYY HH:MM:SS AM/PM","Mon DD, YYYY HH:MM:SS","yyMMdd HH:mm:ss","ddMMyy HH:mm:ss","MMM d hh:mm:ss","dd/MMM/yyyy:HH:mm:ss zzz","yyyy-MM-ddTHH:mm:ssK"]}}

### spec.dataSources.logFiles[].streams

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.dataSources.prometheusForwarders

`[]AzureMonitorDataCollectionRulePrometheusForwarder`

### spec.dataSources.prometheusForwarders[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1","maxLen":"32"}}

### spec.dataSources.prometheusForwarders[].labelIncludeFilters

`[]AzureMonitorDataCollectionRuleLabelIncludeFilter`

### spec.dataSources.prometheusForwarders[].labelIncludeFilters[].label

`string` · required

- rule: {"required":true,"string":{"in":["microsoft_metrics_include_label"]}}

### spec.dataSources.prometheusForwarders[].labelIncludeFilters[].value

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.dataSources.prometheusForwarders[].streams

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"in":["Microsoft-PrometheusMetrics"]}}}}

### spec.dataSources.windowsFirewallLogs

`[]AzureMonitorDataCollectionRuleWindowsFirewallLog`

### spec.dataSources.windowsFirewallLogs[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1","maxLen":"32"}}

### spec.dataSources.windowsFirewallLogs[].streams

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.dataSources.platformTelemetries

`[]AzureMonitorDataCollectionRulePlatformTelemetry`

### spec.dataSources.platformTelemetries[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1","maxLen":"32"}}

### spec.dataSources.platformTelemetries[].streams

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.dataSources.dataImport

`AzureMonitorDataCollectionRuleDataImport`

### spec.dataSources.dataImport.eventHubDataSource

`AzureMonitorDataCollectionRuleDataImportEventHub` · required

- rule: {"required":true}

### spec.dataSources.dataImport.eventHubDataSource.name

`string` · required

- rule: {"required":true,"string":{"minLen":"1","maxLen":"32"}}

### spec.dataSources.dataImport.eventHubDataSource.stream

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.dataSources.dataImport.eventHubDataSource.consumerGroup

`string`

### spec.destinations

`AzureMonitorDataCollectionRuleDestinations` · required

- rule: {"required":true}
- rule: configure at least one destination (log_analytics, azure_monitor_metrics, event_hub, event_hub_direct, monitor_accounts, storage_blobs, storage_blob_directs or storage_table_directs)

### spec.destinations.logAnalytics

`[]AzureMonitorDataCollectionRuleLogAnalytics`

### spec.destinations.logAnalytics[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.destinations.logAnalytics[].workspaceResourceId

`string | valueFrom` · required

- references: AzureLogAnalyticsWorkspace (`status.outputs.workspace_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureLogAnalyticsWorkspace, name: <that resource's name>, fieldPath: status.outputs.workspace_id}} -- a bare string does not parse

### spec.destinations.azureMonitorMetrics

`AzureMonitorDataCollectionRuleAzureMonitorMetrics`

### spec.destinations.azureMonitorMetrics.name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.destinations.eventHub

`AzureMonitorDataCollectionRuleEventHubDestination`

### spec.destinations.eventHub.name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.destinations.eventHub.eventHubId

`string | valueFrom` · required

- references: AzureEventHub (`status.outputs.event_hub_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureEventHub, name: <that resource's name>, fieldPath: status.outputs.event_hub_id}} -- a bare string does not parse

### spec.destinations.eventHubDirect

`AzureMonitorDataCollectionRuleEventHubDestination`

### spec.destinations.eventHubDirect.name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.destinations.eventHubDirect.eventHubId

`string | valueFrom` · required

- references: AzureEventHub (`status.outputs.event_hub_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureEventHub, name: <that resource's name>, fieldPath: status.outputs.event_hub_id}} -- a bare string does not parse

### spec.destinations.monitorAccounts

`[]AzureMonitorDataCollectionRuleMonitorAccount`

### spec.destinations.monitorAccounts[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.destinations.monitorAccounts[].monitorAccountId

`string | valueFrom` · required

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.destinations.storageBlobs

`[]AzureMonitorDataCollectionRuleStorageBlobDestination`

### spec.destinations.storageBlobs[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.destinations.storageBlobs[].containerName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.destinations.storageBlobs[].storageAccountId

`string | valueFrom` · required

- references: AzureStorageAccount (`status.outputs.storage_account_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.storage_account_id}} -- a bare string does not parse

### spec.destinations.storageBlobDirects

`[]AzureMonitorDataCollectionRuleStorageBlobDestination`

### spec.destinations.storageBlobDirects[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.destinations.storageBlobDirects[].containerName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.destinations.storageBlobDirects[].storageAccountId

`string | valueFrom` · required

- references: AzureStorageAccount (`status.outputs.storage_account_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.storage_account_id}} -- a bare string does not parse

### spec.destinations.storageTableDirects

`[]AzureMonitorDataCollectionRuleStorageTableDirect`

### spec.destinations.storageTableDirects[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.destinations.storageTableDirects[].tableName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.destinations.storageTableDirects[].storageAccountId

`string | valueFrom` · required

- references: AzureStorageAccount (`status.outputs.storage_account_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.storage_account_id}} -- a bare string does not parse

### spec.dataFlows

`[]AzureMonitorDataCollectionRuleDataFlow` · required

- rule: {"repeated":{"minItems":"1"}}

### spec.dataFlows[].streams

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.dataFlows[].destinations

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.dataFlows[].builtInTransform

`string`

### spec.dataFlows[].outputStream

`string`

### spec.dataFlows[].transformKql

`string`

### spec.streamDeclarations

`[]AzureMonitorDataCollectionRuleStreamDeclaration`

### spec.streamDeclarations[].streamName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.streamDeclarations[].columns

`[]AzureMonitorDataCollectionRuleStreamDeclarationColumn` · required

- rule: {"repeated":{"minItems":"1"}}

### spec.streamDeclarations[].columns[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.streamDeclarations[].columns[].type

`string` · required

- rule: {"required":true,"string":{"in":["boolean","datetime","dynamic","int","long","real","string"]}}

### spec.tags

`map<string, string>`

## Validation Rules

- `dcr_stream_declaration_names_unique`: stream_declarations must carry unique stream_name values

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureMonitorDataCollectionRule, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.data_collection_rule_id` | `string` |  |
| `status.outputs.data_collection_rule_name` | `string` |  |
| `status.outputs.immutable_id` | `string` |  |
| `status.outputs.identity_principal_id` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.identity.identityIds` | AzureUserAssignedIdentity | `status.outputs.identity_id` |
| `spec.destinations.logAnalytics[].workspaceResourceId` | AzureLogAnalyticsWorkspace | `status.outputs.workspace_id` |
| `spec.destinations.eventHub.eventHubId` | AzureEventHub | `status.outputs.event_hub_id` |
| `spec.destinations.eventHubDirect.eventHubId` | AzureEventHub | `status.outputs.event_hub_id` |
| `spec.destinations.storageBlobs[].storageAccountId` | AzureStorageAccount | `status.outputs.storage_account_id` |
| `spec.destinations.storageBlobDirects[].storageAccountId` | AzureStorageAccount | `status.outputs.storage_account_id` |
| `spec.destinations.storageTableDirects[].storageAccountId` | AzureStorageAccount | `status.outputs.storage_account_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureLogAnalyticsWorkspace | `spec.dataCollectionRuleId` | `status.outputs.data_collection_rule_id` |
| AzureMonitorDataCollectionRuleAssociation | `spec.dataCollectionRuleId` | `status.outputs.data_collection_rule_id` |

## See Also

- [Overview](../README.md)
