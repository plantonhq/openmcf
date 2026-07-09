# Private, Hardened Workspace

This preset creates a Log Analytics Workspace in the regulated-estate posture: no public ingestion or query paths, no shared-key authentication, centralized query permissions, and no post-retention grace store.

## When to Use

- Regulated workloads (finance, health, government) with data-residency and access-control mandates
- Estates that route all monitoring traffic through Azure Monitor Private Link Scope (AMPLS)
- Security-team-owned central workspaces where resource-context access is too permissive

## Key Configuration Choices

- **Private-link-only** (`internetIngestionEnabled/internetQueryEnabled: false`) -- both paths require AMPLS private endpoints; agents outside the private network stop ingesting, so wire AMPLS first
- **Entra-only** (`localAuthenticationEnabled: false`) -- workspace shared keys stop working as credentials; agents must authenticate with managed identities
- **Workspace-context access** (`allowResourceOnlyPermissions: false`) -- every query requires explicit workspace-level permission instead of riding resource RBAC
- **Immediate purge** (`immediateDataPurgeOn30DaysEnabled: true`) -- removes Azure's recovery grace store beyond the retention window (right-to-erasure compliance)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-observability-rg` | Resource group holding the workspace | `AzureResourceGroup` status outputs |
| `my-regulated-logs` | Workspace name | Your naming convention |

## Related Presets

- **01-pay-as-you-go** -- The default open posture for everyday environments
- **02-commitment-tier** -- Combine with this posture for high-volume regulated estates
