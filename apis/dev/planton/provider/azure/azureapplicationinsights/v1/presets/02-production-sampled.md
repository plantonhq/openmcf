# Production Application Insights (Sampled, Cost-Controlled)

This preset creates a production APM resource with the cost levers engaged: 50% sampling, a 10 GB daily cap with notification, one-year retention, and Entra-only ingestion.

## When to Use

- High-traffic production services where full-fidelity telemetry costs more than it informs
- Estates standardizing on keyless (managed-identity) SDK authentication
- Compliance environments needing a one-year telemetry trail

## Key Configuration Choices

- **50% sampling** (`samplingPercentage: 50`) -- statistically representative telemetry at half the volume; APM percentiles stay accurate
- **10 GB/day cap + notification** -- ingestion stops at the cap until the next UTC day; the email tells you telemetry is being dropped rather than letting it fail silently
- **Entra-only** (`localAuthenticationEnabled: false`) -- a bare instrumentation key no longer authorizes ingestion; SDKs authenticate with managed identities
- **365-day retention** -- the workspace's retention still governs workspace-based queries; this aligns the classic-experience tables

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-observability-rg` | Resource group holding the resource | `AzureResourceGroup` status outputs |
| `my-platform-logs` | Workspace storing the telemetry | `AzureLogAnalyticsWorkspace` status outputs |
| `my-prod-app-insights` | Resource name | Your naming convention |

## Related Presets

- **01-standard** -- Full-fidelity default for lower-traffic environments
