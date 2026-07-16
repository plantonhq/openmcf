---
title: "Commitment-Tier Workspace"
description: "This preset creates a Log Analytics Workspace on the CapacityReservation sku -- Azure's discounted commitment billing for estates with sustained high ingestion. The tier discount grows with the..."
type: "preset"
rank: "02"
presetSlug: "02-commitment-tier"
componentSlug: "log-analytics-workspace"
componentTitle: "Log Analytics Workspace"
provider: "azure"
icon: "package"
order: 2
---

# Commitment-Tier Workspace

This preset creates a Log Analytics Workspace on the CapacityReservation sku -- Azure's discounted commitment billing for estates with sustained high ingestion. The tier discount grows with the reservation level.

## When to Use

- Sustained ingestion at or above ~100 GB/day (the smallest tier)
- Central workspaces aggregating a whole estate's telemetry (SIEM, platform logging)
- Anywhere the pay-as-you-go bill has become predictable enough to commit

## Key Configuration Choices

- **CapacityReservation sku + 100 GB/day tier** -- Azure only sells fixed tiers (100, 200, 300, 400, 500, 1000, 2000, 5000, 10000, 25000, 50000); pick the tier just below observed daily volume -- overage bills at the pay-as-you-go rate
- **31-day commitment** -- raising the tier (or entering the sku) restarts Azure's commitment period; the sku can drop back to PerGB2018 in place after it lapses
- **180-day retention** -- central workspaces usually carry longer compliance windows

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-observability-rg` | Resource group holding the workspace | `AzureResourceGroup` status outputs |
| `my-central-logs` | Workspace name | Your naming convention |

## Related Presets

- **01-pay-as-you-go** -- The everyday shape below the commitment threshold
- **03-private-hardened** -- Add the locked-down network/auth posture on top
