# Linux Syslog and Performance to Workspace

This preset deploys the canonical Azure Monitor Agent shape for a Linux fleet: security-relevant syslog plus a CPU/memory/disk performance baseline, landing in one Log Analytics workspace. Attach machines with `AzureMonitorDataCollectionRuleAssociation` resources.

## When to Use

- The baseline rule for any Linux fleet feeding a central operations workspace
- As the starting shape for stricter collection -- narrow the facilities or severities further, or add a `transformKql` on the flow

## Key Configuration Choices

- **Filtered at the source** -- auth/authpriv/daemon facilities at Warning and above, not `*`: everything a flow lands bills per GB at the workspace, and the rule is the cheapest place to drop noise
- **60-second sampling** -- the fleet-baseline sweet spot; tighten per-workload rules separately rather than sampling the whole fleet aggressively
- **One flow, two streams** -- `Microsoft-Syslog` lands in the workspace's `Syslog` table and `Microsoft-Perf` in `Perf`; the flow only routes, the streams decide the tables

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-region>` | The rule's region (co-locate with the workspace) | The workspace's overview page |
| `<your-log-analytics-workspace>` | The Planton name of your `AzureLogAnalyticsWorkspace` resource | Planton console (or replace `valueFrom` with `value:` and a literal workspace ARM ID) |
