# Pay-As-You-Go Workspace

This preset creates the everyday Log Analytics Workspace: PerGB2018 (pay-as-you-go) billing, 90-day retention, and a daily ingestion cap as a cost guard. It is the right starting point for nearly every environment -- diagnostic settings, Container Insights, Application Insights, and query alerts all build on it.

## When to Use

- The first (often only) workspace of an environment or region
- Ingestion below ~100 GB/day, where commitment tiers do not yet pay off
- Anywhere you want logs queryable and alertable without billing commitments

## Key Configuration Choices

- **PerGB2018 sku** -- Azure's recommended tier; you pay only for what you ingest
- **90-day retention** (`retentionInDays: 90`) -- a compliance-friendly middle ground; the first 31 days are free on this tier
- **25 GB/day cap** (`dailyQuotaGb: 25`) -- ingestion stops at the cap until the next UTC day, so size it above normal daily volume; it is a data-loss dial, not just a cost dial

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-observability-rg` | Resource group holding the workspace | `AzureResourceGroup` status outputs |
| `my-platform-logs` | Workspace name (4-63 letters/digits/hyphens) | Your naming convention |

## Related Presets

- **02-commitment-tier** -- Switch when sustained ingestion crosses ~100 GB/day
- **03-private-hardened** -- Private-link-only, Entra-only posture for regulated estates
