# Custom JSON Log to Workspace

This preset teaches the custom-log shape: a JSON application log file gets a declared schema (the stream declaration), an ingestion-time KQL filter, and lands in a Log Analytics custom table (`MyAppLogs_CL`).

## When to Use

- Application logs that live in files the Azure Monitor Agent can read (JSON here; switch `format: text` plus `settings.text.recordStartTimestampFormat` for text logs)
- Whenever you need schema and filtering on logs that no built-in stream covers

## Key Configuration Choices

- **A Data Collection Endpoint is mandatory** -- custom streams ingest through a DCE; the rule is rejected without `dataCollectionEndpointId`. Provide the literal ARM id of an endpoint (the DCE is not yet a catalog kind)
- **The stream declaration is a contract** -- its columns become the custom table's schema; include `TimeGenerated` or the workspace stamps arrival time and skews time-based queries. Version stream names instead of mutating a live schema
- **`transformKql` drops noise before it bills** -- the `where Level != 'debug'` filter runs in the ingestion pipeline; `outputStream` names the destination custom table (`_CL` suffix)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-region>` | The rule's region (co-locate with the workspace) | The workspace's overview page |
| `<your-data-collection-endpoint-arm-id>` | The Data Collection Endpoint's ARM resource ID | The DCE's Properties page in the Azure portal |
| `<your-log-analytics-workspace>` | The Planton name of your `AzureLogAnalyticsWorkspace` resource | Planton console (or replace `valueFrom` with `value:` and a literal workspace ARM ID) |
