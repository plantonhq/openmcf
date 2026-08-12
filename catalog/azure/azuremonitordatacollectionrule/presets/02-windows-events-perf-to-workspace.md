# Windows Events and Performance to Workspace

This preset deploys the canonical Azure Monitor Agent shape for a Windows fleet: XPath-filtered event logs (errors, warnings, and failed logons) plus a performance baseline, landing in one Log Analytics workspace. Attach machines with `AzureMonitorDataCollectionRuleAssociation` resources.

## When to Use

- The baseline rule for any Windows fleet feeding a central operations workspace
- As the starting shape for security-focused collection -- extend the Security channel's XPath with the event IDs your SOC actually consumes

## Key Configuration Choices

- **XPath is the billing filter** -- `Level=1 or Level=2 or Level=3` takes Critical/Error/Warning and skips Information, which is routinely >90% of raw event volume; the filter runs on the machine, before ingestion bills
- **Security stays surgical** -- collecting the whole Security channel is a SIEM decision with SIEM costs; this preset takes only failed logons (4625) as the universally useful signal
- **Streams pick the tables** -- `Microsoft-Event` lands in the workspace's `Event` table and `Microsoft-Perf` in `Perf`

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-region>` | The rule's region (co-locate with the workspace) | The workspace's overview page |
| `<your-log-analytics-workspace>` | The Planton name of your `AzureLogAnalyticsWorkspace` resource | Planton console (or replace `valueFrom` with `value:` and a literal workspace ARM ID) |
