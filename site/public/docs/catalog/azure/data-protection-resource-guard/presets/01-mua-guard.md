---
title: "MUA Guard"
description: "This preset creates the strongest-posture Resource Guard: no exclusions, so every critical vault operation (disabling soft delete, deleting protection, reducing retention) requires an approval..."
type: "preset"
rank: "01"
presetSlug: "01-mua-guard"
componentSlug: "data-protection-resource-guard"
componentTitle: "Data Protection Resource Guard"
provider: "azure"
icon: "package"
order: 1
---

# MUA Guard

This preset creates the strongest-posture Resource Guard: no exclusions, so every critical vault operation (disabling soft delete, deleting protection, reducing retention) requires an approval through the guard. Deploy it in a scope a DIFFERENT administrator controls than the vaults it protects.

## When to Use

- Ransomware-resistant backup estates -- pair with immutable vaults so neither data nor retention can be quietly destroyed
- Separation of duties: backup operations owned by one team, destructive-operation approval by another
- Compliance regimes mandating dual authorization on data-destruction paths

## Key Configuration Choices

- **A security-team resource group** -- the guard's value comes entirely from living in someone else's scope; same-scope guards are speed bumps, not controls
- **No `vaultCriticalOperationExclusionList`** -- empty guards everything; add exclusions only as deliberate, documented exceptions
- **Vaults opt in by reference** -- creating the guard changes nothing until a vault references its `resource_guard_id` output

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<security-team-resource-group>` | An AzureResourceGroup a different administrator controls | The security/governance team's resource group name |

The guard is free -- it is a pure governance object.
